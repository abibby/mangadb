package jobs

import (
	"context"
	"log/slog"
	"math/rand/v2"

	"abibby.com/mangadb/app/events"
	"abibby.com/mangadb/services/datasource"
	"abibby.com/mangadb/services/datasource/anilist"
	"abibby.com/mangadb/services/datasource/mangadex"
	"abibby.com/mangadb/services/datasource/mangaplus"
	"abibby.com/mangadb/services/datasource/viz"
	"github.com/jmoiron/sqlx"
	"gosalusa.com/cache"
	"gosalusa.com/event"
)

type FetchSeries struct {
	Logger *slog.Logger `inject:""`
	DB     *sqlx.DB     `inject:""`

	AnilistClient *anilist.Client   `inject:""`
	MDClient      *mangadex.Client  `inject:""`
	MPClient      *mangaplus.Client `inject:""`
	VizClient     *viz.Client       `inject:""`
	Dispatch      event.Dispatch    `inject:""`
	Cache         cache.MemoryCache `inject:""`
}

func (m *FetchSeries) Handle(ctx context.Context, e *events.FetchSeriesEvent) error {
	idb := make([]byte, 8)
	for i := range idb {
		idb[i] = byte(rand.IntN(26) + 'a')
	}
	id := string(idb)
	m.Logger = m.Logger.With("run_id", id)
	for _, s := range e.Series {
		err := m.fetchSeries(ctx, id, s)
		if err != nil {
			return err
		}
	}
	err := m.Dispatch(ctx, &events.UpdateViewsEvent{})
	if err != nil {
		return err
	}
	return m.Dispatch(ctx, &events.FetchImagesEvent{})
}
func (m *FetchSeries) fetchSeries(ctx context.Context, id string, e *datasource.Series) error {
	cachedID, err := m.Cache.GetOrCreate("fetchSeries:"+e.Source+":"+e.ID, func() []byte {
		return []byte(id)
	})
	if err != nil {
		return err
	}
	if string(cachedID) != id {
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
