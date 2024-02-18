package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jcleira/encinitas-collector-go/internal/app/manager/aggregates"
)

// emailCreator defines the methods needed to publish emails.
type emailCreator interface {
	Create(context.Context, string) error
}

// EmailsCreatorHandler defines the dependencies to create emails.
type EmailsCreatorHandler struct {
	emailCreator emailCreator
}

// NewEmailsCreatorHandler initializes a new EmailsCreatorHandler.
func NewEmailsCreatorHandler(
	emailCreator emailCreator) *EmailsCreatorHandler {
	return &EmailsCreatorHandler{
		emailCreator: emailCreator,
	}
}

type httpEmailCreateRequest struct {
	Email string `json:"email"`
}

// Handle is the handler function to create emails
func (ech *EmailsCreatorHandler) Handle(c *gin.Context) {
	var httpEmailCreateRequest httpEmailCreateRequest
	if err := c.ShouldBindJSON(&httpEmailCreateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ech.emailCreator.Create(
		c.Request.Context(), httpEmailCreateRequest.Email); err != nil {
		switch err {
		case aggregates.ErrEmailAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{})
}
