package models

import (
	"context"

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
}

func init() {
	providers.Add(modeldi.Register[*APIResponse])
}

func ApiResponseQuery(ctx context.Context) *builder.ModelBuilder[*APIResponse] {
	return builder.From[*APIResponse]().WithContext(ctx)
}
