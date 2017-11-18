// Copyright © 2017 Tom Whiston <tom.whiston@gmail.com>
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
	"fmt"
	"os"

	"log"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "tlint",
	Short: "runs the linters recursively within cwd",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		log.SetFlags(0)

		dir, err := os.Getwd()
		checkError(err)

		//goimports
		files, err := globExt(dir, ".go")
		checkError(err)
		for _, v := range files {
			runcmd("goimports", "-w", v)
		}

		//gometalinter
		runcmd("gometalinter", "./...")

		//hadolint
		files, err = glob(dir, "Dockerfile")
		checkError(err)
		for _, v := range files {
			runcmd("hadolint", v)
		}

		//shellcheck
		//TODO - what about other files without an extension, as in s2i?
		files, err = globExt(dir, ".sh")
		checkError(err)
		for _, v := range files {
			runcmd("shellcheck", v)
		}

		//checkmake
		files, err = glob(dir, "Makefile")
		checkError(err)
		for _, v := range files {
			runcmd("checkmake", v)
		}
	},
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tlint.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

		// Search config in home directory with name ".tlint" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".tlint")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
