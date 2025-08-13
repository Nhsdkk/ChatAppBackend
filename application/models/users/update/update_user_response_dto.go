package update

import (
	"chat_app_backend/internal/extensions"
	"chat_app_backend/internal/sqlc/db_queries"
	"time"
)

type UpdateUserResponseDto struct {
	ID                 extensions.UUID     `json:"id"`
	FullName           string              `json:"full_name"`
	Birthday           time.Time           `json:"birthday"`
	Gender             db_queries.Gender   `json:"gender"`
	Email              string              `json:"email"`
	Online             bool                `json:"online"`
	EmailVerified      bool                `json:"email_verified"`
	LastSeen           time.Time           `json:"last_seen"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
	Role               db_queries.RoleType `json:"role"`
	AccessToken        string              `json:"access_token"`
	RefreshToken       string              `json:"refresh_token"`
	AvatarDownloadLink string              `json:"avatar_download_link"`
}
