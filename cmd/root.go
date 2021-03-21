/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nsxtea",
	Short: "A Command Line Tool to interact with the NSX-T API",
	Long: `To configure the CLI, set the following environment variables:

NSXTEA_URL (required)
NSXTEA_USERNAME (required)
NSXTEA_PASSWORD (required)
NSXTEA_INSECURE (optional, default false)
`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	nsxtInsecure, err := strconv.ParseBool(os.Getenv("NSXTEA_INSECURE"))
	if err != nil {
		nsxtInsecure = false
	}

	if nsxtInsecure {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
}
