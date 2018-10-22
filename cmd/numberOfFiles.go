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
	"time"

	"github.com/spf13/cobra"
)

// numberOfFilesCmd represents the numberOfFiles command
var numberOfFilesCmd = &cobra.Command{
	Use:   "numberOfFiles",
	Short: "Returns the number of backup files available for the org",
	Long: `Returns the number of backup files available for the requested org. 	
For example:
go-salesforce-backup-downloader.exe testConnection -u sadmin@atyourcrazyorg -p mypasswordwithtoken
go-salesforce-backup-downloader.exe testConnection --user sadmin@atyourcrazyorg --password mypasswordwithtoken`,
	Run: func(cmd *cobra.Command, args []string) {
		loginData := login()
		paths := getPaths(loginData)
		if len(paths[0]) == 0 {
			fmt.Printf("Number of backup files for %s -> 0", salesForceUserName)
		}
		fmt.Printf("%v\n", len(paths))
	},
}

func init() {
	rootCmd.AddCommand(numberOfFilesCmd)

	httpClient = &http.Client{
		Timeout: time.Minute * 10,
	}

	numberOfFilesCmd.PersistentFlags().StringVarP(&salesForceUserName, "user", "u", "", "Salesforce username e.g something@somewhere.there")
	numberOfFilesCmd.PersistentFlags().StringVarP(&salesForceUserPassword, "password", "p", "", "Salesforce password+token e.g supersecretpasswordwithtoken")
	numberOfFilesCmd.MarkPersistentFlagRequired("user")
	numberOfFilesCmd.MarkPersistentFlagRequired("password")
}
