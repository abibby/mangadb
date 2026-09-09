package models

import (
	"context"
	"time"

	"abibby.com/mangadb/app/providers"
	"abibby.com/salusa/database/builder"
	"abibby.com/salusa/database/jsoncolumn"
	"abibby.com/salusa/database/model/mixins"
	"abibby.com/salusa/database/model/modeldi"
)

type Series struct {
	BaseModel

	mixins.SoftDelete
	mixins.Timestamps

	Title       string                   `json:"title"       db:"title"`
	Aliases     jsoncolumn.Slice[string] `json:"aliases"     db:"aliases"`
	Description string                   `json:"description" db:"description"`
	StartDate   time.Time                `json:"start_date"  db:"start_date"`
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
