package migrations

import (
	"abibby.com/salusa/database/migrate"
	"abibby.com/salusa/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20260831_081823-Staff",
		Up: schema.Create("staffs", func(table *schema.Blueprint) {
			table.Blob("id").Primary()
		}),
		Down: schema.DropIfExists("staffs"),
	})
}
