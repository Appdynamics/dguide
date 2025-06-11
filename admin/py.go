package admin

import (
	"fmt"
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

var pyCmd = &cobra.Command{
	Use:   "py",
	Short: "Collect diag information for python agent",
	Run: func(cmd *cobra.Command, args []string) {
		if osType == "" {
			osType = runtime.GOOS
		}
		dlog := log.GetLogger()
		strings := []string{
			tools.GetOSInfo(osType),
			tools.GetPyVersion(),
			tools.GetpsInfo("py"),
			tools.GetEnvDetails(),
		}
		// output, err := exec.Command("sh", "-c", "cmd").CombinedOutput()
		// if err != nil {
		// 	fmt.Printf("Failed to run command: %s\n", err)
		// 	return
		// }
		//fmt.Printf("%s\n", output)

		for _, str := range strings {
			err := tools.WriteOutput([]byte(str+"\n"), prettyPrint, outputPath, "collect-py", true)
			if err != nil {
				fmt.Printf("Error writing output: %s\n", err)
			}
		}
		err := tools.ZipFile(logPath, enableZip)
		if err != nil {
			fmt.Printf("Error zipping agent log dir %s \n", err)
		}
		if prettyPrint {
			dlog.Info("Console output redirected to %s", outputPath)
		}
		dlog.Success("Ok!")
	},
}

func init() {
	collectCmd.AddCommand(pyCmd)
}
