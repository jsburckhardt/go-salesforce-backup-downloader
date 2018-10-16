package main

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/go-playground/log"

	"github.com/spf13/viper"
)

func exportResultsToCsv(consolidateResults []DownloadResult) {
	fileFolderValidated := folderValidator(viper.GetString("sf.backuppath"))
	exportResultsToCsvFileName := fileFolderValidated + "/" + viper.GetString("sf.username") + ".csv"
	file, err := os.Create(exportResultsToCsvFileName)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	csvheaders := []string{"FileName", "FileSize", "Result", "Attempt", "Duration"}
	err = writer.Write(csvheaders)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, value := range consolidateResults {
		arrayvalue := []string{value.FileName, value.FileSize, value.Result, strconv.Itoa(value.Attempt), (value.Duration).String()}
		err := writer.Write(arrayvalue)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}
