package http

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/sanika-farm/sanika-farm-be/configs"
	"github.com/sanika-farm/sanika-farm-be/infras"
	"github.com/sanika-farm/sanika-farm-be/pkg/logger"
	"github.com/sanika-farm/sanika-farm-be/transports/http/response"
	"github.com/sanika-farm/sanika-farm-be/transports/http/router"
	"github.com/shopspring/decimal"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/swaggo/swag/example/basic/docs"
)

// ServerState is an indicator if this server's state.
type ServerState int

const (
	// ServerStateReady indicates that the server is ready to serve.
	ServerStateReady ServerState = iota + 1
	// ServerStateInGracePeriod indicates that the server is in its grace period and will shut down after it is done cleaning up.
	ServerStateInGracePeriod
	// ServerStateInCleanupPeriod indicates that the server no longer responds to any requests, is cleaning up its internal state, and will shut down shortly.
	ServerStateInCleanupPeriod
)

// HTTP is the HTTP server.
type HTTP struct {
	Config *configs.Config
	DB     *infras.PostgresConn
	Router router.Router
	State  ServerState
	mux    *gin.Engine
}

// ProvideHTTP is the provider for HTTP.
func ProvideHTTP(db *infras.PostgresConn, config *configs.Config, router router.Router) *HTTP {
	return &HTTP{
		DB:     db,
		Config: config,
		Router: router,
	}
}

// SetupAndServe sets up the server and gets it up and running.
func (h *HTTP) SetupAndServe() {
	h.mux = gin.New()
	h.setupMiddleware()
	h.setupSwaggerDocs()
	h.setupRoutes()
	h.setupGracefulShutdown()
	h.State = ServerStateReady

	h.logServerInfo()

	log.Info().Str("port", h.Config.Server.Port).Msg("Starting up HTTP server.")

	err := h.mux.Run(":" + h.Config.Server.Port)
	if err != nil {
		logger.ErrorWithStack(err)
	}
}

func (h *HTTP) setupSwaggerDocs() {
	if h.Config.Server.Env == "development" {
		docs.SwaggerInfo.Title = h.Config.App.Name
		docs.SwaggerInfo.Version = h.Config.App.Revision
		swaggerURL := fmt.Sprintf("%s/swagger/doc.json", h.Config.App.URL)
		// Use `httpSwagger.WrapHandler` directly without additional parameters
		h.mux.GET("/swagger/*any", gin.WrapH(httpSwagger.WrapHandler))
		log.Info().Str("url", swaggerURL).Msg("Swagger documentation enabled.")
	}
}

func (h *HTTP) setupRoutes() {
	decimal.MarshalJSONWithoutQuotes = true
	h.Router.SetupRoutes(h.mux)
}

func (h *HTTP) setupGracefulShutdown() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	go h.respondToSigterm(done)
}

func (h *HTTP) respondToSigterm(done chan os.Signal) {
	<-done
	defer os.Exit(0)

	shutdownConfig := h.Config.Server.Shutdown

	log.Info().Msg("Received SIGTERM.")
	log.Info().Int64("seconds", shutdownConfig.GracePeriodSeconds).Msg("Entering grace period.")
	h.State = ServerStateInGracePeriod
	time.Sleep(time.Duration(shutdownConfig.GracePeriodSeconds) * time.Second)

	log.Info().Int64("seconds", shutdownConfig.CleanupPeriodSeconds).Msg("Entering cleanup period.")
	h.State = ServerStateInCleanupPeriod
	time.Sleep(time.Duration(shutdownConfig.CleanupPeriodSeconds) * time.Second)

	log.Info().Msg("Cleaning up completed. Shutting down now.")
}

func (h *HTTP) setupMiddleware() {
	h.mux.Use(gin.Logger())
	h.mux.Use(gin.Recovery())
	h.mux.Use(h.serverStateMiddleware())
	h.setupCORS()
}

func (h *HTTP) logServerInfo() {
	h.logCORSConfigInfo()
}

func (h *HTTP) logCORSConfigInfo() {
	corsConfig := h.Config.App.CORS
	if corsConfig.Enable {
		log.Info().Msg("CORS Headers and Handlers are enabled.")
		log.Info().Str("CORS Header", fmt.Sprintf("Access-Control-Allow-Credentials: %t", corsConfig.AllowCredentials)).Msg("")
		log.Info().Str("CORS Header", strings.Join(corsConfig.AllowedHeaders, ",")).Msg("Access-Control-Allow-Headers")
		log.Info().Str("CORS Header", strings.Join(corsConfig.AllowedMethods, ",")).Msg("Access-Control-Allow-Methods")
		log.Info().Str("CORS Header", strings.Join(corsConfig.AllowedOrigins, ",")).Msg("Access-Control-Allow-Origin")
		log.Info().Str("CORS Header", fmt.Sprintf("Access-Control-Max-Age: %d", corsConfig.MaxAgeSeconds)).Msg("")
	} else {
		log.Info().Msg("CORS Headers are disabled.")
	}
}

func (h *HTTP) serverStateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch h.State {
		case ServerStateReady:
			// Server is ready to serve, proceed as normal.
			c.Next()
		case ServerStateInGracePeriod:
			// Server is in grace period. Issue a warning message and continue serving.
			log.Warn().Msg("SERVER IS IN GRACE PERIOD")
			c.Next()
		case ServerStateInCleanupPeriod:
			// Server is in cleanup period. Stop request and respond appropriately.
			response.WithPreparingShutdown(c)
			c.Abort()
		}
	}
}

func (h *HTTP) setupCORS() {
	corsConfig := h.Config.App.CORS
	if corsConfig.Enable {
		h.mux.Use(func(c *gin.Context) {
			c.Header("Access-Control-Allow-Credentials", fmt.Sprintf("%t", corsConfig.AllowCredentials))
			c.Header("Access-Control-Allow-Headers", strings.Join(corsConfig.AllowedHeaders, ","))
			c.Header("Access-Control-Allow-Methods", strings.Join(corsConfig.AllowedMethods, ","))
			c.Header("Access-Control-Allow-Origin", strings.Join(corsConfig.AllowedOrigins, ","))
			c.Header("Access-Control-Max-Age", fmt.Sprintf("%d", corsConfig.MaxAgeSeconds))
			c.Next()
		})
	}
}
