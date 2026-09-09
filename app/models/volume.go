package models

import (
	"context"

	"abibby.com/mangadb/app/providers"
	"abibby.com/salusa/database/builder"
	"abibby.com/salusa/database/model/modeldi"
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
