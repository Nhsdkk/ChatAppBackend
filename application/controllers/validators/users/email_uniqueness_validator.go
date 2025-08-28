package user_validators

import (
	"chat_app_backend/internal/sqlc/db"
	"context"
)

type EmailUniquenessValidator struct {
	Db db.IDbConnection
}

func (e EmailUniquenessValidator) Validate(email *string, ctx context.Context) bool {
	if exists, err := e.Db.GetQueries().EmailExists(ctx, *email); err != nil || exists {
		return false
	}

	return true
}
