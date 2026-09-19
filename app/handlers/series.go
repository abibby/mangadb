package handlers

import (
	"context"
	"fmt"

	"abibby.com/mangadb/app/jobs"
	"abibby.com/mangadb/app/models"
	"abibby.com/mangadb/services/datasource"
	"github.com/jmoiron/sqlx"
	"gosalusa.com/database"
	"gosalusa.com/request"
)

type SeriesListRequest struct {
	Read database.Read   `inject:""`
	Ctx  context.Context `inject:""`
}
type SeriesListResponse struct {
	Series []*models.Series `json:"series"`
}

var SeriesList = request.Handler(func(r *SeriesListRequest) (*SeriesListResponse, error) {
	series, err := database.Value(r.Read, func(tx *sqlx.Tx) ([]*models.Series, error) {
		return models.SeriesQuery(r.Ctx).With("IDMaps").Get(tx)
	})
	if err != nil {
		return nil, err
	}
	return &SeriesListResponse{
		Series: series,
	}, nil
})

type SeriesViewRequest struct {
	ID   string          `path:"series_id"`
	Read database.Read   `inject:""`
	Ctx  context.Context `inject:""`
}
type SeriesViewResponse struct {
	Series *models.Series `json:"series"`
}

var SeriesView = request.Handler(func(r *SeriesViewRequest) (*SeriesViewResponse, error) {
	series, err := database.Value(r.Read, func(tx *sqlx.Tx) (*models.Series, error) {
		return models.SeriesQuery(r.Ctx).With("IDMaps").Find(tx, r.ID)
	})
	if err != nil {
		return nil, err
	}
	return &SeriesViewResponse{
		Series: series,
	}, nil
})

type SeriesImportRequest struct {
	URL string `json:"url"`

	Update database.Update `inject:""`
	Ctx    context.Context `inject:""`

	jobs.FetchSeries `inject:""`
}
type SeriesImportResponse struct {
	SeriesID string `json:"series_id"`
}

var SeriesImport = request.Handler(func(r *SeriesImportRequest) (*SeriesImportResponse, error) {
	e, ok := datasource.Get(r.URL)
	if !ok {
		return nil, fmt.Errorf("site not supported: %s", r.URL)
	}
	err := r.FetchSeries.Handle(r.Ctx, e)
	if err != nil {
		return nil, err
	}

	m, err := database.Value(r.Update, func(tx *sqlx.Tx) (*models.IDMap, error) {
		return models.IDMapQuery(r.Ctx).Where("source", "=", e.Source).Where("source_series_id", "=", e.ID).First(tx)
	})
	if err != nil {
		return nil, err
	}

	return &SeriesImportResponse{
		SeriesID: m.SeriesID,
	}, nil
})
