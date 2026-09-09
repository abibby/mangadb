package jobs

import (
	"context"
	"log/slog"

	"abibby.com/mangadb/app/events"
	"abibby.com/mangadb/services/datasource/anilist"
	"github.com/jmoiron/sqlx"
)

type Anilist struct {
	Logger        *slog.Logger    `inject:""`
	AnilistClient *anilist.Client `inject:""`
	DB            *sqlx.DB        `inject:""`
}

func (m *Anilist) Handle(ctx context.Context, e *events.AnilistSeries) error {
	m.Logger.Info("starting sync")
	err := m.AnilistClient.Series(ctx, m.DB, e.ID)
	if err != nil {
		return err
	}
	m.Logger.Info("starting finished")
	return nil
}
