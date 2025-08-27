package update

import (
	"chat_app_backend/internal/extensions"
	"mime/multipart"
)

type UpdateInterestRequestDto struct {
	ID          extensions.UUID       `uri:"id" validator:"not_empty"`
	Icon        *multipart.FileHeader `form:"icon" mapper:"exclude"`
	Description *string               `form:"description" mapper:"exclude"`
}
