package models

import (
	"context"

	"github.com/abibby/icbmdb/app/providers"
	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/model/modeldi"
)

//go:generate spice generate:migration
type Staff struct {
	BaseModel
}

func init() {
	providers.Add(modeldi.Register[*Staff])
}

func StaffQuery(ctx context.Context) *builder.ModelBuilder[*Staff] {
	return builder.From[*Staff]().WithContext(ctx)
}
