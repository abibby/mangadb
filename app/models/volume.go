package models

import (
	"context"

	"abibby.com/mangadb/app/providers"
	"gosalusa.com/database/builder"
	"gosalusa.com/database/model/modeldi"
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
