package interests

import (
	shared_interests "chat_app_backend/application/handlers/shared/interests"
	"chat_app_backend/application/models/interests/assign"
	"chat_app_backend/application/models/interests/get"
	"chat_app_backend/internal/exceptions"
	"chat_app_backend/internal/exceptions/common_exceptions"
	"chat_app_backend/internal/mapper"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/service_wrapper"
	"chat_app_backend/internal/sqlc/db_queries"

	"github.com/gin-gonic/gin"
)

type AssignInterestHandler struct{}

func (a AssignInterestHandler) Handle(
	request *assign.AssignInterestRequestDto,
	services service_wrapper.IServiceWrapper,
	ctx *gin.Context,
	requestEnvironment *request_env.RequestEnv,
) (*assign.AssignInterestResponseDto, exceptions.ITrackableException) {
	if requestEnvironment.User.Role != db_queries.RoleTypeADMIN && request.UserID != requestEnvironment.User.ID {
		return nil, &common_exceptions.ForbiddenException{
			BaseRestException: exceptions.BaseRestException{
				ITrackableException: exceptions.CreateTrackableExceptionFromStringF("user role assignment forbidden"),
				Message:             "not enough privileges to update this user",
			},
		}
	}

	params := db_queries.GetManyInterestsByFiltersParams{Ids: request.InterestIds}

	interests, getError := services.GetDbConnection().
		GetQueries().
		GetManyInterestsByFilters(
			ctx,
			params,
		)

	if getError != nil {
		return nil, exceptions.WrapErrorWithTrackableException(getError)
	}

	if len(interests) != len(request.InterestIds) {
		return nil, &common_exceptions.ResourceNotFoundException{
			BaseRestException: exceptions.BaseRestException{
				ITrackableException: exceptions.CreateTrackableExceptionFromStringF("some interests not found"),
				Message:             "some interests not found",
			},
		}
	}

	assignError := services.GetDbConnection().
		CreateTransaction(
			ctx,
			func(queries *db_queries.Queries) exceptions.ITrackableException {
				if removeError := queries.RemoveUserInterests(ctx, request.UserID); removeError != nil {
					return exceptions.WrapErrorWithTrackableException(removeError)
				}

				var assignParams db_queries.AssignInterestsToUserParams

				paramMappingError := mapper.Mapper{}.Map(&assignParams, *request)

				if paramMappingError != nil {
					return exceptions.WrapErrorWithTrackableException(paramMappingError)
				}

				if assignError := queries.AssignInterestsToUser(ctx, assignParams); assignError != nil {
					return exceptions.WrapErrorWithTrackableException(assignError)
				}

				return nil
			},
		)

	if assignError != nil {
		return nil, assignError
	}

	var response assign.AssignInterestResponseDto

	interestsWithIcons, iconsGetError := shared_interests.GetInterestIcons(interests, services.GetS3Client(), ctx)

	if iconsGetError != nil {
		return nil, iconsGetError
	}

	responseMappingError := mapper.Mapper{}.Map(&response, struct {
		Interests []get.GetInterestResponseDto
	}{
		Interests: interestsWithIcons,
	})

	if responseMappingError != nil {
		return nil, exceptions.WrapErrorWithTrackableException(responseMappingError)
	}

	return &response, nil
}
