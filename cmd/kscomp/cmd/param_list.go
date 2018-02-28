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
	"os"

	"github.com/bryanl/woowoo/component"
	"github.com/bryanl/woowoo/ksplugin"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var paramListCmd = &cobra.Command{
	Use:   "list",
	Short: "param list",
	Long:  `param list`,
	Run: func(cmd *cobra.Command, args []string) {
		pluginEnv, err := ksplugin.Read()
		if err != nil {
			logrus.Fatal(err)
		}

		name := viper.GetString("ns")

		ns, err := component.GetNamespace(fs, pluginEnv.AppDir, name)
		if err != nil {
			logrus.WithError(err).Fatal("could not find namespace")
		}

		paramData, err := ns.Params()
		if err != nil {
			logrus.WithError(err).Fatal("could not list params")
		}

		// TODO: we can do better than this
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"COMPONENT", "KEY", "VALUE"})
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetRowLine(false)
		table.SetBorder(false)
		for _, data := range paramData {
			table.Append([]string{data.Component, data.Key, data.Value})
		}

		table.Render()
	},
}

func init() {
	paramCmd.AddCommand(paramListCmd)

	paramListCmd.Flags().String("ns", "", "Namespace")
	viper.BindPFlag("ns", paramListCmd.Flags().Lookup("ns"))
}
