package migrations

import (
	"context"

	"gosalusa.com/database"
	"gosalusa.com/database/migrate"
	"gosalusa.com/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20260910_165435-make_chapter_view",
		Up: schema.Raw(`
CREATE MATERIALIZED VIEW chapters as
WITH stacked_sources AS (
	SELECT int_viz_chapter.series_id AS source_series_id,
		'viz'::text AS source,
		int_viz_chapter.title,
		int_viz_chapter.volume,
		int_viz_chapter.chapter,
		50 AS quality
	FROM int_viz_chapter
	UNION ALL
	SELECT int_mangaplus_chapter.series_id AS source_series_id,
		'mangaplus'::text AS source,
		int_mangaplus_chapter.title,
		NULL::integer AS volume,
		int_mangaplus_chapter.chapter,
		51 AS quality
	FROM int_mangaplus_chapter
	UNION ALL
	SELECT int_mangadex_chapter.series_id AS source_series_id,
		'mangadex'::text AS source,
		int_mangadex_chapter.title,
		int_mangadex_chapter.volume,
		int_mangadex_chapter.chapter,
		1 AS quality
	FROM int_mangadex_chapter
	WHERE (int_mangadex_chapter.language = 'en'::text)
),
stacked_sources_with_id AS (
	SELECT id_maps.series_id,
		stacked_sources.source_series_id,
		stacked_sources.source,
		stacked_sources.title,
		stacked_sources.volume,
		stacked_sources.chapter,
		stacked_sources.quality
	FROM stacked_sources
	JOIN id_maps ON stacked_sources.source_series_id = id_maps.source_series_id AND stacked_sources.source = id_maps.source
)
SELECT
	series_id,
	(array_remove(array_agg(volume ORDER BY quality DESC), NULL::integer))[1] AS volume,
	chapter,
	(array_remove(array_agg(title ORDER BY quality DESC), NULL::text))[1] AS title
FROM stacked_sources_with_id
GROUP BY series_id, chapter;

CREATE UNIQUE INDEX chapter_series_chapter_idx ON chapters (series_id, chapter);
`),
		Down: schema.Run(func(ctx context.Context, tx database.DB) error {
			return nil
		}),
	})
}
