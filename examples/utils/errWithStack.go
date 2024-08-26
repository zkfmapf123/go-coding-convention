package utils

import (
	"github.com/pkg/errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrUnknown        = errors.New("unknown error")
)

func withErrStack() {

	errors.WithStack(ErrRecordNotFound)
}
