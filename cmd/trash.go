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

// trashCmd represents the trash command
var trashCmd = &cobra.Command{
	Use:   "trash",
	Short: "Empty trash can",
	Long:  `Empty trash can`,
	Args: func(cmd *cobra.Command, args []string) error {
		if url == "" {
			return errors.New("requires Artifactory's url")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		sEmptyTrashCanURL := fmt.Sprintf(sEmptyTrashCanURLTemplate, url)
		fmt.Println(sEmptyTrashCanURL)
		req, err := http.NewRequest("POST", sEmptyTrashCanURL, nil)
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

		fmt.Println("trash called")
	},
}

func init() {
	RootCmd.AddCommand(trashCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	trashCmd.PersistentFlags().StringVarP(&user, "user", "u", "", "username:password")
	trashCmd.PersistentFlags().StringVar(&url, "url", "", "Artifactory's url")
	trashCmd.PersistentFlags().UintVarP(&sr, "show_response", "r", 0, "Show data response")
}

const sEmptyTrashCanURLTemplate = "%s/api/trash/empty"
