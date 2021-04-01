package config

import "github.com/spf13/cobra"

// NewConfigCmd creates config command
func NewConfigCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "config",
		Short: "Manage the config item related to KubeSphere Jenkins",
	}

	cmd.AddCommand(newUpdateCmd())
	return
}
