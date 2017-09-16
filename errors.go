package instrumenter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/rjeczalik/instrumenter/intercept"
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

func defaultInterceptor(e error) {
	req := &intercept.NewRequest{
		Error: intercept.NewError(e),
	}
	p, err := json.Marshal(req)
	if err != nil {
		log.Println("failed to send error:", err)
		return
	}

	resp, err := http.Post("http://127.0.0.1:8484/new", "application/json", bytes.NewReader(p))
	if err != nil {
		log.Println("failed to send error:", err)
		return
	}
	resp.Body.Close()
}

func init() {
	Register(defaultInterceptor)
}
