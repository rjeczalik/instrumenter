package templates

import (
	"fmt"

	"github.com/rjeczalik/instrumenter"
)

func before(msg string, arg interface{}) error {
	return fmt.Errorf(msg, arg)
}

func after(msg string, arg interface{}) error {
	return instrumenter.Errorf(msg, arg)
}
