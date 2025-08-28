package create

import (
	"mime/multipart"
)

type CreateInterestRequestDto struct {
	Title       string                `form:"title" validator:"not_empty;length lt 255"`
	Icon        *multipart.FileHeader `form:"icon" validator:"not_empty"`
	Description string                `form:"description" validator:"not_empty"`
}
