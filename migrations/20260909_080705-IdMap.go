package migrations

import (
	"abibby.com/salusa/database/migrate"
	"abibby.com/salusa/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20260909_080705-IdMap",
		Up: schema.Create("id_maps", func(table *schema.Blueprint) {
			table.String("source_series_id")
			table.String("source")
			table.String("series_id").Index()
			table.Int("match_quality")
			table.String("title")
			table.String("author")
			table.Int("spine_quality")
			table.PrimaryKey("source_series_id", "source")
		}),
		Down: schema.DropIfExists("id_maps"),
	})
}
