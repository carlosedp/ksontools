package cmd

import (
	"github.com/bryanl/woowoo/action"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	vComponentListNamespace = "component-list-ns"
)

var componentListCmd = &cobra.Command{
	Use:   "list",
	Short: "component list",
	Long:  "component list",
	RunE: func(cmd *cobra.Command, args []string) error {
		namespace := viper.GetString(vComponentListNamespace)
		return action.ComponentList(fs, namespace)
	},
}

func init() {
	componentCmd.AddCommand(componentListCmd)

	componentListCmd.Flags().String(flagNamespace, "", "Component namespace")
	viper.BindPFlag(vComponentListNamespace, componentListCmd.Flags().Lookup(flagNamespace))
}
