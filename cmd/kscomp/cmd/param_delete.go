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
	"strings"

	"github.com/bryanl/woowoo/component"
	"github.com/bryanl/woowoo/ksplugin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var paramDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete param",
	Long:  `delete param`,
	Run: func(cmd *cobra.Command, args []string) {
		pluginEnv, err := ksplugin.Read()
		if err != nil {
			logrus.Fatal(err)
		}

		componentName := args[0]

		path := strings.Split(args[1], ".")

		c, err := component.ExtractComponent(fs, pluginEnv.AppDir, componentName)
		if err != nil {
			logrus.WithError(err).Fatal("could not find component")
		}

		if err := c.DeleteParam(path, component.ParamOptions{}); err != nil {
			logrus.WithError(err).Fatal("delete param")
		}
	},
}

func init() {
	paramCmd.AddCommand(paramDeleteCmd)

}
