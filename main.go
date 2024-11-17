package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Erwanph/be-wan-central-lab/internal/config"
)

func main() {
	viperConfig, err := config.NewViper()
	if err != nil {
		log.Fatalf("Failed to initialize viper config: %v", err)
	}
	log := config.NewLogger(viperConfig)
	mongo_1 := config.NewMongoDatabase(viperConfig, "MONGODB_URI_1")
	validate := config.NewValidator()
	app := config.NewFiber(viperConfig)

	config.Bootstrap(&config.BootstrapConfig{
		MongoDB1: mongo_1, // user
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprintf("%d", viperConfig.GetInt("WEB_PORT"))
	}

	err = app.Listen(fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
