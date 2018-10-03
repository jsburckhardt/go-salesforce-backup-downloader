package main

import (
	"net/http"
	"time"

	"github.com/go-playground/log"
	"github.com/go-playground/log/handlers/console"
	"github.com/spf13/viper"
)

var httpClient *http.Client

func init() {
	// logger
	cLog := console.New(true)
	log.AddHandler(cLog, log.AllLevels...)

	// config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error parsing config file: %v", err)
	} else {
		log.Infof("Using configuration file %s", viper.ConfigFileUsed())
	}

	// initialise client
	httpClient = &http.Client{
		Timeout: time.Minute * 10,
	}
}

func main() {
	loginData := login()

	export(loginData)
}
