package templates

import (
	"errors"

	"github.com/rjeczalik/instrumenter"
)

func before(msg string) error {
	return errors.New(msg)
}

func after(msg string) error {
	return instrumenter.Error(msg)
}
