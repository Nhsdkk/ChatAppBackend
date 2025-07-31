package users

import (
	"chat_app_backend/application/models/exception"
	delete2 "chat_app_backend/application/models/users/delete"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/service_wrapper"
	"github.com/gin-gonic/gin"
)

type DeleteUserHandler struct{}

func (d DeleteUserHandler) Handle(
	_ *delete2.DeleteUserRequestDto,
	service service_wrapper.IServiceWrapper,
	ctx *gin.Context,
	requestEnvironment *request_env.RequestEnv,
) (*delete2.DeleteUserResponseDto, error) {
	user := *requestEnvironment.User
	userDeletionError := service.GetDbConnection().GetQueries().RemoveUser(ctx, user.ID)
	if userDeletionError != nil {
		return nil, exception.ServerException{
			Err: userDeletionError,
		}
	}
	return &delete2.DeleteUserResponseDto{}, nil
}
