package eks

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var binaryNotFoundError = &tracer.Error{
	Kind: "binaryNotFoundError",
}

func IsBinaryNotFound(err error) bool {
	return errors.Is(err, binaryNotFoundError)
}

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

var invalidFlagError = &tracer.Error{
	Kind: "invalidFlagError",
}

func IsInvalidFlag(err error) bool {
	return errors.Is(err, invalidFlagError)
}
