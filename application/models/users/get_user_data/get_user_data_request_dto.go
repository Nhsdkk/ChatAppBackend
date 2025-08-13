package get_user_data

import (
	"chat_app_backend/internal/extensions"
)

type GetUserDataRequestDto struct {
	ID extensions.UUID `uri:"id" validator:"not_empty"`
}
