package update

import (
	"chat_app_backend/internal/extensions"
	"chat_app_backend/internal/sqlc/db_queries"
	"mime/multipart"
	"time"
)

type UpdateUserRequestDto struct {
	ID             extensions.UUID       `uri:"id"  validator:"not_empty"`
	FullName       *string               `form:"full_name"  validation:"length gt 10; length lt 255"`
	Birthday       *time.Time            `form:"birthday" binder:"body,birthday"`
	Gender         *db_queries.Gender    `form:"gender"  mapper:"exclude"`
	Email          *string               `form:"email"  validation:"length gt 10; length lt 255"`
	PasswordString *string               `form:"password"  validation:"length gt 10; length lt 255"`
	Role           *db_queries.RoleType  `form:"role"  mapper:"exclude"`
	Avatar         *multipart.FileHeader `form:"avatar"`
}
