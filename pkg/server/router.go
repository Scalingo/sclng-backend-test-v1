package server

import (
	"fmt"
	"net/http"

	"github.com/Scalingo/go-handlers"
	"github.com/Scalingo/sclng-backend-test-v1/config"
	"github.com/sirupsen/logrus"
)

func InitRouter(cfg *config.Config, log logrus.FieldLogger) error {
	log.Info("Initializing routes")
	router := handlers.NewRouter(log)
	router.HandleFunc("/ping", PongHandler)
	router.HandleFunc("/listRepo", ListLastHandredPublicRepo)
	router.HandleFunc("/listAgragateRepo", AgregateLastHandredPublicRepo)
	// Initialize web server and configure the following routes:
	// GET /repos
	// GET /stats

	log = log.WithField("port", cfg.Port)
	log.Info("Listening...")
	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)
	if err != nil {
		log.WithError(err).Error("Fail to listen to the given port")
		return err
	}
	return nil
}
