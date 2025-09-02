package interests_validators

import (
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/s3"
	"context"
)

type IconFileTypeValidator struct{}

func (i IconFileTypeValidator) Validate(iconFileName *string, _ context.Context, _ request_env.RequestEnv) bool {
	_, fileType, fileTypeExtractionError := s3.DeconstructFileName(*iconFileName)
	if fileTypeExtractionError != nil {
		return false
	}

	switch fileType {
	case s3.Png, s3.Svg:
		return true
	default:
		return false
	}
}
