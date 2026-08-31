package app

import (
	"context"

	"abibby.com/salusa/auth"
	"abibby.com/salusa/clog"
	"abibby.com/salusa/database"
	"abibby.com/salusa/email"
	"abibby.com/salusa/event"
	"abibby.com/salusa/kernel"
	"abibby.com/salusa/openapidoc"
	"abibby.com/salusa/openapidoc/openapidocdi"
	"abibby.com/salusa/pubsub/channelpubsub"
	"abibby.com/salusa/request"
	"abibby.com/salusa/view"
	"github.com/abibby/icbmdb/app/jobs"
	"github.com/abibby/icbmdb/app/models"
	"github.com/abibby/icbmdb/app/providers"
	"github.com/abibby/icbmdb/config"
	"github.com/abibby/icbmdb/migrations"
	"github.com/abibby/icbmdb/resources"
	"github.com/abibby/icbmdb/routes"
	"github.com/abibby/icbmdb/services/datasource/mangadex"
	"github.com/go-openapi/spec"
)

var Kernel = kernel.New(
	kernel.Config(config.Load),
	kernel.Bootstrap(
		view.Register(resources.Content, "**/*.html"),
		providers.Register,
		kernel.Register(func(ctx context.Context, c *config.Config) {
			database.Register(ctx, c.Database, migrations.Use())
			email.Register(ctx, c.Mail)
			channelpubsub.Register(ctx)

			clog.RegisterDefault(ctx)
			request.Register(ctx)
			auth.Register[*models.User](ctx)
			event.Register(ctx)
			openapidocdi.Register(ctx)

			mangadex.Register(ctx)
		}),
	),
	kernel.Services(
		event.Service(
			event.NewListener[*jobs.Mangadex](),
		),
	),
	kernel.InitRoutes(routes.InitRoutes),
	kernel.APIDocumentation(
		openapidoc.Info(spec.InfoProps{
			Title:       "Salusa Example API",
			Description: `This is the API documentaion for the example Salusa application`,
		}),
		openapidoc.BasePath("/api"),
	),
)
