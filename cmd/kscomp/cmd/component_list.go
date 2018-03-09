package cmd

import (
	"github.com/bryanl/woowoo/action"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	vComponentListNamespace = "component-list-ns"
	vComponentListOutput    = "component-list-output"
)

var componentListCmd = &cobra.Command{
	Use:   "list",
	Short: "component list",
	Long:  "component list",
	RunE: func(cmd *cobra.Command, args []string) error {
		namespace := viper.GetString(vComponentListNamespace)
		output := viper.GetString(vComponentListOutput)
		return action.ComponentList(fs, namespace, output)
	},
}

func init() {
	componentCmd.AddCommand(componentListCmd)

	componentListCmd.Flags().String(flagNamespace, "", "Component namespace")
	viper.BindPFlag(vComponentListNamespace, componentListCmd.Flags().Lookup(flagNamespace))

	componentListCmd.Flags().StringP(flagOutput, "o", "", "Output format. Valid options: wide")
	viper.BindPFlag(vComponentListOutput, componentListCmd.Flags().Lookup(flagOutput))
}
