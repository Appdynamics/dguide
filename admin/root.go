package admin

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	prettyPrint bool
	outputPath  string
	logPath     string
	enableZip   bool
)

var rootCmd = &cobra.Command{
	Use:   "dguide",
	Short: "Diagnose agent issues as quick and simple as possible",
	Long: `dguide is an agent automation tool.
	Its goal is to simplify the process of collecting diagnostic information for language Agents. 

		dguide collect node --write-to-file (# Writes output to default path /tmp/dguide)
		dguide collect node --write-to-file --output-path <custom-path> (# Writes output to <custom-path>)
		dguide collect node -z -w (# Writes the output to /tmp/dguide and zip the /tmp/appd)`,
	SuggestionsMinimumDistance: 3, // auto-complete
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	// Persistent Flags, global for application.
	rootCmd.PersistentFlags().BoolVarP(&prettyPrint, "write-to-file", "w", false, "Write output to a file, default is stdout")
	rootCmd.PersistentFlags().StringVarP(&outputPath, "output-path", "o", "/tmp/dguide", "Path to output directory, default is /tmp/dguide")
	rootCmd.PersistentFlags().StringVarP(&logPath, "log-path", "l", "/tmp/appd", "Agent log directory dguide uses to zip, default is /tmp/appd")
	rootCmd.PersistentFlags().BoolVarP(&enableZip, "zip", "z", false, "Enable zipping of agent logs")
	// rootCmd.Flags().BoolP("toggle", "t", false, "Some flag to toggle")
}
