package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	serviceConfig "github.com/dibrito/ennismore-weather-app/config"
	"github.com/dibrito/ennismore-weather-app/internal/clients/openstreetmap"
	"github.com/dibrito/ennismore-weather-app/internal/clients/weather"
	"github.com/dibrito/ennismore-weather-app/internal/controller"
	httpHandler "github.com/dibrito/ennismore-weather-app/internal/handler"
	repository "github.com/dibrito/ennismore-weather-app/internal/repository"

	"github.com/dibrito/ennismore-weather-app/pkg/logging"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const serviceName = "ennismore-weather-app"

func main() {
	// define logger
	logger, _ := zap.NewProduction()
	defer logger.Sugar()
	logger = logger.With(zap.String(logging.FieldService, serviceName))

	// read config
	f, err := os.Open("config.yaml")
	if err != nil {
		logger.Fatal("unable to open configuration", zap.Error(err))
	}
	defer f.Close()

	var cfg serviceConfig.ServiceConfig
	if err = yaml.NewDecoder(f).Decode(&cfg); err != nil {
		logger.Fatal("unable to parse configuration", zap.Error(err))
	}

	// setup cache
	cache := repository.New()

	// setup open streat map client
	openstreetmapClient := openstreetmap.New(cfg.OpenstreetmapConfig)

	// setup forecast client
	weatherClient := weather.New(cfg.WeatherConfig)

	// set up controller
	controller := controller.New(openstreetmapClient, weatherClient, cache)
	// set up handler
	handler := httpHandler.New(controller)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.APIConfig.Port),
		Handler: handler.Routes(logger),
	}

	// start server
	go func() {
		logger.Info("started weather-app service",
			zap.String("servicename", serviceName),
			zap.Int("port", cfg.APIConfig.Port))
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("listen and service", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("received shutdown signal ...")

	// Close all e.g. database connection etc....

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.APIConfig.ShutdownTimeout)*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("server shutdown fail", zap.Error(err))
	} else {
		logger.Info("server exiting OK!")
	}

	<-shutdownCtx.Done()
	logger.Info("shutdown timeout is DONE!", zap.Int("timeout", cfg.APIConfig.ShutdownTimeout))

	logger.Info("done!")
}
