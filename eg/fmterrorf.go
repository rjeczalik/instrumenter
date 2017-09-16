package instrumenter

import (
	"fmt"

	"github.com/rjeczalik/instrumenter"
)

func before(msg string, args ...interface{}) error {
	return fmt.Errorf(msg, args...)
}

func after(msg string, args ...interface{}) error {
	return instrumenter.Errorf(msg, args...)
}
