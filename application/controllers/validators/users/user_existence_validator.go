package user_validators

import (
	"chat_app_backend/internal/extensions"
	"chat_app_backend/internal/sqlc/db"
	"context"
)

type UserExistenceValidator struct {
	Db db.IDbConnection
}

func (u UserExistenceValidator) Validate(id *extensions.UUID, ctx context.Context) bool {
	if exists, err := u.Db.GetQueries().UserExists(ctx, *id); err != nil || !exists {
		return false
	}

	return true
}
