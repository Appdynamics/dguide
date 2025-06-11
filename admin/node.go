package admin

import (
	"runtime"

	"github.com/Appdynamics/dguide/log"
	tools "github.com/Appdynamics/dguide/util"
	"github.com/spf13/cobra"
)

// var(
//
//	output []byte
//	err    error
//
// )
var osType string

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Collect diag information for Nodejs agent",
	Run: func(cmd *cobra.Command, args []string) {
		dlog := log.GetLogger()
		if osType == "" {
			osType = runtime.GOOS
		}

		strings := []string{
			tools.GetOSInfo(osType),
			tools.GetNodeversion(),
			tools.GetpsInfo("py"),
			tools.GetEnvDetails(),
			tools.GetOptions(),
		}
		// output, err := exec.Command("sh", "-c", "cmd").CombinedOutput()
		// if err != nil {
		// 	fmt.Printf("Failed to run command: %s\n", err)
		// 	return
		// }
		//fmt.Printf("%s\n", output)

		for _, str := range strings {
			err := tools.WriteOutput([]byte(str+"\n"), prettyPrint, outputPath, "collect-node", true)
			if err != nil {
				//fmt.Printf("Error writing output: %s\n", err)
				dlog.Error("Error writing output: %s", err)
			}
		}
		err := tools.ZipFile(logPath, enableZip)
		if err != nil {
			//fmt.Printf("Error zipping agent log dir %s \n", err)
			dlog.Error("Error zipping agent log dir %s", err)
		}
		if prettyPrint {
			//fmt.Printf("Console output redirected to %s\n", outputPath)
			dlog.Info("Console output redirected to %s", outputPath)
		}
		//fmt.Printf("\u001B[32mSUCCESS\u001B[0m!\n")
		dlog.Success("Ok!")
	},
}

func init() {
	collectCmd.AddCommand(nodeCmd)
}
