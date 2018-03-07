package cmd

import "github.com/spf13/cobra"

var componentCmd = &cobra.Command{
	Use:   "component",
	Short: "component",
	Long:  "component",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(componentCmd)
}
