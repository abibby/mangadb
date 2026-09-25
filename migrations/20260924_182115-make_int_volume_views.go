package migrations

import (
	"context"

	"gosalusa.com/database"
	"gosalusa.com/database/migrate"
	"gosalusa.com/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20260924_182115-make_int_volume_views",
		Up: schema.Raw(`
CREATE MATERIALIZED VIEW int_viz_volumes as
SELECT
	chapter.element ->> 'id' as "id",
	api_responses.source_series_id as "series_id",
	(chapter.element -> 'number')::numeric::int as "number",
	chapter.element ->> 'cover_image' as "cover_url"
FROM
	api_responses,
	jsonb_array_elements(raw_payload->'volumes') AS chapter(element)
where
	"source" = 'viz'
	and "data_type" = 'series';

CREATE UNIQUE INDEX int_viz_volumes_id_idx ON int_viz_volumes (id);

CREATE MATERIALIZED VIEW int_mangadex_volumes as
WITH ParsedCovers AS (
    SELECT 
        chapter.element #>> '{id}' as raw_id,
        api_responses.source_series_id as series_id,
        (chapter.element #>> '{attributes,volume}')::numeric as raw_volume,
        chapter.element #>> '{attributes,locale}' as locale,
        chapter.element #>> '{attributes,fileName}' as file_name
    FROM
        api_responses,
        jsonb_array_elements(raw_payload->'data') AS chapter(element)
    WHERE
        "source" = 'mangadex'
        and "data_type" = 'covers'
        and (chapter.element #>> '{attributes,volume}') ~ '^[0-9]+([.][0-9]+)?$'
)
SELECT DISTINCT ON (series_id, FLOOR(raw_volume))
    raw_id as "id",
    series_id,
    FLOOR(raw_volume)::int as "number",
    'https://mangadex.org/covers/' || series_id || '/' || file_name as "cover_url"
FROM
    ParsedCovers
ORDER BY
    series_id,
    FLOOR(raw_volume),
    -- 1. Prioritize whole numbers
    CASE WHEN raw_volume = FLOOR(raw_volume) THEN 1 ELSE 2 END ASC,
    -- 2. Prioritize 'ja' locale, then 'uk', then everything else
    CASE locale 
        WHEN 'ja' THEN 1 
        WHEN 'uk' THEN 2 
        ELSE 3 
    END ASC;

CREATE UNIQUE INDEX int_mangadex_volumes_id_idx ON int_mangadex_volumes (id);
`),
		Down: schema.Run(func(ctx context.Context, tx database.DB) error {
			return nil
		}),
	})
}
