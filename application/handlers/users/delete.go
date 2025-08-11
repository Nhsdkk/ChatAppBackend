package users

import (
	delete2 "chat_app_backend/application/models/users/delete"
	"chat_app_backend/internal/exceptions"
	"chat_app_backend/internal/exceptions/common_exceptions"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/service_wrapper"
	"chat_app_backend/internal/sqlc/db_queries"
	"fmt"

	"github.com/gin-gonic/gin"
)

type DeleteUserHandler struct{}

func (d DeleteUserHandler) Handle(
	request *delete2.DeleteUserRequestDto,
	service service_wrapper.IServiceWrapper,
	ctx *gin.Context,
	requestEnvironment *request_env.RequestEnv,
) (*delete2.DeleteUserResponseDto, exceptions.ITrackableException) {
	user := *requestEnvironment.User
	if user.Role == db_queries.RoleTypeUSER && user.ID != request.ID {
		message := fmt.Sprintf("can't delete user with id %s", request.ID)
		return nil, common_exceptions.ForbiddenException{
			BaseRestException: exceptions.BaseRestException{
				ITrackableException: exceptions.CreateTrackableExceptionFromStringF(message),
				Message:             message,
			},
		}
	}

	userDeletionError := service.GetDbConnection().GetQueries().RemoveUser(ctx, request.ID)
	if userDeletionError != nil {
		return nil, exceptions.WrapErrorWithTrackableException(userDeletionError)
	}
	return &delete2.DeleteUserResponseDto{}, nil
}
