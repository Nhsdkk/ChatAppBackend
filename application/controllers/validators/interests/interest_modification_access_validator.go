package interests_validators

import (
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/sqlc/db_queries"
	"context"
)

type InterestModificationAccessValidator struct{}

func (i InterestModificationAccessValidator) Validate(_ *interface{}, _ context.Context, env request_env.RequestEnv) bool {
	return env.User != nil && env.User.Role == db_queries.RoleTypeADMIN
}
