package models

import (
	"context"

	"gosalusa.com/auth"
	"gosalusa.com/database/builder"
)

//go:generate spice generate:migration
type User struct {
	auth.EmailVerifiedUser
}

func UserQuery(ctx context.Context) *builder.ModelBuilder[*User] {
	return builder.From[*User]().WithContext(ctx)
}
