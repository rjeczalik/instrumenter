package intercept

import (
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"runtime"
	"time"
)

type Error struct {
	ID         string    `json:"id,omitempty"`
	Path       string    `json:"path"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"createdAt"`
	Stacktrace []byte    `json:"stacktrace"`
	Type       string    `json:"type"`
	Message    string    `json:"message"`
}

func NewError(err error) *Error {
	e := &Error{
		Path:       os.Args[0],
		CreatedAt:  time.Now(),
		Stacktrace: make([]byte, 4096),
		Type:       reflect.TypeOf(err).String(),
		Message:    err.Error(),
	}

	if u, err := user.Current(); err == nil {
		e.Username = u.Username
	}

	if abs, err := filepath.Abs(e.Path); err == nil {
		e.Path = abs
	}

	n := runtime.Stack(e.Stacktrace, false)
	e.Stacktrace = e.Stacktrace[:n]

	return e
}

type NewRequest struct {
	Error *Error `json:"error"`
}

type NewResponse struct {
	ID string `json:"id"`
}
