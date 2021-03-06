package docker

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var executionFailedError = &tracer.Error{
	Kind: "executionFailedError",
}

func IsExecutionFailed(err error) bool {
	return errors.Is(err, executionFailedError)
}

var invalidConfigError = &tracer.Error{
	Kind: "invalidConfigError",
}

func IsInvalidConfig(err error) bool {
	return errors.Is(err, invalidConfigError)
}
