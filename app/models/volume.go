package models

import (
	"context"

	"github.com/abibby/icbmdb/app/providers"
	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/model/modeldi"
)

//go:generate spice generate:migration
type Volume struct {
	BaseModel
	DatasourceIDs
}

func init() {
	providers.Add(modeldi.Register[*Volume])
}

func CollectionQuery(ctx context.Context) *builder.ModelBuilder[*Volume] {
	return builder.From[*Volume]().WithContext(ctx)
}
