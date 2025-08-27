package get

import "chat_app_backend/internal/extensions"

type GetInterestsRequestDto struct {
	Name *string           `json:"name" validator:"length lt 255"`
	Ids  []extensions.UUID `json:"ids"`
}
