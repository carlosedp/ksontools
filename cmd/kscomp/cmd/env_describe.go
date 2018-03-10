package cmd

import (
	"github.com/bryanl/woowoo/action"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// envDescribeCmd represents the env describe command
var envDescribeCmd = &cobra.Command{
	Use:   "describe <env>",
	Short: "describe",
	Long:  `describe`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("env describe <environment>")
		}

		environment := args[0]

		return action.EnvDescribe(fs, environment)
	},
}

func init() {
	envCmd.AddCommand(envDescribeCmd)
}
