package migrations

import (
	"context"

	"gosalusa.com/database"
	"gosalusa.com/database/migrate"
	"gosalusa.com/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20260831_091138-make_int_chapter_views",
		Up: schema.Raw(`
CREATE MATERIALIZED VIEW int_mangadex_chapter as
SELECT distinct on (chapter.element #>> '{id}')
	chapter.element #>> '{id}' as "id",
	source_series_id as "series_id",
	chapter.element #>> '{attributes,title}' as "title",
	NULLIF(chapter.element #>> '{attributes,volume}', 'null')::integer as "volume",
	NULLIF(chapter.element #>> '{attributes,chapter}', 'null')::float as "chapter",
	chapter.element #>> '{attributes,translatedLanguage}' as "language",
	(chapter.element #>> '{attributes,createdAt}')::timestamptz as "created_at",
	(chapter.element #>> '{attributes,updatedAt}')::timestamptz as "updated_at"
FROM
    api_responses,
    jsonb_array_elements(raw_payload->'data') AS chapter(element)
where
	"source" = 'mangadex'
	and "data_type" = 'chapter_list';

CREATE UNIQUE INDEX int_mangadex_chapter_id_idx ON int_mangadex_chapter (id);

CREATE MATERIALIZED VIEW int_mangaplus_chapter as
select
	raw_payload ->> 'chapterId' as "id",
	series_id,
	REGEXP_REPLACE(REGEXP_REPLACE(raw_payload ->> 'name', '^#', ''), '-', '.')::float as "chapter",
	REGEXP_REPLACE(raw_payload ->> 'subTitle', '^[^:]+\d+:', '') as "title",
	to_timestamp((raw_payload ->> 'startTimeStamp')::integer) as "created_at"
from (
	select
		source_series_id as "series_id",
		chapter.element as raw_payload
	from
		api_responses,
		jsonb_array_elements(raw_payload #> '{chapterListGroup,0,firstChapterList}') as chapter(element)
	where
		"source" = 'mangaplus'
		and "data_type" = 'series'
	union all
	select
		source_series_id as "series_id",
		chapter.element as raw_payload
	from
		api_responses,
		jsonb_array_elements(raw_payload #> '{chapterListGroup,1,lastChapterList}') as chapter(element)
	where
		"source" = 'mangaplus'
		and "data_type" = 'series'

);

CREATE UNIQUE INDEX int_mangaplus_chapter_id_idx ON int_mangaplus_chapter (id);

CREATE MATERIALIZED VIEW int_viz_chapter as
SELECT
	chapter.element ->> 'id' as "id",
	source_series_id as "series_id",
	chapter.element ->> 'title' as "title",
	NULLIF(chapter.element ->> 'volume', 'null')::integer as "volume",
	(chapter.element ->> 'chapter')::float as "chapter"
FROM
	api_responses,
	jsonb_array_elements(raw_payload->'chapters') AS chapter(element)
where
	"source" = 'viz'
	and "data_type" = 'series';

CREATE UNIQUE INDEX int_viz_chapter_id_idx ON int_viz_chapter (id);
`),
		Down: schema.Run(func(ctx context.Context, tx database.DB) error {
			return nil
		}),
	})
}
