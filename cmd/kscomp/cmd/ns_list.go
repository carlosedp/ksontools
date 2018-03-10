package cmd

import (
	"github.com/bryanl/woowoo/action"
	"github.com/spf13/cobra"
)

// nsListCmd represents the ns list command
var nsListCmd = &cobra.Command{
	Use:   "list",
	Short: "list",
	Long:  `list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return action.NsList(fs)
	},
}

func init() {
	nsCmd.AddCommand(nsListCmd)
}
