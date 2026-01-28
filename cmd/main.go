package main

import (
	"net"
	"net/http"

	"github.com/dhiyazhar/voucher-redemption-api/internal/config"
	"github.com/dhiyazhar/voucher-redemption-api/internal/delivery/handler/middleware"
	"github.com/go-playground/validator/v10"
)

func main() {
	viperConfig := config.NewViper()
	logger := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, logger)
	validate := validator.New()
	mux := http.NewServeMux()

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		Config:   viperConfig,
		Logger:   logger,
		Mux:      mux,
		Validate: validate,
	})

	host := viperConfig.GetString("WEB_HOST")
	port := viperConfig.GetString("WEB_PORT")
	server := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: middleware.Logging(logger, mux),
	}

	logger.Info("starting server", "addr", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server stopped", "err", err)
		panic(err)
	}

}
