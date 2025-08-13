package interests

import (
	"chat_app_backend/application/handlers/interests"
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
			&router.AuthorizedRoute[interests2.GetInterestsRequestDto, interests2.GetInterestResponseDto]{
				Route: router.CreateBaseRoute(
					wrapper,
					"/",
					interests.GetInterestsHandler{}.Handle,
					validator.Validator[interests2.GetInterestsRequestDto]{},
					router.POST,
				),
			},
		},
	)

	return ic
}
