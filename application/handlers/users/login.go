package users

import (
	interests2 "chat_app_backend/application/models/interests/get_many_by_ids"
	"chat_app_backend/application/models/jwt_claims"
	"chat_app_backend/application/models/users/login"
	"chat_app_backend/internal/exceptions"
	"chat_app_backend/internal/exceptions/common_exceptions"
	"chat_app_backend/internal/mapper"
	"chat_app_backend/internal/password"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/service_wrapper"
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

	interests := make([]interests2.GetInterestsDto, len(interestsRaw))
	for idx, interestRaw := range interestsRaw {
		mapperError := mapper.Mapper{}.Map(
			&interests[idx],
			interestRaw,
		)

		if mapperError != nil {
			return nil, exceptions.WrapErrorWithTrackableException(mapperError)
		}
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
			Interests    []interests2.GetInterestsDto
			AccessToken  string
			RefreshToken string
		}{
			Interests:    interests,
			AccessToken:  accessToken.GetToken(),
			RefreshToken: refreshToken.GetToken(),
		},
	)

	if mappingErr != nil {
		return nil, exceptions.WrapErrorWithTrackableException(mappingErr)
	}

	return &response, nil
}
