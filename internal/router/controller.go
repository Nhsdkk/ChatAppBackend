package router

import (
	"github.com/gin-gonic/gin"
)

type IController interface {
	ConfigureGroup()
}

type Controller struct {
	router         *gin.Engine
	routes         []IRoute
	controllerPath string
}

func (controller Controller) ConfigureGroup() {
	group := controller.router.Group(controller.controllerPath)
	for _, route := range controller.routes {
		RegisterRoute(group, route)
	}
}

func CreateController(router *gin.Engine, controllerPath string, routes []IRoute) (controller Controller) {
	controller.router = router
	controller.controllerPath = controllerPath
	controller.routes = routes
	return controller
}
