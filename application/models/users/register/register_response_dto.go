package register

import (
	interests "chat_app_backend/application/models/interests/get_many_by_ids"
	"chat_app_backend/internal/extensions"
	"chat_app_backend/internal/sqlc/db_queries"
	"time"
)

type RegisterResponseDto struct {
	ID             extensions.UUID             `json:"id"`
	FullName       string                      `json:"full_name"`
	Birthday       time.Time                   `json:"birthday"`
	Gender         db_queries.Gender           `json:"gender"`
	Email          string                      `json:"email"`
	AvatarFileName string                      `json:"avatar_file_name"`
	Online         bool                        `json:"online"`
	EmailVerified  bool                        `json:"email_verified"`
	LastSeen       time.Time                   `json:"last_seen"`
	CreatedAt      time.Time                   `json:"created_at"`
	UpdatedAt      time.Time                   `json:"updated_at"`
	Interests      []interests.GetInterestsDto `json:"interests"`
	Role           db_queries.RoleType         `json:"role"`
	AccessToken    string                      `json:"access_token"`
	RefreshToken   string                      `json:"refresh_token"`
}
