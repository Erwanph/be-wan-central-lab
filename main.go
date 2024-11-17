package main

import (
	"net/http"

	"github.com/Erwanph/be-wan-central-lab/internal/config"
	"github.com/gofiber/adaptor/v2"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	viperConfig, err := config.NewViper()
	if err != nil {
		http.Error(w, "Failed to initialize config", http.StatusInternalServerError)
		return
	}

	log := config.NewLogger(viperConfig)
	mongo_1 := config.NewMongoDatabase(viperConfig, "MONGODB_URI_1")
	validate := config.NewValidator()
	app := config.NewFiber(viperConfig)

	config.Bootstrap(&config.BootstrapConfig{
		MongoDB1: mongo_1,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
	})

	adaptor.FiberApp(app)(w, r)
}
