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

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nsxtea",
	Short: "A Command Line Tool to interact with the NSX-T Policy API",
	Long: `To configure the CLI, set following Environment Variables:

NSXTEA_URL (required)
NSXTEA_USERNAME (required)
NSXTEA_PASSWORD (required)
NSXTEA_INSECURE (optional, default false)
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	nsxtInsecure, err := strconv.ParseBool(os.Getenv("NSXTEA_INSECURE"))
	if err != nil {
		nsxtInsecure = false
	}

	if nsxtInsecure {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// nsxtHost := os.Getenv("NSXTEA_URL")
	// nsxtUsername := os.Getenv("NSXTEA_USERNAME")
	// nsxtPassword := os.Getenv("NSXTEA_PASSWORD")

	// if len(nsxtHost) == 0 {
	// 	fmt.Println("NSXTEA_URL not set")
	// 	os.Exit(1)
	// }

	// if len(nsxtUsername) == 0 {
	// 	fmt.Println("NSXTEA_USERNAME not set")
	// 	os.Exit(1)
	// }

	// if len(nsxtPassword) == 0 {
	// 	fmt.Println("NSXTEA_PASSWORD not set")
	// 	os.Exit(1)
	// }

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nsxtea.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".nsxtea" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".nsxtea")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
