package users

import (
	interests2 "chat_app_backend/application/models/interests/get_many_by_ids"
	"chat_app_backend/application/models/users/get_user_data"
	"chat_app_backend/internal/exceptions"
	"chat_app_backend/internal/exceptions/common_exceptions"
	"chat_app_backend/internal/mapper"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/service_wrapper"
	"chat_app_backend/internal/sqlc/db_queries"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type GetUserDataHandler struct{}

func (g GetUserDataHandler) Handle(
	request *get_user_data.GetUserDataRequestDto,
	services service_wrapper.IServiceWrapper,
	ctx *gin.Context,
	requestEnvironment *request_env.RequestEnv,
) (*get_user_data.GetUserDataResponseDto, exceptions.ITrackableException) {
	requestingUser := *requestEnvironment.User
	if requestingUser.Role == db_queries.RoleTypeUSER && requestingUser.ID != request.ID {
		message := fmt.Sprintf("can't get user with id %s", request.ID)
		return nil, common_exceptions.ForbiddenException{
			BaseRestException: exceptions.BaseRestException{
				ITrackableException: exceptions.CreateTrackableExceptionFromStringF(message),
				Message:             message,
			},
		}
	}

	user, userQueryError := services.GetDbConnection().GetQueries().GetUserById(ctx, request.ID)

	switch {
	case errors.Is(userQueryError, pgx.ErrNoRows):
		return nil, common_exceptions.ResourceNotFoundException{
			BaseRestException: exceptions.BaseRestException{
				ITrackableException: exceptions.WrapErrorWithTrackableException(userQueryError),
				Message:             "user not found",
			},
		}
	case userQueryError != nil:
		return nil, exceptions.WrapErrorWithTrackableException(userQueryError)
	}

	interestsRaw, interestsQueryError := services.GetDbConnection().GetQueries().GetUserInterests(ctx, request.ID)
	if interestsQueryError != nil {
		return nil, exceptions.WrapErrorWithTrackableException(interestsQueryError)
	}

	interests := make([]interests2.GetInterestsDto, len(interestsRaw))
	for idx, interestRaw := range interestsRaw {
		mapperError := mapper.Mapper{}.Map(
			&interests[idx],
			interestRaw,
		)

		if mapperError != nil {
			return nil, exceptions.WrapErrorWithTrackableException(mapperError)
		}
	}

	var response get_user_data.GetUserDataResponseDto
	_ = mapper.Mapper{}.Map(
		&response,
		user,
		struct {
			Interests []interests2.GetInterestsDto
		}{
			Interests: interests,
		},
	)
	return &response, nil
}
