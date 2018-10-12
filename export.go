package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

var wg sync.WaitGroup
var mutex = &sync.Mutex{}

func export(lr loginRes, consolidateResults *[]DownloadResult) {
	paths := getPaths(lr)

	if len(paths[0]) == 0 {
		err := errors.New("Getting files to download failed, number of files zero")
		log.Fatalln(err)
	}
	fmt.Printf("NUMBER OF URLS TO DOWNLOAD: %v\n", len(paths))
	fmt.Printf("USING %v THREADS\n", viper.GetInt("sf.workers"))

	wg.Add(viper.GetInt("sf.workers"))
	go pool(&wg, viper.GetInt("sf.workers"), paths, lr, consolidateResults)
	wg.Wait()
}

func pool(wg *sync.WaitGroup, workers int, paths []string, lr loginRes, consolidateResults *[]DownloadResult) {
	tasksCh := make(chan string)

	for i := 0; i < workers; i++ {
		fmt.Printf("Creating pool of workers. Worker: %v\n", i+1)
		go worker(tasksCh, wg, lr, consolidateResults)
	}

	for _, v := range paths {
		tasksCh <- v
	}
	close(tasksCh)
}

func worker(tasksCh <-chan string, wg *sync.WaitGroup, lr loginRes, consolidateResults *[]DownloadResult) {
	defer wg.Done()

	for {
		task, ok := <-tasksCh
		if !ok {
			return
		}

		//Download details for consolidate results
		var downloadResultTemp DownloadResult

		// getting file's expected size
		url := viper.GetString("sf.baseUrl") + task
		expectecSize := getDownloadSize(lr, url)

		// downloading file
		fn := fileName(url)
		filePath := viper.GetString("sf.backuppath") + fn + ".zip"
		log.Printf("Staring download: %s", fn)

		startDownloadTime := time.Now()
		err := DownloadFile(lr, filePath, url)
		if err != nil {
			log.Fatalln(err)
			return
		}
		endDownloadTime := time.Now()

		//verifying file size
		fi, err := os.Stat(filePath)
		if err != nil {
			log.Fatalln(err)
		}

		attempt := 0
		downloadResultTemp.Duration = endDownloadTime.Sub(startDownloadTime)
		downloadResultTemp.FileName = fn

		if expectedSizeInt, _ := strconv.ParseInt(expectecSize, 10, 64); expectedSizeInt == fi.Size() {
			completedfn := viper.GetString("sf.backuppath") + fn
			_, err := os.Create(completedfn)
			if err != nil {
				log.Fatalln(err)
			}
			attempt++
			downloadResultTemp.Attempt = attempt
			downloadResultTemp.FileSize = expectecSize
			downloadResultTemp.Result = "Successful"
			log.Printf("Successful download: %s", fn)
		} else {
			attempt++
			downloadResultTemp.Attempt = attempt
			downloadResultTemp.FileSize = string(0)
			downloadResultTemp.Result = "Fail"
			os.Remove(filePath)
			errorMessage := "Downloading file " + fn + " failed"
			err := errors.New(errorMessage)
			log.Fatalln(err)
		}
		mutex.Lock()
		*consolidateResults = append(*consolidateResults, downloadResultTemp)
		mutex.Unlock()
	}
}

func getPaths(lr loginRes) []string {
	url := viper.GetString("sf.baseUrl") + "/servlet/servlet.OrgExport"

	headers := map[string]string{
		"Cookie":         fmt.Sprintf("oid=%s;sid=%s", lr.orgID, lr.sID),
		"X-SFDC-Session": lr.sID,
	}

	status, res, err, _ := sendRequest("POST", url, headers, nil)
	if err != nil || status >= 400 {
		log.Fatalf("Error upon login: %v", err)
	}

	return strings.Split(strings.TrimSpace(string(res)), "\n")
}

func fileName(url string) string {
	var re = regexp.MustCompile(`.*fileName=(.*)\.ZIP.*`)
	result := re.FindAllStringSubmatch(url, -1)
	return result[0][1]
}

func getDownloadSize(lr loginRes, url string) string {
	headers := map[string]string{
		"Cookie":         fmt.Sprintf("oid=%s;sid=%s", lr.orgID, lr.sID),
		"X-SFDC-Session": lr.sID,
	}
	status, _, err, head := sendRequest("HEAD", url, headers, nil)
	if err != nil || status >= 400 {
		log.Fatalf("Error upon login: %v", err)
	}
	return string(head[0])
}

//DownloadFile downloads file into ...
func DownloadFile(lr loginRes, filepath string, url string) error {

	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	headers := map[string]string{
		"Cookie":         fmt.Sprintf("oid=%s;sid=%s", lr.orgID, lr.sID),
		"X-SFDC-Session": lr.sID,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	// Write body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
