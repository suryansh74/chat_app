package main

import (
	"github.com/suryansh74/chat_app/config"
	"github.com/suryansh74/chat_app/pkg/logger"
	"github.com/suryansh74/chat_app/server"
)

func main() {
	// logger setup
	logger.Init()
	defer logger.Sync()
	logger.Log.Info("Starting the Go Chat Backend...")

	// load config
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// start server
	srv := server.NewServer(&cfg)
	if err := srv.Run(); err != nil {
		panic(err)
	}
	srv.Run()
}
