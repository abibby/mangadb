package anilist

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"abibby.com/mangadb/app/models"
	"abibby.com/mangadb/services/datasource"
	"abibby.com/mangadb/version"
	"abibby.com/salusa/database"
	"abibby.com/salusa/di"
	"abibby.com/salusa/jsonio"
	"go.uber.org/ratelimit"
)

const seriesQuery = `
query Series($mediaId: Int) {
  Media(id: $mediaId) {
    id
    idMal
    title {
      romaji
      english
      native
      userPreferred
    }
    type
    format
    status
    description
    startDate {
      year
      month
      day
    }
    endDate {
      year
      month
      day
    }
    season
    seasonYear
    episodes
    duration
    chapters
    volumes
    countryOfOrigin
    isLicensed
    source
    hashtag
    trailer {
      id
      site
      thumbnail
    }
    updatedAt
    coverImage {
      extraLarge
      large
      medium
      color
    }
    bannerImage
    genres
    synonyms
    averageScore
    meanScore
    popularity
    isLocked
    trending
    favourites
    tags {
      id
      name
      description
      category
      rank
      isGeneralSpoiler
      isMediaSpoiler
      isAdult
      userId
    }
    isFavourite
    isFavouriteBlocked
    isAdult
    siteUrl
    autoCreateForumThread
    isRecommendationBlocked
    isReviewBlocked
    modNotes
    externalLinks {
      id
      url
      site
      siteId
      type
      language
      color
      icon
      notes
      isDisabled
    }
    staff {
      edges {
        node {
          id
          name {
            first
            middle
            last
            full
            native
            alternative
            userPreferred
          }
        }
        role
        id
      }
    }
  }
}
`

type Client struct {
	httpClient *http.Client
	buffer     bytes.Buffer

	mtx     sync.Mutex
	limiter ratelimit.Limiter
}

type GraphQLResponse struct {
	Data struct {
		Media Media `json:"media"`
	} `json:"data"`
}

type Media struct {
	ID            int            `json:"id"`
	ExternalLinks []ExternalLink `json:"externalLinks"`
}

type ExternalLink struct {
	URL      string `json:"url"`
	Language string `json:"language"`
}

func New() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
		limiter: ratelimit.New(20, ratelimit.Per(time.Minute)),
	}
}

func Register(ctx context.Context) {
	di.RegisterLazySingleton(ctx, func() (*Client, error) {
		return New(), nil
	})
}

// Series implements [datasource.Datasource].
func (m *Client) Series(ctx context.Context, tx database.DB, id string) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.buffer.Reset()

	err := m.request(seriesQuery, map[string]any{"mediaId": id}, &m.buffer)
	if err != nil {
		return err
	}

	r := &GraphQLResponse{}

	err = json.Unmarshal(m.buffer.Bytes(), r)
	if err != nil {
		return err
	}

	s := r.Data.Media
	ids := map[string]string{
		datasource.AnilistSource: fmt.Sprint(s.ID),
	}
	for _, l := range s.ExternalLinks {
		if l.Language != "English" {
			continue
		}
		datasource.AddFromURL(l.URL, ids)
	}
	err = models.IDMapCreate(ctx, tx, datasource.AnilistQuality, ids)
	if err != nil {
		return err
	}
	return models.ApiResponseCreateOrUpdate(ctx, tx, &models.APIResponse{
		Source:         datasource.AnilistSource,
		SourceSeriesID: id,
		DataType:       "series",
		URL:            fmt.Sprintf("https://graphql.anilist.co Media(id: %s)", id),
		Page:           0,
		SyncJobID:      "",
		RawPayload:     m.buffer.Bytes(),
	})
}

type GraphqlError struct {
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func (e *GraphqlError) Error() string {
	if len(e.Errors) == 0 {
		return "unknown error"
	}
	return e.Errors[0].Message
}

func (m *Client) request(query string, variables map[string]any, w io.Writer) error {
	m.limiter.Take()

	r, err := http.NewRequest(http.MethodPost, "https://graphql.anilist.co/", jsonio.NewReader(map[string]any{
		"query":     query,
		"variables": variables,
	}))
	if err != nil {
		return fmt.Errorf("failed to create request: query %s: %w", query, err)
	}

	r.Header.Add("User-Agent", "mangadb "+version.Version)
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Origin", "https://studio.apollographql.com")
	r.Header.Add("Referer", "https://studio.apollographql.com")

	resp, err := m.httpClient.Do(r)
	if err != nil {
		return fmt.Errorf("request failed: query %s: %w", query, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		e := &GraphqlError{}
		jsonErr := json.NewDecoder(resp.Body).Decode(e)
		if jsonErr != nil {
			return fmt.Errorf("failed to unmarshal graphql error: %w", jsonErr)
		}
		return e
	}

	_, err = io.Copy(w, resp.Body)
	return err
}
