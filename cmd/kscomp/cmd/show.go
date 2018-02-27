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
	"github.com/bryanl/woowoo/ksutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show a component",
	Long:  `show a component`,
	Run: func(cmd *cobra.Command, args []string) {
		pluginEnv, err := ksplugin.Read()
		if err != nil {
			logrus.Fatal(err)
		}

		env := viper.GetString("env")

		app := ksutil.NewApp(fs, pluginEnv.AppDir)
		namespaces, err := component.NamespacesFromEnv(fs, app, pluginEnv.AppDir, env)
		if err != nil {
			logrus.WithError(err).Fatal("find namespaces")
		}

		var objects []*unstructured.Unstructured
		for _, ns := range namespaces {
			members, err := ns.Components()
			if err != nil {
				logrus.WithError(err).Fatal("find components")
			}
			for _, c := range members {
				o, err := c.Objects()
				if err != nil {
					logrus.WithError(err).Fatal("get objects")
				}
				objects = append(objects, o...)
			}
		}

		ksutil.Fprint(os.Stdout, objects, "yaml")
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	showCmd.Flags().String("env", "default", "Environment")
	viper.BindPFlag("env", showCmd.Flags().Lookup("env"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
