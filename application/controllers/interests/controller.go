package interests

import (
	"chat_app_backend/application/handlers/interests"
	"chat_app_backend/application/models/interests/assign"
	"chat_app_backend/application/models/interests/create"
	"chat_app_backend/application/models/interests/delete"
	"chat_app_backend/application/models/interests/get"
	"chat_app_backend/application/models/interests/update"
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
					validator.Validator[create.CreateInterestRequestDto]{},
					router.POST,
				),
			},
			&router.AuthorizedRoute[delete.DeleteInterestRequestDto, delete.DeleteInterestResponseDto]{
				Route: router.CreateBaseRoute(
					wrapper,
					"/:id",
					interests.DeleteInterestHandler{}.Handle,
					validator.Validator[delete.DeleteInterestRequestDto]{},
					router.DELETE,
				),
			},
			&router.AuthorizedRoute[update.UpdateInterestRequestDto, update.UpdateInterestResponseDto]{
				Route: router.CreateBaseRoute(
					wrapper,
					"/:id",
					interests.UpdateInterestsHandler{}.Handle,
					validator.Validator[update.UpdateInterestRequestDto]{},
					router.PUT,
				),
			},
			&router.AuthorizedRoute[assign.AssignInterestRequestDto, assign.AssignInterestResponseDto]{
				Route: router.CreateBaseRoute(
					wrapper,
					"/assign",
					interests.AssignInterestHandler{}.Handle,
					validator.Validator[assign.AssignInterestRequestDto]{},
					router.PUT,
				),
			},
		},
	)

	return ic
}
