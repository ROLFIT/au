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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// purgeCmd represents the purge command
var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Delete old artifacts",
	Long:  `Delete old artifacts`,
	Args: func(cmd *cobra.Command, args []string) error {
		if url == "" {
			return errors.New("requires Artifactory's url")
		}
		if repo == "" {
			return errors.New("requires Artifactory's repository name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		bgDate := time.Now().Add(time.Duration(-(100 + days)) * time.Hour * 24)
		fnDate := bgDate.Add(100 * time.Hour * 24)
		bg := bgDate.UnixNano() / 1000000
		fn := fnDate.UnixNano() / 1000000

		sCandidatesURL := fmt.Sprintf(sCandidatesURLTemplate, url, bg, fn, repo)
		fmt.Println(sCandidatesURL)
		req, err := http.NewRequest("GET", sCandidatesURL, nil)
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

		uris := struct {
			Results []struct {
				URI     string
				Created time.Time
			}
		}{Results: []struct {
			URI     string
			Created time.Time
		}{}}
		if err := json.NewDecoder(r.Body).Decode(&uris); err != nil {
			log.Fatal(err)
		}
		minDt := time.Now()
		maxDt := time.Time{}

		regExp, err := regexp.Compile(fmt.Sprintf("%s/%s/[A-Za-z]_\\d*.\\d*/\\d*/", url, repo))
		if err != nil {
			log.Fatal(err)
		}

		folders := make(map[string]bool, 0)
		for _, v := range uris.Results {
			if v.Created.After(maxDt) {
				maxDt = v.Created
			}
			if v.Created.Before(minDt) {
				minDt = v.Created
			}
			delURI := strings.Replace(v.URI, "/api/storage", "", 1)
			if regExp.MatchString(delURI) {
				delURI = regExp.FindString(delURI)
				delURI = delURI[:len(delURI)-1]
				if _, ok := folders[delURI]; !ok {
					folders[delURI] = true
				}
			}
		}
		for k := range folders {
			reqDel, err := http.NewRequest("DELETE", k, nil)
			if err != nil {
				log.Fatal(err)
			}
			reqDel.SetBasicAuth(u[0], u[1])
			respDel, err := cli.Do(reqDel)
			if err != nil {
				log.Fatal(err)
			}
			defer respDel.Body.Close()
			if respDel.StatusCode != 204 {
				log.Fatal(fmt.Errorf(fmt.Sprintf("Invalid status code %d", respDel.StatusCode)))
			}
			fmt.Println(k, "deleted")
		}
		fmt.Println(minDt, maxDt)
		fmt.Println("purge called")
	},
}

func init() {
	RootCmd.AddCommand(purgeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	purgeCmd.PersistentFlags().StringVarP(&user, "user", "u", "", "username:password")
	purgeCmd.PersistentFlags().StringVar(&url, "url", "", "Artifactory's url")
	purgeCmd.PersistentFlags().StringVar(&repo, "repo", "", "Artifactory's repository name")
	purgeCmd.PersistentFlags().UintVarP(&days, "days", "d", 7, "Number of days")
	purgeCmd.PersistentFlags().UintVarP(&sr, "show_response", "r", 0, "Show data response")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// purgeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

const sCandidatesURLTemplate = "%s/api/search/dates?dateFields=created,lastModified,lastDownloaded&from=%d&to=%d&repos=%s"
