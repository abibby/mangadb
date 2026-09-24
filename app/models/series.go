package models

import (
	"context"

	"abibby.com/mangadb/app/providers"
	"github.com/google/uuid"
	"gosalusa.com/database"
	"gosalusa.com/database/builder"
	"gosalusa.com/database/jsoncolumn"
	"gosalusa.com/database/model/modeldi"
	"gosalusa.com/di"
	"gosalusa.com/nulls/v2"
	"gosalusa.com/router"
)

type Series struct {
	ViewModel

	// mixins.SoftDelete
	// mixins.Timestamps

	ID            string                   `json:"id"          db:"id"`
	Title         string                   `json:"title"       db:"title"`
	Aliases       jsoncolumn.Slice[string] `json:"aliases"     db:"aliases"`
	Description   string                   `json:"description" db:"description"`
	StartDate     any                      `json:"start_date"  db:"start_date"`
	Genre         jsoncolumn.Slice[string] `json:"genre"       db:"genre"`
	Tags          jsoncolumn.Slice[string] `json:"tags"        db:"tags"`
	CoverImageID  uuid.NullUUID            `json:"-"           db:"cover_image_id"`
	BannerImageID uuid.NullUUID            `json:"-"           db:"banner_image_id"`

	CoverImageURL  nulls.Null[string] `json:"cover_image_url"  db:"-"`
	BannerImageURL nulls.Null[string] `json:"banner_image_url" db:"-"`

	IDMaps   *builder.HasMany[*IDMap]   `json:"ids"                foreign:"series_id"`
	Chapters *builder.HasMany[*Chapter] `json:"chapters,omitempty" foreign:"series_id" `
	Volumes  *builder.HasMany[*Volume]  `json:"volumes,omitempty"  foreign:"series_id" `
}

func init() {
	providers.Add(modeldi.Register[*Series])
}

func SeriesQuery(ctx context.Context) *builder.ModelBuilder[*Series] {
	return builder.From[*Series]().WithContext(ctx)
}

func (s *Series) AfterLoad(ctx context.Context, tx database.DB) error {
	url, err := di.Resolve[router.URLResolver](ctx)
	if err != nil {
		return nil
	}

	if s.CoverImageID.Valid {
		s.CoverImageURL.Valid = true
		s.CoverImageURL.V = url.Resolve("image.get", "image_id", s.CoverImageID.UUID)
	}
	if s.BannerImageID.Valid {
		s.BannerImageURL.Valid = true
		s.BannerImageURL.V = url.Resolve("image.get", "image_id", s.BannerImageID.UUID)
	}
	return nil
}
