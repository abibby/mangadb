package models

import (
	"context"
	"fmt"
	"time"

	"abibby.com/mangadb/app/providers"
	"gosalusa.com/database"
	"gosalusa.com/database/builder"
	"gosalusa.com/database/model"
	"gosalusa.com/database/model/mixins"
	"gosalusa.com/database/model/modeldi"
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
	return model.InsertMany(
		tx,
		stream.Of(imageURLs).
			Filter(func(u string) bool {
				fmt.Println("url", u)
				return u != ""
			}).
			Map(func(u string) *Image {
				return &Image{SourceURL: u}
			}).
			Slice(),
	)
}
