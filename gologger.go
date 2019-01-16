package gologger

import (
	"fmt"
	"log"
	"runtime/debug"
	"strings"

	bugsnag "github.com/bugsnag/bugsnag-go"
)

// Config Go logger config struct
// Bugsnag true/false - turn on or off logging to bugsnag
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
	logLevel("ERROR", err.Error(), a...)
	// append error class so bugsnag can group errors using this

	if loggerConfig.Bugsnag {
		fmt.Println(getErrorClass(err))
		a = append([]interface{}{bugsnag.ErrorClass{getErrorClass(err)}}, a...)
		err = bugsnag.Notify(err, a...)

		if err != nil {
			log.Panicln(err)
		}
	}
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

func logLevel(logLevel string, message interface{}, a ...interface{}) {
	stringMessage := fmt.Sprintf(interfaceToString(message), a...)

	// Add a stack trace if it is an error
	if logLevel == "ERROR" {
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
