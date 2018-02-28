package cmd

import (
	"github.com/bryanl/woowoo/action"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var paramSetGlobalCmd = &cobra.Command{
	Use:   "set-global",
	Short: "param set-global",
	Long:  "param set-global",
	RunE: func(cmd *cobra.Command, args []string) error {
		var nsName, key, value string
		switch len(args) {
		case 2:
			key = args[0]
			value = args[1]
		case 3:
			nsName = args[0]
			key = args[1]
			value = args[2]
		default:
			return errors.New("set-global [namespace] <param-key> <param-value>")
		}

		apsg, err := action.NewParamSetGlobal(fs, nsName, key, value)
		if err != nil {
			return errors.Wrap(err, "unable to initialize param set global action")
		}

		if err := apsg.Run(); err != nil {
			return errors.Wrap(err, "set global param")
		}

		return nil
	},
}

func init() {
	paramCmd.AddCommand(paramSetGlobalCmd)

	// TODO: support global in env
}
