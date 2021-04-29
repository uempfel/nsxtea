/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
)

var filePath string

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Interact with the Hierarchical Policy API",
	Long: `Decalaratively apply configurations via yaml or json
files.

Examples:
nsxtea apply -f infra.yaml
nsxtea apply -f infra.json`,
	Run: handleApply,
}

func handleApply(cmd *cobra.Command, args []string) {

	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error parsing file\n", err)
		os.Exit(1)
	}

	fileContent, err = yaml.YAMLToJSON(fileContent)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	endpoint := "/policy/api/v1/infra"
	body := bytes.NewBuffer(fileContent)
	req, err := http.NewRequest("PATCH", "https://"+url+endpoint, body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(userName, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		bodyText, err := ioutil.ReadAll(resp.Body)
		printErrIfNotNil(err)
		s := string(bodyText)
		fmt.Println(s)
		os.Exit(1)
	}
	fmt.Println(resp.Status)
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.PersistentFlags().StringVarP(&filePath, "filepath", "f", "", "Path to the file that contains the configuration to apply")
	cobra.MarkFlagRequired(applyCmd.PersistentFlags(), "filepath")
}
