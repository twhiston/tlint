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

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"log"
)

var cfgFile string
var localConfPath string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "tlint",
	Short: "runs linters recursively within cwd",
	Long: `
All linters enabled in a .tlint.yml file will be run over the current cwd and all directories below.
If a .tlint.yml file cannot be found tlint will look for $HOME/.tlint.yml

Config file will be overriden by options. For example
	tlint -m=false
ensures that gometalinter is never run and
	tlint -d
ensures that hadolint is always run
	`,
	Run: func(cmd *cobra.Command, args []string) {

		log.SetFlags(0)

		dir, err := os.Getwd()
		checkError(err)

		//goimports
		if viper.GetBool("status.imports") {
			files, err := globExt(dir, ".go")
			checkError(err)
			runArgsCmd("goimports", files, []string{"-w"})
		}

		//gofmt
		if viper.GetBool("status.fmt") {
			files, err := globExt(dir, ".go")
			checkError(err)
			runArgsCmd("go", files, []string{"fmt"})
		}

		//gometalinter
		if viper.GetBool("status.gometalinter") {
			cmdArgs := []string{
				"./...",
			}
			if _, err := os.Stat(localConfPath + ".gometalinter.json"); !os.IsNotExist(err) {
				cmdArgs = append(cmdArgs, "--config", localConfPath+"gometalinter.json")
			}
			runcmd("gometalinter", cmdArgs...)
		}

		//hadolint
		if viper.GetBool("status.hadolint") {
			files, err := glob(dir, "Dockerfile")
			checkError(err)
			cmdArgs := buildCmdArgs("hadolint.ignore", "--ignore")
			runArgsCmd("hadolint", files, cmdArgs)
		}

		//shellcheck
		if viper.GetBool("status.shellcheck") {
			files, err := globExt(dir, ".sh")
			checkError(err)
			cmdArgs := buildCmdArgs("shellcheck.ignore", "-e")
			runArgsCmd("shellcheck", files, cmdArgs)
		}

		//shellcheck every file in a bin directory in the image, as these are usually shell scripts
		//Could be better!
		if viper.GetBool("status.shellcheck_bin") {
			files, err := getDirContents(dir, "bin")
			checkError(err)
			cmdArgs := buildCmdArgs("shellcheck.ignore", "-e")
			runArgsCmd("shellcheck", files, cmdArgs)
		}

		//checkmake
		if viper.GetBool("status.checkmake") {
			files, err := glob(dir, "Makefile")
			checkError(err)

			var cmdArgs []string
			if _, err := os.Stat(localConfPath + ".checkmake.ini"); !os.IsNotExist(err) {
				cmdArgs = append(cmdArgs, "--config="+localConfPath+".checkmake.ini")
			}
			runArgsCmd("checkmake", files, cmdArgs)
		}
	},
}

func buildCmdArgs(sliceKey string, optionKey string) []string {
	var cmdArgs []string
	cmdData := viper.GetStringSlice(sliceKey)
	if len(cmdData) > 0 {
		for _, v := range cmdData {
			cmdArgs = append(cmdArgs, optionKey, v)
		}
	}
	return cmdArgs
}

func runArgsCmd(cmd string, files []string, cmdArgs []string) {
	for _, v := range files {
		fileArgs := append(cmdArgs, v)
		runcmd(cmd, fileArgs...)
	}
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
	log.SetFlags(0)
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tlint.yml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("imports", "i", true, "fix go imports")
	RootCmd.Flags().BoolP("fmt", "f", true, "fix go formatting")
	RootCmd.Flags().BoolP("gometalinter", "m", true, "run gometalinter")
	RootCmd.Flags().BoolP("hadolint", "d", true, "run hadolint")
	RootCmd.Flags().BoolP("shellcheck", "s", true, "run shellcheck")
	RootCmd.Flags().BoolP("shellcheck-bin", "b", true, "run shellcheck on ANY file in a folder called bin, useful for s2i image linting")
	RootCmd.Flags().BoolP("checkmake", "c", true, "run checkmake")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	configName := ".tlint"

	dir, err := os.Getwd()
	checkError(err)
	localConfPath = dir + "/"

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		localConf := localConfPath + configName + ".yml"
		if _, err := os.Stat(localConf); os.IsNotExist(err) {
			// Find home directory.
			home, err := homedir.Dir()
			checkError(err)
			// Search config in home directory with name ".tlint" (without extension).
			viper.AddConfigPath(home)
			viper.SetConfigName(configName)
		} else {
			viper.SetConfigFile(localConf)
		}
	}

	viper.SetEnvPrefix("TL_")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		viper.Set("status.imports", true)
		viper.Set("status.fmt", true)
		viper.Set("status.gometalinter", true)
		viper.Set("status.hadolint", true)
		viper.Set("status.shellcheck", true)
		viper.Set("status.shellcheck_bin", true)
		viper.Set("status.checkmake", true)
	}

	ftests := []string{"imports", "fmt", "gometalinter", "hadolint", "shellcheck", "shellcheck-bin", "checkmake"}

	for _, v := range ftests {
		if RootCmd.Flags().Changed(v) {
			log.Println("Config overriden by option:", v)
			val, err := RootCmd.Flags().GetBool(v)
			checkError(err)
			viper.Set("status."+v, val)
		}
	}

}
