package migrations

import (
	"context"

	"abibby.com/salusa/database"
	"abibby.com/salusa/database/migrate"
	"abibby.com/salusa/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20260831_091137-make_int_series_views",
		Up: schema.Raw(`
CREATE MATERIALIZED VIEW int_mangadex_series as
select
	raw_payload #>> '{data,0,id}' as "id",
	jsonb_path_query_first(raw_payload #> '{data,0,attributes,title}', '$.*') #>> '{}' as "title",
	to_jsonb(ARRAY(
		SELECT distinct jsonb_path_query(raw_payload #> '{data,0,attributes,altTitles}', '$[*].*') #>> '{}'
	)) AS "aliases",
	raw_payload #>> '{data,0,attributes,description,en}' as "description",
	make_date((raw_payload #>> '{data,0,attributes,year}')::int, 1, 1) as "start_date",
	(raw_payload #>> '{data,0,attributes,updatedAt}')::timestamptz as "updated_at",
	(raw_payload #>> '{data,0,attributes,createdAt}')::timestamptz as "created_at"
from api_responses
where
	"source" = 'mangadex'
	and "data_type" = 'series';

CREATE UNIQUE INDEX int_mangadex_series_id_idx ON int_mangadex_series (id);

CREATE MATERIALIZED VIEW int_mangaplus_series as
select
	(raw_payload #>> '{title,titleId}')::varchar as "id",
	raw_payload #>> '{title,name}' as "title",
	raw_payload ->> 'overview' as "description"
from
	api_responses
where
	"source" = 'mangaplus'
	and "data_type" = 'series';

CREATE UNIQUE INDEX int_mangaplus_series_id_idx ON int_mangaplus_series (id);

CREATE MATERIALIZED VIEW int_viz_series as
select
  "source_series_id" as "id",
  "raw_payload" ->> 'title' as "title",
  "raw_payload" ->> 'description' as "description",
  "raw_payload" ->> 'author' as "author"
from
  api_responses
where
  source = 'viz'
  and data_type = 'series';

CREATE UNIQUE INDEX int_viz_series_id_idx ON int_viz_series (id);

CREATE MATERIALIZED VIEW int_anilist_series as
select
	(raw_payload #>> '{data,Media,id}')::text as "id",
	(raw_payload #>> '{data,Media,title,userPreferred}')::text as "title",
	to_jsonb(ARRAY(
		SELECT jsonb_array_elements_text(raw_payload #> '{data,Media,synonyms}')
	)) as "aliases",
	(raw_payload #>> '{data,Media,description}')::text as "description",
	to_jsonb(ARRAY(
		SELECT jsonb_array_elements_text(raw_payload #> '{data,Media,genres}')
	)) AS "genres",
	to_jsonb(ARRAY(
		SELECT distinct jsonb_path_query(raw_payload #> '{data,Media,tags}', '$.name') #>> '{}'
	)) AS "tags",
	to_jsonb(ARRAY(
		SELECT edge #>> '{node,name,userPreferred}'
		FROM jsonb_array_elements(raw_payload #> '{data,Media,staff,edges}') AS edge
		WHERE edge #>> '{role}' ILIKE '%Story & Art%'
			OR edge #>> '{role}' ILIKE '%Original Creator%'
			OR edge #>> '{role}' ILIKE '%Supervisor%'
    )) AS "authors",
	raw_payload
from
	api_responses
where
	"source" = 'anilist'
	and "data_type" = 'series';

CREATE UNIQUE INDEX int_anilist_series_id_idx ON int_anilist_series (id);

`),
		/*
			select
				(raw_payload #>> '{data,Media,id}')::varchar as "id",
				(raw_payload #>> '{data,Media,title,userPreferred}')::varchar as "title",
				(raw_payload #>> '{data,Media,description}')::varchar as "description"
			from
				api_responses
			where
				"source" = 'anilist'
				and "data_type" = 'series';
		*/
		Down: schema.Run(func(ctx context.Context, tx database.DB) error {
			return nil
		}),
	})
}
