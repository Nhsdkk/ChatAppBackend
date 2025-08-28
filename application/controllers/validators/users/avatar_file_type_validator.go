package user_validators

import (
	"chat_app_backend/internal/s3"
	"context"
)

type AvatarFileTypeValidator struct{}

func (a AvatarFileTypeValidator) Validate(avatarFileName *string, _ context.Context) bool {
	_, fileType, fileTypeExtractionError := s3.DeconstructFileName(*avatarFileName)
	if fileTypeExtractionError != nil {
		return false
	}

	switch fileType {
	case s3.Png, s3.Jpeg:
		return true
	default:
		return false
	}
}
