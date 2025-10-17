package ethToMvx

import (
	"fmt"
	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	convertCommon "github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/common"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/eth"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/mvx"
	"strconv"
)

func EthereumToMultiversXArguments(arguments ethAbi.Arguments) (convertCommon.Arguments, error) {
	convertedArguments := make(convertCommon.Arguments, len(arguments))
	for position, argument := range arguments {
		convertedArgument, err := ethereumToMultiversXArgument(&argument.Type)
		if err != nil {
			return nil, err
		}
		convertedArguments[position] = convertedArgument
	}
	return convertedArguments, nil
}

func ethereumToMultiversXArgument(argument *ethAbi.Type) (*convertCommon.Argument, error) {
	switch argument.T {
	case ethAbi.IntTy:
		return ethereumIntToMultiversXArgument(argument)
	case ethAbi.UintTy:
		return ethereumUintToMultiversXArgument(argument)
	case ethAbi.BoolTy:
		return &convertCommon.Argument{Type: mvx.Bool}, nil
	case ethAbi.StringTy:
		return &convertCommon.Argument{Type: mvx.String}, nil
	case ethAbi.SliceTy:
		return ethereumSliceToMultiversXArgument(argument)
	case ethAbi.ArrayTy:
		return ethereumArrayToMultiversXArgument(argument)
	case ethAbi.TupleTy:
		return ethereumTupleToMultiversXArgument(argument)
	case ethAbi.AddressTy:
		return &convertCommon.Argument{Type: mvx.Address}, nil
	case ethAbi.FixedBytesTy, ethAbi.FunctionTy:
		return ethereumFixedBytesToMultiversXArgument(argument)
	case ethAbi.BytesTy:
		return &convertCommon.Argument{Type: mvx.Bytes}, nil
	default:
		return nil, fmt.Errorf("%w: %v", ErrEthToMvxUnhandledAbiType, argument.String())
	}
}

func ethereumIntToMultiversXArgument(argument *ethAbi.Type) (*convertCommon.Argument, error) {
	switch size := argument.Size; {
	case size <= eth.Size8:
		return &convertCommon.Argument{Type: mvx.I8}, nil
	case size <= eth.Size16:
		return &convertCommon.Argument{Type: mvx.I16}, nil
	case size <= eth.Size32:
		return &convertCommon.Argument{Type: mvx.I32}, nil
	case size <= eth.Size64:
		return &convertCommon.Argument{Type: mvx.I64}, nil
	default:
		return &convertCommon.Argument{Type: mvx.BigInt}, nil
	}
}

func ethereumUintToMultiversXArgument(argument *ethAbi.Type) (*convertCommon.Argument, error) {
	switch size := argument.Size; {
	case size <= eth.Size8:
		return &convertCommon.Argument{Type: mvx.U8}, nil
	case size <= eth.Size16:
		return &convertCommon.Argument{Type: mvx.U16}, nil
	case size <= eth.Size32:
		return &convertCommon.Argument{Type: mvx.U32}, nil
	case size <= eth.Size64:
		return &convertCommon.Argument{Type: mvx.U64}, nil
	default:
		return &convertCommon.Argument{Type: mvx.BigUint}, nil
	}
}

func ethereumSliceToMultiversXArgument(argument *ethAbi.Type) (*convertCommon.Argument, error) {
	nestedArgument, err := ethereumToMultiversXArgument(argument.Elem)
	if err != nil {
		return nil, err
	}
	return &convertCommon.Argument{Type: mvx.List, Arguments: convertCommon.Arguments{nestedArgument}}, nil
}

func ethereumArrayToMultiversXArgument(argument *ethAbi.Type) (*convertCommon.Argument, error) {
	nestedArgument, err := ethereumToMultiversXArgument(argument.Elem)
	if err != nil {
		return nil, err
	}
	return &convertCommon.Argument{Type: mvx.Array + strconv.Itoa(argument.Size), Arguments: convertCommon.Arguments{nestedArgument}}, nil
}

func ethereumFixedBytesToMultiversXArgument(argument *ethAbi.Type) (*convertCommon.Argument, error) {
	uint8Type := ethAbi.Type{T: ethAbi.UintTy, Size: eth.Size8}
	sizedArrayType := ethAbi.Type{T: ethAbi.ArrayTy, Size: argument.Size, Elem: &uint8Type}
	return ethereumArrayToMultiversXArgument(&sizedArrayType)
}

func ethereumTupleToMultiversXArgument(argument *ethAbi.Type) (*convertCommon.Argument, error) {
	nestedArguments := make(convertCommon.Arguments, len(argument.TupleElems))
	for position, element := range argument.TupleElems {
		nestedArgument, err := ethereumToMultiversXArgument(element)
		if err != nil {
			return nil, err
		}
		nestedArguments[position] = nestedArgument
	}
	return &convertCommon.Argument{Type: mvx.Tuple, Arguments: nestedArguments}, nil
}
