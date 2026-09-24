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
SELECT distinct on (chapter.element #>> '{id}')
	chapter.element #>> '{id}' as "id",
	api_responses.source_series_id as "series_id",
	(chapter.element #>> '{attributes,volume}')::numeric::int as "number",
    'https://mangadex.org/covers/' || api_responses.source_series_id || '/' || (chapter.element #>> '{attributes,fileName}')::text as "cover_url"
FROM
    api_responses,
    jsonb_array_elements(raw_payload->'data') AS chapter(element)
WHERE
	"source" = 'mangadex'
	and "data_type" = 'covers'
	and (chapter.element #>> '{attributes,volume}') ~ '^[0-9]+([.][0-9]+)?$'
	and (chapter.element #>> '{attributes,volume}')::numeric % 1 = 0;

CREATE UNIQUE INDEX int_mangadex_volumes_id_idx ON int_mangadex_volumes (id);
`),
		Down: schema.Run(func(ctx context.Context, tx database.DB) error {
			return nil
		}),
	})
}
