package controllers

import (
	"fmt"
	"goMailer/internal/domain"
	"goMailer/internal/services/sender_mail"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			errorMessages := translateValidationErrors(validationErrs)
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

func translateValidationErrors(err validator.ValidationErrors) []string {
	var errorMessages []string

	for _, err := range err {
		var message string

		switch err.Tag() {
		case "required":
			message = fmt.Sprintf("The field '%s' is required.", err.Field())
		case "email":
			message = fmt.Sprintf("The field '%s' must be a valid email address.", err.Field())
		case "nefield":
			message = fmt.Sprintf("The field '%s' must not be the same as field '%s'.", err.Field(), err.Param())
		default:
			message = fmt.Sprintf("Validation error on field '%s' with tag '%s'.", err.Field(), err.Tag())
		}

		errorMessages = append(errorMessages, message)
	}

	return errorMessages
}
