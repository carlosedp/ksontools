package cmd

import "github.com/spf13/cobra"

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "env",
	Long:  `env`,
}

func init() {
	rootCmd.AddCommand(envCmd)
}
