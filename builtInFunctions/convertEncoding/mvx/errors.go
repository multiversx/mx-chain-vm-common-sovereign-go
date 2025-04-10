package mvx

import (
	"errors"
)

var ErrExpectedTypeEnd = errors.New("expected type end")

var ErrInvalidSignatureAbiType = errors.New("invalid signature abi type provided")

var ErrUnhandledAbiArgumentForClone = errors.New("unhandled abi argument for clone")
