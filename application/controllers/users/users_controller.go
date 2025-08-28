package users

import (
	"chat_app_backend/application/controllers/users/validators"
	"chat_app_backend/application/handlers/users"
	delete2 "chat_app_backend/application/models/users/delete"
	"chat_app_backend/application/models/users/get_user_data"
	"chat_app_backend/application/models/users/login"
	"chat_app_backend/application/models/users/refresh_token"
	"chat_app_backend/application/models/users/register"
	"chat_app_backend/application/models/users/update"
	"chat_app_backend/internal/router"
	"chat_app_backend/internal/service_wrapper"
	"chat_app_backend/internal/validator"
	"time"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	router.Controller
}

func CreateUserController(engine *gin.Engine, serviceWrapper service_wrapper.IServiceWrapper) (uc UserController) {
	uc.Controller = router.CreateController(
		engine,
		"/users",
		[]router.IRoute{
			router.CreateBaseRoute(
				serviceWrapper,
				"/register",
				users.RegisterHandler{}.Handle,
				validator.
					Validator[register.RegisterRequestDto]{}.
					AttachValidator(
						validator.ExternalValidator[register.RegisterRequestDto, string]{}.
							RuleFor(
								func(data *register.RegisterRequestDto) *string {
									return &data.Email
								},
							).
							Must(validators.EmailValidator{}).
							WithMessage("email is of wrong format").
							Validate,
					).
					AttachValidator(
						validator.ExternalValidator[register.RegisterRequestDto, time.Time]{}.
							RuleFor(
								func(data *register.RegisterRequestDto) *time.Time {
									return &data.Birthday
								},
							).
							Must(validators.BirthDateValidator{}).
							WithMessage("you are not old enough to register").
							Validate,
					).
					AttachValidator(
						validator.ExternalValidator[register.RegisterRequestDto, string]{}.
							RuleFor(
								func(data *register.RegisterRequestDto) *string {
									return &data.Password
								},
							).
							Must(validators.PasswordValidator{}).
							WithMessage("password should have at least one of each of this characters (special characters, upper and lowercase letters, digits)").
							Validate,
					),
				router.POST,
			),
			router.CreateBaseRoute(
				serviceWrapper,
				"/login",
				users.LoginHandler{}.Handle,
				validator.
					Validator[login.LoginRequestDto]{},
				router.POST,
			),
			&router.AuthorizedRoute[get_user_data.GetUserDataRequestDto, get_user_data.GetUserDataResponseDto]{
				Route: router.CreateBaseRoute(
					serviceWrapper,
					"/:id",
					users.GetUserDataHandler{}.Handle,
					validator.
						Validator[get_user_data.GetUserDataRequestDto]{},
					router.GET,
				),
			},
			router.CreateBaseRoute(
				serviceWrapper,
				"/refresh_token",
				users.RefreshTokenHandler{}.Handle,
				validator.
					Validator[refresh_token.RefreshTokenRequestDto]{},
				router.POST,
			),
			&router.AuthorizedRoute[delete2.DeleteUserRequestDto, delete2.DeleteUserResponseDto]{
				Route: router.CreateBaseRoute(
					serviceWrapper,
					"/:id",
					users.DeleteUserHandler{}.Handle,
					validator.
						Validator[delete2.DeleteUserRequestDto]{},
					router.DELETE,
				),
			},
			&router.AuthorizedRoute[update.UpdateUserRequestDto, update.UpdateUserResponseDto]{
				Route: router.CreateBaseRoute(
					serviceWrapper,
					"/:id",
					users.UpdateUserHandler{}.Handle,
					validator.
						Validator[update.UpdateUserRequestDto]{}.
						AttachValidator(
							validator.ExternalValidator[update.UpdateUserRequestDto, string]{}.
								RuleFor(
									func(data *update.UpdateUserRequestDto) *string {
										return data.Email
									},
								).
								Must(validators.EmailValidator{}).
								WithMessage("email is of wrong format").
								Optional().
								Validate,
						).
						AttachValidator(
							validator.ExternalValidator[update.UpdateUserRequestDto, time.Time]{}.
								RuleFor(
									func(data *update.UpdateUserRequestDto) *time.Time {
										return data.Birthday
									},
								).
								Must(validators.BirthDateValidator{}).
								WithMessage("you are not old enough to register").
								Optional().
								Validate,
						).
						AttachValidator(
							validator.ExternalValidator[update.UpdateUserRequestDto, string]{}.
								RuleFor(
									func(data *update.UpdateUserRequestDto) *string {
										return data.PasswordString
									},
								).
								Must(validators.PasswordValidator{}).
								WithMessage("password should have at least one of each of this characters (special characters, upper and lowercase letters, digits)").
								Optional().
								Validate,
						),
					router.PUT,
				),
			},
		},
	)

	return uc
}
