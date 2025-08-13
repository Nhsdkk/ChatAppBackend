package interests

import (
	"chat_app_backend/internal/extensions"
	"time"
)

type GetInterestsDto struct {
	ID           extensions.UUID `json:"id,omitempty"`
	Title        string          `json:"title,omitempty"`
	IconFileName string          `json:"icon_file_name,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}
