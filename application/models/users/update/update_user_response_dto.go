package update

import (
	"chat_app_backend/internal/sqlc/db_queries"
	"github.com/google/uuid"
	"time"
)

type UpdateUserResponseDto struct {
	ID             uuid.UUID           `json:"id"`
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
	Role           db_queries.RoleType `json:"role"`
	AccessToken    string              `json:"access_token"`
	RefreshToken   string              `json:"refresh_token"`
}
