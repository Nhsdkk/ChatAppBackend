package user_validators

import (
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/sqlc/db"
	"context"
)

type NameUniquenessValidator struct {
	Db db.IDbConnection
}

func (n NameUniquenessValidator) Validate(name *string, ctx context.Context, _ request_env.RequestEnv) bool {
	if exists, err := n.Db.GetQueries().NameExists(ctx, *name); err != nil || exists {
		return false
	}

	return true
}
