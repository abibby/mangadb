package migrations

import (
	"abibby.com/salusa/database/migrate"
	"abibby.com/salusa/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20260831_045946-User",
		Up: schema.Create("users", func(table *schema.Blueprint) {
			table.Blob("id").Primary()
			table.String("email").Unique()
			table.Blob("password")
			table.Bool("validated")
			table.String("lookup_token")
		}),
		Down: schema.DropIfExists("users"),
	})
}
