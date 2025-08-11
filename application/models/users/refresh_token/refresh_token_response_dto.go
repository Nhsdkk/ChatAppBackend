package refresh_token

import (
	"chat_app_backend/internal/extensions"
	"chat_app_backend/internal/sqlc/db_queries"
	"time"
)

type RefreshTokenResponseDto struct {
	ID             extensions.UUID     `json:"id"`
	FullName       string              `json:"full_name"`
	Birthday       time.Time           `json:"birthday"`
	Gender         db_queries.Gender   `json:"gender"`
	Email          string              `json:"email"`
	AvatarFileName string              `json:"avatar_file_name"`
	Online         bool                `json:"online"`
	EmailVerified  bool                `json:"email_verified"`
	LastSeen       time.Time           `json:"last_seen"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	AccessToken    string              `json:"access_token"`
	Role           db_queries.RoleType `json:"role"`
}
