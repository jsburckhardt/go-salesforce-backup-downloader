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
	"github.com/go-playground/log/handlers/console"
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

// rootCmd -> download org files (main purpose)
var rootCmd = &cobra.Command{
	Use:   "go-salesforce-backup-downloader",
	Short: "Commandline application for downloading SalesForce backup files.",
	Long: `A commandline application for downloading multiple files concurrently.
You should execute giving Username (u) and Password (p). For example:
go-salesforce-backup-downloader.exe -u sadmin@atyourcrazyorg -p mypasswordwithtoken
go-salesforce-backup-downloader.exe --user sadmin@atyourcrazyorg --password mypasswordwithtoken
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Future development -> accept config file
	//cobra.OnInitialize(initConfig)

	cLog := console.New(true)
	log.AddHandler(cLog, log.AllLevels...)

	httpClient = &http.Client{
		Timeout: time.Minute * 10,
	}

	// global parameters for the application.
	// Future development -> accept config file
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-salesforce-backup-downloader.yaml)")
	rootCmd.PersistentFlags().StringVarP(&salesForceUserName, "user", "u", "", "Salesforce username e.g something@somewhere.there")
	rootCmd.PersistentFlags().StringVarP(&salesForceUserPassword, "password", "p", "", "Salesforce password+token e.g supersecretpasswordwithtoken")
	rootCmd.PersistentFlags().IntVarP(&maxWorkers, "maxworkers", "m", 5, "Maximum number of workers for concurrency. (default is 5)")
	rootCmd.MarkPersistentFlagRequired("user")
	rootCmd.MarkPersistentFlagRequired("password")

	// debug flag should be use for verbosity
	//rootCmd.PersistentFlags().BoolVarP(&debugStatus, "debug", "d", false, "debug?")
}

// initConfig reads in config file and ENV variables if set.
// currently no config file is being passed.
// func initConfig() {
// 	if cfgFile != "" {
// 		// Use config file from the flag.
// 		viper.SetConfigFile(cfgFile)
// 	} else {
// 		// Find home directory.
// 		home, err := homedir.Dir()
// 		if err != nil {
// 			fmt.Println(err)
// 			os.Exit(1)
// 		}

// 		// Search config in home directory with name ".go-salesforce-backup-downloader" (without extension).
// 		viper.AddConfigPath(home)
// 		viper.SetConfigName(".go-salesforce-backup-downloader")
// 	}

// 	viper.AutomaticEnv() // read in environment variables that match

// 	// If a config file is found, read it in.
// 	if err := viper.ReadInConfig(); err == nil {
// 		fmt.Println("Using config file:", viper.ConfigFileUsed())
// 	}
//}
