package models

import (
	"context"

	"abibby.com/mangadb/app/events"
	"abibby.com/mangadb/app/providers"
	"abibby.com/salusa/database"
	"abibby.com/salusa/database/builder"
	"abibby.com/salusa/database/model"
	"abibby.com/salusa/database/model/modeldi"
	"abibby.com/salusa/di"
	"abibby.com/salusa/event"
	"github.com/google/uuid"
)

//go:generate spice generate:migration
type IDMap struct {
	model.BaseModel

	SourceSeriesID string `json:"source_series_id" db:"source_series_id,primary"`
	Source         string `json:"source"           db:"source,primary"`
	SeriesID       string `json:"series_id"        db:"series_id,index"`
	MatchQuality   int    `json:"match_quality"    db:"match_quality"`
}

func init() {
	providers.Add(modeldi.Register[*IDMap])
}

func IDMapQuery(ctx context.Context) *builder.ModelBuilder[*IDMap] {
	return builder.From[*IDMap]().WithContext(ctx)
}

func IDMapCreate(ctx context.Context, tx database.DB, quality int, m map[string]string) error {
	if len(m) == 0 {
		return nil
	}
	q := IDMapQuery(ctx)

	for k, v := range m {
		q.Or(func(q *builder.Conditions) {
			q.Where("source", "=", k).Where("source_series_id", "=", v)
		})
	}
	idMaps, err := q.Get(tx)
	if err != nil {
		return err
	}

	seriesID := ""

	idMapMap := map[string]*IDMap{}
	for _, idMap := range idMaps {
		idMapMap[idMap.Source] = idMap
		if seriesID == "" {
			seriesID = idMap.SeriesID
		} else if seriesID != idMap.SeriesID {
			panic("do something here")
		}
	}
	if seriesID == "" {
		seriesID = uuid.NewString()
	}

	for k, v := range m {
		idMap, ok := idMapMap[k]
		if ok {
			if idMap.MatchQuality < quality {
				idMap.SeriesID = seriesID
				idMap.MatchQuality = quality
			}
		} else {
			idMap = &IDMap{
				SourceSeriesID: v,
				Source:         k,
				SeriesID:       seriesID,
				MatchQuality:   quality,
			}
		}
		err = model.SaveContext(ctx, tx, idMap)
		if err != nil {
			return err
		}
	}
	dispatch, err := di.Resolve[event.Dispatch](ctx)
	if err != nil {
		return err
	}
	for k, v := range m {
		dispatch(ctx, &events.FetchSeriesEvent{
			Source: k,
			ID:     v,
		})
	}

	return nil
}
