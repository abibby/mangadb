package main

import (
	"context"
	"os"
	"time"

	"abibby.com/mangadb/app"
	"abibby.com/mangadb/app/events"
	"gosalusa.com/clog"
	"gosalusa.com/di"
	"gosalusa.com/event"
)

func main() {

	ctx := di.ContextWithDependencyProvider(
		context.Background(),
		di.NewDependencyProvider(),
	)

	err := app.Kernel.Bootstrap(ctx)
	if err != nil {
		clog.Use(ctx).Error("error bootstrapping", "error", err)
		os.Exit(1)
	}
	go func() {
		time.Sleep(100 * time.Millisecond)

		dispatch, err := di.Resolve[event.Dispatch](ctx)
		if err != nil {
			panic(err)
		}
		es := []event.Event{
			// &events.MangadexSeries{ID: "32505911-558f-4ce0-9eed-0b4538f29efc"}, // Days Off in the Dragon's Stomach
			// &events.MangadexSeries{ID: "68112dc1-2b80-4f20-beb8-2f2a8716a430"}, // Dandadan
			// &events.MangadexSeries{ID: "801513ba-a712-498c-8f57-cae55b38cc92"}, // Berserk
			// &events.MangadexSeries{ID: "a1c7c817-4e59-43b7-9365-09675a149a6f"}, // One Piece
			// &events.FetchSeriesEvent{Source: datasource.MangadexSource, ID: "239d6260-d71f-43b0-afff-074e3619e3de"}, // Bleach

			// &events.MangaPlusSeries{ID: "100171"}, // Dandadan
			// &events.MangaPlusSeries{ID: "100020"}, // One Piece
			// &events.MangaPlusSeries{ID: "100056"}, // SPY x FAMILY
			// &events.MangaPlusSeries{ID: "100269"}, // Boruto: Two Blue Vortex
			// &events.VizSeries{ID: "dandadan"},
			// &events.VizSeries{ID: "one-piece"},
			// &events.FetchSeriesEvent{Source: datasource.AnilistSource, ID: "30002"}, // Berserk
			// &events.FetchSeriesEvent{Source: datasource.AnilistSource, ID: "132029"}, // Dandadan
			// &events.FetchSeriesEvent{Source: datasource.AnilistSource, ID: "30013"}, // One Piece
			&events.UpdateViewsEvent{},
			&events.FetchImagesEvent{},
		}
		for _, e := range es {
			err = dispatch(ctx, e)
			if err != nil {
				panic(err)
			}
		}
	}()

	err = app.Kernel.Run(ctx)
	if err != nil {
		clog.Use(ctx).Error("error running", "error", err)
		os.Exit(1)
	}
}
