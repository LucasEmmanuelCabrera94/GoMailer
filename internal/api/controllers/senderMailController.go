package controllers

import (
	"goMailer/internal/domain"
	"goMailer/internal/services/sender_mail"
	"goMailer/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type SenderMailController struct {
	SenderMailService sender_mail.Service
}

func NewSenderMailController(service sender_mail.Service) *SenderMailController {
	return &SenderMailController{SenderMailService: service}
}

func (ctrl *SenderMailController) SenderMail(c *gin.Context) {
	var params []domain.SenderMailParams

	if err := c.ShouldBindJSON(&params); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			errorMessages := utils.TranslateValidationErrors(validationErrs)
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request payload",
				"details": strings.Join(errorMessages, ", "),
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		}
		return
	}

	if err := ctrl.SenderMailService.SendEmail(c, params); err != nil {
		switch err := err.(type) {
		case *domain.ValidationError:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case *domain.ExternalServiceError:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}
