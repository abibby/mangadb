package models

import (
	"context"

	"abibby.com/mangadb/app/providers"
	"abibby.com/salusa/database/builder"
	"abibby.com/salusa/database/model/modeldi"
)

type Staff struct {
	BaseModel
}

func init() {
	providers.Add(modeldi.Register[*Staff])
}

func StaffQuery(ctx context.Context) *builder.ModelBuilder[*Staff] {
	return builder.From[*Staff]().WithContext(ctx)
}
