package admin

import (
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs a subcommand under dguide",
	Long: `Run custom commands to gather troubleshooting details
	e.g.
		dguide run <mycustom command>
	`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
