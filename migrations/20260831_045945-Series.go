package migrations

import (
	"abibby.com/salusa/database/migrate"
	"abibby.com/salusa/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20260831_045945-Series",
		Up: schema.Create("series", func(table *schema.Blueprint) {
			table.Blob("id").Primary()
			table.String("title")
			table.JSON("aliases")
			table.String("description")
			table.DateTime("start_date")
			table.JSON("genre")
			table.JSON("tags")
			table.String("mangadex_id")
			table.DateTime("mangadex_updated")
			table.String("manga_plus_id")
			table.DateTime("manga_plus_updated")
		}),
		Down: schema.DropIfExists("series"),
	})
}
