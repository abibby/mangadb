package models

import (
	"context"

	"abibby.com/mangadb/app/providers"
	"gosalusa.com/database/builder"
	"gosalusa.com/database/jsoncolumn"
	"gosalusa.com/database/model/modeldi"
)

type Series struct {
	ViewModel

	// mixins.SoftDelete
	// mixins.Timestamps

	ID          string                   `json:"id"          db:"id"`
	Title       string                   `json:"title"       db:"title"`
	Aliases     jsoncolumn.Slice[string] `json:"aliases"     db:"aliases"`
	Description string                   `json:"description" db:"description"`
	StartDate   any                      `json:"start_date"  db:"start_date"`
	Genre       jsoncolumn.Slice[string] `json:"genre"       db:"genre"`
	Tags        jsoncolumn.Slice[string] `json:"tags"        db:"tags"`

	Staff    *builder.HasMany[*Staff]   `json:"staff"`
	Chapters *builder.HasMany[*Chapter] `json:"chapters"`
	Volumes  *builder.HasMany[*Volume]  `json:"volumes"`
}

func init() {
	providers.Add(modeldi.Register[*Series])
}

func SeriesQuery(ctx context.Context) *builder.ModelBuilder[*Series] {
	return builder.From[*Series]().WithContext(ctx)
}
