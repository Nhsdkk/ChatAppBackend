package interests

import (
	"chat_app_backend/application/handlers/interests"
	"chat_app_backend/application/models/interests/create"
	delete2 "chat_app_backend/application/models/interests/delete"
	interests2 "chat_app_backend/application/models/interests/get"
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
			&router.AuthorizedRoute[interests2.GetInterestsRequestDto, interests2.GetInterestsResponseDto]{
				Route: router.CreateBaseRoute(
					wrapper,
					"/",
					interests.GetInterestsHandler{}.Handle,
					validator.Validator[interests2.GetInterestsRequestDto]{},
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
			&router.AuthorizedRoute[delete2.DeleteInterestRequestDto, delete2.DeleteInterestResponseDto]{
				Route: router.CreateBaseRoute(
					wrapper,
					"/:id",
					interests.DeleteInterestHandler{}.Handle,
					validator.Validator[delete2.DeleteInterestRequestDto]{},
					router.DELETE,
				),
			},
		},
	)

	return ic
}
