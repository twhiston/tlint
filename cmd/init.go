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
	"log"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "install the linting tools",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFlags(0)
		runcmd("brew", "install", "shellcheck", "hadolint")
		runcmd("go", "get", "-u", "github.com/mrtazz/checkmake")
		runcmd("go", "get", "-u", "github.com/alecthomas/gometalinter")
		runcmd("gometalinter", "--install")
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
