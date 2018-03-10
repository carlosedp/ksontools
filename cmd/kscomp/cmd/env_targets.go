package cmd

import (
	"github.com/bryanl/woowoo/action"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	vEnvTargetComponents = "env-target-components"
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

		components := viper.GetStringSlice(vEnvTargetComponents)
		return action.EnvTargets(fs, environment, components)
	},
}

func init() {
	envCmd.AddCommand(envTargetsCmd)

	envTargetsCmd.Flags().StringSliceP(flagComponent, "c", nil, "Components to include")
	viper.BindPFlag(vEnvTargetComponents, envTargetsCmd.Flags().Lookup(flagComponent))
}
