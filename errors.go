package instrumenter

import (
	"fmt"
	"os"
	"sync/atomic"

	_ "golang.org/x/tools/refactor/eg"
)

var interceptor atomic.Value

// Error
func Error(s string) error {
	var err error = errorString(s)
	Intercept(err)
	return err
}

// Errorf
func Errorf(s string, args ...interface{}) error {
	var err error = errorString(fmt.Sprintf(s, args...))
	Intercept(err)
	return err
}

// Register
func Register(intercept func(error)) {
	interceptor.Store(intercept)
}

func Intercept(err error) {
	interceptor.Load().(func(error))(err)
}

type errorString string

func (e errorString) Error() string { return string(e) }

func defaultInterceptor(err error) {
	fmt.Fprintln(os.Stderr, "INTERCEPTED:", err)
}

func init() {
	Register(defaultInterceptor)
}
