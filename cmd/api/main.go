package main

import (
	"goMailer/internal/api/controllers"
	"goMailer/internal/api/http"
	repository "goMailer/internal/repositories/sender_mail"
	service "goMailer/internal/services/sender_mail"
)

func main() {
	senderMailController := initSenderMail()

	apiRouter := http.NewAPIRouter(senderMailController)
	apiRouter.Run()
}

func initSenderMail() *controllers.SenderMailController {
	senderMailRepository := repository.NewSenderMailRepository()
	senderMailService := service.NewSenderMailService(senderMailRepository)
	return controllers.NewSenderMailController(senderMailService)
}
