package models

import (
	"context"

	"abibby.com/mangadb/app/providers"
	"github.com/google/uuid"
	"gosalusa.com/database"
	"gosalusa.com/database/builder"
	"gosalusa.com/database/model/modeldi"
	"gosalusa.com/di"
	"gosalusa.com/nulls/v2"
	"gosalusa.com/router"
)

type Volume struct {
	ViewModel

	SeriesID     string        `json:"series_id" db:"series_id"`
	Number       float32       `json:"number"    db:"number"`
	CoverImageID uuid.NullUUID `json:"-"         db:"cover_image_id"`

	CoverImageURL nulls.Null[string] `json:"cover_image_url"  db:"-"`
}

func init() {
	providers.Add(modeldi.Register[*Volume])
}

func CollectionQuery(ctx context.Context) *builder.ModelBuilder[*Volume] {
	return builder.From[*Volume]().WithContext(ctx)
}

func (v *Volume) AfterLoad(ctx context.Context, tx database.DB) error {
	url, err := di.Resolve[router.URLResolver](ctx)
	if err != nil {
		return nil
	}

	if v.CoverImageID.Valid {
		v.CoverImageURL.Valid = true
		v.CoverImageURL.V = url.Resolve("image.get", "image_id", v.CoverImageID.UUID)
	}
	return nil
}
