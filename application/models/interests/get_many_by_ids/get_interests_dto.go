package interests

import (
	"github.com/google/uuid"
	"time"
)

type GetInterestsDto struct {
	ID           uuid.UUID `json:"id,omitempty"`
	Title        string    `json:"title,omitempty"`
	IconFileName string    `json:"icon_file_name,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
