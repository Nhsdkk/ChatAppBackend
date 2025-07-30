package users

import (
	"chat_app_backend/application/models/exception"
	"chat_app_backend/application/models/jwt_claims"
	"chat_app_backend/application/models/users/login"
	"chat_app_backend/internal/mapper"
	"chat_app_backend/internal/password"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/service_wrapper"
	"chat_app_backend/internal/sqlc/db_queries"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"time"
)

type LoginHandler struct{}

func (l LoginHandler) Handle(
	request *login.LoginRequestDto,
	service service_wrapper.IServiceWrapper,
	ctx *gin.Context,
	_ *request_env.RequestEnv,
) (*login.LoginResponseDto, error) {
	user, err := service.GetDbConnection().GetQueries().GetUserByEmail(ctx, request.Email)

	switch {
	case err != nil && errors.Is(err, pgx.ErrNoRows):
		return nil, exception.InvalidBodyException{
			Err: errors.New("invalid credentials"),
		}
	case err != nil:
		return nil, err
	}

	if !password.ComparePassword(request.Password, user.Password) {
		return nil, exception.InvalidBodyException{
			Err: errors.New("invalid credentials"),
		}
	}

	var userClaims jwt_claims.UserClaims

	mappingErr := mapper.Mapper{}.Map(user, &userClaims)
	if mappingErr != nil {
		return nil, mappingErr
	}

	accessToken, refreshToken, tokenGenerationError := service.GetJwtHandler().GenerateJwtPair(userClaims)
	if tokenGenerationError != nil {
		return nil, tokenGenerationError
	}

	var response login.LoginResponseDto

	mappingErr = mapper.Mapper{}.Map(
		struct {
			ID             uuid.UUID
			FullName       string
			Birthday       time.Time
			Gender         db_queries.Gender
			Email          string
			Password       []byte
			AvatarFileName string
			Online         bool
			EmailVerified  bool
			LastSeen       time.Time
			CreatedAt      time.Time
			UpdatedAt      time.Time
			AccessToken    string
			RefreshToken   string
		}{
			ID:             user.ID,
			FullName:       user.FullName,
			Birthday:       user.Birthday,
			Gender:         user.Gender,
			Email:          user.Email,
			Password:       user.Password,
			AvatarFileName: user.AvatarFileName,
			Online:         user.Online,
			EmailVerified:  user.EmailVerified,
			LastSeen:       user.LastSeen,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
			AccessToken:    accessToken.GetToken(),
			RefreshToken:   refreshToken.GetToken(),
		},
		&response,
	)

	if mappingErr != nil {
		return nil, mappingErr
	}

	return &response, nil
}
