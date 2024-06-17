package controllers

import (
	"goMailer/internal/domain"
	"goMailer/internal/services/sender_mail"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SenderMailController struct {
	SenderMailService *sender_mail.SenderMailService
}

func NewSenderMailController(service *sender_mail.SenderMailService) *SenderMailController {
	return &SenderMailController{SenderMailService: service}
}

func (ctrl *SenderMailController) SenderMail(c *gin.Context) {
	var params domain.SenderMailParams

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error to bind json": err.Error()})
		return
	}

	if err := ctrl.SenderMailService.SendEmail(c, params); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error to send email": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})

}
