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

const (
	vShowEnv       = "show-env"
	vShowComponent = "show-component"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show a component",
	Long:  `show a component`,
	RunE: func(cmd *cobra.Command, args []string) error {
		env := viper.GetString(vShowEnv)
		components := viper.GetStringSlice(vShowComponent)

		return action.Show(fs, env, action.ShowWithComponents(components...))
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	showCmd.Flags().String(flagEnv, "default", "Environment")
	viper.BindPFlag(vShowEnv, showCmd.Flags().Lookup(flagEnv))

	showCmd.Flags().StringSliceP(flagComponent, "c", nil, "Components to include")
	viper.BindPFlag(vShowComponent, showCmd.Flags().Lookup(flagComponent))
}
