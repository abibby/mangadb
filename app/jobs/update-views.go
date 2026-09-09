package jobs

import (
	"context"
	"log/slog"
	"sync"

	"github.com/abibby/icbmdb/app/events"
	"github.com/abibby/icbmdb/services/datasource/viz"
	"github.com/jmoiron/sqlx"
)

type UpdateViews struct {
	Logger    *slog.Logger `inject:""`
	VizClient *viz.Client  `inject:""`
	DB        *sqlx.DB     `inject:""`
}

func (m *UpdateViews) Handle(ctx context.Context, e *events.UpdateViewsEvent) error {
	m.Logger.Info("Started updating views")

	views := []string{
		"int_mangadex_series",
		"int_mangaplus_series",
		"int_viz_series",
		"int_anilist_series",
		"int_mangadex_chapter",
		"int_mangaplus_chapter",
		"int_viz_chapter",
	}

	wg := sync.WaitGroup{}
	for _, v := range views {
		wg.Go(func() {
			_, err := m.DB.ExecContext(ctx, "REFRESH MATERIALIZED VIEW CONCURRENTLY "+v)
			if err != nil {
				m.Logger.Warn("Failed to update view", "view", v, "error", err)
			}
		})
	}
	wg.Wait()

	m.Logger.Info("Finished updating views")
	return nil
}
