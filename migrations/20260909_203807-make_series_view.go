package migrations

import (
	"context"

	"gosalusa.com/database"
	"gosalusa.com/database/migrate"
	"gosalusa.com/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20260909_203807-make_series_view",
		Up: schema.Raw(`
CREATE MATERIALIZED VIEW series as
WITH "stacked_sources" AS (
    SELECT
        "id" as "source_id",
        'viz' AS "source",
        "title",
        null as "aliases",
        "description",
        null as "start_date",
        null as "genre",
        null as "tags",
        50 AS "quality"
    FROM "int_viz_series"
    union all
    SELECT
        "id" as "source_id",
        'anilist' AS "source",
        "title",
        "aliases",
        "description",
        null as "start_date",
        null as "genre",
        "tags",
        25 AS "quality"
    FROM "int_anilist_series"
    union all
    SELECT
        "id" as "source_id",
        'mangadex' AS "source",
        "title",
        "aliases",
        "description",
        null as "start_date",
        null as "genre",
        null as "tags",
        1 AS "quality"
    FROM "int_mangadex_series"
    union all
    SELECT
        "id" as "source_id",
        'mangaplus' AS "source",
        "title",
        null as "aliases",
        "description",
        null as "start_date",
        null as "genre",
        null as "tags",
        50 AS "quality"
    FROM "int_mangaplus_series"
),
"stacked_sources_with_id" AS (
    SELECT
        "id_maps"."series_id" as "id",
        "stacked_sources".*
    from "stacked_sources"
    join "id_maps" on "stacked_sources"."source_id" = "id_maps"."source_series_id" and "stacked_sources"."source" = "id_maps"."source"
)
SELECT 
    "id",
    (array_remove(array_agg("title" ORDER BY "quality" DESC), NULL))[1] AS "title",
    (array_remove(array_agg("aliases" ORDER BY "quality" DESC), NULL))[1] AS "aliases",
    (array_remove(array_agg("description" ORDER BY "quality" DESC), NULL))[1] AS "description",
    (array_remove(array_agg("start_date" ORDER BY "quality" DESC), NULL))[1] AS "start_date",
    (array_remove(array_agg("genre" ORDER BY "quality" DESC), NULL))[1] AS "genre",
    (array_remove(array_agg("tags" ORDER BY "quality" DESC), NULL))[1] AS "tags"
FROM "stacked_sources_with_id"
GROUP BY "id";

CREATE UNIQUE INDEX series_id_idx ON series (id);
`),
		Down: schema.Run(func(ctx context.Context, tx database.DB) error {
			return nil
		}),
	})
}
