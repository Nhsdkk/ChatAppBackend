package create

import (
	"chat_app_backend/internal/extensions"
	"time"
)

type CreateInterestResponseDto struct {
	ID               extensions.UUID `json:"id"`
	Title            string          `json:"title"`
	IconDownloadLink string          `json:"icon_download_link"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}
