package ethToMvx

import "errors"

var ErrUnhandledAbiArgumentForAttach = errors.New("unhandled abi argument for attach")

var ErrInvalidValuesSizeForAttach = errors.New("invalid values size for attach")

var ErrInvalidValueSizeForArrayAttach = errors.New("invalid value size for array attach")

var ErrAttachFailed = errors.New("attach failed")

var ErrEthToMvxUnhandledAbiType = errors.New("unhandled abi type provided for ethereum to multiversX arguments conversion")

var ErrEthToMvxExpectedOneArgument = errors.New("expected one argument for ethereum to multiversX conversion")
