package test

import (
	"context"
	"log"

	"abibby.com/mangadb/app"
	"abibby.com/mangadb/config"
	"abibby.com/mangadb/migrations"
	"abibby.com/salusa/database/dbtest"
	"abibby.com/salusa/database/dialects/sqlite"
	"abibby.com/salusa/email/emailtest"
	"abibby.com/salusa/testing/kerneltest"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var runner = dbtest.NewRunner(func() (*sqlx.DB, error) {
	sqlite.UseSQLite()

	db, err := sqlx.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}

	err = migrations.Use().Up(context.Background(), db)
	if err != nil {
		return nil, err
	}

	log.Print("db loaded")

	return db, nil
})

var Run = runner.Run
var RunBenchmark = runner.RunBenchmark

var Kernel = kerneltest.NewTestKernelFactory(app.Kernel, &config.Config{
	Port:     443,
	BasePath: "https://example.test",

	Database: sqlite.NewConfig(":memory:"),
	Mail:     emailtest.NewTestMailerConfig(),
})
