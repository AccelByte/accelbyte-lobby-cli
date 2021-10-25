package main

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type config struct {
	LobbyBaseURL    string `env:"LOBBY_BASE_URL" envDocs:"The base URL for lobby server"`
	IAMBaseURL      string `env:"IAM_BASE_URL" envDocs:"The base URL where the IAM located"`
	QOSBaseURL      string `env:"QOS_BASE_URL" envDocs:"The base URL where the QOS located"`
	IAMClientID     string `env:"IAM_CLIENT_ID" envDocs:"The client ID for the Lobby Server"`
	IAMClientSecret string `env:"IAM_CLIENT_SECRET" envDocs:"The client secret"`
}

func getConfig() *config {
	err := godotenv.Load()
	if err != nil {
		logrus.Info("no .env file found")
	}

	config := &config{}
	err = env.Parse(config)
	if err != nil {
		log.Fatal("failed to load config")
	}
	return config
}
