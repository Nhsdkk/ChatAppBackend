package users

import (
	interests2 "chat_app_backend/application/models/interests/get_many_by_ids"
	"chat_app_backend/application/models/jwt_claims"
	"chat_app_backend/application/models/users/register"
	"chat_app_backend/internal/mapper"
	"chat_app_backend/internal/password"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/service_wrapper"
	"chat_app_backend/internal/sqlc/db_queries"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

type RegisterHandler struct{}

func (r RegisterHandler) Handle(
	request *register.RegisterRequestDto,
	services service_wrapper.IServiceWrapper,
	ctx *gin.Context,
	_ *request_env.RequestEnv,
) (*register.RegisterResponseDto, error) {
	var response register.RegisterResponseDto

	transactionError := services.
		GetDbConnection().
		CreateTransaction(ctx, func(queries *db_queries.Queries) error {
			createUserParams := db_queries.CreateUserParams{
				FullName: request.FullName,
				Birthday: request.Birthday,
				Gender:   request.Gender,
				Email:    request.Email,
				Password: password.HashPassword(request.Password),
				// TODO: fill with real avatar
				AvatarFileName: "avatar.png",
			}

			user, createUserError := queries.CreateUser(ctx, createUserParams)
			if createUserError != nil {
				return createUserError
			}

			assignInterestsParams := db_queries.AssignInterestsToUserParams{
				UserID:      user.ID,
				InterestIds: request.Interests,
			}

			if assignInterestError := queries.AssignInterestsToUser(ctx, assignInterestsParams); assignInterestError != nil {
				return assignInterestError
			}

			interests, getInterestsError := queries.GetManyInterestsById(ctx, request.Interests)
			if getInterestsError != nil {
				return getInterestsError
			}

			interestsMapped := make([]interests2.GetInterestsDto, len(interests))

			for idx, interest := range interests {
				err := mapper.Mapper{}.Map(
					interest,
					&interestsMapped[idx],
				)

				if err != nil {
					return err
				}
			}

			var claims jwt_claims.UserClaims
			mappingErr := mapper.Mapper{}.Map(user, &claims)
			if mappingErr != nil {
				return mappingErr
			}

			accessToken, refreshToken, tokenGenerationError := services.
				GetJwtHandler().
				GenerateJwtPair(claims)

			if tokenGenerationError != nil {
				return tokenGenerationError
			}

			mappingError := mapper.Mapper{}.Map(
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
					Interests      []interests2.GetInterestsDto
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
					Interests:      interestsMapped,
					AccessToken:    accessToken.GetToken(),
					RefreshToken:   refreshToken.GetToken(),
				},
				&response,
			)

			if mappingError != nil {
				return mappingError
			}

			return nil
		})

	if transactionError != nil {
		return nil, transactionError
	}

	return &response, nil
}
