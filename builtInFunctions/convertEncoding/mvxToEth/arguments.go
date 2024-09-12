package mvxToEth

import (
	"fmt"
	convertCommon "github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/common"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/eth"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/mvx"
	"strconv"
	"strings"
)

func MultiversXToEthereumArguments(arguments convertCommon.Arguments) (convertCommon.Arguments, error) {
	convertedArguments := make(convertCommon.Arguments, len(arguments))
	for position, argument := range arguments {
		convertedArgument, err := multiversXToEthereumArgument(argument)
		if err != nil {
			return nil, err
		}
		convertedArguments[position] = convertedArgument
	}
	return convertedArguments, nil
}

func multiversXToEthereumArgument(argument *convertCommon.Argument) (*convertCommon.Argument, error) {
	if strings.HasPrefix(argument.Type, mvx.Array) {
		return multiversXArrayToEthereumArgument(argument)
	}
	if strings.HasPrefix(argument.Type, mvx.Tuple) {
		return multiversXTupleToEthereumArgument(argument)
	}
	switch abiType := argument.Type; abiType {
	case mvx.U8:
		return &convertCommon.Argument{Type: eth.Uint8}, nil
	case mvx.U16:
		return &convertCommon.Argument{Type: eth.Uint16}, nil
	case mvx.U32:
		return &convertCommon.Argument{Type: eth.Uint32}, nil
	case mvx.U64:
		return &convertCommon.Argument{Type: eth.Uint64}, nil
	case mvx.BigUint:
		return &convertCommon.Argument{Type: eth.Uint256}, nil
	case mvx.I8:
		return &convertCommon.Argument{Type: eth.Int8}, nil
	case mvx.I16:
		return &convertCommon.Argument{Type: eth.Int16}, nil
	case mvx.I32:
		return &convertCommon.Argument{Type: eth.Int32}, nil
	case mvx.I64:
		return &convertCommon.Argument{Type: eth.Int64}, nil
	case mvx.BigInt:
		return &convertCommon.Argument{Type: eth.Int256}, nil
	case mvx.Bool:
		return &convertCommon.Argument{Type: eth.Bool}, nil
	case mvx.Bytes, mvx.TokenIdentifier:
		return &convertCommon.Argument{Type: eth.Bytes}, nil
	case mvx.Address:
		return &convertCommon.Argument{Type: eth.Address}, nil
	case mvx.String:
		return &convertCommon.Argument{Type: eth.String}, nil
	case mvx.List, mvx.Variadic:
		return multiversXListToEthereumArgument(argument)
	case mvx.Option, mvx.Optional:
		return multiversXOptionToEthereumArgument(argument)
	case mvx.Multi:
		return multiversXTupleToEthereumArgument(argument)
	default:
		return nil, fmt.Errorf("%w: %v", ErrMvxToEthUnhandledAbiType, abiType)
	}
}

func multiversXListToEthereumArgument(argument *convertCommon.Argument) (*convertCommon.Argument, error) {
	nestedArgument, err := multiversXToEthereumNestedArgument(argument)
	if err != nil {
		return nil, err
	}
	return &convertCommon.Argument{Type: nestedArgument.Type + eth.BeginArray + eth.EndArray, Arguments: nestedArgument.Arguments}, nil
}

func multiversXArrayToEthereumArgument(argument *convertCommon.Argument) (*convertCommon.Argument, error) {
	arraySize, err := mvx.ExtractMvxArraySize(argument.Type)
	if err != nil {
		return nil, err
	}
	nestedArgument, err := multiversXToEthereumNestedArgument(argument)
	if err != nil {
		return nil, err
	}
	return &convertCommon.Argument{Type: nestedArgument.Type + eth.BeginArray + strconv.Itoa(arraySize) + eth.EndArray, Arguments: nestedArgument.Arguments}, nil
}

func multiversXTupleToEthereumArgument(argument *convertCommon.Argument) (*convertCommon.Argument, error) {
	convertedArguments := make(convertCommon.Arguments, len(argument.Arguments))
	for position, nestedArgument := range argument.Arguments {
		convertedArgument, err := multiversXToEthereumArgument(nestedArgument)
		if err != nil {
			return nil, err
		}
		convertedArguments[position] = convertedArgument
	}
	return &convertCommon.Argument{Type: eth.Tuple, Arguments: convertedArguments}, nil
}

func multiversXOptionToEthereumArgument(argument *convertCommon.Argument) (*convertCommon.Argument, error) {
	nestedArgument, err := convertCommon.ExtractNestedArgument(argument)
	if err != nil {
		return nil, err
	}
	return multiversXTupleToEthereumArgument(&convertCommon.Argument{Arguments: convertCommon.Arguments{{Type: mvx.Bool}, nestedArgument}})
}

func multiversXToEthereumNestedArgument(argument *convertCommon.Argument) (*convertCommon.Argument, error) {
	nestedArgument, err := convertCommon.ExtractNestedArgument(argument)
	if err != nil {
		return nil, err
	}
	return multiversXToEthereumArgument(nestedArgument)
}
