package app

import (
	"context"

	"abibby.com/mangadb/app/events"
	"abibby.com/mangadb/app/jobs"
	"abibby.com/mangadb/app/models"
	"abibby.com/mangadb/app/providers"
	"abibby.com/mangadb/config"
	"abibby.com/mangadb/migrations"
	"abibby.com/mangadb/resources"
	"abibby.com/mangadb/routes"
	"abibby.com/mangadb/services/datasource/anilist"
	"abibby.com/mangadb/services/datasource/mangadex"
	"abibby.com/mangadb/services/datasource/mangaplus"
	"abibby.com/mangadb/services/datasource/viz"
	"abibby.com/salusa/auth"
	"abibby.com/salusa/clog"
	"abibby.com/salusa/database"
	"abibby.com/salusa/email"
	"abibby.com/salusa/event"
	"abibby.com/salusa/event/cron"
	"abibby.com/salusa/kernel"
	"abibby.com/salusa/openapidoc"
	"abibby.com/salusa/openapidoc/openapidocdi"
	"abibby.com/salusa/pubsub/channelpubsub"
	"abibby.com/salusa/request"
	"abibby.com/salusa/view"
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
			mangaplus.Register(ctx)
			viz.Register(ctx)
			anilist.Register(ctx)
		}),
	),
	kernel.Services(
		cron.Service().
			Schedule("0 * * * *", &events.UpdateViewsEvent{}),
		event.Service(
			event.NewListener[*jobs.FetchSeries](),
			event.NewListener[*jobs.UpdateViews](),
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
