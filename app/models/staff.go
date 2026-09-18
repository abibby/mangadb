package models

import (
	"context"

	"abibby.com/mangadb/app/providers"
	"gosalusa.com/database/builder"
	"gosalusa.com/database/model/modeldi"
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
