package users

import (
	"chat_app_backend/application/models/exception"
	"chat_app_backend/application/models/users/get_user_data"
	"chat_app_backend/internal/mapper"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/service_wrapper"
	"errors"
	"github.com/gin-gonic/gin"
)

type GetUserDataHandler struct{}

func (g GetUserDataHandler) Handle(
	_ *get_user_data.GetUserDataRequestDto,
	_ service_wrapper.IServiceWrapper,
	_ *gin.Context,
	requestEnvironment *request_env.RequestEnv,
) (*get_user_data.GetUserDataResponseDto, error) {
	if requestEnvironment.User == nil {
		return nil, exception.ServerException{
			Err: errors.New("user is null on authorized handler"),
		}
	}

	var response get_user_data.GetUserDataResponseDto
	_ = mapper.Mapper{}.Map(*requestEnvironment.User, &response)
	return &response, nil
}
