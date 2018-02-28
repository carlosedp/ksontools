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
	"github.com/bryanl/woowoo/action"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show a component",
	Long:  `show a component`,
	RunE: func(cmd *cobra.Command, args []string) error {
		env := viper.GetString("env")
		components := viper.GetStringSlice("component")

		showAction, err := action.NewShow(fs, env, action.ShowWithComponents(components...))
		if err != nil {
			return err
		}

		if err := showAction.Run(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	showCmd.Flags().String("env", "default", "Environment")
	viper.BindPFlag("env", showCmd.Flags().Lookup("env"))

	showCmd.Flags().StringSliceP("component", "c", nil, "Components to include")
	viper.BindPFlag("component", showCmd.Flags().Lookup("component"))
}
