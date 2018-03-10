package cmd

import "github.com/spf13/cobra"

// nsCmd represents the ns command
var nsCmd = &cobra.Command{
	Use:   "ns",
	Short: "ns",
	Long:  `ns`,
}

func init() {
	rootCmd.AddCommand(nsCmd)
}
