package main

import (
	"os"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/sclng-backend-test-v1/config"
	"github.com/Scalingo/sclng-backend-test-v1/internal/server"
)

func main() {
	log := logger.Default()
	log.Info("Initializing app")
	cfg, err := config.NewConfig()
	if err != nil {
		log.WithError(err).Error("Fail to initialize configuration")
		os.Exit(1)
	}

	err = server.InitRouter(cfg, log)
	if err != nil {
		os.Exit(2)
	}

}
