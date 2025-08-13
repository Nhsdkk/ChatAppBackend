package users

import (
	sharedinterests "chat_app_backend/application/handlers/shared/interests"
	interests "chat_app_backend/application/models/interests/get"
	"chat_app_backend/application/models/users/get_user_data"
	"chat_app_backend/internal/exceptions"
	"chat_app_backend/internal/exceptions/common_exceptions"
	"chat_app_backend/internal/mapper"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/s3"
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

	rawInterests, interestsQueryError := services.GetDbConnection().GetQueries().GetUserInterests(ctx, request.ID)
	if interestsQueryError != nil {
		return nil, exceptions.WrapErrorWithTrackableException(interestsQueryError)
	}

	avatarDownloadLink, downloadLinkGenerationError := services.GetS3Client().
		GetDownloadUrl(ctx, user.AvatarFileName, s3.AvatarsBucket)

	if downloadLinkGenerationError != nil {
		return nil, exceptions.WrapErrorWithTrackableException(downloadLinkGenerationError)
	}

	mappedInterests, err := sharedinterests.GetInterestIcons(rawInterests, services.GetS3Client(), ctx)
	if err != nil {
		return nil, err
	}

	var response get_user_data.GetUserDataResponseDto
	_ = mapper.Mapper{}.Map(
		&response,
		user,
		struct {
			Interests          []interests.GetInterestResponseDto
			AvatarDownloadLink string
		}{
			Interests:          mappedInterests,
			AvatarDownloadLink: avatarDownloadLink,
		},
	)
	return &response, nil
}
