package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/log"
	"github.com/go-playground/log/handlers/console"
	"github.com/spf13/viper"
)

//DownloadResult struct used for catching results from downloaded file attempt
type DownloadResult struct {
	FileName, FileSize, Result string
	Attempt                    int
	Duration                   time.Duration
}

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
	start := time.Now()
	fmt.Printf("Start time -> %s\n", start.Format(time.ANSIC))
	loginData := login()

	var consolidateResults []DownloadResult

	//Download files
	export(loginData, &consolidateResults)

	//Export results
	exportResultsToCsv(consolidateResults)

	t := time.Now()
	fmt.Printf("End time -> %s\n", t.Format(time.ANSIC))
	fmt.Printf("total time -> %s\n", t.Sub(start))
}
