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

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// testCredentialsCmd represents the testCredentials command
var testCredentialsCmd = &cobra.Command{
	Use:   "testCredentials",
	Short: "Test credentials for org",
	Long: `use testConnection for testing Salesforce Credentials.
For example:
go-salesforce-backup-downloader.exe testConnection -u sadmin@atyourcrazyorg -p mypasswordwithtoken
go-salesforce-backup-downloader.exe testConnection --user sadmin@atyourcrazyorg --password mypasswordwithtoken`,
	Run: func(cmd *cobra.Command, args []string) {
		loginResult := login()
		if len(loginResult.sID) > 0 {
			color.Set(color.FgHiGreen)
			fmt.Printf("*** Login successful %s ***\n", salesForceUserName)
			color.Unset()
		}
	},
}

func init() {
	rootCmd.AddCommand(testCredentialsCmd)

	httpClient = &http.Client{
		Timeout: time.Minute * 10,
	}

	testCredentialsCmd.PersistentFlags().StringVarP(&salesForceUserName, "user", "u", "", "Salesforce username e.g something@somewhere.there")
	testCredentialsCmd.PersistentFlags().StringVarP(&salesForceUserPassword, "password", "p", "", "Salesforce password+token e.g supersecretpasswordwithtoken")
	testCredentialsCmd.MarkPersistentFlagRequired("user")
	testCredentialsCmd.MarkPersistentFlagRequired("password")
}
