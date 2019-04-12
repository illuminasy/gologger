package gologger

import (
	"fmt"
	"log"
	"runtime/debug"
	"strings"
)

// Config Go logger config struct
// CustomErrorClass mapping errors to string, used in bugsnag to group by
// Requires manually setting group classes in bugsnag website
type Config struct {
	Bugsnag          bool
	CustomErrorClass map[string]error
}

var loggerConfig Config

// Configure Configure gologger
func Configure(config Config) {
	loggerConfig = config
}

// Info Logs an info
func Info(message interface{}, a ...interface{}) {
	logLevel("INFO", message, a...)
}

// Warn Logs a warning
func Warn(err error, a ...interface{}) {
	logLevel("WARN", err.Error(), a...)
}

// Error Logs an Error
func Error(err error, a ...interface{}) {
	logWithMiddleware("ERROR", err, a...)
}

// Fatal Logs an Fatal
func Fatal(err error, a ...interface{}) {
	logWithMiddleware("FATAL", err, a...)
}

// WarnIfErr Warn only if there is an error
func WarnIfErr(err error) {
	if err != nil {
		Warn(err)
	}
}

// LogIfErr Log only if there is an error
func LogIfErr(err error) {
	if err != nil {
		Error(err)
	}
}

// Wrap wraps error with custom error
// so we can group them by error class in bugsnag
func Wrap(err1 error, err2 error) error {
	if err1 == nil {
		return nil
	}
	return fmt.Errorf("%v: %v", err2, err1)
}

func logWithMiddleware(level string, err error, a ...interface{}) {
	logLevel(level, err.Error(), a...)

	if loggerConfig.Bugsnag {
		sendErrorToBugsnag(err)
	}
}

func logLevel(logLevel string, message interface{}, a ...interface{}) {
	stringMessage := fmt.Sprintf(interfaceToString(message), a...)

	// Add a stack trace if it is an error
	if logLevel == "ERROR" || logLevel == "FATAL" {
		stringMessage += "\nStacktrace:\n" + string(debug.Stack())
	}

	log.Printf("[%s] %s", logLevel, stringMessage)
}

func interfaceToString(message interface{}) string {
	if message == nil {
		return "nil"
	}

	if v, ok := message.(error); ok {
		return v.Error()
	}

	// Fall back to what Sprintf provides
	return fmt.Sprintf("%v", message)
}

func getErrorClass(err error) string {
	for key, value := range loggerConfig.CustomErrorClass {
		if strings.Contains(err.Error(), value.Error()) {
			return key
		}
	}

	return "GenericError"
}
