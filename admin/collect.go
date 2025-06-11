package admin

import (
	"github.com/spf13/cobra"
)

var collectCmd = &cobra.Command{
	Use:   "collect [command]",
	Short: "Collects various types of informations for different agents",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(collectCmd)
}
