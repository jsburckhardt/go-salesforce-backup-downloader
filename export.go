package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/dustin/go-humanize"
	"github.com/spf13/viper"
)

func export(lr loginRes) {
	//paths := []string{"/servlet/servlet.OrgExport?fileName=WE_00D90000000JJYKEA4_1.ZIP&id=0920W00000t5oFk "}
	paths := getPaths(lr)
	fmt.Println(paths)

	fmt.Printf("NUMBER OF URLS TO DOWNLOAD: %v\n", len(paths))
	if len(paths) == 0 {
		log.Fatalln("NO BACKUP FILES!!!")
	}

	var wg sync.WaitGroup
	fmt.Printf("Creating workers\n")
	wg.Add(3)
	fmt.Printf("Creating pool\n")
	go pool(&wg, 3, paths, lr)
	wg.Wait()
}

// /////////////////////////////////////////////////////////////////////////////////////

// 	done := make(chan bool, len(paths))
// 	errch := make(chan error, len(paths))

// 	for _, path := range paths {

// 		go func(path string) {
// 			url := viper.GetString("sf.baseUrl") + path

// 			fmt.Printf("Working on url: %s\n", url)
// 			expectecSize := getDownloadSize(lr, url)
// 			fmt.Printf("size: %s\n", expectecSize)
// 			fn := fileName(url)
// 			filePath := viper.GetString("sf.backuppath") + fn + ".zip"
// 			//retry := 0

// 			err := DownloadFile(lr, filePath, url)

// 			//panic(err)
// 			if err != nil {
// 				errch <- err
// 				done <- false
// 				return
// 			}
// 			done <- true
// 			errch <- nil
// 		}(path)
// 	}

// 	var result []bool
// 	var errStr string
// 	for i := 0; i < len(paths); i++ {
// 		//bytesArray = append(bytesArray, <-done)
// 		result = append(result, <-done)
// 		if err := <-errch; err != nil {
// 			errStr = errStr + " " + err.Error()
// 		}
// 	}
// 	fmt.Println(errStr)
// 	fmt.Println(result)

// }

func pool(wg *sync.WaitGroup, workers int, paths []string, lr loginRes) {
	tasksCh := make(chan string)

	for i := 0; i < workers; i++ {
		fmt.Printf("Creating pool of workers. Worker: %v\n", i+1)
		go worker(tasksCh, wg, lr)
	}

	for _, v := range paths {
		tasksCh <- v
	}
	close(tasksCh)
}

func worker(tasksCh <-chan string, wg *sync.WaitGroup, lr loginRes) {
	defer wg.Done()

	for {
		task, ok := <-tasksCh
		if !ok {
			return
		}

		// getting file's expected size
		url := viper.GetString("sf.baseUrl") + task
		expectecSize := getDownloadSize(lr, url)
		fmt.Printf("downloading file -> %s with size %v\n", task, expectecSize)

		// downloading file
		fn := fileName(url)
		filePath := viper.GetString("sf.backuppath") + fn + ".zip"
		err := DownloadFile(lr, filePath, url)
		if err != nil {
			log.Fatalln(err)
			return
		}

		//verifying file size
		fi, err := os.Stat(filePath)
		if err != nil {
			log.Fatalln(err)
		}
		expectedSizeInt, err := strconv.ParseInt(expectecSize, 10, 64)
		if expectedSizeInt == fi.Size() {
			completedfn := viper.GetString("sf.backuppath") + fn
			_, err := os.Create(completedfn)
			if err != nil {
				log.Fatalln(err)
			}
		}
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
	fmt.Println("Getting download size...")
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
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	// resp, err := http.Get(url)
	// if err != nil {
	// 	return err
	// }
	// defer resp.Body.Close()

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

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Print("\n")

	out.Close()
	err = os.Rename(filepath+".tmp", filepath)
	if err != nil {
		return err
	}

	return nil
}

// WriteCounter ...
type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}
