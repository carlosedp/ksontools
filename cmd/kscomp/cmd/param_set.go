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
)

// setCmd represents the set command
var paramSetCmd = &cobra.Command{
	Use:   "set",
	Short: "param set",
	Long:  `param set`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 3 {
			logrus.Fatal("set <component-name> <param-key> <param-value>")
		}

		actionParamSet, err := action.NewParamSet(fs, args[0], args[1], args[2])
		if err != nil {
			return errors.Wrap(err, "unable to initialize param set action")
		}

		if err := actionParamSet.Run(); err != nil {
			return errors.Wrap(err, "set param")
		}

		return nil
	},
}

func init() {
	paramCmd.AddCommand(paramSetCmd)

	// TODO: support env
}
