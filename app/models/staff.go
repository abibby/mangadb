package models

import (
	"context"

	"abibby.com/salusa/database/builder"
	"abibby.com/salusa/database/model/modeldi"
	"github.com/abibby/icbmdb/app/providers"
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
