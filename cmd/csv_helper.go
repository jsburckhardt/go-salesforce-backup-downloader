package cmd

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"

	"github.com/go-playground/log"
)

func exportResultsToCsv(consolidateResults []DownloadResult) {
	fileFolderValidated := strings.Split(salesForceUserName, "@")[len(strings.Split(salesForceUserName, "@"))-1]
	exportResultsToCsvFileName := fileFolderValidated + "/" + salesForceUserName + ".csv"
	file, err := os.Create(exportResultsToCsvFileName)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	csvheaders := []string{"FileName", "FileSize", "Result", "Error", "Attempt", "Duration"}
	err = writer.Write(csvheaders)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, value := range consolidateResults {
		arrayvalue := []string{value.FileName, value.FileSize, value.Result, value.Error, strconv.Itoa(value.Attempt), (value.Duration).String()}
		err := writer.Write(arrayvalue)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}
