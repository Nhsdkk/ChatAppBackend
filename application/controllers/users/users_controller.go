package controllers

import (
	"chat_app_backend/application/controllers/users/validators"
	"chat_app_backend/application/handlers/users"
	"chat_app_backend/application/models/users/get_user_data"
	"chat_app_backend/application/models/users/login"
	"chat_app_backend/application/models/users/refresh_token"
	"chat_app_backend/application/models/users/register"
	"chat_app_backend/internal/router"
	"chat_app_backend/internal/service_wrapper"
	"chat_app_backend/internal/validator"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	controller router.Controller
}

func (uc UserController) ConfigureGroup() {
	uc.controller.ConfigureGroup()
}

func CreateUserController(engine *gin.Engine, serviceWrapper service_wrapper.IServiceWrapper) (uc UserController) {
	uc.controller = router.CreateController(
		engine,
		"/users",
		[]router.IRoute{
			router.RouteFactory(
				router.Base,
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
			router.RouteFactory(
				router.Base,
				serviceWrapper,
				"/login",
				users.LoginHandler{}.Handle,
				validator.
					Validator[login.LoginRequestDto]{},
				router.POST,
			),
			router.RouteFactory(
				router.Authorized,
				serviceWrapper,
				"/",
				users.GetUserDataHandler{}.Handle,
				validator.
					Validator[get_user_data.GetUserDataRequestDto]{},
				router.GET,
			),
			router.RouteFactory(
				router.Base,
				serviceWrapper,
				"/refresh_token",
				users.RefreshTokenHandler{}.Handle,
				validator.
					Validator[refresh_token.RefreshTokenRequestDto]{},
				router.POST,
			),
		},
	)

	return uc
}
