// admin/dotnet.go
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

// dotnetCmd represents the dotnet command
var dotnetCmd = &cobra.Command{
    Use:   "dotnet",
    Short: "Gather environment details for DotNet agent",
    Run: func(cmd *cobra.Command, args []string) {
		dlog := log.GetLogger()
		zInB64 := "H4sIAAAAAAAA/7xY3W7bOhK+Fp9iKsuw3a7i9KYXaV1sUWe7BdIkWBdb7DaFQUsjm1uKVEnKbZrkAfbyPMF5xfMIByQlW7KVn3NOUSIXMjl/33BmOJPeo/GCibFeEXI2m0TDUtAcIdYjIktTlGaeMY6TcJzieqxNKksTEi1LleA8ZWoSjk1ejGlRpCEhPXgrmGGUs+8IFBJZCgOZVICcLdmCIxRKJqg1alJvzR3V5HC7saGZhCFJSqVQmHnB0kkUAfTgDRowK4TqpBYJb6fELpqm8xXSFNVwBFck6MFrjlQ5Fo8ILCJgGTAD+I1po0nwEsKogTckAACYrCSEkz+8Qnh5hzzYWQenaODV0kKZGWp086wmvkvej7bvDQpU1GAKUhxBNEypQXgS9v8T9/O4n0L/n0f9d0f9GfT/C/3v4ehuaT/auq7TG0ISWphS4VxfaoP5HMV6c/v+BPxJzHGNHFCsmZIit05fU8XogqMmwe0qqqOZEwInTshxQ8i/ayFH3dxifbvQ+P51C+oNbKnnKRrKuPaw7wXy2C440zD1bOA2OhkSamCMJhlLHSvkSDXeLrf3oHWHYfchzRTN8atUn/8U4Kk0Ag38qxSG5fgA8KlniGPOtImV59M/H38Nvyp2XeA9wvOqGt4PrQcnUhZgVkqWyxVQzhulVJPA1u2CpcAERMNCQ4wQS7szGT2HVJIg8H89OP6W8DLFrqr8N6DiEhZUr0Aq0KvtA2BP0iort7swPH87hRcTeHp4OCJBwDL4CGFUsDSEGL9AGDUehBA+PbdKBQmCIJHCMFEiCYKMEeIMe73C5LOt9GOrYGzFjKvU31R+ryLGJk0i87wh24pnGbhE2CG6hqXCAuJjiL/AwOK81qsBdJpl7Qq8tIZp1mm13y0xZULDwMfdwDqNS5piCu9YoqSWmYHpyYmuTLrTpi9bMdfXu5Q5LXSTkrNFrhOpOFscaHlxzdkikQoTrtzPjfbBFlrg3+LaeNs3uO06Ru3GJBp2mDhydJWMV0WBIt1rEYCJTIKRzkF1jXYsXY1CtL8ZNe1wcXUE1oLRhQjbcnwHEg2H7R14Ak9HI7K19EHh1BlQNVUjXmv33fIUOQPFer6mSk+ioVEwuDgcwOBCDOBFh+hRBUnjDm94Krt1AF1Txu2nd4fNmq1ZVdTlMi1re6rv/SutYol+/QyD6Bk8mtgydlUoJgxEz24GANegpTIQl6NaxwfFDNrrbfRk7uiqst/VtLqandIc7e01bjRs0dWX29pt/WiXx3ZwdL/kvnZCS0pU+/V2RR1aT6oU9g60qIsuOyqNbYWV11ubD+gYHrzuA3LT8XaQZqhtI612rpBmJ7gyBuALc5BKgX/9tbRPoiwu51wu3Y7etHubHK3nknGVm+69sUzAjC+1KIyuK0z10GeSp6hIlcMphNF21mk+N0kBsWqdjsPHEEZX5x+mNxayB/haFgzTrbJMyRwaTFa75zmwHYB3qWc9kUtIJXpvOgSWJGMWukutuS7znKrLGvm5y7f9ScvhNtJQXg1k1M9PKNIH903HtdTzjdTbG4s6U/Yq8h3U7519+2qOIGoX5S4ZZKda2JdD5dQwKSCRnGNipxnbx7TZbuzE+o9SJI7S1iKVrFAbO/04F+lEscJ6H5PS0hD/hc7lbijZTpvQg6kUBhTmco3E9s27I4nj2G/Zt6R7ze32aKfxI0Hv7/+jDI6g1AhLLheUQ8bpEmI+jGMulwU1q5EFhcKmoY986sZMLpeQIbVyD0jQaycSCVrR9YAQ+bkD4PHp9GcovCE9mJVFIZXRcMJE+Q2k4JekakjPZiFMIHQHm7oAJKgChLhU3tr0Xl0ysbS3oUoBUgCFUsTay8e09UB/jM5mnw7gQvz26y//n21IzmbuPylVnfIXySrTaupzP51ZHY3wlQK2mgpOjc2ORyHJ2O8BAAD//2aUPTf8EQAA"
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
		
        err = tools.WriteOutput(output, prettyPrint, outputPath, "dotnet", false)
        if err != nil {
         fmt.Printf("\033[31mERROR\033[0m Error writing output: %s\n", err)
        }
		// caveat for dotnet agent , only enabled for java & dotnet
		cmdName := "dotnet"
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
    runCmd.AddCommand(dotnetCmd)
}
