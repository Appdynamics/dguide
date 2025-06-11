package util

import (
	"fmt"
	"os/exec"
)

type agents string

const (
	node   agents = "Node.Js"
	php    agents = "PHP"
	py     agents = "Python"
	websrv agents = "Webserver"
	sdk    agents = "CPP~GO~IIB~ABAP"
	java   agents = "Java"
	Dnet   agents = "DotNet"
)

var (
	output []byte
	err    error
)

var currAgent agents

var hdrRuntime string = fmt.Sprintf("%s     VERSION    ", currAgent)

const (
	HeaderTemplate       string = "============= %s ==============\n"
	hdrOSinfo            string = "Operating System"
	hdrPSinfo            string = " *Process Info  "
	hdrEnvVarInfo        string = "Environment var "
	hdrHttpdPs           string = " *httpd Process "
	hdrFlagOrOptionsInfo string = "Runtime Flag/Opt"
)

func FormatString(header, content string) string {
	return fmt.Sprintf(HeaderTemplate, header) + content + "\n"
}

func GetOSInfo(osType string) string {

	switch osType {
	case "windows":
		output, err = exec.Command("cmd", "/C", "ver").CombinedOutput()
	case "linux":
		output, err = exec.Command("sh", "-c", "cat /etc/*release").CombinedOutput()
	case "darwin": // macOS is reported as "darwin" by Go
		output, err = exec.Command("sw_vers").CombinedOutput()
	default:
		formattedOutput := FormatString(hdrOSinfo, "GetOSInfo(osType string) :: Unsupported operating system.\n")
		return formattedOutput
	}

	if err != nil {
		formattedOutput := FormatString(hdrOSinfo, fmt.Sprintf("Failed to get OS info! Error: %s\n", err))
		return formattedOutput
	}
	formattedOutput := FormatString(hdrOSinfo, string(output))
	return formattedOutput

}

func GetPyVersion() string {

	var pythonCmd string
	var err error

	// Check for various python commands
	if _, err = exec.LookPath("python"); err == nil {
		pythonCmd = "python"
	} else if _, err = exec.LookPath("python3"); err == nil {
		pythonCmd = "python3"
	} else if _, err = exec.LookPath("python2"); err == nil {
		pythonCmd = "python2"
	} else {
		// No error , just gracefull exits.
		formattedOutput := FormatString(hdrRuntime, string("Python is not available on this system."))
		return formattedOutput
	}
	cmd := exec.Command("sh", "-c", fmt.Sprintf(`
        echo "Found!" && %s --version &&
        echo "\n[ site information]" &&
        %s -m site &&
		echo "\n[pip version]" && %s -m pip --version &&
		echo "\n" && %s -m pip list
    `, pythonCmd, pythonCmd, pythonCmd, pythonCmd))

	output, err = cmd.CombinedOutput()
	if err != nil {
		formattedOutput := FormatString(hdrRuntime, fmt.Sprintf("GetPyVersion()::Failed to get Python version! Error %s\n", err))
		return formattedOutput
	}
	currAgent = py

	formattedOutput := FormatString(hdrRuntime, string(output))
	return formattedOutput
}

func GetNodeversion() string {
	output, err := exec.Command("sh", "-c", `echo "NodeJs ver:" && node -v && echo "\nNPM ver:" && npm -v`).CombinedOutput()
	if err != nil {
		formattedOutput := FormatString(hdrRuntime, fmt.Sprintf("GetNodeversion()::Failed to get Node version! Error %s\n", err))
		return formattedOutput
	}
	currAgent = node
	formattedOutput := FormatString(hdrRuntime, string(output))
	return formattedOutput

}

func GetWebsrvVersion() string {
	output, err := exec.Command("sh", "-c", `echo "HTTPD version:" && httpd -V`).CombinedOutput()
	if err != nil {
		formattedOutput := FormatString(hdrRuntime, fmt.Sprintf("GetWebsrvVersion()::Failed to get Node version! Error %s\n", err))
		return formattedOutput
	}
	currAgent = websrv
	formattedOutput := FormatString(hdrRuntime, string(output))
	return formattedOutput
}
func GetpsInfo(agentType string) string {
	var cmd string = ""
	switch agentType {
	case "websrv":
		cmd = "ps -aef | grep proxy"
	case "py":
		cmd = "ps -aef | grep proxy"
	case "php":
		cmd = "ps -aef | grep proxy"
	case "node":
		cmd = "ps -aef | grep node"

	default:
		cmd = "ps -aef"
	}
	cmd += "| grep -v grep"
	output, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		formattedOutput := FormatString(hdrPSinfo, fmt.Sprintf("GetpsInfo()::Failed to getpsInfo! Error %s\n", err))
		return formattedOutput
	}
	formattedOutput := FormatString(hdrPSinfo, string(output))
	return formattedOutput

}

/*
Fetch env details specific to APPDYNAMICS
*/
func GetEnvDetails() string {
	output, err := exec.Command("sh", "-c", "env | grep APPD").CombinedOutput()
	if err != nil {
		formattedOutput := FormatString(hdrEnvVarInfo, fmt.Sprintf("GetEnvDetails()::Failed to get env variables! Error %s\n", err))
		return formattedOutput
	}
	formattedOutput := FormatString(hdrEnvVarInfo, string(output))
	return formattedOutput
}

/*
Runtime Flag Options for Various Programming Languages
*/
func GetOptions() string {
	output, err := exec.Command("sh", "-c", "env | grep NODE_OPTIONS").CombinedOutput()
	if err != nil {
		formattedOutput := FormatString(hdrFlagOrOptionsInfo, fmt.Sprintf("GetOptions()::Failed to get runtime OPTIONS variables! Error %s\n", err))
		return formattedOutput
	}
	formattedOutput := FormatString(hdrFlagOrOptionsInfo, string(output))
	return formattedOutput
}

func GethttpdProcess() string {
	output, err := exec.Command("sh", "-c", "ps -ae | grep httpd |  grep -v 'sh -c ps -ae | grep httpd'").CombinedOutput()
	if err != nil {
		formattedOutput := FormatString(hdrHttpdPs, fmt.Sprintf("GethttpdProcess()::Failed to get httpd info! Error %s\n", err))
		return formattedOutput
	}
	formattedOutput := FormatString(hdrHttpdPs, string(output))
	return formattedOutput
}
