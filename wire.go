//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/sanika-farm/sanika-farm-be/configs"
	"github.com/sanika-farm/sanika-farm-be/infras"
	usersRepository "github.com/sanika-farm/sanika-farm-be/internal/domain/users/repository"
	usersService "github.com/sanika-farm/sanika-farm-be/internal/domain/users/services"
	usersHandlers "github.com/sanika-farm/sanika-farm-be/internal/handlers"
	"github.com/sanika-farm/sanika-farm-be/transports/http"
	"github.com/sanika-farm/sanika-farm-be/transports/http/router"
)

// Wiring for configurations.
var configurationsService = wire.NewSet(
	configs.Get,
)

// Wiring for persistences.
var persistencesService = wire.NewSet(
	infras.ProvidePostgresConn,
)

// Wiring for domain users
var domainUsersService = wire.NewSet(
	usersService.ProvideUsersService,
	wire.Bind(new(usersService.UsersService), new(*usersService.UsersServiceImpl)),

	usersRepository.ProvideUsersRepository,
	wire.Bind(new(usersRepository.UsersRepository), new(*usersRepository.UsersRepositoryImpl)),
)

// Wiring for all domains
var domainsServices = wire.NewSet(
	domainUsersService,
)

// Wiring for HTTP routing
var httpRouting = wire.NewSet(
	wire.Struct(new(router.DomainHandlers), "*"),
	usersHandlers.ProvideUsersHandler,
	router.ProvideRouter,
)

// Wiring for everything.
func InitializeApp() *http.HTTP {
	wire.Build(
		configurationsService,
		persistencesService,
		domainsServices,
		httpRouting,
		http.ProvideHTTP,
	)
	return &http.HTTP{}
}
