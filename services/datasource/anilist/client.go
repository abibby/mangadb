package anilist

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"abibby.com/salusa/database"
	"abibby.com/salusa/di"
	"abibby.com/salusa/jsonio"
	"github.com/abibby/icbmdb/app/models"
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

	u := fmt.Sprintf("https://graphql.anilist.co Media(id: %s)", id)
	err := m.request(seriesQuery, map[string]any{
		"mediaId": id,
	}, &m.buffer)
	if err != nil {
		return err
	}

	return models.ApiResponseCreateOrUpdate(ctx, tx, &models.APIResponse{
		Source:         "anilist",
		SourceSeriesID: id,
		DataType:       "series",
		URL:            u,
		Page:           0,
		SyncJobID:      "",
		RawPayload:     m.buffer.Bytes(),
	})
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

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")

	resp, err := m.httpClient.Do(r)
	if err != nil {
		return fmt.Errorf("request failed: query %s: %w", query, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("request failed: query %s: status %s", query, resp.Status)
	}

	_, err = io.Copy(w, resp.Body)
	return err
}
