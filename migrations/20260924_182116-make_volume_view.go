package migrations

import (
	"context"

	"gosalusa.com/database"
	"gosalusa.com/database/migrate"
	"gosalusa.com/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20260924_182116-make_volume_view",
		Up: schema.Raw(`
CREATE MATERIALIZED VIEW volumes as
WITH stacked_sources AS (
	SELECT int_viz_volumes.series_id AS source_series_id,
		'viz'::text AS source,
		int_viz_volumes.number AS number,
		50 AS quality,
		int_viz_volumes.cover_url AS cover_url
	FROM int_viz_volumes
	UNION ALL
	SELECT int_mangadex_volumes.series_id AS source_series_id,
		'mangadex'::text AS source,
		int_mangadex_volumes.number,
		1 AS quality,
		int_mangadex_volumes.cover_url
	FROM int_mangadex_volumes
),
stacked_sources_with_id AS (
	SELECT id_maps.series_id,
		stacked_sources.source_series_id,
		stacked_sources.source,
		stacked_sources.number,
		stacked_sources.quality,
		cover_image.id as "cover_image_id"
	FROM stacked_sources
	JOIN id_maps ON stacked_sources.source_series_id = id_maps.source_series_id AND stacked_sources.source = id_maps.source
	LEFT JOIN images as cover_image ON stacked_sources.cover_url = cover_image.source_url
	WHERE stacked_sources.number IS NOT NULL
)
SELECT
	series_id,
	number,
	(array_remove(array_agg(cover_image_id ORDER BY quality DESC), NULL))[1] AS cover_image_id
FROM stacked_sources_with_id
GROUP BY series_id, number;

CREATE UNIQUE INDEX volumes_series_number_idx ON volumes (series_id, number);
`),
		Down: schema.Run(func(ctx context.Context, tx database.DB) error {
			return nil
		}),
	})
}
