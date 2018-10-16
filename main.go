package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/log"
	"github.com/go-playground/log/handlers/console"
	"github.com/spf13/viper"
)

//DownloadResult struct used for catching results from downloaded file attempt
type DownloadResult struct {
	FileName, FileSize, Result, Error string
	Attempt                           int
	Duration                          time.Duration
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
		Timeout: time.Minute * 1,
	}
}

func main() {
	start := time.Now()
	fmt.Printf("Start time -> %s\n", start.Format(time.ANSIC))
	loginData := login()

	var consolidateResults []DownloadResult

	//Verifying target folder
	fileFolderValidated := folderValidator(viper.GetString("sf.backuppath"))
	if _, err := os.Stat(fileFolderValidated); os.IsNotExist(err) {
		log.Infof("backuppath folder doesn't exist. Creating: %s", fileFolderValidated)
		err := os.Mkdir(fileFolderValidated, 0777)
		if err != nil {
			log.Fatal(err)
		}
	}

	export(loginData, &consolidateResults)

	//Export results
	exportResultsToCsv(consolidateResults)

	t := time.Now()
	fmt.Printf("End time -> %s\n", t.Format(time.ANSIC))
	fmt.Printf("total time -> %s\n", t.Sub(start))
}
