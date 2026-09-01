package migrations

import (
	"context"

	"abibby.com/salusa/database"
	"abibby.com/salusa/database/migrate"
	"abibby.com/salusa/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20260831_091138-make_int_mangadex_chapter_view",
		Up: schema.Run(func(ctx context.Context, tx database.DB) error {
			_, err := tx.ExecContext(ctx, `
CREATE MATERIALIZED VIEW IF NOT EXISTS int_mangadex_chapter as 
SELECT 
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
	and "data_type" = 'chapter_list'
`)

			/*
				select
					raw_payload ->> 'chapterId' as "id",
					series_id,
					REGEXP_REPLACE(raw_payload ->> 'name', '^#', '')::float as "chapter",
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

				)
			*/
			return err
		}),
		Down: schema.Run(func(ctx context.Context, tx database.DB) error {
			return nil
		}),
	})
}
