package register

import (
	"chat_app_backend/internal/sqlc/db_queries"
	"github.com/google/uuid"
	"time"
)

type RegisterRequestDto struct {
	FullName  string            `json:"full_name" validator:"not_empty;length gt 10;length lt 255" binder:"body,full_name"`
	Birthday  time.Time         `json:"birthday" validator:"not_empty" binder:"body,birthday"`
	Gender    db_queries.Gender `json:"gender" binder:"body,gender"`
	Email     string            `json:"email" validator:"not_empty;length lt 255" binder:"body,email"`
	Password  string            `json:"password" validator:"not_empty;length gt 10;length lt 255" binder:"body,password"`
	Interests []uuid.UUID       `json:"interests" binder:"body,interests"`
}
