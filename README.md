# GOLogger [![Build Status](https://travis-ci.org/illuminasy/gologger.svg?branch=master)](https://travis-ci.org/illuminasy/gologger) [![GoDoc](https://godoc.org/github.com/illuminasy/gologger?status.svg)](https://godoc.org/github.com/illuminasy/gologger) [![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/illuminasy/gologger/blob/master/LICENSE.md)

GOLogger
A Simple logging package, prefixes logging level.
Also supports error logging to bugsnag.

Currently supported middlewares:
 
# Usage

Get the library:

    $ go get -v github.com/illuminasy/gologger

Without Bugsnag
```go
func startApp() {
	err := doSomething()
	gologger.LogIfErr(err)
}

```

With Bugsnag
There are 2 ways to use this package with bugsnag
1) With out gorouter package, which has a middleware which does error reporting to external service
checkout https://github.com/illuminasy/gorouter for more info and configuration. No need to reconfigure bugsnag here
2) With out gorouter package, you will need to configure bugsnag here to make it work

1) With out gorouter package
```go
func startApp() {
	errTest = errors.New("test error class")
	lc := gologger.Config{
		Bugsnag: true,
		CustomErrorClass: map[string]error{
			"ErrTest": errTest,
		},
	}
	gologger.Configure(lc)

	// Wrap error with custom error for grouping by error class in Bugsnag
	err := gologger.Wrap(errTest, fmt.Errorf("Error something should not have happened"))
	gologger.Error(err)
}

```

2) With out gorouter package
```go
func startApp() {
	bc := gologger.BugsnagConfig {
		APIKey:              "testing",
		ReleaseStage:        "Dev",
		AppType:             "logger",
		AppVersion:          "0.1.0",
		ProjectPackages:     []string{
			"main",
		},
		NotifyReleaseStages: []string{
			"Dev",
		},
		ParamsFilters:       []string{"password", "secret"},
		PanicHandler:        func() {},
		Hostname:            "localhost",
	}
	gologger.ConfigureBugsnag(bc)

	errTest = errors.New("test error class")
	lc := gologger.Config{
		Bugsnag: true,
		CustomErrorClass: map[string]error{
			"ErrTest": errTest,
		},
	}
	gologger.Configure(lc)

	// Wrap error with custom error for grouping by error class in Bugsnag
	err := gologger.Wrap(errTest, fmt.Errorf("Error something should not have happened"))
	gologger.Error(err)
}

```
