package templates

import (
	"fmt"

	"github.com/rjeczalik/instrumenter"
)

func before(msg string, arg1, arg2, arg3, arg4, arg5, arg6 interface{}) error {
	return fmt.Errorf(msg, arg1, arg2, arg3, arg4, arg5, arg6)
}

func after(msg string, arg1, arg2, arg3, arg4, arg5, arg6 interface{}) error {
	return instrumenter.Errorf(msg, arg1, arg2, arg3, arg4, arg5, arg6)
}
