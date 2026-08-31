package jobs

import (
	"context"
	"log/slog"

	"abibby.com/salusa/database"
	"github.com/abibby/icbmdb/app/events"
	"github.com/abibby/icbmdb/services/datasource/mangadex"
	"github.com/jmoiron/sqlx"
)

type Mangadex struct {
	Logger   *slog.Logger     `inject:""`
	MDClient *mangadex.Client `inject:""`
	Update   database.Update  `inject:""`
}

func (m *Mangadex) Handle(ctx context.Context, e *events.MangadexSeries) error {
	m.Logger.Info("starting sync")
	return m.Update(func(tx *sqlx.Tx) error {
		err := m.MDClient.Series(ctx, tx, e.ID)
		if err != nil {
			return err
		}
		err = m.MDClient.Chapters(ctx, tx, e.ID)
		if err != nil {
			return err
		}
		return nil
	})
}
