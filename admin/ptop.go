// admin/ptop.go
package admin

import (
	roller "github.com/Appdynamics/dguide/roller"
	g "github.com/Appdynamics/dguide/globals"
	tools "github.com/Appdynamics/dguide/util"
	"github.com/Appdynamics/dguide/log"
	"context"
	"fmt"
	"os/exec"
	"time"
	"runtime"

	"github.com/spf13/cobra"
)

// ptopCmd represents the ptop command
var ptopCmd = &cobra.Command{
    Use:   "ptop",
    Short: "Fetch top 10 CPU intensive processes",
    Run: func(cmd *cobra.Command, args []string) {
		dlog := log.GetLogger()
		zInB64 := "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"
        script, err := tools.DecodeAndDecompress(zInB64)
        if err != nil {
            fmt.Printf("ERROR: Failed to decode and decompress script: %s\n", err)
            return
        }
		ctx, cancel := context.WithTimeout(context.Background(), g.DEFAULT_COMMAND_WAITING_SECONDS * time.Second)
        defer cancel()
        var output []byte
		
		r := roller.StartNew(g.DEFAULT_WAITING_MESSAGE)
		time.Sleep(g.DEFAULT_WAITING_SECONDS * time.Second)
		//--begin
		switch "" {
        case "curl":
            //output, err = exec.Command("curl", "-v", script).CombinedOutput()
			output, err = exec.CommandContext(ctx, "curl", "-v", script).CombinedOutput()

        default:
            //output, err = exec.Command("sh", "-c", script).CombinedOutput()
			output, err = exec.CommandContext(ctx, "sh", "-c", script).CombinedOutput()
        }
		if ctx.Err() == context.DeadlineExceeded {
			r.Stop()
            fmt.Printf("\033[31mERROR\033[0m The operation has timed out. Try again!.")
            return
        }	
		r.Stop()
		//--end
        if err != nil {
            fmt.Printf("\033[33mWARNING\033[0m Failed to run command: %s\n", err)
            //return
        }
		// @jay : Fix this for cascade commands with splitter && , handle each of them gracefully !
		
        err = tools.WriteOutput(output, prettyPrint, outputPath, "ptop", false)
        if err != nil {
         fmt.Printf("\033[31mERROR\033[0m Error writing output: %s\n", err)
        }
		// caveat for dotnet agent , only enabled for java & dotnet
		cmdName := "ptop"
		isZipEnabled := false
		if (cmdName == "java"){ isZipEnabled = true } 
		if (cmdName == "dotnet"){ 
			if (runtime.GOOS == "linux" ){isZipEnabled = true}
		}

		if (isZipEnabled){
			err = tools.ZipFile(logPath, enableZip)
			if err != nil {
				fmt.Printf("Error zipping agent log dir %s \n", err)
		    }
		}
	
		if (prettyPrint){
			dlog.Info("Output written to %s successfully!", outputPath)
		}
		dlog.Success("Ok!")
    },
}
func init() {
    runCmd.AddCommand(ptopCmd)
}
