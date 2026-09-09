package models

import (
	"context"

	"abibby.com/mangadb/app/providers"
	"abibby.com/salusa/database"
	"abibby.com/salusa/database/builder"
	"abibby.com/salusa/database/model"
	"abibby.com/salusa/database/model/modeldi"
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
}

func init() {
	providers.Add(modeldi.Register[*IDMap])
}

func IDMapQuery(ctx context.Context) *builder.ModelBuilder[*IDMap] {
	return builder.From[*IDMap]().WithContext(ctx)
}

func IDMapCreateOrUpdate(ctx context.Context, tx database.DB, r *IDMap) error {
	existing, err := IDMapQuery(ctx).Where("source", "=", r.Source).Where("source_series_id", "=", r.SourceSeriesID).First(tx)
	if err != nil {
		return err
	}
	changes := false
	if existing == nil {
		existing = r
		changes = true
	} else {
		if existing.MatchQuality < r.MatchQuality {
			existing.SeriesID = r.SeriesID
			existing.MatchQuality = r.MatchQuality
			changes = true
		}
		if existing.Title != r.Title {
			existing.Title = r.Title
			changes = true
		}
		if existing.Author != r.Author {
			existing.Author = r.Author
			changes = true
		}
	}
	if !changes {
		return nil
	}
	return model.SaveContext(ctx, tx, r)
}
