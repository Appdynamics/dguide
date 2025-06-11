// util/output_test.go
// How to run ?
// automation: go test ./util/
// verbose:  go test ./util/ -v
package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLargeShellScriptFile(t *testing.T) {
	sizes := []int{1, 10, 20, 50, 80, 100} // Sizes in MB for testing purposes
	for _, size := range sizes {
		t.Run(fmt.Sprintf("%dMB", size), func(t *testing.T) {
			scriptFile := fmt.Sprintf("../testdata/script_%d.sh", size)
			scriptFilePath, err := filepath.Abs(scriptFile)
			if err != nil {
				t.Fatalf("Failed to resolve path for script file: %s", err)
			}

			if err := os.MkdirAll(filepath.Dir(scriptFilePath), 0755); err != nil {
				t.Fatalf("Failed to create testdata directory: %s", err)
			}

			generateShPath, err := filepath.Abs("../testdata/generate_sh.go")
			if err != nil {
				t.Fatalf("Failed to resolve path to generate_sh.go: %s", err)
			}

			cmd := exec.Command("go", "run", generateShPath, fmt.Sprintf("%d", size), scriptFilePath)

			var stderr bytes.Buffer
			cmd.Stderr = &stderr

			//t.Logf("cmd: %s %s %s\n", cmd.Path, cmd.Args[1], cmd.Args[2])

			if err := cmd.Run(); err != nil {
				t.Fatalf("Failed to generate shell script: %s, stderr: %s", err, stderr.String())
			}

			if _, err := os.Stat(scriptFilePath); os.IsNotExist(err) {
				t.Fatalf("Generated shell script file does not exist: %s", scriptFilePath)
			}

			scriptContent, err := ioutil.ReadFile(scriptFilePath)
			if err != nil {
				t.Fatalf("Failed to read shell script: %s", err)
			}

			// Clean up the generated file
			defer os.Remove(scriptFilePath)

			compressed, err := CompressAndEncode(string(scriptContent))
			if err != nil {
				t.Fatalf("Failed to compress and encode: %s", err)
			}

			decompressed, err := DecodeAndDecompress(compressed)
			if err != nil {
				t.Fatalf("Failed to decode and decompress: %s", err)
			}
			// validate the content
			assert.Equal(t, string(scriptContent), decompressed, "Original and decompressed content should match")
		})
	}
}
