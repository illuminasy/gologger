package gologger

import (
	"fmt"
	"log"
	"strings"

	bugsnag "github.com/bugsnag/bugsnag-go"
	"github.com/illuminasy/gorouter/middleware"
)

var bugsnagConfigured = false

// BugsnagConfig Bugsnag Configuration
type BugsnagConfig struct {
	APIKey              string
	ReleaseStage        string
	AppType             string
	AppVersion          string
	ProjectPackages     []string
	NotifyReleaseStages []string
	ParamsFilters       []string
	PanicHandler        func()
	Hostname            string
}

// ConfigureBugsnag Configure bugsnag
func ConfigureBugsnag(config BugsnagConfig) {
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:              config.APIKey,
		ReleaseStage:        config.ReleaseStage,
		AppType:             config.AppType,
		AppVersion:          config.AppVersion,
		ProjectPackages:     config.ProjectPackages,
		NotifyReleaseStages: config.NotifyReleaseStages,
		ParamsFilters:       []string{"password", "secret"},
		PanicHandler:        func() {},
		Hostname:            config.Hostname,
	})

	bugsnagConfigured = true
}

func sendErrorToBugsnag(err error, a ...interface{}) error {
	errClass := getErrorClass(err)

	// First try with middleware
	err = middleware.ReportErrorToBugsnag(getErrorClass(err), err, a...)

	if err != nil && strings.Contains(err.Error(), "Bugsnag has not been configured") {
		// Try to send using local bugsnag package
		if bugsnagConfigured {
			// append error class so bugsnag can group errors using this
			a = append([]interface{}{bugsnag.ErrorClass{Name: errClass}}, a...)
			return bugsnag.Notify(err, a...)
		}
	} else if err != nil {
		log.Panicln(err)
	}

	if bugsnagConfigured {
		// append error class so bugsnag can group errors using this
		a = append([]interface{}{bugsnag.ErrorClass{Name: errClass}}, a...)
		return bugsnag.Notify(err, a...)
	}

	return fmt.Errorf("Bugsnag has not been configured")
}
