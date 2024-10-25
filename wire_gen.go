// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"github.com/sanika-farm/sanika-farm-be/configs"
	"github.com/sanika-farm/sanika-farm-be/infras"
	"github.com/sanika-farm/sanika-farm-be/internal/domain/users/repository"
	"github.com/sanika-farm/sanika-farm-be/internal/domain/users/services"
	"github.com/sanika-farm/sanika-farm-be/internal/handlers"
	"github.com/sanika-farm/sanika-farm-be/transports/http"
	"github.com/sanika-farm/sanika-farm-be/transports/http/router"
)

// Injectors from wire.go:

// Wiring for everything.
func InitializeApp() *http.HTTP {
	config := configs.Get()
	postgresConn := infras.ProvidePostgresConn(config)
	usersRepositoryImpl := repository.ProvideUsersRepository(postgresConn)
	usersServiceImpl := services.ProvideUsersService(usersRepositoryImpl, config)
	usersHandler := handlers.ProvideUsersHandler(usersServiceImpl)
	domainHandlers := router.DomainHandlers{
		UsersHandler: usersHandler,
	}
	routerRouter := router.ProvideRouter(domainHandlers)
	httpHTTP := http.ProvideHTTP(postgresConn, config, routerRouter)
	return httpHTTP
}

// wire.go:

// Wiring for configurations.
var configurationsService = wire.NewSet(configs.Get)

// Wiring for persistences.
var persistencesService = wire.NewSet(infras.ProvidePostgresConn)

// Wiring for domain users
var domainUsersService = wire.NewSet(services.ProvideUsersService, wire.Bind(new(services.UsersService), new(*services.UsersServiceImpl)), repository.ProvideUsersRepository, wire.Bind(new(repository.UsersRepository), new(*repository.UsersRepositoryImpl)))

// Wiring for all domains
var domainsServices = wire.NewSet(
	domainUsersService,
)

// Wiring for HTTP routing
var httpRouting = wire.NewSet(wire.Struct(new(router.DomainHandlers), "*"), handlers.ProvideUsersHandler, router.ProvideRouter)
