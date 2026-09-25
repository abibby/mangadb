package handlers

import (
	"context"
	"fmt"
	"time"

	"abibby.com/mangadb/app/events"
	"abibby.com/mangadb/app/models"
	"abibby.com/mangadb/services/datasource"
	"github.com/jmoiron/sqlx"
	"gosalusa.com/database"
	"gosalusa.com/event"
	"gosalusa.com/extra/sets"
	"gosalusa.com/request"
)

type SeriesListRequest struct {
	Query string `query:"q"`
	URL   string `query:"url"`

	Read database.Read   `inject:""`
	Ctx  context.Context `inject:""`
}
type SeriesListResponse struct {
	Series []*models.Series `json:"series"`
}

var SeriesList = request.Handler(func(r *SeriesListRequest) (*SeriesListResponse, error) {
	series, err := database.Value(r.Read, func(tx *sqlx.Tx) ([]*models.Series, error) {
		q := models.SeriesQuery(r.Ctx).
			With("IDMaps")
		if r.Query != "" {
			q.Where("title", "%", r.Query).
				OrderByRaw(`similarity("title", $1) DESC`)
		}
		if r.URL != "" {
			e, ok := datasource.Get(r.URL)
			if !ok {
				return nil, fmt.Errorf("site not supported: %s", r.URL)
			}
			q.Where("id", "=", models.IDMapQuery(r.Ctx).Select("series_id").Where("source", "=", e.Source).Where("source_series_id", "=", e.ID))
		}

		return q.Get(tx)
	})
	if err != nil {
		return nil, err
	}
	return &SeriesListResponse{
		Series: series,
	}, nil
})

type SeriesViewRequest struct {
	ID   string   `path:"series_id"`
	With []string `query:"with"`

	Read database.Read   `inject:""`
	Ctx  context.Context `inject:""`
}
type SeriesViewResponse struct {
	Series *models.Series `json:"series"`
}

var SeriesView = request.Handler(func(r *SeriesViewRequest) (*SeriesViewResponse, error) {
	withSet := sets.New(r.With...)
	series, err := database.Value(r.Read, func(tx *sqlx.Tx) (*models.Series, error) {
		q := models.SeriesQuery(r.Ctx).With("IDMaps")
		if withSet.Has("chapters") {
			q.With("Chapters")
		}
		if withSet.Has("volumes") {
			q.With("Volumes")
		}
		return q.Find(tx, r.ID)
	})
	if err != nil {
		return nil, err
	}
	if series == nil {
		return nil, request.ErrStatusNotFound
	}
	return &SeriesViewResponse{
		Series: series,
	}, nil
})

type SeriesImportRequest struct {
	URL string `json:"url"`

	Read database.Read   `inject:""`
	Ctx  context.Context `inject:""`

	Dispatch event.Dispatch `inject:""`
}
type SeriesImportResponse struct {
	SeriesID string `json:"series_id"`
}

var SeriesImport = request.Handler(func(r *SeriesImportRequest) (*SeriesImportResponse, error) {
	e, ok := datasource.Get(r.URL)
	if !ok {
		return nil, fmt.Errorf("site not supported: %s", r.URL)
	}
	err := r.Dispatch(r.Ctx, &events.FetchSeriesEvent{
		Series: []*datasource.Series{e},
	})
	if err != nil {
		return nil, err
	}

	for i := 0; i < 10; i++ {
		m, err := database.Value(r.Read, func(tx *sqlx.Tx) (*models.IDMap, error) {
			return models.IDMapQuery(r.Ctx).Where("source", "=", e.Source).Where("source_series_id", "=", e.ID).First(tx)
		})
		if err != nil {
			return nil, err
		}
		if m != nil {
			return &SeriesImportResponse{
				SeriesID: m.SeriesID,
			}, nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return &SeriesImportResponse{}, nil
})

type SeriesLookupRequest struct {
	URL string `query:"url"`

	Read database.Read   `inject:""`
	Ctx  context.Context `inject:""`
}
type SeriesLookupResponse struct {
	SeriesID *string `json:"series_id"`
}

var SeriesLookup = request.Handler(func(r *SeriesLookupRequest) (*SeriesLookupResponse, error) {
	e, ok := datasource.Get(r.URL)
	if !ok {
		return &SeriesLookupResponse{
			SeriesID: nil,
		}, nil
	}

	m, err := database.Value(r.Read, func(tx *sqlx.Tx) (*models.IDMap, error) {
		return models.IDMapQuery(r.Ctx).Where("source", "=", e.Source).Where("source_series_id", "=", e.ID).First(tx)
	})
	if err != nil {
		return nil, err
	}
	if m == nil {
		return &SeriesLookupResponse{
			SeriesID: nil,
		}, nil
	}

	return &SeriesLookupResponse{
		SeriesID: &m.SeriesID,
	}, nil
})
