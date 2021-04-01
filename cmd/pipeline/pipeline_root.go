package pipeline

import "github.com/spf13/cobra"

// NewPipelineRootCommand returns the root command of pipeline
func NewPipelineRootCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:     "pipeline",
		Aliases: []string{"pip"},
	}

	cmd.AddCommand(newBackupCommand(), newRestoreCommand())
	return
}
