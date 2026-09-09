package models

import (
	"context"

	"abibby.com/salusa/database/builder"
	"abibby.com/salusa/database/model"
	"abibby.com/salusa/database/model/modeldi"
	"github.com/abibby/icbmdb/app/providers"
)

//go:generate spice generate:migration
type IDMap struct {
	model.BaseModel

	SourceSeriesID string `json:"source_series_id" db:"source_series_id,primary"`
	Source         string `json:"source"           db:"source,primary"`
	SeriesID       string `json:"series_id"        db:"series_id,index"`
	MatchQuality   int    `json:"match_quality"    db:"match_quality"`
	Title          string `json:"title"            db:"title"`
	Author         string `json:"author"           db:"author"`
	SpineQuality   int    `json:"spine_quality"    db:"spine_quality"`
}

func init() {
	providers.Add(modeldi.Register[*IDMap])
}

func IDMapQuery(ctx context.Context) *builder.ModelBuilder[*IDMap] {
	return builder.From[*IDMap]().WithContext(ctx)
}
