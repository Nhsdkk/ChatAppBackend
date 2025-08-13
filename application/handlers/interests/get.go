package interests

import (
	interests "chat_app_backend/application/models/interests/get"
	"chat_app_backend/internal/exceptions"
	"chat_app_backend/internal/mapper"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/s3"
	"chat_app_backend/internal/service_wrapper"
	"chat_app_backend/internal/sqlc/db_queries"

	"github.com/gin-gonic/gin"
)

type GetInterestsHandler struct{}

func (g GetInterestsHandler) Handle(
	request *interests.GetInterestsRequestDto,
	service service_wrapper.IServiceWrapper,
	ctx *gin.Context,
	_ *request_env.RequestEnv,
) (*interests.GetInterestsResponseDto, exceptions.ITrackableException) {
	rawInterests, dbError := service.GetDbConnection().
		GetQueries().
		GetManyInterestsByFilters(ctx, db_queries.GetManyInterestsByFiltersParams{
			Ids:  request.Ids,
			Name: request.Name,
		})

	if dbError != nil {
		return nil, exceptions.WrapErrorWithTrackableException(dbError)
	}

	mappedInterests := make([]interests.GetInterestResponseDto, len(rawInterests), len(rawInterests))

	for idx, rawInterest := range rawInterests {
		link, s3Error := service.GetS3Client().GetDownloadUrl(ctx, rawInterest.IconFileName, s3.InterestsIconBucket)
		if s3Error != nil {
			return nil, exceptions.WrapErrorWithTrackableException(s3Error)
		}

		mappingErr := mapper.Mapper{}.Map(
			&mappedInterests[idx],
			rawInterest,
			struct {
				IconDownloadLink string
			}{
				IconDownloadLink: link,
			},
		)

		if mappingErr != nil {
			return nil, exceptions.WrapErrorWithTrackableException(mappingErr)
		}
	}

	return &interests.GetInterestsResponseDto{Interests: mappedInterests}, nil
}
