package create

import (
	"chat_app_backend/internal/s3"
	"mime/multipart"
)

type CreateInterestRequestDto struct {
	Title        string                `form:"title" validator:"not_empty;length lt 255"`
	Icon         *multipart.FileHeader `form:"icon" validator:"not_empty"`
	IconFileType s3.FileType           `form:"icon_file_type"`
	Description  string                `form:"description" validator:"not_empty"`
}
