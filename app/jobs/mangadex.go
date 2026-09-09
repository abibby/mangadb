package jobs

import (
	"context"
	"log/slog"

	"abibby.com/mangadb/app/events"
	"abibby.com/mangadb/services/datasource/mangadex"
	"github.com/jmoiron/sqlx"
)

type Mangadex struct {
	Logger   *slog.Logger     `inject:""`
	MDClient *mangadex.Client `inject:""`
	DB       *sqlx.DB         `inject:""`
}

func (m *Mangadex) Handle(ctx context.Context, e *events.MangadexSeries) error {
	m.Logger.Info("starting sync")
	err := m.MDClient.Series(ctx, m.DB, e.ID)
	if err != nil {
		return err
	}
	err = m.MDClient.Chapters(ctx, m.DB, e.ID)
	if err != nil {
		return err
	}

	m.Logger.Info("starting finished")
	return nil
}
