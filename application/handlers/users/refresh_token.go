package users

import (
	"chat_app_backend/application/models/exception"
	"chat_app_backend/application/models/users/refresh_token"
	"chat_app_backend/internal/jwt"
	"chat_app_backend/internal/mapper"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/service_wrapper"
	"chat_app_backend/internal/sqlc/db_queries"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"time"
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

	accessToken, accessTokenGenerationError := validToken.RefreshRelatedAccessToken(service.GetJwtHandler())
	if accessTokenGenerationError != nil {
		return nil, accessTokenGenerationError
	}

	var response refresh_token.RefreshTokenResponseDto
	mappingErr := mapper.Mapper{}.Map(
		struct {
			ID             uuid.UUID
			FullName       string
			Birthday       time.Time
			Gender         db_queries.Gender
			Email          string
			AvatarFileName string
			Online         bool
			EmailVerified  bool
			LastSeen       time.Time
			CreatedAt      time.Time
			UpdatedAt      time.Time
			AccessToken    string
		}{
			ID:             user.ID,
			FullName:       user.FullName,
			Birthday:       user.Birthday,
			Gender:         user.Gender,
			Email:          user.Email,
			AvatarFileName: user.AvatarFileName,
			Online:         user.Online,
			EmailVerified:  user.EmailVerified,
			LastSeen:       user.LastSeen,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
			AccessToken:    accessToken.GetToken(),
		},
		&response,
	)
	if mappingErr != nil {
		return nil, mappingErr
	}

	return &response, nil
}
