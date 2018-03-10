package cmd

import (
	"github.com/bryanl/woowoo/action"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	vEnvTargetNamespaces = "env-target-namespaces"
)

// envTargetsCmd represents the env targets command
var envTargetsCmd = &cobra.Command{
	Use:   "targets",
	Short: "targets",
	Long:  `targets`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("env targets <environment> --component name")
		}

		environment := args[0]

		components := viper.GetStringSlice(vEnvTargetNamespaces)
		return action.EnvTargets(fs, environment, components)
	},
}

func init() {
	envCmd.AddCommand(envTargetsCmd)

	envTargetsCmd.Flags().StringSlice(flagNamespace, nil, "Components to include")
	viper.BindPFlag(vEnvTargetNamespaces, envTargetsCmd.Flags().Lookup(flagNamespace))
}
