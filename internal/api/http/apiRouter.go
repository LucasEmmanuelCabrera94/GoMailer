package http

import (
	"goMailer/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

type ApiRouter struct {
	senderMailController *controllers.SenderMailController
}

func NewAPIRouter(controller *controllers.SenderMailController) *ApiRouter {
	return &ApiRouter{senderMailController: controller}
}

func (api *ApiRouter) Run() {
	r := gin.Default()

	api.configRoutes(r)

	r.Run()
}

func (api *ApiRouter) configRoutes(r *gin.Engine) {
	//r.GET("/ping", func(c *gin.Context) { api.pingCtrl.Ping(c) })
	r.POST("/messages/:id", func(c *gin.Context) { api.senderMailController.SenderMail(c) })
}
