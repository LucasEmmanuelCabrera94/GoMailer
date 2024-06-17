package main

import (
	"goMailer/internal/api/controllers"
	"goMailer/internal/api/http"
	repository "goMailer/internal/repositories/sender_mail"
	service "goMailer/internal/services/sender_mail"
)

func main() {
	senderMailRepository := repository.NewSenderMailRepository()
	senderMailService := service.NewSenderMailService(senderMailRepository)
	senderMailController := controllers.NewSenderMailController(senderMailService)
	apiRouter := http.NewAPIRouter(senderMailController)

	apiRouter.Run()
}
