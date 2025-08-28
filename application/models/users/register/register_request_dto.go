package register

import (
	"chat_app_backend/internal/extensions"
	"chat_app_backend/internal/sqlc/db_queries"
	"mime/multipart"
	"time"
)

type RegisterRequestDto struct {
	FullName  string                `form:"full_name" validator:"not_empty;length gt 10;length lt 255" `
	Birthday  time.Time             `form:"birthday" validator:"not_empty"`
	Gender    db_queries.Gender     `form:"gender"`
	Email     string                `form:"email" validator:"not_empty;length lt 255" `
	Password  string                `form:"password" validator:"not_empty;length gt 10;length lt 255" `
	Interests []extensions.UUID     `form:"interests"`
	Avatar    *multipart.FileHeader `form:"avatar"`
}
