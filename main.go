package main

import (
	"context"
	"os"
	"time"

	"abibby.com/mangadb/app"
	"abibby.com/mangadb/app/events"
	"abibby.com/salusa/clog"
	"abibby.com/salusa/di"
	"abibby.com/salusa/event"
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
			// &events.MangaPlusSeries{ID: "100171"}, // Dandadan
			// &events.MangaPlusSeries{ID: "100020"}, // One Piece
			// &events.MangaPlusSeries{ID: "100056"}, // SPY x FAMILY
			// &events.MangaPlusSeries{ID: "100269"}, // Boruto: Two Blue Vortex
			// &events.VizSeries{ID: "dandadan"},
			// &events.VizSeries{ID: "one-piece"},
			// &events.AnilistSeries{ID: "30002"}, // Berserk
			&events.AnilistSeries{ID: "132029"}, // Dandadan
			&events.UpdateViewsEvent{},
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
