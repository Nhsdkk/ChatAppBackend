package request_env

import "chat_app_backend/internal/sqlc/db_queries"

type RequestEnv struct {
	User *db_queries.User
}
