package cmd

import (
	"github.com/bryanl/woowoo/action"
	"github.com/spf13/cobra"
)

// envListCmd represents the env list command
var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "list",
	Long:  `list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return action.EnvList(fs)
	},
}

func init() {
	envCmd.AddCommand(envListCmd)
}
