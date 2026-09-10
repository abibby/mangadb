package handlers

import (
	"context"

	"abibby.com/mangadb/app/models"
	"abibby.com/salusa/database"
	"abibby.com/salusa/request"
	"github.com/jmoiron/sqlx"
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
		return models.SeriesQuery(r.Ctx).Get(tx)
	})
	if err != nil {
		return nil, err
	}
	return &SeriesListResponse{
		Series: series,
	}, nil
})
