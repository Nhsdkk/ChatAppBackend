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
					AttachValidator(func(data *register.RegisterRequestDto) error { return validators.ValidateEmail(data.Email) }).
					AttachValidator(func(data *register.RegisterRequestDto) error { return validators.ValidateBirthDate(data.Birthday) }).
					AttachValidator(func(data *register.RegisterRequestDto) error { return validators.ValidatePassword(data.Password) }),
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
						AttachValidator(func(data *update.UpdateUserRequestDto) error {
							if data.PasswordString == nil {
								return nil
							}
							return validators.ValidatePassword(*data.PasswordString)
						}).
						AttachValidator(func(data *update.UpdateUserRequestDto) error {
							if data.Email == nil {
								return nil
							}
							return validators.ValidateEmail(*data.Email)
						}).
						AttachValidator(func(data *update.UpdateUserRequestDto) error {
							if data.Birthday == nil {
								return nil
							}
							return validators.ValidateBirthDate(*data.Birthday)
						}),
					router.PUT,
				),
			},
		},
	)

	return uc
}
