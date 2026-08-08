package models

import (
	"context"
	"time"

	"github.com/abibby/icbmdb/app/providers"
	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/jsoncolumn"
	"github.com/abibby/salusa/database/model/modeldi"
)

//go:generate spice generate:migration
type Series struct {
	BaseModel

	Title       string                   `json:"title"       db:"title"`
	Aliases     jsoncolumn.Slice[string] `json:"aliases"     db:"aliases"`
	Description string                   `json:"description" db:"description"`
	StartDate   time.Time                `json:"start_date"  db:"start_date"`
	Genre       jsoncolumn.Slice[string] `json:"genre"       db:"genre"`
	Tags        jsoncolumn.Slice[string] `json:"tags"        db:"tags"`

	DatasourceIDs

	Staff   *builder.HasMany[*Staff]  `json:"staff"`
	Issues  *builder.HasMany[*Issue]  `json:"issues"`
	Volumes *builder.HasMany[*Volume] `json:"volumes"`
}

func init() {
	providers.Add(modeldi.Register[*Series])
}

func SeriesQuery(ctx context.Context) *builder.ModelBuilder[*Series] {
	return builder.From[*Series]().WithContext(ctx)
}
