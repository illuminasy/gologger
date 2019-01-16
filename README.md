# GOLogger [![Build Status](https://travis-ci.org/Illuminasy/gologger.svg?branch=master)](https://travis-ci.org/Illuminasy/gologger) [![GoDoc](https://godoc.org/github.com/Illuminasy/gologger?status.svg)](https://godoc.org/github.com/Illuminasy/gologger) [![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Illuminasy/gologger/blob/master/LICENSE.md)

GOLogger
A Simple logging package, prefixes logging level.
Also supports error logging to bugsnag.

Currently supported middlewares:
 
# Usage

Get the library:

    $ go get -v github.com/Illuminasy/gologger

Without Bugsnag
```go
func startApp() {
	err := doSomething()
	gologger.LogIfErr(err)
}

```

With Bugsnag
```go
func startApp() {
	errTest = errors.New("test error class")
	lc := gologger.Config{
		Bugsnag: false,
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
