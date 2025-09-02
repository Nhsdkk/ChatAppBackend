package user_validators

import (
	"chat_app_backend/internal/extensions"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/sqlc/db_queries"
	"context"
)

type UserModificationAccessValidator struct{}

func (u UserModificationAccessValidator) Validate(userId *extensions.UUID, _ context.Context, env request_env.RequestEnv) bool {
	return env.User != nil && (env.User.Role == db_queries.RoleTypeADMIN || env.User.ID == *userId)
}
