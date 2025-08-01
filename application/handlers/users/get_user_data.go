package users

import (
	"chat_app_backend/application/models/exception"
	interests2 "chat_app_backend/application/models/interests/get_many_by_ids"
	"chat_app_backend/application/models/users/get_user_data"
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
) (*get_user_data.GetUserDataResponseDto, error) {
	requestingUser := *requestEnvironment.User
	if requestingUser.Role == db_queries.RoleTypeUSER && requestingUser.ID != request.ID {
		return nil, exception.ForbiddenException{
			Err: errors.New(fmt.Sprintf("can't get user with id %s", request.ID)),
		}
	}

	user, userQueryError := services.GetDbConnection().GetQueries().GetUserById(ctx, request.ID)

	switch {
	case userQueryError != nil && errors.Is(userQueryError, pgx.ErrNoRows):
		return nil, exception.ResourceNotFoundException{
			Err: errors.New("user not found"),
		}
	case userQueryError != nil:
		return nil, exception.ServerException{
			Err: userQueryError,
		}
	}

	interestsRaw, interestsQueryError := services.GetDbConnection().GetQueries().GetUserInterests(ctx, request.ID)
	if interestsQueryError != nil {
		return nil, exception.ServerException{
			Err: interestsQueryError,
		}
	}

	interests := make([]interests2.GetInterestsDto, len(interestsRaw))
	for idx, interestRaw := range interestsRaw {
		mapperError := mapper.Mapper{}.Map(
			&interests[idx],
			interestRaw,
		)

		if mapperError != nil {
			return nil, exception.ServerException{
				Err: mapperError,
			}
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
