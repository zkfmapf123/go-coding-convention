package utils

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

func Test_errWithStack(t *testing.T) {

	err := func() error {
		return func() error {
			err := errors.Wrap(ErrRecordNotFound, "err-1")
			err = errors.Wrap(err, "err-2")
			err = errors.Wrap(err, "err-3")
			err = errors.Wrap(err, "err-4")
			err = errors.Wrap(err, "err-5")
			err = errors.Wrap(err, "err-6")
			return err
		}()
	}()

	fmt.Printf("%+v\n", err)
}
