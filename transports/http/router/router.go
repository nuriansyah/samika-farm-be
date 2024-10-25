package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sanika-farm/sanika-farm-be/internal/handlers"
)

// DomainHandlers is a struct that contains all domain-specific handlers.
type DomainHandlers struct {
	UsersHandler handlers.UsersHandler
}

// Router is the router struct containing handlers.
type Router struct {
	DomainHandlers DomainHandlers
}

// ProvideRouter is the provider function for this router.
func ProvideRouter(domainHandlers DomainHandlers) Router {
	return Router{
		DomainHandlers: domainHandlers,
	}
}

// SetupRoutes sets up all routing for this server.
func (r *Router) SetupRoutes(router *gin.Engine) {
	v1 := router.Group("/v1")
	{
		r.DomainHandlers.UsersHandler.Router(v1)
	}
}
