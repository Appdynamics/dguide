package admin

import (
	"fmt"
	"runtime"

	"github.com/Appdynamics/dguide/log"
	tools "github.com/Appdynamics/dguide/util"
	"github.com/spf13/cobra"
)

var webSrvType string

var websrvCmd = &cobra.Command{
	Use:   "websrv",
	Short: "Collect diag information for webserver agent (auto-detects type from httpd -V)",
	Run: func(cmd *cobra.Command, args []string) {
		if osType == "" {
			osType = runtime.GOOS
		}
		dlog := log.GetLogger()
		results := []string{
			tools.GetOSInfo(osType),
			tools.GetWebsrvVersion(),
			tools.GetpsInfo("websrv"),
			tools.GethttpdProcess(),
			tools.GetEnvDetails(),
			tools.CollectHttpdLogs(webSrvType, outputPath),
		}

		for _, str := range results {
			err := tools.WriteOutput([]byte(str+"\n"), prettyPrint, outputPath, "collect-websrv", true)
			if err != nil {
				fmt.Printf("Error writing output: %s\n", err)
			}
		}
		if enableZip && !cmd.Root().PersistentFlags().Changed("log-path") {
			dlog.Error("ERROR: -z requires -l <log-dir> for websrv. " +
				"Apache agent logs are typically at /opt/appdynamics-sdk-native/logs. " +
				"Re-run with: dguide collect websrv -z -l /opt/appdynamics-sdk-native/logs")
			return
		}
		err := tools.ZipFile(logPath, enableZip, outputPath)
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
	websrvCmd.Flags().StringVarP(&webSrvType, "server-type", "s", "auto",
		"Web server type for log collection (auto, apache, apache-rhel, weblogic, ibm); 'auto' detects from httpd -V")
	collectCmd.AddCommand(websrvCmd)
}
