package templates

import (
	"fmt"

	"github.com/rjeczalik/instrumenter"
)

func before(msg string, arg1, arg2 interface{}) error {
	return fmt.Errorf(msg, arg1, arg2)
}

func after(msg string, arg interface{}) error {
	return instrumenter.Errorf(msg, arg1, arg2)
}
