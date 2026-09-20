package migrations

import (
	"gosalusa.com/database/migrate"
	"gosalusa.com/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20260920_115349-Images",
		Up: schema.Create("images", func(table *schema.Blueprint) {
			table.Blob("id").Primary()
			table.DateTime("created_at")
			table.DateTime("updated_at")
			table.String("source_url").Unique()
			table.DateTime("fetched_at").Nullable()
			table.String("path")
			table.String("format")
			table.Int("width")
			table.Int("height")
		}),
		Down: schema.DropIfExists("images"),
	})
}
