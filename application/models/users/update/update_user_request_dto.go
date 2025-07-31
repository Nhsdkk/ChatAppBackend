package update

import (
	"chat_app_backend/internal/sqlc/db_queries"
	"time"
)

type UpdateUserRequestDto struct {
	FullName       *string            `json:"full_name" binder:"body,full_name" validation:"length gt 10; length lt 255"`
	Birthday       *time.Time         `json:"birthday" binder:"body,birthday"`
	Gender         *db_queries.Gender `json:"gender" binder:"body,gender" mapper:"exclude"`
	Email          *string            `json:"email" binder:"body,email" validation:"length gt 10; length lt 255"`
	PasswordString *string            `json:"password" binder:"body,password" validation:"length gt 10; length lt 255"`
}
