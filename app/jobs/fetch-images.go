package jobs

import (
	"bytes"
	"context"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "golang.org/x/image/webp"

	"abibby.com/mangadb/app/events"
	"abibby.com/mangadb/app/models"
	"abibby.com/mangadb/services/datasource/viz"
	"github.com/jmoiron/sqlx"
	"gosalusa.com/database"
	"gosalusa.com/database/builder"
	"gosalusa.com/database/model"
)

type FetchImages struct {
	Logger    *slog.Logger    `inject:""`
	VizClient *viz.Client     `inject:""`
	Update    database.Update `inject:""`
}

func (m *FetchImages) Handle(ctx context.Context, e *events.FetchImagesEvent) error {
	m.Logger.Info("Started fetching images")
	var errNoRows = errors.New("no rows")
	var err error
	for err == nil {
		err = m.Update(func(tx *sqlx.Tx) error {
			images, err := models.ImagesQuery(ctx).
				Where(
					"id",
					"in",
					models.ImagesQuery(ctx).
						Select("id").
						Where("fetched_at", "=", nil).
						Limit(1),
				).
				ForUpdateSkipLocked().
				UpdateReturning(tx, builder.Updates{
					"fetched_at": time.Now(),
				})
			if err != nil {
				return err
			}

			if len(images) == 0 {
				return errNoRows
			}

			img := images[0]

			r, err := http.Get(img.SourceURL)
			if err != nil {
				return err
			}
			defer r.Body.Close()

			b, err := io.ReadAll(r.Body)
			if err != nil {
				return err
			}

			i, format, err := image.Decode(bytes.NewBuffer(b))
			if err != nil {
				return err
			}

			p := "images/" + img.ID.String() + "." + format
			err = os.WriteFile(p, b, 0644)
			if err != nil {
				return err
			}

			img.Path = p
			img.Format = format
			img.Width = i.Bounds().Dx()
			img.Height = i.Bounds().Dy()

			return model.SaveContext(ctx, tx, img)
		})
	}
	if err == errNoRows {
	} else if err != nil {
		return err
	}

	m.Logger.Info("Finished fetchiong images")
	return nil
}
