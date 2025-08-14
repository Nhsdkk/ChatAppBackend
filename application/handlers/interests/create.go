package interests

import (
	"chat_app_backend/application/models/interests/create"
	"chat_app_backend/internal/exceptions"
	"chat_app_backend/internal/exceptions/common_exceptions"
	"chat_app_backend/internal/mapper"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/s3"
	"chat_app_backend/internal/service_wrapper"
	"chat_app_backend/internal/sqlc/db_queries"

	"github.com/gin-gonic/gin"
)

type CreateInterestHandler struct{}

func (c CreateInterestHandler) Handle(
	request *create.CreateInterestRequestDto,
	services service_wrapper.IServiceWrapper,
	ctx *gin.Context,
	requestEnvironment *request_env.RequestEnv,
) (*create.CreateInterestResponseDto, exceptions.ITrackableException) {
	if requestEnvironment.User.Role != db_queries.RoleTypeADMIN {
		return nil, common_exceptions.ForbiddenException{
			BaseRestException: exceptions.BaseRestException{
				ITrackableException: exceptions.CreateTrackableExceptionFromStringF("interest creation forbidden"),
				Message:             "only admins can access this route",
			},
		}
	}

	filename := s3.ConstructFilenameFromFileType(request.IconFileType)
	iconDownloadLink, fileUploadError := services.GetS3Client().
		UploadFile(ctx, request.Icon, filename, s3.InterestsIconBucket)
	if fileUploadError != nil {
		return nil, exceptions.WrapErrorWithTrackableException(fileUploadError)
	}

	dbRequest := db_queries.CreateInterestParams{}
	dbRequestMapperError := mapper.Mapper{}.Map(
		&dbRequest,
		*request,
		struct {
			IconFileName string
		}{
			IconFileName: filename,
		},
	)

	if dbRequestMapperError != nil {
		return nil, exceptions.WrapErrorWithTrackableException(dbRequestMapperError)
	}

	interest, creationError := services.GetDbConnection().
		GetQueries().
		CreateInterest(ctx, dbRequest)
	if creationError != nil {
		return nil, exceptions.WrapErrorWithTrackableException(creationError)
	}

	response := create.CreateInterestResponseDto{}
	responseMappingError := mapper.Mapper{}.Map(
		&response,
		interest,
		struct {
			IconDownloadLink string
		}{
			IconDownloadLink: iconDownloadLink,
		},
	)

	if responseMappingError != nil {
		return nil, exceptions.WrapErrorWithTrackableException(responseMappingError)
	}

	return &response, nil
}
