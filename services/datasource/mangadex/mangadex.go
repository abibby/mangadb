package mangadex

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"abibby.com/salusa/database"
	"abibby.com/salusa/database/model"
	"abibby.com/salusa/di"
	"github.com/abibby/icbmdb/app/models"
)

type Client struct {
	httpClient *http.Client
	buffer     bytes.Buffer

	mtx sync.Mutex
}

func New() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
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

	u := fmt.Sprintf("https://api.mangadex.org/manga?ids[]=%s", id)
	err := m.request(http.MethodGet, u, &m.buffer)
	if err != nil {
		return err
	}

	r, err := models.ApiResponseQuery(ctx).Where("url", "=", u).First(tx)
	if err != nil {
		return err
	}

	if r == nil {
		r = &models.APIResponse{
			Source:         "mangadex",
			SourceSeriesID: id,
			DataType:       "series",
			URL:            u,
			Page:           0,
		}
	}
	r.SyncJobID = ""
	r.RawPayload = m.buffer.Bytes()
	return model.SaveContext(ctx, tx, r)
}

// Series implements [datasource.Datasource].
func (m *Client) Chapters(ctx context.Context, tx database.DB, id string) error {
	return m.chapters(ctx, tx, id, 0)
}
func (m *Client) chapters(ctx context.Context, tx database.DB, id string, page int) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.buffer.Reset()

	u := fmt.Sprintf("https://api.mangadex.org/manga/%s/feed", id)
	err := m.request(http.MethodGet, u, &m.buffer)
	if err != nil {
		return err
	}

	r, err := models.ApiResponseQuery(ctx).Where("url", "=", u).First(tx)
	if err != nil {
		return err
	}

	if r == nil {
		r = &models.APIResponse{
			Source:         "mangadex",
			SourceSeriesID: id,
			DataType:       "chapter_list",
			URL:            u,
			Page:           0,
		}
	}
	r.SyncJobID = ""
	r.RawPayload = m.buffer.Bytes()
	err = model.SaveContext(ctx, tx, r)
	if err != nil {
		return err
	}
	return nil
}

func (m *Client) request(method, url string, w io.Writer) error {
	r, err := http.NewRequest(method, url, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// if c.token != nil {
	// 	r.Header.Add("Authorization", "Bearer "+c.token.Session)
	// }
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")

	resp, err := m.httpClient.Do(r)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("request failed: status %s", resp.Status)
	}

	_, err = io.Copy(w, resp.Body)
	return err
}
