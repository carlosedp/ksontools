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
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	vRootVerbose = "root-verbose"
)

var fs = afero.NewOsFs()

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kscomp",
	Short: "ks plugin for new component work",
	Long:  `ks plugin for new component work`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		verbosity := viper.GetInt(vRootVerbose)
		logrus.SetLevel(logLevel(verbosity))

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().IntP(flagVerbose, "v", 0, "Verbosity level")
	viper.BindPFlag(vRootVerbose, rootCmd.PersistentFlags().Lookup(flagVerbose))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}

func logLevel(verbosity int) logrus.Level {
	switch verbosity {
	case 0:
		return logrus.InfoLevel
	default:
		return logrus.DebugLevel
	}
}
