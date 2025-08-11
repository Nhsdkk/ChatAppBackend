package users

import (
	"chat_app_backend/application/models/jwt_claims"
	"chat_app_backend/application/models/users/login"
	"chat_app_backend/internal/exceptions"
	"chat_app_backend/internal/exceptions/common_exceptions"
	"chat_app_backend/internal/mapper"
	"chat_app_backend/internal/password"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/service_wrapper"
	"chat_app_backend/internal/sqlc/db_queries"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type LoginHandler struct{}

func (l LoginHandler) Handle(
	request *login.LoginRequestDto,
	service service_wrapper.IServiceWrapper,
	ctx *gin.Context,
	_ *request_env.RequestEnv,
) (*login.LoginResponseDto, exceptions.ITrackableException) {
	user, userExistenceError := service.GetDbConnection().GetQueries().GetUserByEmail(ctx, request.Email)

	switch {
	case errors.Is(userExistenceError, pgx.ErrNoRows):
		return nil, common_exceptions.InvalidBodyException{
			BaseRestException: exceptions.BaseRestException{
				ITrackableException: exceptions.WrapErrorWithTrackableException(userExistenceError),
				Message:             "invalid credentials",
			},
		}
	case userExistenceError != nil:
		return nil, exceptions.WrapErrorWithTrackableException(userExistenceError)
	}

	if !password.ComparePassword(request.Password, user.Password) {
		message := "invalid credentials"
		return nil, common_exceptions.InvalidBodyException{
			BaseRestException: exceptions.BaseRestException{
				ITrackableException: exceptions.CreateTrackableExceptionFromStringF(message),
				Message:             message,
			},
		}
	}

	interestsRaw, interestsQueryError := service.GetDbConnection().GetQueries().GetUserInterests(ctx, user.ID)
	if interestsQueryError != nil {
		return nil, exceptions.WrapErrorWithTrackableException(interestsQueryError)
	}

	var userClaims jwt_claims.UserClaims

	mappingErr := mapper.Mapper{}.Map(&userClaims, user)
	if mappingErr != nil {
		return nil, exceptions.WrapErrorWithTrackableException(mappingErr)
	}

	accessToken, refreshToken, tokenGenerationError := service.GetJwtHandler().GenerateJwtPair(userClaims)
	if tokenGenerationError != nil {
		return nil, exceptions.WrapErrorWithTrackableException(tokenGenerationError)
	}

	var response login.LoginResponseDto

	mappingErr = mapper.Mapper{}.Map(
		&response,
		user,
		struct {
			Interests    []db_queries.Interest
			AccessToken  string
			RefreshToken string
		}{
			Interests:    interestsRaw,
			AccessToken:  accessToken.GetToken(),
			RefreshToken: refreshToken.GetToken(),
		},
	)

	if mappingErr != nil {
		return nil, exceptions.WrapErrorWithTrackableException(mappingErr)
	}

	return &response, nil
}
