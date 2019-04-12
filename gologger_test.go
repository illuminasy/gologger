package gologger

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	errAPIConfig = errors.New("Invalid API config")
	errTest      = errors.New("testing error class")

	customErrors = map[string]error{
		"ErrAPIConfig": errAPIConfig,
		"ErrTest":      errTest,
	}
)

func TestInfo(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	Info(fmt.Errorf("something happened"))
	assert.Contains(t, buf.String(), "[INFO] something happened")
}

func TestWarn(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	Warn(fmt.Errorf("Error handling something"))
	assert.Contains(t, buf.String(), "[WARN] Error handling something")
}

func TestError(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	Error(fmt.Errorf("Error something should not have happened"))
	assert.Contains(t, buf.String(), "[ERROR] Error something should not have happened")
}

func ExampleWrap_output() {
	errTest = errors.New("test error class")
	lc := Config{
		CustomErrorClass: map[string]error{
			"ErrTest": errTest,
		},
	}
	Configure(lc)

	err := Wrap(errTest, fmt.Errorf("Error something should not have happened"))
	fmt.Println(err.Error())
	// Output: Error something should not have happened: test error class
}

func TestErrorWrap(t *testing.T) {
	lc := Config{
		CustomErrorClass: customErrors,
	}
	Configure(lc)

	// Wrap was intended for bugsnag, can also be used without bugsnag wont cause any issue.
	err := Wrap(errAPIConfig, fmt.Errorf("Error something should not have happened"))
	assert.Contains(t, err.Error(), "Error something should not have happened: Invalid API config")
}

func TestGetErrorClass(t *testing.T) {
	lc := Config{
		CustomErrorClass: customErrors,
	}
	Configure(lc)
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ErrGenericError",
			args: args{
				err: errors.New("generic error"),
			},
			want: "GenericError",
		},
		{
			name: "ErrAPIConfig",
			args: args{
				err: errAPIConfig,
			},
			want: "ErrAPIConfig",
		},
		{
			name: "ErrTest",
			args: args{
				err: errTest,
			},
			want: "ErrTest",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getErrorClass(tt.args.err))
		})
	}
}
