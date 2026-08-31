package models

import (
	"context"

	"abibby.com/salusa/database/builder"
	"abibby.com/salusa/database/model"
	"abibby.com/salusa/database/model/modeldi"
	"github.com/abibby/icbmdb/app/providers"
	"github.com/google/uuid"
)

//go:generate spice generate:migration
type IdMap struct {
	model.BaseModel

	SourceID string    `json:"source_id" db:"source_id,primary"`
	LocalID  uuid.UUID `json:"local_id" db:"local_id,primary"`
}

func init() {
	providers.Add(modeldi.Register[*IdMap])
}

func IdMapQuery(ctx context.Context) *builder.ModelBuilder[*IdMap] {
	return builder.From[*IdMap]().WithContext(ctx)
}
