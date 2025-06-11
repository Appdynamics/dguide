package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// Log levels
const (
	INFO    = "INFO"
	ERROR   = "ERROR"
	WARN    = "WARN"
	SUCCESS = "SUCCESS"
)

// Colors for log levels
var Colors = map[string]string{
	INFO:    "\033[37m", // White
	ERROR:   "\033[31m", // Red
	WARN:    "\033[33m", // Yellow
	SUCCESS: "\033[32m", // Green
	"CLR":   "\033[0m",  // Reset
}

// Global logger instance
var (
	instance *Logger
	once     sync.Once
)

// Logger represents a custom logger
type Logger struct {
	logger *log.Logger
}

// Logger with output file (redirect to a file)
func GetLoggerOutput(output io.Writer) *Logger {
	once.Do(func() {
		instance = &Logger{
			logger: log.New(output, "", 0),
		}
	})
	return instance
}

func GetLogger() *Logger {
	once.Do(func() {
		instance = &Logger{
			logger: log.New(os.Stdout, "", 0),
		}
	})
	return instance
}

// Log logs a message with the given level and message
func (l *Logger) Log(level, message string) {
	color, exists := Colors[level]
	if !exists {
		color = Colors[INFO]
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	if !strings.HasSuffix(message, "\n") {
		message += "\n"
	}
	logMessage := fmt.Sprintf("%s %s[%s]%s %s", timestamp, color, level, Colors["CLR"], message)
	l.logger.Print(logMessage)
}

// Logf formats its arguments according to the format, analogous to Printf
func (l *Logger) Logf(level, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.Log(level, message)
}

// Info logs an info message with formatting
func (l *Logger) Info(format string, args ...interface{}) {
	l.Logf(INFO, format, args...)
}

// Error logs an error message with formatting
func (l *Logger) Error(format string, args ...interface{}) {
	l.Logf(ERROR, format, args...)
}

// Warn logs a warning message with formatting
func (l *Logger) Warn(format string, args ...interface{}) {
	l.Logf(WARN, format, args...)
}

// Success logs a success message with formatting
func (l *Logger) Success(format string, args ...interface{}) {
	l.Logf(SUCCESS, format, args...)
}
