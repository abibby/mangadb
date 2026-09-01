package jobs

import (
	"context"
	"log/slog"

	"github.com/abibby/icbmdb/app/events"
	"github.com/abibby/icbmdb/services/datasource/viz"
	"github.com/jmoiron/sqlx"
)

type Viz struct {
	Logger    *slog.Logger `inject:""`
	VizClient *viz.Client  `inject:""`
	DB        *sqlx.DB     `inject:""`
}

func (m *Viz) Handle(ctx context.Context, e *events.VizSeries) error {
	m.Logger.Info("starting sync")
	err := m.VizClient.Series(ctx, m.DB, e.ID)
	if err != nil {
		return err
	}
	// err = m.MPClient.Chapters(ctx, m.DB, e.ID)
	// if err != nil {
	// 	return err
	// }
	m.Logger.Info("starting finished")
	return nil
}
