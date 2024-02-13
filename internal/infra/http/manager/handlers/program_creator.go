package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jcleira/encinitas-collector-go/internal/app/manager/aggregates"
)

// programCreator defines the methods needed to publish programs.
type programCreator interface {
	Create(context.Context, aggregates.Program) error
}

// ProgramsCreatorHandler defines the dependencies to create programs.
type ProgramsCreatorHandler struct {
	programCreator programCreator
}

// NewProgramsCreatorHandler initializes a new ProgramsCreatorHandler.
func NewProgramsCreatorHandler(
	programCreator programCreator) *ProgramsCreatorHandler {
	return &ProgramsCreatorHandler{
		programCreator: programCreator,
	}
}

// Handle is the handler function to create programs
func (ech *ProgramsCreatorHandler) Handle(c *gin.Context) {
	var httpProgramCreateRequest httpProgramCreateRequest
	if err := c.ShouldBindJSON(&httpProgramCreateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ech.programCreator.Create(
		c.Request.Context(), httpProgramCreateRequest.ToAggregate()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}
