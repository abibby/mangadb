package mangadex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sync"
	"time"

	"abibby.com/mangadb/app/models"
	"abibby.com/mangadb/services/datasource"
	"github.com/abibby/mangadexv5"
	"go.uber.org/ratelimit"
	"gosalusa.com/database"
	"gosalusa.com/di"
)

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

type PaginatedResponse struct {
	Result   string `json:"result"`
	Response string `json:"response"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
	Total    int    `json:"total"`
}

type SeriesResponse struct {
	Data []MDSeries `json:"data"`
}

type MDSeries struct {
	ID         string             `json:"id"`
	Attributes MDSeriesAttributes `json:"attributes"`
}

type MDSeriesAttributes struct {
	Links map[string]string `json:"links"`
}

// Series implements [datasource.Datasource].
func (m *Client) Series(ctx context.Context, tx database.DB, id string) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.buffer.Reset()

	u := fmt.Sprintf("https://api.mangadex.org/manga?ids[]=%s", id)
	err := m.request(http.MethodGet, u, &m.buffer)
	if err != nil {
		return err
	}

	r := &SeriesResponse{}

	err = json.Unmarshal(m.buffer.Bytes(), r)
	if err != nil {
		return err
	}

	if len(r.Data) == 0 {
		return fmt.Errorf("mangadex: no series in response")
	}
	s := r.Data[0]
	ids := map[string]string{
		datasource.MangadexSource: s.ID,
	}
	for k, v := range s.Attributes.Links {
		switch k {
		case "al":
			ids[datasource.AnilistSource] = v
		case "engtl":
			e, ok := datasource.Get(v)
			if ok {
				ids[e.Source] = e.ID
			}
		}
	}
	err = models.IDMapCreate(ctx, tx, datasource.MangadexQuality, ids)
	if err != nil {
		return err
	}

	return models.ApiResponseCreateOrUpdate(ctx, tx, &models.APIResponse{
		Source:         datasource.MangadexSource,
		SourceSeriesID: id,
		DataType:       "series",
		URL:            u,
		Page:           0,
		SyncJobID:      "",
		RawPayload:     m.buffer.Bytes(),
	})
}

// Series implements [datasource.Datasource].
func (m *Client) Chapters(ctx context.Context, tx database.DB, id string) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	mangadexv5.NewClient()
	limit := 100
	offset := 0
	total := math.MaxInt

	fails := 0

	for offset < total {
		m.buffer.Reset()

		u := fmt.Sprintf("https://api.mangadex.org/manga/%s/feed?limit=%d&offset=%d", id, limit, offset)
		err := m.request(http.MethodGet, u, &m.buffer)
		if err != nil {
			fails++
			if fails > 5 {
				return err
			}
			time.Sleep(time.Second ^ time.Duration(fails))
			fmt.Println(time.Second ^ time.Duration(fails))
			continue
		}
		fails = 0

		err = models.ApiResponseCreateOrUpdate(ctx, tx, &models.APIResponse{
			Source:         datasource.MangadexSource,
			SourceSeriesID: id,
			DataType:       "chapter_list",
			URL:            u,
			SyncJobID:      "",
			RawPayload:     m.buffer.Bytes(),
			Page:           offset / limit,
		})
		if err != nil {
			return err
		}

		p := &PaginatedResponse{}

		err = json.Unmarshal(m.buffer.Bytes(), p)
		if err != nil {
			return err
		}

		total = p.Total
		offset += limit
	}
	return nil
}

func (m *Client) request(method, url string, w io.Writer) error {
	m.limiter.Take()

	fmt.Println("REQUEST " + url)

	r, err := http.NewRequest(method, url, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create request: url %s: %w", url, err)
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")

	resp, err := m.httpClient.Do(r)
	if err != nil {
		return fmt.Errorf("request failed: url %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("request failed: url %s: status %s", url, resp.Status)
	}

	_, err = io.Copy(w, resp.Body)
	return err
}
