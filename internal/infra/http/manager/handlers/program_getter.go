package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jcleira/encinitas-collector-go/internal/app/manager/aggregates"
)

// programGetter defines the methods needed to publish programs.
type programGetter interface {
	GetPrograms(context.Context) ([]aggregates.Program, error)
}

// ProgramGetterHandler defines the dependencies to create programs.
type ProgramGetterHandler struct {
	programGetter programGetter
}

// NewProgramGetterHandler initializes a new ProgramGetterHandler.
func NewProgramGetterHandler(
	programGetter programGetter) *ProgramGetterHandler {
	return &ProgramGetterHandler{
		programGetter: programGetter,
	}
}

// Handle is the handler function to create programs
func (ech *ProgramGetterHandler) Handle(c *gin.Context) {
	programs, err := ech.programGetter.GetPrograms(
		c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, httpProgramsFromAggregates(programs))
}
