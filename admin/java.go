// admin/java.go
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

// javaCmd represents the java command
var javaCmd = &cobra.Command{
    Use:   "java",
    Short: "Gather environment details for Java agent",
    Run: func(cmd *cobra.Command, args []string) {
		dlog := log.GetLogger()
		zInB64 := "H4sIAAAAAAAA/8xXX0/jOBB/z6cY0ohSrUIFLyeVK7qK8m/FFcSifTn2IhNPGy+u3bOdLohyn/1kOw1pmnaX1Ul3edhFM+Pxb2Z+M+O2droPTHR1FgQyN7PcJGPGsR92Kc672lCZmzAICKVJhoSi2uvASwAAgGkmIey/+wvhGMKocldYcQe17yOZExhMUBj4ZIjRVd3S+HiLv5+Ct8XfOQpUxCAFKXoQ7VFisLP9yL8NoUn7GgRcTsZlbbhMCQeNqWFSJIYZW9HoIKwomUgTzoSVvxz2YqNyfA0BWjDEMcm5ASPBCoGNQUgDMyXnjCINnA82hj8gjJZOQuh74y9HYDIUwbJGHnL87m9TCsaskoloJcBNR/5DqK9B0Prxsg+vYXR9BycXg9H56c73rIOgBScZpo++R+aoNJMiSK0o+UrmJClEFVJMxhA6a9dLHmbrt6+Egc5nM4VaMzEBIp7BRwNjJac2S2CeZwjxDKxjyww3HUTOOSxznMrplAgK8dwZNSX3TOaCehf4hGluyANHYAJuBncXG7JoPxdO3/7ji89dSWMBYfRx8HmQXFz/fhrCF9jdteKnqtjNNnsytPV+D6bSxTowqCFrvK+AqrFGMJf/kTT+5v0t7CnZa+mbFEHUWUvREMZ1P9orbeKi8HB4vHvQKS0L6UZLWAD59gjxGbTDNrS7ha4LLzPFhIHo8LXdqQUTFddviqK0/Oyd9SCMCr/bjiyDfjOF8GD/l4bwa3ldxsM0TKVCMBkR4E5uwcc1lr5+wDFHrb1jqQD/ygm3DfGdS4q5VcwEuGLa+L6dKZmi1qiDCTEZKt+6pbSpeW+WSu9+piEmOIYFTBT6Di3+jOf+f1/Y9k5XZxCn0Nq5tzS97+qs225eJlWExG3fZpxOtx2tX94l5t5G0M5XDxZV6D+D/IYojW5oefhqkk8tBDucUsk5psY39+D8dHSXXI4+3Q2urpLh5W1DaExoQzhPKFP14BpdwLBoiKBk06ZQl2V5GUu1x/oHR+zX/ujsiH340LENsBcx+Bu6f8blgW4Hik5kr21YQJobiGm714Z4fAgLeCJqooEyJcgUYQHfMsYRFBJqhUdA5UqD+QFKmQrXVmG1wa1+e/OsHPPJOZE5p+7ZQNGgmjKBbq1UGHHpMwtDpjA1Uj3vww1HohHcCrMQ32roTjANKheCicl+uHLng0LyWG81KgV6QvgtObMwtG1ibcc7kjTbRIPKDl1jQWL9lFT4P9TW4x9e3jZUmEIYlfrmOvuC2Xa3u7+WpVoF6LJWPai4XXWnIebkdvXed/OnJAVQidoRCZ+YNs33rlX8LBfuYWgHs1RphtrYR7sLR6eKzUyx7O2Lyf+FZU3ffurYx7AUBhRO5Rz9Pl5/Xzlx4/Be09TG5QZ9hW31K5vpGFTexPfiVNANPxGKSIN/AgAA//+32M3n8w0AAA=="
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
		
        err = tools.WriteOutput(output, prettyPrint, outputPath, "java", false)
        if err != nil {
         fmt.Printf("\033[31mERROR\033[0m Error writing output: %s\n", err)
        }
		// caveat for dotnet agent , only enabled for java & dotnet
		cmdName := "java"
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
    runCmd.AddCommand(javaCmd)
}
