// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

// optimizeCmd represents the optimize command
var optimizeCmd = &cobra.Command{
	Use:   "optimize",
	Short: "Optimize system storage",
	Long:  `Optimize system storage`,
	Args: func(cmd *cobra.Command, args []string) error {
		if url == "" {
			return errors.New("requires Artifactory's url")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		sOptimizeURL := fmt.Sprintf(sOptimizeURLTemplate, url)
		fmt.Println(sOptimizeURL)
		req, err := http.NewRequest("POST", sOptimizeURL, nil)
		if err != nil {
			log.Fatal(err)
		}
		u := strings.Split(user, ":")
		req.SetBasicAuth(u[0], u[1])

		cli := &http.Client{}
		r, err := cli.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer r.Body.Close()

		fmt.Println("response Status:", r.Status)
		if sr == 1 {
			fmt.Println("response Headers:", r.Header)
			body, _ := ioutil.ReadAll(r.Body)
			fmt.Println("response Body:", string(body))
		}

		fmt.Println("optimize called")
	},
}

func init() {
	RootCmd.AddCommand(optimizeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	optimizeCmd.PersistentFlags().StringVarP(&user, "user", "u", "", "username:password")
	optimizeCmd.PersistentFlags().StringVar(&url, "url", "", "Artifactory's url")
	optimizeCmd.PersistentFlags().UintVarP(&sr, "show_response", "r", 0, "Show data response")
}

const sOptimizeURLTemplate = "%s/api/system/storage/optimize"
