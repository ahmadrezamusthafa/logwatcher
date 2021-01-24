package main

import (
	"fmt"
	"github.com/ahmadrezamusthafa/logwatcher/common/di/container"
	"github.com/ahmadrezamusthafa/logwatcher/common/logger"
	"github.com/ahmadrezamusthafa/logwatcher/config"
	"github.com/ahmadrezamusthafa/logwatcher/domain/service/thirdparty"
	"github.com/ahmadrezamusthafa/logwatcher/pkg/breaker"
	"github.com/ahmadrezamusthafa/logwatcher/pkg/database"
	"github.com/ahmadrezamusthafa/logwatcher/pkg/worker"
	"github.com/ahmadrezamusthafa/logwatcher/server"
	httphealth "github.com/ahmadrezamusthafa/logwatcher/server/http/health"
	httpthirdparty "github.com/ahmadrezamusthafa/logwatcher/server/http/thirdparty"
)

func main() {
	logger.SetupLoggerAuto("", "")
	fmt.Println(`
   __   ___ ____           __   _____ ___    _    _   __  ____ ___ _   _ ____ ____ 
  /__\ / __|_  _)   ___   (  ) (  _  ) __)  ( \/\/ ) /__\(_  _) __| )_( | ___|  _ \
 /(__)\\__ \_)(_   (___)   )(__ )(_)( (_-.   )    ( /(__)\ )(( (__ ) _ ( )__) )   /
(__)(__|___(____)         (____|_____)___/  (__/\__|__)(__|__)\___|_) (_|____|_)\_)

`)

	conf, err := config.New()
	if err != nil {
		logger.Warn("%v", err)
	}

	logger.Info("Starting service registry...")
	registry := container.NewContainer()
	registry.RegisterService("config", *conf)
	registry.RegisterService("worker", new(worker.EngineWorker))
	registry.RegisterService("database", new(database.EngineDatabase))
	registry.RegisterService("breaker", new(breaker.EngineBreaker))
	registry.RegisterService("thirdPartyService", new(thirdparty.Service))
	registry.RegisterService("healthHandler", new(httphealth.Handler))
	registry.RegisterService("thirdPartyHandler", new(httpthirdparty.Handler))

	rootHandler := new(server.RootHandler)
	registry.RegisterService("rootHandler", rootHandler)

	if err := registry.Ready(); err != nil {
		logger.Fatal("Failed to populate services %v", err)
	}

	h := server.HttpServer{Config: *conf, RootHandler: rootHandler}
	h.Serve()
}
