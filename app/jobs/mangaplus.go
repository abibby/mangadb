package jobs

import (
	"context"
	"log/slog"

	"abibby.com/mangadb/app/events"
	"abibby.com/mangadb/services/datasource/mangaplus"
	"github.com/jmoiron/sqlx"
)

type MangaPlus struct {
	Logger   *slog.Logger      `inject:""`
	MPClient *mangaplus.Client `inject:""`
	DB       *sqlx.DB          `inject:""`
}

func (m *MangaPlus) Handle(ctx context.Context, e *events.MangaPlusSeries) error {
	m.Logger.Info("starting sync")
	err := m.MPClient.Series(ctx, m.DB, e.ID)
	if err != nil {
		return err
	}
	m.Logger.Info("starting finished")
	return nil
}
