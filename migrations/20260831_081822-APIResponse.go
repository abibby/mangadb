package migrations

import (
	"abibby.com/salusa/database/dialects"
	"abibby.com/salusa/database/migrate"
	"abibby.com/salusa/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20260831_081822-APIResponse",
		Up: schema.Create("api_responses", func(table *schema.Blueprint) {
			table.Blob("id").Primary()
			table.DateTime("deleted_at").Nullable()
			table.DateTime("created_at")
			table.DateTime("updated_at")
			table.String("sync_job_id")
			table.String("source")
			table.String("source_series_id")
			table.String("data_type")
			table.String("url")
			table.Int("page")
			table.OfType(dialects.DataType{Name: "jsonb"}, "raw_payload")
		}),
		Down: schema.DropIfExists("api_responses"),
	})
}
