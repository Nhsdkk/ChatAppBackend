package users

import (
	"chat_app_backend/application/models/exception"
	"chat_app_backend/application/models/users/refresh_token"
	"chat_app_backend/internal/jwt"
	"chat_app_backend/internal/mapper"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/service_wrapper"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type RefreshTokenHandler struct{}

func (r RefreshTokenHandler) Handle(
	request *refresh_token.RefreshTokenRequestDto,
	service service_wrapper.IServiceWrapper,
	ctx *gin.Context,
	_ *request_env.RequestEnv,
) (*refresh_token.RefreshTokenResponseDto, error) {
	token := jwt.CreateTokenFromHandlerAndString(service.GetJwtHandler(), request.RefreshToken, jwt.RefreshToken)

	validToken, validationErr := token.Validate()
	if validationErr != nil {
		return nil, &exception.UnauthorizedException{
			Err: validationErr,
		}
	}

	user, err := service.
		GetDbConnection().
		GetQueries().
		GetUserById(ctx, validToken.GetClaims().ID)

	switch {
	case err != nil && errors.Is(err, pgx.ErrNoRows):
		return nil, &exception.ResourceNotFoundException{
			Err: errors.New("owner of the token not found"),
		}
	case err != nil:
		return nil, err
	}

	if !validToken.GetClaims().Equals(&user) {
		return nil, exception.UnauthorizedException{
			Err: errors.New("claims and user data does not match"),
		}
	}

	accessToken, accessTokenGenerationError := validToken.RefreshRelatedAccessToken(service.GetJwtHandler())
	if accessTokenGenerationError != nil {
		return nil, accessTokenGenerationError
	}

	var response refresh_token.RefreshTokenResponseDto
	mappingErr := mapper.Mapper{}.Map(
		&response,
		user,
		struct {
			AccessToken string
		}{
			AccessToken: accessToken.GetToken(),
		},
	)
	if mappingErr != nil {
		return nil, mappingErr
	}

	return &response, nil
}
