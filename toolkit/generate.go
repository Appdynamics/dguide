// generate.go
package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"os"
	"text/template"

	"gopkg.in/yaml.v3"
)

type Command struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Execute     string `yaml:"execute"`
	Script      string `yaml:"script"`
	CommandType string `yaml:"commandType"`
}

type Commands struct {
	Commands []Command `yaml:"commands"`
}

func compressAndEncode(data string) (string, error) {
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	_, err := gzipWriter.Write([]byte(data))
	if err != nil {
		return "", err
	}
	gzipWriter.Close()
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

const commandTemplate = `// admin/{{ .Name }}.go
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

// {{ .Name }}Cmd represents the {{ .Name }} command
var {{ .Name }}Cmd = &cobra.Command{
    Use:   "{{ .Name }}",
    Short: "{{ .Description }}",
    Run: func(cmd *cobra.Command, args []string) {
		dlog := log.GetLogger()
		zInB64 := "{{ .Text }}"
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
		switch "{{ .CommandType }}" {
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
		
        err = tools.WriteOutput(output, prettyPrint, outputPath, "{{ .Name }}", false)
        if err != nil {
         fmt.Printf("\033[31mERROR\033[0m Error writing output: %s\n", err)
        }
		// caveat for dotnet agent , only enabled for java & dotnet
		cmdName := "{{ .Name }}"
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
    runCmd.AddCommand({{ .Name }}Cmd)
}
`

func main() {
	yamlFile, err := os.ReadFile("directives.yaml")
	if err != nil {
		panic(err)
	}

	var cmds Commands
	err = yaml.Unmarshal(yamlFile, &cmds)
	if err != nil {
		panic(err)
	}

	for _, c := range cmds.Commands {
		var text string
		if c.Script != "" {
			scriptContent, err := os.ReadFile(c.Script)
			if err != nil {
				panic(err)
			}
			text, err = compressAndEncode(string(scriptContent))
			if err != nil {
				panic(err)
			}
		} else if c.Execute != "" {
			text, err = compressAndEncode(c.Execute)
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Printf("No script or execute field provided for command: %s\n", c.Name)
			continue
		}

		f, err := os.Create("admin/" + c.Name + ".go")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		tmpl, err := template.New("command").Parse(commandTemplate)
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(f, struct {
			Name        string
			Description string
			Text        string
			CommandType string
		}{
			Name:        c.Name,
			Description: c.Description,
			Text:        text,
			CommandType: c.CommandType,
		})
		if err != nil {
			panic(err)
		}
	}
}
