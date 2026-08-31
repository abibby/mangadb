package main

import (
	"context"
	"os"
	"time"

	"abibby.com/salusa/clog"
	"abibby.com/salusa/di"
	"abibby.com/salusa/event"
	"github.com/abibby/icbmdb/app"
	"github.com/abibby/icbmdb/app/events"
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
		time.Sleep(time.Second)

		dispatch, err := di.Resolve[event.Dispatch](ctx)
		if err != nil {
			panic(err)
		}

		err = dispatch(ctx, &events.MangadexSeries{
			ID: "32505911-558f-4ce0-9eed-0b4538f29efc",
		})
		if err != nil {
			panic(err)
		}
	}()

	err = app.Kernel.Run(ctx)
	if err != nil {
		clog.Use(ctx).Error("error running", "error", err)
		os.Exit(1)
	}
}
