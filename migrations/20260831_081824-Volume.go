package migrations

import (
	"abibby.com/salusa/database/migrate"
	"abibby.com/salusa/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20260831_081824-Volume",
		Up: schema.Create("volumes", func(table *schema.Blueprint) {
			table.Blob("id").Primary()
			table.String("mangadex_id")
			table.DateTime("mangadex_updated")
			table.String("manga_plus_id")
			table.DateTime("manga_plus_updated")
		}),
		Down: schema.DropIfExists("volumes"),
	})
}
