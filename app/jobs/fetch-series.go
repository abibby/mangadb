package jobs

import (
	"context"
	"log/slog"
	"time"

	"abibby.com/mangadb/app/events"
	"abibby.com/mangadb/app/models"
	"abibby.com/mangadb/services/datasource"
	"abibby.com/mangadb/services/datasource/anilist"
	"abibby.com/mangadb/services/datasource/mangadex"
	"abibby.com/mangadb/services/datasource/mangaplus"
	"abibby.com/mangadb/services/datasource/viz"
	"github.com/jmoiron/sqlx"
)

type FetchSeries struct {
	Logger *slog.Logger `inject:""`
	DB     *sqlx.DB     `inject:""`

	AnilistClient *anilist.Client   `inject:""`
	MDClient      *mangadex.Client  `inject:""`
	MPClient      *mangaplus.Client `inject:""`
	VizClient     *viz.Client       `inject:""`
}

func (m *FetchSeries) Handle(ctx context.Context, e *events.FetchSeriesEvent) error {
	exists, err := models.ApiResponseQuery(ctx).
		Where("source", "=", e.Source).
		Where("source_series_id", "=", e.ID).
		Where("updated_at", ">", time.Now().Add(-10*time.Minute).Format(time.RFC3339)).
		Count(m.DB)
	if err != nil {
		return err
	}
	if exists > 0 {
		m.Logger.Info("Fetch series skipped", "source", e.Source, "id", e.ID)
		return nil
	}

	m.Logger.Info("Fetch series start", "source", e.Source, "id", e.ID)
	switch e.Source {
	case datasource.AnilistSource:
		err := m.AnilistClient.Series(ctx, m.DB, e.ID)
		if err != nil {
			return err
		}
	case datasource.MangadexSource:
		err := m.MDClient.Series(ctx, m.DB, e.ID)
		if err != nil {
			return err
		}
		err = m.MDClient.Chapters(ctx, m.DB, e.ID)
		if err != nil {
			return err
		}
		err = m.MDClient.Covers(ctx, m.DB, e.ID)
		if err != nil {
			return err
		}
	case datasource.MangaplusSource:
		err := m.MPClient.Series(ctx, m.DB, e.ID)
		if err != nil {
			return err
		}
	case datasource.VizSource:
		err := m.VizClient.Series(ctx, m.DB, e.ID)
		if err != nil {
			return err
		}
	}
	m.Logger.Info("Fetch series finished", "source", e.Source, "id", e.ID)
	return nil
}
