package cmd

import (
	"github.com/bryanl/woowoo/action"
	"github.com/bryanl/woowoo/pkg/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	vDeleteGracePeriod = "delete-grace-period"
)

var (
	deleteClientConfig *client.Config
)

// showCmd represents the show command
var deleteCmd = &cobra.Command{
	Use:   "delete <environment>",
	Short: "delete a component",
	Long:  `delete a component`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("delete <environment>")
		}
		env := args[0]

		components := viper.GetStringSlice(vShowComponent)
		gracePeriod := viper.GetInt64(vDeleteGracePeriod)

		options := client.DeleteOptions{
			GracePeriod: gracePeriod,
			Client:      deleteClientConfig,
		}

		return action.Delete(fs, env, options, action.DeleteWithComponents(components...))
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteClientConfig = client.NewDefaultClientConfig()
	deleteClientConfig.BindClientGoFlags(deleteCmd)

	deleteCmd.Flags().Int64(flagGracePeriod, -1, "Number of seconds given to resources to terminate gracefully. A negative value is ignored")
	viper.BindPFlag(vDeleteGracePeriod, deleteCmd.Flags().Lookup(flagGracePeriod))
}
