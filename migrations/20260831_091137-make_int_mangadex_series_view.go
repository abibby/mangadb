package migrations

import (
	"context"

	"abibby.com/salusa/database"
	"abibby.com/salusa/database/migrate"
	"abibby.com/salusa/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20260831_091137-make_int_mangadex_series_view",
		Up: schema.Run(func(ctx context.Context, tx database.DB) error {
			_, err := tx.ExecContext(ctx, `
CREATE MATERIALIZED VIEW IF NOT EXISTS int_mangadex_series as 
select
	raw_payload #>> '{data,0,id}' as "id",
	jsonb_path_query_first(raw_payload #> '{data,0,attributes,title}', '$.*') #>> '{}' as "title",
	ARRAY(
        SELECT distinct jsonb_path_query(raw_payload #> '{data,0,attributes,altTitles}', '$[*].*') #>> '{}'
    ) AS "aliases",
	raw_payload #>> '{data,0,attributes,description,en}' as "description",
	make_date((raw_payload #>> '{data,0,attributes,year}')::int, 1, 1) as "start_date",
	(raw_payload #>> '{data,0,attributes,updatedAt}')::timestamptz as "updated_at",
	(raw_payload #>> '{data,0,attributes,createdAt}')::timestamptz as "created_at"
from api_responses
where
	"source" = 'mangadex'
	and "data_type" = 'series'
`)
			return err
		}),
		Down: schema.Run(func(ctx context.Context, tx database.DB) error {
			return nil
		}),
	})
}
