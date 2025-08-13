package interests

import (
	interests2 "chat_app_backend/application/handlers/shared/interests"
	interests "chat_app_backend/application/models/interests/get"
	"chat_app_backend/internal/exceptions"
	"chat_app_backend/internal/request_env"
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

	mappedInterests, err := interests2.GetInterestIcons(rawInterests, service.GetS3Client(), ctx)
	if err != nil {
		return nil, err
	}

	return &interests.GetInterestsResponseDto{Interests: mappedInterests}, nil
}
