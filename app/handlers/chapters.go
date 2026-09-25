package handlers

import (
	"context"

	"abibby.com/mangadb/app/models"
	"github.com/jmoiron/sqlx"
	"gosalusa.com/database"
	"gosalusa.com/request"
)

type ChapterListRequest struct {
	SeriesID string `path:"series_id"`

	Read database.Read   `inject:""`
	Ctx  context.Context `inject:""`
}
type ChapterListResponse struct {
	Data []*models.Chapter `json:"data"`
}

var ChapterList = request.Handler(func(r *ChapterListRequest) (*ChapterListResponse, error) {
	series, err := database.Value(r.Read, func(tx *sqlx.Tx) ([]*models.Chapter, error) {
		return models.ChapterQuery(r.Ctx).Where("series_id", "=", r.SeriesID).Get(tx)
	})
	if err != nil {
		return nil, err
	}
	return &ChapterListResponse{
		Data: series,
	}, nil
})
