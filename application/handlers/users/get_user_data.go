package users

import (
	"chat_app_backend/application/models/exception"
	interests2 "chat_app_backend/application/models/interests/get_many_by_ids"
	"chat_app_backend/application/models/users/get_user_data"
	"chat_app_backend/internal/mapper"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/service_wrapper"
	"github.com/gin-gonic/gin"
)

type GetUserDataHandler struct{}

func (g GetUserDataHandler) Handle(
	_ *get_user_data.GetUserDataRequestDto,
	services service_wrapper.IServiceWrapper,
	ctx *gin.Context,
	requestEnvironment *request_env.RequestEnv,
) (*get_user_data.GetUserDataResponseDto, error) {
	interestsRaw, interestsQueryError := services.GetDbConnection().GetQueries().GetUserInterests(ctx, requestEnvironment.User.ID)
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
		*requestEnvironment.User,
		struct {
			Interests []interests2.GetInterestsDto
		}{
			Interests: interests,
		},
	)
	return &response, nil
}
