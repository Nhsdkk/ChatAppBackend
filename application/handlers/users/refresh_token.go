package users

import (
	"chat_app_backend/application/models/users/refresh_token"
	"chat_app_backend/internal/exceptions"
	"chat_app_backend/internal/exceptions/common_exceptions"
	"chat_app_backend/internal/jwt"
	"chat_app_backend/internal/mapper"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/s3"
	"chat_app_backend/internal/service_wrapper"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type RefreshTokenHandler struct{}

func (r RefreshTokenHandler) Handle(
	request *refresh_token.RefreshTokenRequestDto,
	services service_wrapper.IServiceWrapper,
	ctx *gin.Context,
	_ *request_env.RequestEnv,
) (*refresh_token.RefreshTokenResponseDto, exceptions.ITrackableException) {
	token := jwt.CreateTokenFromHandlerAndString(services.GetJwtHandler(), request.RefreshToken, jwt.RefreshToken)

	validToken, validationErr := token.Validate()
	if validationErr != nil {
		return nil, common_exceptions.UnauthorizedException{
			BaseRestException: exceptions.BaseRestException{
				ITrackableException: exceptions.WrapErrorWithTrackableException(validationErr),
				Message:             "",
			},
		}
	}

	user, err := services.
		GetDbConnection().
		GetQueries().
		GetUserById(ctx, validToken.GetClaims().ID)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, common_exceptions.ResourceNotFoundException{
			BaseRestException: exceptions.BaseRestException{
				ITrackableException: exceptions.WrapErrorWithTrackableException(err),
				Message:             "owner of the token not found",
			},
		}
	case err != nil:
		return nil, exceptions.WrapErrorWithTrackableException(err)
	}

	if !validToken.GetClaims().Equals(&user) {
		message := "claims and user data does not match"
		return nil, common_exceptions.UnauthorizedException{
			BaseRestException: exceptions.BaseRestException{
				ITrackableException: exceptions.CreateTrackableExceptionFromStringF(message),
				Message:             message,
			},
		}
	}

	accessToken, accessTokenGenerationError := validToken.RefreshRelatedAccessToken(services.GetJwtHandler())
	if accessTokenGenerationError != nil {
		return nil, exceptions.WrapErrorWithTrackableException(accessTokenGenerationError)
	}

	avatarDownloadLink, downloadLinkGenerationError := services.GetS3Client().
		GetDownloadUrl(ctx, user.AvatarFileName, s3.AvatarBucket)

	if downloadLinkGenerationError != nil {
		return nil, exceptions.WrapErrorWithTrackableException(downloadLinkGenerationError)
	}

	var response refresh_token.RefreshTokenResponseDto
	mappingErr := mapper.Mapper{}.Map(
		&response,
		user,
		struct {
			AccessToken        string
			AvatarDownloadLink string
		}{
			AvatarDownloadLink: avatarDownloadLink,
			AccessToken:        accessToken.GetToken(),
		},
	)
	if mappingErr != nil {
		return nil, exceptions.WrapErrorWithTrackableException(mappingErr)
	}

	return &response, nil
}
