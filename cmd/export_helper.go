package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/log"
)

var wg sync.WaitGroup
var mutex = &sync.Mutex{}

func export(lr loginRes, consolidateResults *[]DownloadResult) {
	paths := getPaths(lr)

	if len(paths[0]) == 0 {
		err := errors.New("Getting list of files to download failed. Number of files can't be zero")
		log.Fatalf("Getting paths failed: %s", err)
	}
	fmt.Printf("NUMBER OF URLS TO DOWNLOAD: %v\n", len(paths))
	workers := 0
	if len(paths) < maxWorkers {
		workers = len(paths)
	} else {
		workers = maxWorkers
	}
	fmt.Printf("USING %v THREADS\n", workers)

	wg.Add(workers)
	go pool(&wg, workers, paths, lr, consolidateResults)
	wg.Wait()
}

func pool(wg *sync.WaitGroup, workers int, paths []string, lr loginRes, consolidateResults *[]DownloadResult) {
	tasksCh := make(chan string)

	for i := 0; i < workers; i++ {
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

		var downloadResultTemp DownloadResult
		url := lr.orgURL + task
		attempt := 0
		validateDownloadResult := true

		expectecSize := getDownloadSize(lr, url)
		fn := fileName(url)
		fileFolderValidated := strings.Split(salesForceUserName, "@")[len(strings.Split(salesForceUserName, "@"))-1]
		filePath := fileFolderValidated + "/" + fn + ".zip"

		startDownloadTime := time.Now()
		log.Infof("Downloading file %s. Attempt: %v", filePath, attempt+1)

		for validateDownloadResult {
			filestat, err := os.Stat(filePath)
			if err != nil {
				attempt++
				err := DownloadFile(lr, filePath, url)
				if err != nil {
					downloadResultTemp.Error = err.Error()
					log.Info(err)
				}
			} else if attempt == 3 {
				log.Infof("Fail to download: %s.", fn)
				validateDownloadResult = false
				downloadResultTemp.Result = "Fail"
				os.Remove(filePath)
			} else if expectedSizeInt, _ := strconv.ParseInt(expectecSize, 10, 64); expectedSizeInt != filestat.Size() {
				attempt++
				log.Infof("The file is corrupted. Retry download attempt: %v", attempt)
				err := DownloadFile(lr, filePath, url)
				if err != nil {
					downloadResultTemp.Error = err.Error()
					log.Info(err)
				}
			} else {
				log.Infof("Successful download: %s", fn)
				downloadResultTemp.Error = "nil"
				validateDownloadResult = false
				downloadResultTemp.Result = "Successful"
			}
		}
		endDownloadTime := time.Now()
		downloadResultTemp.Duration = endDownloadTime.Sub(startDownloadTime)
		downloadResultTemp.FileName = fn
		downloadResultTemp.Attempt = attempt
		downloadResultTemp.FileSize = expectecSize

		mutex.Lock()
		*consolidateResults = append(*consolidateResults, downloadResultTemp)
		mutex.Unlock()
	}
}

func getPaths(lr loginRes) []string {
	url := lr.orgURL + "/servlet/servlet.OrgExport"
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
