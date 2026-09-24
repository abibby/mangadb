package models

import (
	"context"
	"time"

	"abibby.com/mangadb/app/providers"
	"gosalusa.com/clog"
	"gosalusa.com/database"
	"gosalusa.com/database/builder"
	"gosalusa.com/database/model"
	"gosalusa.com/database/model/mixins"
	"gosalusa.com/database/model/modeldi"
	"gosalusa.com/extra/sets"
	"gosalusa.com/stream"
)

//go:generate spice generate:migration
type Image struct {
	BaseModel

	mixins.Timestamps

	SourceURL string     `json:"source_url" db:"source_url,unique"`
	FetchedAt *time.Time `json:"fetched_at" db:"fetched_at"`
	Path      string     `json:"path"       db:"path"`
	Format    string     `json:"format"     db:"format"`
	Width     int        `json:"width"      db:"width"`
	Height    int        `json:"height"     db:"height"`
}

func init() {
	providers.Add(modeldi.Register[*Image])
}

func ImagesQuery(ctx context.Context) *builder.ModelBuilder[*Image] {
	return builder.From[*Image]().WithContext(ctx)
}

func CreateImages(ctx context.Context, tx database.DB, imageURLs ...string) error {
	urls := stream.Of(imageURLs).
		Filter(func(u string) bool {
			return u != ""
		})

	existing, err := ImagesQuery(ctx).
		WhereIn(
			"source_url",
			urls.Map(func(u string) any {
				return u
			}).Slice(),
		).
		Get(tx)
	if err != nil {
		return err
	}
	existingSet := sets.New(stream.Of(existing).Map(func(i *Image) string {
		return i.SourceURL
	}).Slice()...)

	images := urls.
		Filter(func(s string) bool {
			return !existingSet.Has(s)
		}).
		Map(func(u string) *Image {
			return &Image{SourceURL: u}
		})

	for i := range images.All() {
		err := model.SaveContext(ctx, tx, i)
		if err != nil {
			clog.Use(ctx).Warn("Add image failed", "error", err)
		}
	}
	return nil
}
