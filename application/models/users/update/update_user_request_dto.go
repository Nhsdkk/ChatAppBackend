package update

import (
	"chat_app_backend/internal/extensions"
	"chat_app_backend/internal/sqlc/db_queries"
	"time"
)

type UpdateUserRequestDto struct {
	ID             extensions.UUID      `uri:"id"  validator:"not_empty"`
	FullName       *string              `json:"full_name"  validation:"length gt 10; length lt 255"`
	Birthday       *time.Time           `json:"birthday" binder:"body,birthday"`
	Gender         *db_queries.Gender   `json:"gender"  mapper:"exclude"`
	Email          *string              `json:"email"  validation:"length gt 10; length lt 255"`
	PasswordString *string              `json:"password"  validation:"length gt 10; length lt 255"`
	Role           *db_queries.RoleType `json:"role"  mapper:"exclude"`
}
