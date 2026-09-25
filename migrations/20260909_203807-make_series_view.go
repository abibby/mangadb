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
        "banner_url",
        "cover_url",
        51 AS "quality"
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
        "banner_url",
        "cover_url",
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
        null as "banner_url",
        "cover_url",
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
        "banner_url",
        "cover_url",
        50 AS "quality"
    FROM "int_mangaplus_series"
),
"stacked_sources_with_id" AS (
    SELECT
        "id_maps"."series_id" as "id",
        "stacked_sources".*,
        "banner_image".id as "banner_image_id",
        "cover_image".id as "cover_image_id"
    from "stacked_sources"
    join "id_maps" on "stacked_sources"."source_id" = "id_maps"."source_series_id" and "stacked_sources"."source" = "id_maps"."source"
    left join "images" as "banner_image" on "stacked_sources"."banner_url" = "banner_image"."source_url"
    left join "images" as "cover_image" on "stacked_sources"."cover_url" = "cover_image"."source_url"
)
SELECT 
    "id",
    (array_remove(array_agg("title" ORDER BY "quality" DESC), NULL))[1] AS "title",
    (array_remove(array_agg("aliases" ORDER BY "quality" DESC), NULL))[1] AS "aliases",
    (array_remove(array_agg("description" ORDER BY "quality" DESC), NULL))[1] AS "description",
    (array_remove(array_agg("start_date" ORDER BY "quality" DESC), NULL))[1] AS "start_date",
    (array_remove(array_agg("genre" ORDER BY "quality" DESC), NULL))[1] AS "genre",
    (array_remove(array_agg("tags" ORDER BY "quality" DESC), NULL))[1] AS "tags",
    (array_remove(array_agg("banner_image_id" ORDER BY "quality" DESC), NULL))[1] AS "banner_image_id",
    (array_remove(array_agg("cover_image_id" ORDER BY "quality" DESC), NULL))[1] AS "cover_image_id"
FROM "stacked_sources_with_id"
GROUP BY "id";

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE UNIQUE INDEX series_id_idx ON series (id);
CREATE INDEX series_title_trgm_idx ON series USING gin (title gin_trgm_ops);
`),
		Down: schema.Run(func(ctx context.Context, tx database.DB) error {
			return nil
		}),
	})
}
