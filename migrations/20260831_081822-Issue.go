package migrations

import (
	"abibby.com/salusa/database/migrate"
	"abibby.com/salusa/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20260831_081822-Issue",
		Up: schema.Create("issues", func(table *schema.Blueprint) {
			table.Blob("id").Primary()
			table.String("title")
			table.Blob("series_id")
			table.Float("number")
			table.Int("volume")
			table.String("summary")
			table.String("notes")
			table.DateTime("release_date")
			table.String("imprint")
			table.String("web")
			table.String("language_iso")
			table.String("format")
			table.String("mangadex_id")
			table.DateTime("mangadex_updated")
			table.String("manga_plus_id")
			table.DateTime("manga_plus_updated")
		}),
		Down: schema.DropIfExists("issues"),
	})
}
