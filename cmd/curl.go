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
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

var bodyData string
var method string
var override bool

var curlCmd = &cobra.Command{
	Use:   "curl <endpoint>",
	Short: "Interact with NSX-T's REST-API",
	Long: `Interact with NSX-T's REST-API.

Examples:
nsxtea curl -X DELETE  /policy/api/v1/
nsxtea curl -X PUT /api/v1/certificate -d @path-to-body-file`,
	Run: handleCurl,
	Args: cobra.MinimumNArgs(1),
}

func handleCurl(cmd *cobra.Command, args []string) {
	var body io.Reader = nil
	if bodyData != "" {
		body = strings.NewReader(bodyData)
	}
    
	isFile := strings.HasPrefix(bodyData,"@")
	if isFile {
		filePath = strings.TrimPrefix(bodyData,"@")
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

		body = bytes.NewReader(fileContent)
	}
	


	endpoint := args[0]

	req, err := http.NewRequest(method, "https://"+url+endpoint, body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")
	if override {
		req.Header.Set("X-Allow-Overwrite", "true")
	}
	req.SetBasicAuth(userName, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	printErrIfNotNil(err)
	s := string(bodyText)
	fmt.Println(s)
}

func init() {
	rootCmd.AddCommand(curlCmd)
	curlCmd.PersistentFlags().StringVarP(&bodyData, "data", "d", "", "Body data. You can specifiy a path to a yaml or json file with the '@' prefix")
	curlCmd.PersistentFlags().StringVarP(&method, "method", "X", "GET", "HTTP Method")
	curlCmd.PersistentFlags().BoolVarP(&override, "override", "o", false, "Add the 'X-Allow-Overwrite: true' header to mutate protected objects")
}
