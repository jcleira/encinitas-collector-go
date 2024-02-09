package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jcleira/encinitas-collector-go/internal/app/agent/aggregates"
)

// eventsPublisher defines the methods needed to publish events.
type eventsPublisher interface {
	Publish(context.Context, aggregates.Event) error
}

// EventsCreatorHandler defines the dependencies to create events.
type EventsCreatorHandler struct {
	eventsPublisher eventsPublisher
}

// NewEventsCreatorHandler initializes a new EventsCreatorHandler.
func NewEventsCreatorHandler(
	eventsPublisher eventsPublisher) *EventsCreatorHandler {
	return &EventsCreatorHandler{
		eventsPublisher: eventsPublisher,
	}
}

// Handle is the handler function to create events
func (ech *EventsCreatorHandler) Handle(c *gin.Context) {
	var httpEventRequest httpEventRequest
	if err := c.ShouldBindJSON(&httpEventRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ech.eventsPublisher.Publish(
		c.Request.Context(), httpEventRequest.ToAggregate()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}
