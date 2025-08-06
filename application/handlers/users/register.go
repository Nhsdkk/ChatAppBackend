package users

import (
	interests2 "chat_app_backend/application/models/interests/get_many_by_ids"
	"chat_app_backend/application/models/jwt_claims"
	"chat_app_backend/application/models/users/register"
	"chat_app_backend/internal/exceptions"
	"chat_app_backend/internal/mapper"
	"chat_app_backend/internal/password"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/service_wrapper"
	"chat_app_backend/internal/sqlc/db_queries"

	"github.com/gin-gonic/gin"
)

type RegisterHandler struct{}

func (r RegisterHandler) Handle(
	request *register.RegisterRequestDto,
	services service_wrapper.IServiceWrapper,
	ctx *gin.Context,
	_ *request_env.RequestEnv,
) (*register.RegisterResponseDto, exceptions.ITrackableException) {
	var response register.RegisterResponseDto

	transactionError := services.
		GetDbConnection().
		CreateTransaction(ctx, func(queries *db_queries.Queries) exceptions.ITrackableException {
			createUserParams := db_queries.CreateUserParams{
				FullName: request.FullName,
				Birthday: request.Birthday,
				Gender:   request.Gender,
				Email:    request.Email,
				Password: password.HashPassword(request.Password),
				// TODO(issue #5): fill with real avatar
				AvatarFileName: "avatar.png",
			}

			user, createUserError := queries.CreateUser(ctx, createUserParams)
			if createUserError != nil {
				return exceptions.WrapErrorWithTrackableException(createUserError)
			}

			assignInterestsParams := db_queries.AssignInterestsToUserParams{
				UserID:      user.ID,
				InterestIds: request.Interests,
			}

			if assignInterestError := queries.AssignInterestsToUser(ctx, assignInterestsParams); assignInterestError != nil {
				return exceptions.WrapErrorWithTrackableException(assignInterestError)
			}

			interests, getInterestsError := queries.GetUserInterests(ctx, user.ID)
			if getInterestsError != nil {
				return exceptions.WrapErrorWithTrackableException(getInterestsError)
			}

			interestsMapped := make([]interests2.GetInterestsDto, len(interests))

			for idx, interest := range interests {
				err := mapper.Mapper{}.Map(
					&interestsMapped[idx],
					interest,
				)

				if err != nil {
					return exceptions.WrapErrorWithTrackableException(err)
				}
			}

			var claims jwt_claims.UserClaims
			mappingErr := mapper.Mapper{}.Map(&claims, user)
			if mappingErr != nil {
				return exceptions.WrapErrorWithTrackableException(mappingErr)
			}

			accessToken, refreshToken, tokenGenerationError := services.
				GetJwtHandler().
				GenerateJwtPair(claims)

			if tokenGenerationError != nil {
				return exceptions.WrapErrorWithTrackableException(tokenGenerationError)
			}

			mappingError := mapper.Mapper{}.Map(
				&response,
				user,
				struct {
					Interests    []interests2.GetInterestsDto
					AccessToken  string
					RefreshToken string
				}{
					Interests:    interestsMapped,
					AccessToken:  accessToken.GetToken(),
					RefreshToken: refreshToken.GetToken(),
				},
			)

			if mappingError != nil {
				return exceptions.WrapErrorWithTrackableException(mappingError)
			}

			return nil
		})

	if transactionError != nil {
		return nil, transactionError
	}

	return &response, nil
}
