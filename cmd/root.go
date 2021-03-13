package cmd

import (
	"github.com/linuxsuren/jcli-ks-plugin/cmd/config"
	"github.com/spf13/cobra"
)

// NewKSPlugin returns the command of jcli ks
func NewKSPlugin() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "jcli ks",
		Short: "jcli plugin for KubeSphere",
	}

	cmd.AddCommand(config.NewConfigCmd())
	return
}
