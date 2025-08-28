package interests_validators

import (
	"chat_app_backend/internal/extensions"
	"chat_app_backend/internal/sqlc/db"
	"context"
)

type InterestsExistenceValidator struct {
	Db db.IDbConnection
}

func (i InterestsExistenceValidator) Validate(ids *[]extensions.UUID, ctx context.Context) bool {
	if count, err := i.Db.GetQueries().ExistenceCheck(ctx, *ids); err != nil || int64(len(*ids)) != count {
		return false
	}

	return true
}
