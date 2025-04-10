package mvx

import (
	"fmt"
	convertCommon "github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/common"
	mvxAbi "github.com/multiversx/mx-sdk-abi-go/abi"
	"strings"
)

func BuildMultiversXAbi(arguments convertCommon.Arguments) (AbiArguments, error) {
	abiArguments := make(AbiArguments, len(arguments))
	for position, argument := range arguments {
		abiArgument, err := buildAbiArgument(argument)
		if err != nil {
			return nil, err
		}
		abiArguments[position] = abiArgument
	}
	return abiArguments, nil
}

func buildAbiArgument(argument *convertCommon.Argument) (AbiArgument, error) {
	switch argument.Type {
	case Variadic:
		return argumentToVariadicValues(argument)
	case Optional:
		return argumentToOptionalValue(argument)
	case Multi:
		return argumentToMultiValue(argument)
	default:
		return buildSingleValue(argument)
	}
}

func buildSingleValue(argument *convertCommon.Argument) (mvxAbi.SingleValue, error) {
	if strings.HasPrefix(argument.Type, Array) {
		return argumentToArrayValue(argument)
	}
	if strings.HasPrefix(argument.Type, Tuple) {
		return argumentToStructValue(argument)
	}
	switch abiType := argument.Type; abiType {
	case U8:
		return &mvxAbi.U8Value{}, nil
	case U16:
		return &mvxAbi.U16Value{}, nil
	case U32:
		return &mvxAbi.U32Value{}, nil
	case U64:
		return &mvxAbi.U64Value{}, nil
	case BigUint:
		return &mvxAbi.BigUIntValue{}, nil
	case I8:
		return &mvxAbi.I8Value{}, nil
	case I16:
		return &mvxAbi.I16Value{}, nil
	case I32:
		return &mvxAbi.I32Value{}, nil
	case I64:
		return &mvxAbi.I64Value{}, nil
	case BigInt:
		return &mvxAbi.BigIntValue{}, nil
	case Bool:
		return &mvxAbi.BoolValue{}, nil
	case Bytes:
		return &mvxAbi.BytesValue{}, nil
	case Address:
		return &mvxAbi.AddressValue{}, nil
	case String:
		return &mvxAbi.StringValue{}, nil
	case TokenIdentifier:
		return &mvxAbi.BytesValue{}, nil
	case List:
		return argumentToListValue(argument)
	case Option:
		return argumentToOptionValue(argument)
	default:
		return nil, fmt.Errorf("%w: %v", ErrInvalidSignatureAbiType, abiType)
	}
}

func argumentToListValue(argument *convertCommon.Argument) (mvxAbi.SingleValue, error) {
	singleValue, err := extractNestedSingleValue(argument)
	if err != nil {
		return nil, err
	}
	return &mvxAbi.ListValue{ItemCreator: singleValueItemCreator(singleValue)}, nil
}

func argumentToArrayValue(argument *convertCommon.Argument) (mvxAbi.SingleValue, error) {
	arraySize, err := ExtractMvxArraySize(argument.Type)
	if err != nil {
		return nil, err
	}
	singleValue, err := extractNestedSingleValue(argument)
	if err != nil {
		return nil, err
	}
	return &mvxAbi.ArrayValue{Size: uint32(arraySize), ItemCreator: singleValueItemCreator(singleValue)}, nil
}

func argumentToStructValue(argument *convertCommon.Argument) (mvxAbi.SingleValue, error) {
	fields := make([]mvxAbi.Field, len(argument.Arguments))
	for position, component := range argument.Arguments {
		singleValue, err := buildSingleValue(component)
		if err != nil {
			return nil, err
		}
		fields[position] = mvxAbi.Field{Value: singleValue}
	}
	return &mvxAbi.StructValue{Fields: fields}, nil
}

func argumentToOptionValue(argument *convertCommon.Argument) (mvxAbi.SingleValue, error) {
	singleValue, err := extractNestedSingleValue(argument)
	if err != nil {
		return nil, err
	}
	return &mvxAbi.OptionValue{Value: singleValue}, nil
}

func argumentToOptionalValue(argument *convertCommon.Argument) (AbiArgument, error) {
	abiArgument, err := extractNestedAbiArgument(argument)
	if err != nil {
		return nil, err
	}
	return &mvxAbi.OptionalValue{Value: abiArgument}, nil
}

func argumentToVariadicValues(argument *convertCommon.Argument) (AbiArgument, error) {
	abiArgument, err := extractNestedAbiArgument(argument)
	if err != nil {
		return nil, err
	}
	return &mvxAbi.VariadicValues{ItemCreator: abiArgumentItemCreator(abiArgument)}, nil
}

func argumentToMultiValue(argument *convertCommon.Argument) (AbiArgument, error) {
	items := make(AbiArguments, len(argument.Arguments))
	for position, component := range argument.Arguments {
		item, err := buildAbiArgument(component)
		if err != nil {
			return nil, err
		}
		items[position] = item
	}
	return &mvxAbi.MultiValue{Items: items}, nil
}

func extractNestedSingleValue(argument *convertCommon.Argument) (mvxAbi.SingleValue, error) {
	nestedArgument, err := convertCommon.ExtractNestedArgument(argument)
	if err != nil {
		return nil, err
	}
	return buildSingleValue(nestedArgument)
}

func extractNestedAbiArgument(argument *convertCommon.Argument) (AbiArgument, error) {
	nestedArgument, err := convertCommon.ExtractNestedArgument(argument)
	if err != nil {
		return nil, err
	}
	return buildAbiArgument(nestedArgument)
}

func singleValueItemCreator(toClone mvxAbi.SingleValue) func() mvxAbi.SingleValue {
	return func() mvxAbi.SingleValue {
		cloned, unexpectedErr := cloneSingleValue(toClone)
		if unexpectedErr != nil {
			panic(unexpectedErr)
		}
		return cloned
	}
}

func abiArgumentItemCreator(toClone AbiArgument) func() AbiArgument {
	return func() AbiArgument {
		cloned, unexpectedErr := cloneAbiArgument(toClone)
		if unexpectedErr != nil {
			panic(unexpectedErr)
		}
		return cloned
	}
}
