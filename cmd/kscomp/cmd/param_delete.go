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
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteCmd represents the delete command
var paramDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete param",
	Long:  `delete param`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			logrus.Fatal("delete <component-name> <param-key> ")
		}

		indexOpt := action.ParamDeleteWithIndex(viper.GetInt("index"))
		actionParamDelete, err := action.NewParamDelete(fs, args[0], args[1], indexOpt)
		if err != nil {
			return errors.Wrap(err, "unable to initialize param delete action")
		}

		if err := actionParamDelete.Run(); err != nil {
			return errors.Wrap(err, "delete param")
		}

		return nil
	},
}

func init() {
	paramCmd.AddCommand(paramDeleteCmd)

	paramDeleteCmd.Flags().IntP("index", "i", 0, "Index in manifest")
	viper.BindPFlag("index", paramDeleteCmd.Flags().Lookup("index"))
}
