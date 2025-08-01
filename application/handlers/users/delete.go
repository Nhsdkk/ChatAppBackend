package users

import (
	"chat_app_backend/application/models/exception"
	delete2 "chat_app_backend/application/models/users/delete"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/service_wrapper"
	"chat_app_backend/internal/sqlc/db_queries"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type DeleteUserHandler struct{}

func (d DeleteUserHandler) Handle(
	request *delete2.DeleteUserRequestDto,
	service service_wrapper.IServiceWrapper,
	ctx *gin.Context,
	requestEnvironment *request_env.RequestEnv,
) (*delete2.DeleteUserResponseDto, error) {
	user := *requestEnvironment.User
	if user.Role == db_queries.RoleTypeUSER && user.ID != request.ID {
		return nil, exception.ForbiddenException{
			Err: errors.New(fmt.Sprintf("can't delete user with id %s", request.ID)),
		}
	}

	userDeletionError := service.GetDbConnection().GetQueries().RemoveUser(ctx, request.ID)
	if userDeletionError != nil {
		return nil, exception.ServerException{
			Err: userDeletionError,
		}
	}
	return &delete2.DeleteUserResponseDto{}, nil
}
