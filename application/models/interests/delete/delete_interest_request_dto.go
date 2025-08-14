package delete

import "chat_app_backend/internal/extensions"

type DeleteInterestRequestDto struct {
	ID extensions.UUID `uri:"id" validator:"not_empty"`
}
