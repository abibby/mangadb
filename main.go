package main

import (
	"context"
	"os"

	"abibby.com/salusa/clog"
	"abibby.com/salusa/di"
	"github.com/abibby/icbmdb/app"
	"github.com/abibby/icbmdb/services/seriesupdate"
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

	err = seriesupdate.Update()
	if err != nil {
		panic(err)
	}
	// err = app.Kernel.Run(ctx)
	// if err != nil {
	// 	clog.Use(ctx).Error("error running", "error", err)
	// 	os.Exit(1)
	// }
}
