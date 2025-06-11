package util

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Appdynamics/dguide/log"
)

func WriteOutput(output []byte, prettyPrint bool, outputPath string, commandName string, appendToFile bool) error {
	if prettyPrint {
		outputFile := filepath.Join(outputPath, fmt.Sprintf("%s.log", commandName))
		err := os.MkdirAll(filepath.Dir(outputFile), 0755)
		if err != nil {
			return fmt.Errorf("\033[31mERROR\033[0m failed to create directory: %w", err)
		}
		if appendToFile {
			f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return fmt.Errorf("\033[31mERROR\033[0m failed to open file: %w", err)
			}
			defer f.Close()
			if _, err := f.Write(output); err != nil {
				return fmt.Errorf("\033[31mERROR\033[0m failed to write to file: %w", err)
			}
		} else {
			err = ioutil.WriteFile(outputFile, output, 0644)
			if err != nil {
				return fmt.Errorf("\033[31mERROR\033[0m failed to write to file: %w", err)
			}
		}
		//fmt.Printf("\u001B[32mSUCCESS\u001B[0m!\n")
		//fmt.Printf("Console output redirected to %s\n", outputFile)
	} else {
		fmt.Printf("%s\n", output)
	}
	return nil
}

func ZipFile(logPath string, enableZip bool) error {
	if enableZip {
		err := ZipAndMove(logPath, "/tmp/dguide")
		if err != nil {
			return fmt.Errorf("\033[31mERROR\033[0m Failed to zip and move logs: %s", err)
		}
		//fmt.Printf("\u001B[32mSUCCESS\u001B[0m!\n")
		//fmt.Printf("Logs zipped and moved to /tmp/dguide successfully.\n")
		log.GetLogger().Info("Logs zipped and moved to /tmp/dguide successfully.")
	} else {
		//fmt.Printf("Zipping not enabled; logs not moved.\n")
	}
	return nil
}

func ZipAndMove(srcDir, destDir string) error {
	zipFilePath := filepath.Join(destDir, "agentlogs.zip")

	// Create the destination directory if it doesn't exist
	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Create a zip file
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	// Create a new zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk the source directory and add files to the zip
	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Create a zip file entry
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		zipEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		// Open the file
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Copy the file content to the zip entry
		_, err = io.Copy(zipEntry, file)
		return err
	})

	if err != nil {
		return fmt.Errorf("failed to zip directory: %w", err)
	}

	return nil
}

func DecodeAndDecompress(encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer reader.Close()

	decompressed, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(decompressed), nil
}

func CompressAndEncode(data string) (string, error) {
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	_, err := gzipWriter.Write([]byte(data))
	if err != nil {
		return "", err
	}
	gzipWriter.Close()
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
