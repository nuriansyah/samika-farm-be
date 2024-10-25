package main

//go:generate go run github.com/swaggo/swag/cmd/swag init
//go:generate go run github.com/google/wire/cmd/wire

import (
	"github.com/sanika-farm/sanika-farm-be/configs"
	"github.com/sanika-farm/sanika-farm-be/pkg/logger"
)

var configSvc *configs.Config

func main() {
	logger.InitLogger()
	configSvc = configs.Get()
	logger.SetLogLevel(configSvc)

	// InitializeApp()
	httpSvc := InitializeApp()
	httpSvc.SetupAndServe()
}
