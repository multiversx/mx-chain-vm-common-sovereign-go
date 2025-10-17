package mvxToEth

import "errors"

var ErrUnhandledAbiArgumentForDetach = errors.New("unhandled abi argument for detach")

var ErrInvalidValueSizeForArrayDetach = errors.New("invalid value size for array detach")

var ErrDetachFailed = errors.New("detach failed")

var ErrMvxToEthUnhandledAbiType = errors.New("unhandled abi type provided for multiversX to ethereum arguments conversion")
