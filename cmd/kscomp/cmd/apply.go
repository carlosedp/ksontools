package cmd

import (
	"github.com/bryanl/woowoo/action"
	"github.com/ksonnet/ksonnet/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	vApplyCreate = "apply-create"
	vApplyDryRun = "apply-dru-run"
	vApplyGcTag  = "apply-gc-tag"
	vApplySkipGc = "apply-skip-gc"
)

var (
	applyClientConfig *client.Config
)

// showCmd represents the show command
var applyCmd = &cobra.Command{
	Use:   "apply <environment>",
	Short: "apply a component",
	Long:  `apply a component`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("apply <environment>")
		}

		env := args[0]

		options := action.ApplyOptions{
			Create: viper.GetBool(vApplyCreate),
			SkipGc: viper.GetBool(vApplySkipGc),
			GcTag:  viper.GetString(vApplyGcTag),
			DryRun: viper.GetBool(vApplyDryRun),
			Client: applyClientConfig,
		}

		return action.Apply(fs, env, options)
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyClientConfig = client.NewDefaultClientConfig()
	applyClientConfig.BindClientGoFlags(applyCmd)

	applyCmd.Flags().Bool(flagCreate, true, "Option to create resources if they do not already exist on the cluster")
	viper.BindPFlag(vApplyCreate, applyCmd.Flags().Lookup(flagCreate))

	applyCmd.Flags().Bool(flagSkipGc, false, "Option to skip garbage collection, even with --"+flagGcTag+" specified")
	viper.BindPFlag(vApplySkipGc, applyCmd.Flags().Lookup(flagSkipGc))

	applyCmd.Flags().String(flagGcTag, "", "A tag that's (1) added to all updated objects (2) used to garbage collect existing objects that are no longer in the manifest")
	viper.BindPFlag(vApplyGcTag, applyCmd.Flags().Lookup(flagGcTag))

	applyCmd.Flags().Bool(flagDryRun, false, "Option to preview the list of operations without changing the cluster state")
	viper.BindPFlag(vApplyDryRun, applyCmd.Flags().Lookup(flagDryRun))
}
