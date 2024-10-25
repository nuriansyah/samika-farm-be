package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sanika-farm/sanika-farm-be/pkg/failure"
)

// Base is the base object of all responses
type Base struct {
	Data     *interface{} `json:"data,omitempty"`
	Metadata *interface{} `json:"metadata,omitempty"`
	Error    *string      `json:"error,omitempty"`
	Message  *string      `json:"message,omitempty"`
}

// NoContent sends a response without any content
func NoContent(c *gin.Context) {
	respond(c, http.StatusNoContent, nil)
}

// WithMessage sends a response with a simple text message
func WithMessage(c *gin.Context, code int, message string) {
	respond(c, code, Base{Message: &message})
}

// WithJSON sends a response containing a JSON object
func WithJSON(c *gin.Context, code int, jsonPayload interface{}) {
	respond(c, code, Base{Data: &jsonPayload})
}

// WithMetadata sends a response containing a JSON object with metadata
func WithMetadata(c *gin.Context, code int, jsonPayload interface{}, metadata interface{}) {
	respond(c, code, Base{Data: &jsonPayload, Metadata: &metadata})
}

// WithError sends a response with an error message
func WithError(c *gin.Context, err error) {
	code := failure.GetCode(err)
	errMsg := err.Error()
	respond(c, code, Base{Error: &errMsg})
}

// WithPreparingShutdown sends a default response for when the server is preparing to shut down
func WithPreparingShutdown(c *gin.Context) {
	WithMessage(c, http.StatusServiceUnavailable, "SERVER PREPARING TO SHUT DOWN")
}

// WithUnhealthy sends a default response for when the server is unhealthy
func WithUnhealthy(c *gin.Context) {
	WithMessage(c, http.StatusServiceUnavailable, "SERVER UNHEALTHY")
}

func respond(c *gin.Context, code int, payload interface{}) {
	c.JSON(code, payload)
}
