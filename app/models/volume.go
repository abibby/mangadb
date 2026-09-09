package models

import (
	"context"

	"abibby.com/salusa/database/builder"
	"abibby.com/salusa/database/model/modeldi"
	"github.com/abibby/icbmdb/app/providers"
)

type Volume struct {
	BaseModel
}

func init() {
	providers.Add(modeldi.Register[*Volume])
}

func CollectionQuery(ctx context.Context) *builder.ModelBuilder[*Volume] {
	return builder.From[*Volume]().WithContext(ctx)
}
