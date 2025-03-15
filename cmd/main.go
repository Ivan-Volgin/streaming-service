package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"streaming-service/internal/repo"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	"streaming-service/internal/api"
	"streaming-service/internal/config"
	customLogger "streaming-service/internal/logger"
	"streaming-service/internal/service"
)

func main() {
	if err := godotenv.Load("local.env"); err != nil {
		log.Fatal(errors.Wrap(err, "Error loading .env file"))
	}

	var cfg config.AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(errors.Wrap(err, "failed to process configuration"))
	}
	
	logger, err := customLogger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to initialize logger"))
	}

	repository, err := repo.NewRepository(context.Background(), cfg.PostgreSQL)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to initialize repository"))
	}

	serviceInstance := service.NewService(repository.MovieRepo, repository.OwnerRepo, logger)

	app := api.NewRouters(&api.Routers{
		MovieService: serviceInstance,
		OwnerService: serviceInstance,
	}, cfg.Rest.Token)

	go func() {
		logger.Infof("Starting server on %s", cfg.Rest.ListenAddress)
		if err := app.Listen(cfg.Rest.ListenAddress); err != nil {
			log.Fatal(errors.Wrap(err, "failed to start server"))
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Infof("Shutting down server...")
}
