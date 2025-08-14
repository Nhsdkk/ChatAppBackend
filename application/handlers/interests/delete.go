package interests

import (
	delete2 "chat_app_backend/application/models/interests/delete"
	"chat_app_backend/internal/exceptions"
	"chat_app_backend/internal/exceptions/common_exceptions"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/service_wrapper"
	"chat_app_backend/internal/sqlc/db_queries"

	"github.com/gin-gonic/gin"
)

type DeleteInterestHandler struct{}

func (d DeleteInterestHandler) Handle(
	request *delete2.DeleteInterestRequestDto,
	services service_wrapper.IServiceWrapper,
	ctx *gin.Context,
	requestEnvironment *request_env.RequestEnv,
) (*delete2.DeleteInterestResponseDto, exceptions.ITrackableException) {
	if requestEnvironment.User.Role != db_queries.RoleTypeADMIN {
		return nil, common_exceptions.ForbiddenException{
			BaseRestException: exceptions.BaseRestException{
				ITrackableException: exceptions.CreateTrackableExceptionFromStringF("interest deletion forbidden"),
				Message:             "only admins can access this route",
			},
		}
	}

	deletionError := services.GetDbConnection().
		GetQueries().
		DeleteInterest(ctx, request.ID)

	if deletionError != nil {
		return nil, exceptions.WrapErrorWithTrackableException(deletionError)
	}

	return &delete2.DeleteInterestResponseDto{}, nil
}
