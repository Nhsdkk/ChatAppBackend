package users

import (
	interests_validators "chat_app_backend/application/controllers/validators/interests"
	"chat_app_backend/application/controllers/validators/users"
	"chat_app_backend/application/handlers/users"
	"chat_app_backend/application/models/users/delete"
	"chat_app_backend/application/models/users/get_user_data"
	"chat_app_backend/application/models/users/login"
	"chat_app_backend/application/models/users/refresh_token"
	"chat_app_backend/application/models/users/register"
	"chat_app_backend/application/models/users/update"
	"chat_app_backend/internal/exceptions"
	"chat_app_backend/internal/exceptions/common_exceptions"
	"chat_app_backend/internal/extensions"
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
						validator.CreateValidatorGroup[register.RegisterRequestDto]().
							AttachValidation(
								validator.ExternalValidator[register.RegisterRequestDto, string]{}.
									RuleFor(
										func(data *register.RegisterRequestDto) *string {
											return &data.Avatar.Filename
										},
									).
									Must(user_validators.AvatarFileTypeValidator{}).
									WithMessage("avatar file type is invalid").
									Validate,
							).
							AttachValidation(
								validator.ExternalValidator[register.RegisterRequestDto, string]{}.
									RuleFor(
										func(data *register.RegisterRequestDto) *string {
											return &data.Email
										},
									).
									Must(user_validators.EmailFormatValidator{}).
									WithMessage("email is of wrong format").
									Validate,
							).
							AttachValidation(
								validator.ExternalValidator[register.RegisterRequestDto, time.Time]{}.
									RuleFor(
										func(data *register.RegisterRequestDto) *time.Time {
											return &data.Birthday
										},
									).
									Must(user_validators.BirthDateValidator{}).
									WithMessage("you are not old enough to register").
									Validate,
							).
							AttachValidation(
								validator.ExternalValidator[register.RegisterRequestDto, string]{}.
									RuleFor(
										func(data *register.RegisterRequestDto) *string {
											return &data.Password
										},
									).
									Must(user_validators.PasswordValidator{}).
									WithMessage("password should have at least one of each of this characters (special characters, upper and lowercase letters, digits)").
									Validate,
							).Validate,
					).
					AttachValidator(
						validator.ExternalValidator[register.RegisterRequestDto, string]{}.
							RuleFor(
								func(data *register.RegisterRequestDto) *string {
									return &data.FullName
								},
							).
							Must(
								user_validators.NameUniquenessValidator{
									Db: serviceWrapper.GetDbConnection(),
								},
							).
							WithExceptionFactory(
								func(message string) error {
									return &common_exceptions.InvalidBodyException{
										BaseRestException: exceptions.BaseRestException{
											ITrackableException: exceptions.CreateTrackableExceptionFromStringF(message),
											Message:             message,
										},
									}
								},
							).
							WithMessage("that full name is already taken").
							Validate,
					).
					AttachValidator(
						validator.ExternalValidator[register.RegisterRequestDto, string]{}.
							RuleFor(
								func(data *register.RegisterRequestDto) *string {
									return &data.Email
								},
							).
							Must(
								user_validators.EmailUniquenessValidator{
									Db: serviceWrapper.GetDbConnection(),
								},
							).
							WithExceptionFactory(
								func(message string) error {
									return &common_exceptions.InvalidBodyException{
										BaseRestException: exceptions.BaseRestException{
											ITrackableException: exceptions.CreateTrackableExceptionFromStringF(message),
											Message:             message,
										},
									}
								},
							).
							WithMessage("that email is already used").
							Validate,
					).
					AttachValidator(
						validator.ExternalValidator[register.RegisterRequestDto, []extensions.UUID]{}.
							RuleFor(
								func(data *register.RegisterRequestDto) *[]extensions.UUID {
									return &data.Interests
								},
							).
							Must(
								interests_validators.InterestsExistenceValidator{
									Db: serviceWrapper.GetDbConnection(),
								},
							).
							WithExceptionFactory(
								func(message string) error {
									return &common_exceptions.ResourceNotFoundException{
										BaseRestException: exceptions.BaseRestException{
											ITrackableException: exceptions.CreateTrackableExceptionFromStringF(message),
											Message:             message,
										},
									}
								},
							).
							WithMessage("some interests not found").
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
						Validator[get_user_data.GetUserDataRequestDto]{}.
						AttachValidator(
							validator.ExternalValidator[get_user_data.GetUserDataRequestDto, extensions.UUID]{}.
								RuleFor(
									func(data *get_user_data.GetUserDataRequestDto) *extensions.UUID {
										return &data.ID
									},
								).
								Must(
									user_validators.UserExistenceValidator{
										Db: serviceWrapper.GetDbConnection(),
									},
								).
								WithMessage("user with this id does not exist").
								Validate,
						),
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
			&router.AuthorizedRoute[delete.DeleteUserRequestDto, delete.DeleteUserResponseDto]{
				Route: router.CreateBaseRoute(
					serviceWrapper,
					"/:id",
					users.DeleteUserHandler{}.Handle,
					validator.
						Validator[delete.DeleteUserRequestDto]{}.
						AttachValidator(
							validator.ExternalValidator[delete.DeleteUserRequestDto, extensions.UUID]{}.
								RuleFor(
									func(data *delete.DeleteUserRequestDto) *extensions.UUID {
										return &data.ID
									},
								).
								Must(user_validators.UserModificationAccessValidator{}).
								WithExceptionFactory(
									func(message string) error {
										return &common_exceptions.ForbiddenException{
											BaseRestException: exceptions.BaseRestException{
												ITrackableException: exceptions.CreateTrackableExceptionFromStringF(message),
												Message:             message,
											},
										}
									},
								).
								WithMessage("you dont have access to delete this user").
								Validate,
						).
						AttachValidator(
							validator.ExternalValidator[delete.DeleteUserRequestDto, extensions.UUID]{}.
								RuleFor(
									func(data *delete.DeleteUserRequestDto) *extensions.UUID {
										return &data.ID
									},
								).
								Must(
									user_validators.UserExistenceValidator{
										Db: serviceWrapper.GetDbConnection(),
									},
								).
								WithMessage("user with this id does not exist").
								Validate,
						),
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
							validator.ExternalValidator[update.UpdateUserRequestDto, extensions.UUID]{}.
								RuleFor(
									func(data *update.UpdateUserRequestDto) *extensions.UUID {
										return &data.ID
									},
								).
								Must(user_validators.UserModificationAccessValidator{}).
								WithExceptionFactory(
									func(message string) error {
										return &common_exceptions.ForbiddenException{
											BaseRestException: exceptions.BaseRestException{
												ITrackableException: exceptions.CreateTrackableExceptionFromStringF(message),
												Message:             message,
											},
										}
									},
								).
								WithMessage("you dont have access to modify this user").
								Validate,
						).
						AttachValidator(
							validator.CreateValidatorGroup[update.UpdateUserRequestDto]().
								AttachValidation(
									validator.ExternalValidator[update.UpdateUserRequestDto, string]{}.
										RuleFor(
											func(data *update.UpdateUserRequestDto) *string {
												return &data.Avatar.Filename
											},
										).
										Must(user_validators.AvatarFileTypeValidator{}).
										WithMessage("avatar file type is invalid").
										Optional().
										Validate,
								).
								AttachValidation(
									validator.ExternalValidator[update.UpdateUserRequestDto, string]{}.
										RuleFor(
											func(data *update.UpdateUserRequestDto) *string {
												return data.Email
											},
										).
										Must(user_validators.EmailFormatValidator{}).
										WithMessage("email is of wrong format").
										Optional().
										Validate,
								).
								AttachValidation(
									validator.ExternalValidator[update.UpdateUserRequestDto, time.Time]{}.
										RuleFor(
											func(data *update.UpdateUserRequestDto) *time.Time {
												return data.Birthday
											},
										).
										Must(user_validators.BirthDateValidator{}).
										WithMessage("you are not old enough to register").
										Optional().
										Validate,
								).
								AttachValidation(
									validator.ExternalValidator[update.UpdateUserRequestDto, string]{}.
										RuleFor(
											func(data *update.UpdateUserRequestDto) *string {
												return data.PasswordString
											},
										).
										Must(user_validators.PasswordValidator{}).
										WithMessage("password should have at least one of each of this characters (special characters, upper and lowercase letters, digits)").
										Optional().
										Validate,
								).
								Validate,
						).
						AttachValidator(
							validator.ExternalValidator[update.UpdateUserRequestDto, string]{}.
								RuleFor(
									func(data *update.UpdateUserRequestDto) *string {
										return data.FullName
									},
								).
								Must(
									user_validators.NameUniquenessValidator{
										Db: serviceWrapper.GetDbConnection(),
									},
								).
								Optional().
								WithMessage("that full name is already taken").
								Validate,
						).
						AttachValidator(
							validator.ExternalValidator[update.UpdateUserRequestDto, string]{}.
								RuleFor(
									func(data *update.UpdateUserRequestDto) *string {
										return data.Email
									},
								).
								Must(
									user_validators.EmailUniquenessValidator{
										Db: serviceWrapper.GetDbConnection(),
									},
								).
								Optional().
								WithMessage("that email is already used").
								Validate,
						).
						AttachValidator(
							validator.ExternalValidator[update.UpdateUserRequestDto, extensions.UUID]{}.
								RuleFor(
									func(data *update.UpdateUserRequestDto) *extensions.UUID {
										return &data.ID
									},
								).
								Must(
									user_validators.UserExistenceValidator{
										Db: serviceWrapper.GetDbConnection(),
									},
								).
								WithExceptionFactory(
									func(message string) error {
										return &common_exceptions.ResourceNotFoundException{
											BaseRestException: exceptions.BaseRestException{
												ITrackableException: exceptions.CreateTrackableExceptionFromStringF(message),
												Message:             message,
											},
										}
									},
								).
								WithMessage("user with this id does not exist").
								Validate,
						),
					router.PUT,
				),
			},
		},
	)

	return uc
}
