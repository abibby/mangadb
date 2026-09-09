package models

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"abibby.com/mangadb/app/providers"
	"abibby.com/salusa/database"
	"abibby.com/salusa/database/builder"
	"abibby.com/salusa/database/dialects"
	"abibby.com/salusa/database/model/mixins"
	"abibby.com/salusa/database/model/modeldi"
	"github.com/google/uuid"
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
	URL            string          `db:"url,unique"             json:"url"`
	Page           int             `db:"page"                   json:"page"`
	RawPayload     json.RawMessage `db:"raw_payload" json:"raw_payload"`
}

func init() {
	providers.Add(modeldi.Register[*APIResponse])
}

func ApiResponseQuery(ctx context.Context) *builder.ModelBuilder[*APIResponse] {
	return builder.From[*APIResponse]().WithContext(ctx)
}

func ApiResponseCreateOrUpdate(ctx context.Context, tx database.DB, r *APIResponse) error {
	existing, err := ApiResponseQuery(ctx).Where("url", "=", r.URL).First(tx)
	if err != nil {
		return err
	}

	if existing == nil {
		existing = &APIResponse{}
		*existing = *r
	}
	existing.SyncJobID = r.SyncJobID
	existing.RawPayload = r.RawPayload

	d, err := dialects.New(tx.DriverName())
	if err != nil {
		return err
	}

	now := time.Now().Format(time.RFC3339)
	m := map[string]any{
		"id":               uuid.NewString(),
		"created_at":       now,
		"updated_at":       now,
		"sync_job_id":      r.SyncJobID,
		"source":           r.Source,
		"source_series_id": r.SourceSeriesID,
		"data_type":        r.DataType,
		"url":              r.URL,
		"page":             r.Page,
		"raw_payload":      r.RawPayload,
	}

	insertResult, err := d.EncodeInsertQuery(&dialects.InsertQuery{
		Table:  database.GetTable(r),
		Values: []map[string]any{m},
	})
	if err != nil {
		return err
	}

	delete(m, "id")
	delete(m, "created_at")

	updateResult, err := d.EncodeUpdateQuery(&dialects.UpdateQuery{
		Values: m,
	})
	if err != nil {
		return err
	}

	q := dialects.RawQuery{
		SQL:      insertResult.SQL + ` ON CONFLICT ("url") DO ` + strings.ReplaceAll(updateResult.SQL, `"" `, ""),
		Bindings: append(insertResult.Bindings, updateResult.Bindings...),
	}

	_, err = tx.ExecContext(ctx, q.SQL, q.Bindings...)
	return err
}
