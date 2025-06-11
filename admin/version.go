package admin

import (
	"fmt"

	"github.com/spf13/cobra"
)

//	type VersionInfo struct {
//	    Version string `json:"version"`
//	    Build   int    `json:"build"`
//	    BuildDate    string `json:"date"`
//	}
//
// version represents current version and build details
var (
	Version   string
	Build     string
	BuildDate string
)

// type VersionInfo struct {
// 	Version   string `json:"version"`
// 	Build     string `json:"build"`
// 	BuildDate string `json:"date"`
// }

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show dguide version details",
	Run: func(cmd *cobra.Command, args []string) {
		//versionInfo, err := readVersionInfo()
		// if err != nil {
		// 	fmt.Printf("Failed to read version info: %s\n", err)
		// 	return
		// }
		// if Version == "" {
		// 	Version = versionInfo.Version
		// }
		// if Build == "" {
		// 	Build = versionInfo.Build
		// }
		// if BuildDate == "" {
		// 	BuildDate = versionInfo.BuildDate
		// }
		fmt.Printf("dguide version: %s\n", Version)
		fmt.Printf("Build number: %s\n", Build)
		fmt.Printf("Build date: %s\n", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// func readVersionInfo() (*VersionInfo, error) {
// 	file, err := os.Open("version.json")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	var versionInfo VersionInfo
// 	decoder := json.NewDecoder(file)
// 	err = decoder.Decode(&versionInfo)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &versionInfo, nil
// }
