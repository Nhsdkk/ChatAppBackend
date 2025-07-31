package users

import (
	"chat_app_backend/application/models/exception"
	"chat_app_backend/application/models/jwt_claims"
	"chat_app_backend/application/models/users/update"
	"chat_app_backend/internal/mapper"
	"chat_app_backend/internal/password"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/service_wrapper"
	"chat_app_backend/internal/sqlc/db_queries"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateUserHandler struct{}

func (u UpdateUserHandler) Handle(
	request *update.UpdateUserRequestDto,
	service service_wrapper.IServiceWrapper,
	ctx *gin.Context,
	requestEnvironment *request_env.RequestEnv,
) (*update.UpdateUserResponseDto, error) {
	user := *requestEnvironment.User

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

	mapperError := mapper.Mapper{}.Map(
		&updateUserParams,
		*request,
		struct {
			ID             uuid.UUID
			Password       *[]byte
			AvatarFileName *string
			Gender         db_queries.NullGender
			Online         *bool
		}{
			ID:             user.ID,
			Password:       newPasswordBytes,
			AvatarFileName: nil,
			Gender:         nullGender,
			Online:         nil,
		},
	)

	if mapperError != nil {
		return nil, exception.ServerException{
			Err: mapperError,
		}
	}

	newUser, updateUserError := service.GetDbConnection().GetQueries().UpdateUser(ctx, updateUserParams)
	if updateUserError != nil {
		return nil, exception.ServerException{
			Err: updateUserError,
		}
	}

	var claims jwt_claims.UserClaims
	claimsMappingError := mapper.Mapper{}.Map(
		&claims,
		newUser,
	)
	if claimsMappingError != nil {
		return nil, exception.ServerException{
			Err: claimsMappingError,
		}
	}

	accessToken, refreshToken, tokenGenerationError := service.GetJwtHandler().GenerateJwtPair(claims)
	if tokenGenerationError != nil {
		return nil, exception.ServerException{
			Err: tokenGenerationError,
		}
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
		return nil, exception.ServerException{
			Err: responseMappingError,
		}
	}

	return &response, nil

}
