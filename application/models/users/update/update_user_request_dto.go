package update

import (
	"chat_app_backend/internal/sqlc/db_queries"
	"github.com/google/uuid"
	"time"
)

type UpdateUserRequestDto struct {
	ID             uuid.UUID            `json:"id" binder:"path,id" validator:"not_empty"`
	FullName       *string              `json:"full_name" binder:"body,full_name" validation:"length gt 10; length lt 255"`
	Birthday       *time.Time           `json:"birthday" binder:"body,birthday"`
	Gender         *db_queries.Gender   `json:"gender" binder:"body,gender" mapper:"exclude"`
	Email          *string              `json:"email" binder:"body,email" validation:"length gt 10; length lt 255"`
	PasswordString *string              `json:"password" binder:"body,password" validation:"length gt 10; length lt 255"`
	Role           *db_queries.RoleType `json:"role" binder:"body,role" mapper:"exclude"`
}
