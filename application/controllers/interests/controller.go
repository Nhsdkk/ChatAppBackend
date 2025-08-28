package interests

import (
	interests_validators "chat_app_backend/application/controllers/validators/interests"
	user_validators "chat_app_backend/application/controllers/validators/users"
	"chat_app_backend/application/handlers/interests"
	"chat_app_backend/application/models/interests/assign"
	"chat_app_backend/application/models/interests/create"
	"chat_app_backend/application/models/interests/delete"
	"chat_app_backend/application/models/interests/get"
	"chat_app_backend/application/models/interests/update"
	"chat_app_backend/internal/extensions"
	"chat_app_backend/internal/router"
	"chat_app_backend/internal/service_wrapper"
	"chat_app_backend/internal/validator"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	router.Controller
}

func CreateInterestsController(
	r *gin.Engine,
	wrapper service_wrapper.IServiceWrapper,
) (ic Controller) {
	ic.Controller = router.CreateController(
		r,
		"/interests",
		[]router.IRoute{
			&router.AuthorizedRoute[get.GetInterestsRequestDto, get.GetInterestsResponseDto]{
				Route: router.CreateBaseRoute(
					wrapper,
					"/",
					interests.GetInterestsHandler{}.Handle,
					validator.Validator[get.GetInterestsRequestDto]{},
					router.POST,
				),
			},
			&router.AuthorizedRoute[create.CreateInterestRequestDto, create.CreateInterestResponseDto]{
				Route: router.CreateBaseRoute(
					wrapper,
					"/create",
					interests.CreateInterestHandler{}.Handle,
					validator.Validator[create.CreateInterestRequestDto]{}.
						AttachValidator(
							validator.ExternalValidator[create.CreateInterestRequestDto, string]{}.
								RuleFor(
									func(data *create.CreateInterestRequestDto) *string {
										return &data.Icon.Filename
									},
								).
								Must(interests_validators.IconFileTypeValidator{}).
								WithMessage("invalid file type").
								Validate,
						),
					router.POST,
				),
			},
			&router.AuthorizedRoute[delete.DeleteInterestRequestDto, delete.DeleteInterestResponseDto]{
				Route: router.CreateBaseRoute(
					wrapper,
					"/:id",
					interests.DeleteInterestHandler{}.Handle,
					validator.Validator[delete.DeleteInterestRequestDto]{}.
						AttachValidator(
							validator.ExternalValidator[delete.DeleteInterestRequestDto, []extensions.UUID]{}.
								RuleFor(
									func(data *delete.DeleteInterestRequestDto) *[]extensions.UUID {
										return &[]extensions.UUID{data.ID}
									},
								).
								Must(
									interests_validators.InterestsExistenceValidator{
										Db: wrapper.GetDbConnection(),
									},
								).
								WithMessage("interest with provided id does not exist").
								Validate,
						),
					router.DELETE,
				),
			},
			&router.AuthorizedRoute[update.UpdateInterestRequestDto, update.UpdateInterestResponseDto]{
				Route: router.CreateBaseRoute(
					wrapper,
					"/:id",
					interests.UpdateInterestsHandler{}.Handle,
					validator.Validator[update.UpdateInterestRequestDto]{}.
						AttachValidator(
							validator.ExternalValidator[update.UpdateInterestRequestDto, []extensions.UUID]{}.
								RuleFor(
									func(data *update.UpdateInterestRequestDto) *[]extensions.UUID {
										return &[]extensions.UUID{data.ID}
									},
								).
								Must(
									interests_validators.InterestsExistenceValidator{
										Db: wrapper.GetDbConnection(),
									},
								).
								WithMessage("interest with provided id does not exist").
								Validate,
						).
						AttachValidator(
							validator.ExternalValidator[update.UpdateInterestRequestDto, string]{}.
								RuleFor(
									func(data *update.UpdateInterestRequestDto) *string {
										return &data.Icon.Filename
									},
								).
								Must(interests_validators.IconFileTypeValidator{}).
								WithMessage("invalid file type").
								Optional().
								Validate,
						),
					router.PUT,
				),
			},
			&router.AuthorizedRoute[assign.AssignInterestRequestDto, assign.AssignInterestResponseDto]{
				Route: router.CreateBaseRoute(
					wrapper,
					"/assign",
					interests.AssignInterestHandler{}.Handle,
					validator.Validator[assign.AssignInterestRequestDto]{}.
						AttachValidator(
							validator.ExternalValidator[assign.AssignInterestRequestDto, []extensions.UUID]{}.
								RuleFor(
									func(data *assign.AssignInterestRequestDto) *[]extensions.UUID {
										return &data.InterestIds
									},
								).
								Must(
									interests_validators.InterestsExistenceValidator{
										Db: wrapper.GetDbConnection(),
									},
								).
								WithMessage("some interests dont exist").
								Validate,
						).
						AttachValidator(
							validator.ExternalValidator[assign.AssignInterestRequestDto, extensions.UUID]{}.
								RuleFor(
									func(data *assign.AssignInterestRequestDto) *extensions.UUID {
										return &data.UserID
									},
								).
								Must(
									user_validators.UserExistenceValidator{
										Db: wrapper.GetDbConnection(),
									},
								).
								WithMessage("user does not exist").
								Validate,
						),
					router.PUT,
				),
			},
		},
	)

	return ic
}
