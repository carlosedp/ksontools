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
	"bytes"
	"os"
	"path/filepath"
	"strings"

	kscomponent "github.com/ksonnet/ksonnet/component"
	ksparam "github.com/ksonnet/ksonnet/metadata/params"
	"github.com/ksonnet/ksonnet/prototype"
	"github.com/spf13/afero"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import manifest",
	Long:  `Import manifest`,
	Run: func(cmd *cobra.Command, args []string) {
		ksAppDir := os.Getenv("KS_APP_DIR")
		if ksAppDir == "" {
			logrus.Fatal("cannot find ks application directory")
		}

		fileName := viper.GetString("f")
		if fileName == "" {
			logrus.Fatal("-f is required")
		}

		namespace := viper.GetString("ns")

		var name bytes.Buffer
		if namespace != "" {
			name.WriteString(namespace + "/")
		}

		base := filepath.Base(fileName)
		ext := filepath.Ext(base)

		templateType, err := prototype.ParseTemplateType(strings.TrimPrefix(ext, "."))
		if err != nil {
			logrus.WithError(err).Fatal("parse template type")
		}

		name.WriteString(strings.TrimSuffix(base, ext))

		contents, err := afero.ReadFile(fs, fileName)
		if err != nil {
			logrus.WithError(err).Fatal("read manifest")
		}

		params := ksparam.Params{}

		_, err = kscomponent.Create(fs, ksAppDir, name.String(), string(contents), params, templateType)
		if err != nil {
			logrus.WithError(err).Fatal("create component")
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringP("f", "f", "", "Filename for component to import")
	viper.BindPFlag("f", importCmd.Flags().Lookup("f"))
	importCmd.Flags().String("ns", "", "Component namespace")
	viper.BindPFlag("ns", importCmd.Flags().Lookup("ns"))
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
