// Copyright © 2018 Juan Burckhardt <jsburckhardt>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-playground/log"
	"github.com/spf13/cobra"
)

//DownloadResult struct for structuring download results once
//files are successfully downloaded or failed.
type DownloadResult struct {
	FileName, FileSize, Result, Error string
	Attempt                           int
	Duration                          time.Duration
}

var cfgFile string
var debugStatus bool
var salesForceUserName string
var salesForceUserPassword string
var maxWorkers int
var httpClient *http.Client

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "download backup files from SalesForce.",
	Long: `download backup files from SalesForce concurrently.
You should execute giving Username (u) and Password (p). For example:
go-salesforce-backup-downloader.exe download -u sadmin@atyourcrazyorg -p mypasswordwithtoken
go-salesforce-backup-downloader.exe download --user sadmin@atyourcrazyorg --password mypasswordwithtoken
`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		fmt.Printf("Start time -> %s\n", start.Format(time.ANSIC))
		loginData := login()

		var consolidateResults []DownloadResult

		backupfolder := strings.Split(salesForceUserName, "@")[len(strings.Split(salesForceUserName, "@"))-1]

		if _, err := os.Stat(backupfolder); os.IsNotExist(err) {
			log.Infof("backupfolder folder doesn't exist. Creating: %s", backupfolder)
			err := os.Mkdir(backupfolder, 0777)
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

	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	httpClient = &http.Client{
		Timeout: time.Minute * 10,
	}

	// global parameters for the application.
	// Future development -> accept config file
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-salesforce-backup-downloader.yaml)")
	downloadCmd.PersistentFlags().StringVarP(&salesForceUserName, "user", "u", "", "Salesforce username e.g something@somewhere.there")
	downloadCmd.PersistentFlags().StringVarP(&salesForceUserPassword, "password", "p", "", "Salesforce password+token e.g supersecretpasswordwithtoken")
	downloadCmd.PersistentFlags().IntVarP(&maxWorkers, "maxworkers", "m", 5, "Maximum number of workers for concurrency. (default is 5)")
	downloadCmd.MarkPersistentFlagRequired("user")
	downloadCmd.MarkPersistentFlagRequired("password")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
