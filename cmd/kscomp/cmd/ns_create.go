package cmd

import (
	"github.com/bryanl/woowoo/action"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	vNsCreateNamespace = "ns-create-namespace"
)

// nsCreateCmd creates a ns create command.
var nsCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "create",
	Long:  `create`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("ns create <namespace>")
		}

		nsName := args[0]

		return action.NsCreate(fs, nsName)
	},
}

func init() {
	nsCmd.AddCommand(nsCreateCmd)

}
