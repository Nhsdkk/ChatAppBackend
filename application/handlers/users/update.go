package users

import (
	"chat_app_backend/application/models/jwt_claims"
	"chat_app_backend/application/models/users/update"
	"chat_app_backend/internal/exceptions"
	"chat_app_backend/internal/exceptions/common_exceptions"
	"chat_app_backend/internal/mapper"
	"chat_app_backend/internal/password"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/service_wrapper"
	"chat_app_backend/internal/sqlc/db_queries"
	"fmt"

	"github.com/gin-gonic/gin"
)

type UpdateUserHandler struct{}

func (u UpdateUserHandler) Handle(
	request *update.UpdateUserRequestDto,
	service service_wrapper.IServiceWrapper,
	ctx *gin.Context,
	requestEnvironment *request_env.RequestEnv,
) (*update.UpdateUserResponseDto, exceptions.ITrackableException) {
	user := *requestEnvironment.User

	if user.Role == db_queries.RoleTypeUSER && (user.ID != request.ID || request.Role != nil) {
		message := fmt.Sprintf("can't update user with id %s", request.ID)
		return nil, common_exceptions.ForbiddenException{
			BaseRestException: exceptions.BaseRestException{
				ITrackableException: exceptions.CreateTrackableExceptionFromStringF(message),
				Message:             message,
			},
		}
	}

	var updateUserParams db_queries.UpdateUserParams

	var newPasswordBytes *[]byte
	if request.PasswordString != nil {
		newPasswordBytes = new([]byte)
		*newPasswordBytes = password.HashPassword(*request.PasswordString)
	}

	nullGender := db_queries.NullGender{}
	if request.Gender != nil {
		nullGender.Gender = *request.Gender
		nullGender.Valid = true
	}

	nullRole := db_queries.NullRoleType{}
	if request.Role != nil {
		nullRole.RoleType = *request.Role
		nullRole.Valid = true
	}

	mapperError := mapper.Mapper{}.Map(
		&updateUserParams,
		*request,
		struct {
			Password       *[]byte
			AvatarFileName *string
			Gender         db_queries.NullGender
			Role           db_queries.NullRoleType
			Online         *bool
		}{
			Password: newPasswordBytes,
			//TODO(issue #5): add avatar upload
			AvatarFileName: nil,
			Gender:         nullGender,
			Online:         nil,
			Role:           nullRole,
		},
	)

	if mapperError != nil {
		return nil, exceptions.WrapErrorWithTrackableException(mapperError)
	}

	newUser, updateUserError := service.GetDbConnection().GetQueries().UpdateUser(ctx, updateUserParams)
	if updateUserError != nil {
		return nil, exceptions.WrapErrorWithTrackableException(updateUserError)
	}

	var claims jwt_claims.UserClaims
	claimsMappingError := mapper.Mapper{}.Map(
		&claims,
		newUser,
	)
	if claimsMappingError != nil {
		return nil, exceptions.WrapErrorWithTrackableException(claimsMappingError)
	}

	accessToken, refreshToken, tokenGenerationError := service.GetJwtHandler().GenerateJwtPair(claims)
	if tokenGenerationError != nil {
		return nil, exceptions.WrapErrorWithTrackableException(tokenGenerationError)
	}

	var response update.UpdateUserResponseDto
	responseMappingError := mapper.Mapper{}.Map(
		&response,
		newUser,
		struct {
			AccessToken  string
			RefreshToken string
		}{
			AccessToken:  accessToken.GetToken(),
			RefreshToken: refreshToken.GetToken(),
		},
	)
	if responseMappingError != nil {
		return nil, exceptions.WrapErrorWithTrackableException(responseMappingError)
	}

	return &response, nil

}
