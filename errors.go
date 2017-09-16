package instrumenter

import (
	"errors"
	"fmt"
	"os"
	"sync/atomic"
)

var interceptor atomic.Value

// Error
func Error(s string) error {
	err := errors.New(s)
	Intercept(err)
	return err
}

// Errorf
func Errorf(s string, args ...interface{}) error {
	err := fmt.Errorf(s, args...)
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

func defaultInterceptor(err error) {
	fmt.Fprintln(os.Stderr, "INTERCEPTED:", err)
}

func init() {
	Register(defaultInterceptor)
}
