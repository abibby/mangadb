package models

import (
	"context"

	"abibby.com/salusa/auth"
	"abibby.com/salusa/database/builder"
)

//go:generate spice generate:migration
type User struct {
	auth.EmailVerifiedUser
}

func UserQuery(ctx context.Context) *builder.ModelBuilder[*User] {
	return builder.From[*User]().WithContext(ctx)
}
