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
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import manifest",
	Long:  `Import manifest`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fileName := viper.GetString("f")
		if fileName == "" {
			return errors.New("-f is required")
		}

		namespace := viper.GetString("ns")

		importAction, err := action.NewImport(fs, namespace, fileName)
		if err != nil {
			return err
		}

		if err := importAction.Run(); err != nil {
			logrus.Errorf("import failed: %+v", err)
			return errors.Wrap(err, "unable to import file or directory")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringP("f", "f", "", "Filename or directory for component to import")
	viper.BindPFlag("f", importCmd.Flags().Lookup("f"))
	importCmd.Flags().String("ns", "", "Component namespace")
	viper.BindPFlag("ns", importCmd.Flags().Lookup("ns"))
}
