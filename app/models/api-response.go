package models

import (
	"context"
	"encoding/json"

	"abibby.com/salusa/database/builder"
	"abibby.com/salusa/database/model/mixins"
	"abibby.com/salusa/database/model/modeldi"
	"github.com/abibby/icbmdb/app/providers"
)

//go:generate spice generate:migration
type APIResponse struct {
	BaseModel

	mixins.SoftDelete
	mixins.Timestamps

	SyncJobID      string          `db:"sync_job_id"            json:"sync_job_id"`
	Source         string          `db:"source"                 json:"source"`
	SourceSeriesID string          `db:"source_series_id"       json:"source_series_id"`
	DataType       string          `db:"data_type"              json:"data_type"`
	URL            string          `db:"url"                    json:"url"`
	Page           int             `db:"page"                   json:"page"`
	RawPayload     json.RawMessage `db:"raw_payload,type:jsonb" json:"raw_payload"`
}

func init() {
	providers.Add(modeldi.Register[*APIResponse])
}

func ApiResponseQuery(ctx context.Context) *builder.ModelBuilder[*APIResponse] {
	return builder.From[*APIResponse]().WithContext(ctx)
}
