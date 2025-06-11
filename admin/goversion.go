package admin

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// goversionCmd represents the goversion command
var goversionCmd = &cobra.Command{
	Use:   "goversion",
	Short: "Display Go runtime version",
	Run: func(cmd *cobra.Command, args []string) {
		output, err := exec.Command("sh", "-c", "go version").CombinedOutput()
		if err != nil {
			fmt.Printf("Failed to run command: %s\n", err)
			return
		}
		fmt.Printf("%s\n", output)
	},
}

func init() {
	//rootCmd.AddCommand(goversionCmd)
}
