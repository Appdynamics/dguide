package util

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	output []byte
	err    error
)

const (
	hdrRuntime           string = "Runtime     VERSION    "
	HeaderTemplate       string = "============= %s ==============\n"
	hdrOSinfo            string = "Operating System"
	hdrPSinfo            string = " *Process Info  "
	hdrEnvVarInfo        string = "Environment var "
	hdrHttpdPs           string = " *httpd Process "
	hdrFlagOrOptionsInfo string = "Runtime Flag/Opt"
)

type httpdLogCollector struct {
    AccessLogPath string
    ErrorLogPath  string
}

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
	formattedOutput := FormatString(hdrRuntime, string(output))
	return formattedOutput
}

func GetNodeversion() string {
	output, err := exec.Command("sh", "-c", `echo "NodeJs ver:" && node -v && echo "\nNPM ver:" && npm -v`).CombinedOutput()
	if err != nil {
		formattedOutput := FormatString(hdrRuntime, fmt.Sprintf("GetNodeversion()::Failed to get Node version! Error %s\n", err))
		return formattedOutput
	}
	formattedOutput := FormatString(hdrRuntime, string(output))
	return formattedOutput

}

// httpdKnownPaths lists absolute locations to probe when the binary is not on $PATH.
var httpdKnownPaths = []string{
	"/usr/sbin/httpd",
	"/usr/sbin/apache2",
	"/usr/sbin/apache2ctl",
	"/usr/local/sbin/httpd",
	"/usr/local/apache2/bin/httpd",
	"/usr/local/bin/httpd",
	"/usr/bin/httpd",
	"/usr/bin/apache2",
}

// httpdBinary returns the absolute path of the first usable Apache/httpd binary.
// Strategy:
//  1. Resolve via $PATH using LookPath (fast, works when PATH is set).
//  2. Fall back to probing well-known absolute paths (works in minimal/container
//     environments where PATH may not include /usr/sbin).
func httpdBinary() (string, bool) {
	for _, bin := range []string{"httpd", "apache2", "apache2ctl", "httpd2"} {
		if absPath, err := exec.LookPath(bin); err == nil {
			return absPath, true
		}
	}
	for _, absPath := range httpdKnownPaths {
		if _, err := os.Stat(absPath); err == nil {
			return absPath, true
		}
	}
	return "", false
}

func GetWebsrvVersion() string {
	bin, found := httpdBinary()
	if !found {
		return FormatString(hdrRuntime,
			"GetWebsrvVersion()::No httpd/apache2 binary found (tried $PATH and known locations)\n")
	}

	// apache2/apache2ctl use lowercase -v (short version string).
	// httpd/httpd2 use uppercase -V (full build info including compile-time flags).
	// filepath.Base extracts the binary name from the absolute path returned by httpdBinary().
	flag := "-V"
	switch filepath.Base(bin) {
	case "apache2", "apache2ctl":
		flag = "-v"
	}

	// exec.Command requires binary and arguments as separate parameters.
	// Passing "apache2 -v" as a single string would search for a file literally
	// named "apache2 -v" and fail with "executable file not found".
	output, err := exec.Command(bin, flag).CombinedOutput()
	if err != nil {
		return FormatString(hdrRuntime,
			fmt.Sprintf("GetWebsrvVersion()::Failed to run %s %s: %s\n", bin, flag, err))
	}
	return FormatString(hdrRuntime, fmt.Sprintf("[%s %s]\n%s", bin, flag, string(output)))
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
		// Primary query failed — fall back to listing the top 5 running processes.
		// Useful in containers to reveal the entrypoint and co-located services.
		fallback, fbErr := exec.Command("sh", "-c", "ps -eo pid,comm,args | head -6").CombinedOutput()
		if fbErr != nil {
			// ps itself is unavailable; try /proc as a last resort (Linux only).
			fallback, _ = exec.Command("sh", "-c",
				`ls /proc | grep -E '^[0-9]+$' | head -5 | while read p; do cat /proc/$p/cmdline 2>/dev/null | tr '\0' ' '; echo; done`,
			).CombinedOutput()
		}
		result := fmt.Sprintf("Primary query failed (%s)\nFallback — top processes:\n%s", err, string(fallback))
		return FormatString(hdrPSinfo, result)
	}
	return FormatString(hdrPSinfo, string(output))
}

// envShowPlain contains non-APPD variables that are safe and useful to show
// in plain text — host identity and filesystem path context.
var envShowPlain = map[string]struct{}{
	"HOSTNAME": {}, "HOME": {}, "PATH": {},
}

// envSystemVars is the set of well-known OS/shell/container variables that
// carry no diagnostic value for AppDynamics investigations.
var envSystemVars = map[string]struct{}{
	"USER": {}, "LOGNAME": {}, "SHELL": {},
	"TERM": {}, "LANG": {}, "LC_ALL": {}, "LC_CTYPE": {}, "LC_MESSAGES": {},
	"PWD": {}, "OLDPWD": {}, "SHLVL": {}, "_": {},
	"MAIL": {}, "MANPATH": {}, "INFOPATH": {},
	"DISPLAY": {}, "DBUS_SESSION_BUS_ADDRESS": {},
	"GOPATH": {}, "GOROOT": {}, "GOMODCACHE": {}, "GOENV": {},
	"GOCACHE": {}, "GOTOOLCHAIN": {}, "GOLANG_VERSION": {},
	"XDG_RUNTIME_DIR": {}, "XDG_SESSION_TYPE": {}, "XDG_DATA_DIRS": {},
}


// envFallback reads all environment variables via os.Environ().
// Rules applied:
//  1. Known OS/system vars are skipped entirely (no diagnostic value).
//  2. APPD* vars are shown with their actual values (always relevant).
//  3. Every other var is shown with its value redacted — we cannot know
//     whether an org-specific or runtime variable is sensitive, so we
//     surface the name (useful for diagnosis) but hide the value.
func envFallback() string {
	var appdBuf, otherBuf strings.Builder
	appdCount, otherCount := 0, 0

	for _, entry := range os.Environ() {
		idx := strings.IndexByte(entry, '=')
		if idx < 0 {
			continue
		}
		name, value := entry[:idx], entry[idx+1:]

		// Skip known system/OS variables.
		if _, skip := envSystemVars[name]; skip {
			continue
		}
		if strings.HasPrefix(name, "LC_") ||
			strings.HasPrefix(name, "XDG_") ||
			strings.HasPrefix(name, "GNOME_") ||
			strings.HasPrefix(name, "KDE_") ||
			strings.HasPrefix(name, "DBUS_") {
			continue
		}

		if strings.HasPrefix(name, "APPD") {
			// AppDynamics variables: always show actual value.
			appdBuf.WriteString(fmt.Sprintf("  %s=%s\n", name, value))
			appdCount++
		} else if _, plain := envShowPlain[name]; plain {
			// HOSTNAME, HOME, PATH: safe host-context vars, shown as-is.
			otherBuf.WriteString(fmt.Sprintf("  %s=%s\n", name, value))
			otherCount++
		} else {
			// Everything else: show name but redact value — org-specific or
			// runtime vars whose sensitivity we cannot determine.
			otherBuf.WriteString(fmt.Sprintf("  %s=[REDACTED]\n", name))
			otherCount++
		}
	}

	var buf strings.Builder
	if appdCount > 0 {
		buf.WriteString("APPDYNAMICS variables:\n")
		buf.WriteString(appdBuf.String())
		buf.WriteString("\n")
	} else {
		buf.WriteString("(No APPD* variables found)\n\n")
	}
	if otherCount > 0 {
		buf.WriteString("Other variables (values redacted):\n")
		buf.WriteString(otherBuf.String())
	}
	if appdCount == 0 && otherCount == 0 {
		buf.WriteString("(no environment variables found)\n")
	}
	return buf.String()
}

/*
Fetch env details specific to APPDYNAMICS.
Reads os.Environ() directly — no shell subprocess or hardcoded grep pattern.
APPD* vars are shown in full; everything else has its value redacted.
*/
func GetEnvDetails() string {
	return FormatString(hdrEnvVarInfo, envFallback())
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
	// Match both httpd (RHEL/CentOS/Alpine) and apache2 (Debian/Ubuntu) process names.
	output, err := exec.Command("sh", "-c", "ps -ae | grep -E 'httpd|apache2' | grep -v grep").CombinedOutput()
	if err != nil {
		return FormatString(hdrHttpdPs,
			fmt.Sprintf("GethttpdProcess()::Failed to get httpd/apache2 process info: %s\n", err))
	}
	return FormatString(hdrHttpdPs, string(output))
}

// httpdVersionInfo holds details parsed from `httpd -V` output.
type httpdVersionInfo struct {
	ServerType   string // resolved type: "apache", "apache-rhel", "ibm", "weblogic"
	ServerString string // raw value from "Server version:" line, e.g. "Apache/2.4.52 (Ubuntu)"
	HttpdRoot    string // from -D HTTPD_ROOT="..."
	ErrorLogPath string // resolved from -D DEFAULT_ERRORLOG (absolute path)
	RawOutput    string // full httpd -V output
}

// DetectWebServerInfo runs `httpd -V` (falling back to apache2 and httpd2) and
// parses the server type, HTTPD_ROOT, and DEFAULT_ERRORLOG from the output.
func DetectWebServerInfo() httpdVersionInfo {
	info := httpdVersionInfo{ServerType: "apache"} // safe default

	bin, found := httpdBinary()
	if !found {
		return info
	}
	out, cmdErr := exec.Command(bin, "-V").CombinedOutput()
	if cmdErr != nil || len(out) == 0 {
		return info
	}
	info.RawOutput = string(out)

	var rawErrorLog string
	for _, line := range strings.Split(info.RawOutput, "\n") {
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "Server version:"):
			info.ServerString = strings.TrimSpace(strings.TrimPrefix(line, "Server version:"))
			lower := strings.ToLower(info.ServerString)
			switch {
			case strings.Contains(lower, "ibm_http"):
				info.ServerType = "ibm"
			case strings.Contains(lower, "oracle http"):
				info.ServerType = "weblogic"
			case strings.Contains(lower, "red hat"),
				strings.Contains(lower, "centos"),
				strings.Contains(lower, "amazon linux"),
				strings.Contains(lower, "fedora"),
				strings.Contains(lower, "rhel"):
				info.ServerType = "apache-rhel"
			default:
				info.ServerType = "apache"
			}

		case strings.Contains(line, "HTTPD_ROOT="):
			info.HttpdRoot = extractDefineValue(line)

		case strings.Contains(line, "DEFAULT_ERRORLOG="):
			rawErrorLog = extractDefineValue(line)
		}
	}

	// Resolve error log: if relative, join with HTTPD_ROOT
	if rawErrorLog != "" {
		if filepath.IsAbs(rawErrorLog) {
			info.ErrorLogPath = rawErrorLog
		} else if info.HttpdRoot != "" {
			info.ErrorLogPath = filepath.Join(info.HttpdRoot, rawErrorLog)
		}
	}

	return info
}

// extractDefineValue parses the quoted value from httpd -V define lines, e.g.:
//
//	-D HTTPD_ROOT="/etc/httpd"  →  "/etc/httpd"
func extractDefineValue(line string) string {
	start := strings.IndexByte(line, '"')
	end := strings.LastIndexByte(line, '"')
	if start >= 0 && end > start {
		return line[start+1 : end]
	}
	return ""
}

func GetHttpdLogPaths(webServer string) (*httpdLogCollector, error) {
	switch webServer {
	case "apache":
		return &httpdLogCollector{
			AccessLogPath: "/var/log/apache2/access.log",
			ErrorLogPath:  "/var/log/apache2/error.log",
		}, nil
	case "apache-rhel":
		// Red Hat / CentOS / Amazon Linux use httpd path instead of apache2
		return &httpdLogCollector{
			AccessLogPath: "/var/log/httpd/access_log",
			ErrorLogPath:  "/var/log/httpd/error_log",
		}, nil
	case "weblogic":
		// Default WebLogic domain log path; override via WEBLOGIC_DOMAIN_HOME env var
		domainHome := os.Getenv("WEBLOGIC_DOMAIN_HOME")
		if domainHome == "" {
			domainHome = "/u01/app/oracle/user_projects/domains/base_domain/servers/AdminServer"
		}
		return &httpdLogCollector{
			AccessLogPath: filepath.Join(domainHome, "logs", "access.log"),
			ErrorLogPath:  filepath.Join(domainHome, "logs", "AdminServer.log"),
		}, nil
	case "ibm":
		// IBM HTTP Server (IHS) default log path; override via IBM_HTTP_SERVER_HOME env var
		ibmHome := os.Getenv("IBM_HTTP_SERVER_HOME")
		if ibmHome == "" {
			ibmHome = "/opt/IBM/HTTPServer"
		}
		return &httpdLogCollector{
			AccessLogPath: filepath.Join(ibmHome, "logs", "access_log"),
			ErrorLogPath:  filepath.Join(ibmHome, "logs", "error_log"),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported web server type %q (supported: apache, apache-rhel, weblogic, ibm)", webServer)
	}
}

// CollectHttpdLogs auto-detects the web server type from `httpd -V` when
// webServer is "auto", then copies the access and error logs into
// outputPath/httpd-logs/ and returns a formatted status string.
func CollectHttpdLogs(webServer, outputPath string) string {
	const hdr = "Httpd Log Collection"

	if runtime.GOOS != "linux" {
		return FormatString(hdr, fmt.Sprintf(
			"Skipped: log collection only supported on Linux (current OS: %s)\n", runtime.GOOS))
	}

	// Auto-detect from httpd -V when not explicitly overridden
	var detected httpdVersionInfo
	if webServer == "auto" || webServer == "" {
		detected = DetectWebServerInfo()
		webServer = detected.ServerType
	}

	logCollector, err := GetHttpdLogPaths(webServer)
	if err != nil {
		return FormatString(hdr, fmt.Sprintf("Failed to resolve log paths: %s\n", err))
	}

	// Prefer the error log path extracted directly from httpd -V when available
	// (it reflects non-default HTTPD_ROOT installations)
	if detected.ErrorLogPath != "" {
		logCollector.ErrorLogPath = detected.ErrorLogPath
	}

	destDir := filepath.Join(outputPath, "httpd-logs")
	if mkErr := os.MkdirAll(destDir, 0755); mkErr != nil {
		return FormatString(hdr, fmt.Sprintf("Failed to create destination directory %s: %s\n", destDir, mkErr))
	}

	serverLabel := webServer
	if detected.ServerString != "" {
		serverLabel = fmt.Sprintf("%s (%s)", webServer, detected.ServerString)
	}
	result := fmt.Sprintf("Web server: %s  Platform: %s\n", serverLabel, runtime.GOOS)

	for _, src := range []string{logCollector.AccessLogPath, logCollector.ErrorLogPath} {
		dest := filepath.Join(destDir, filepath.Base(src))
		if cpErr := copyLogFile(src, dest); cpErr != nil {
			result += fmt.Sprintf("WARN: could not copy %s: %s\n", src, cpErr)
		} else {
			result += fmt.Sprintf("Copied: %s -> %s\n", src, dest)
		}
	}

	return FormatString(hdr, result)
}

func copyLogFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
