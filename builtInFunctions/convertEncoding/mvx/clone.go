package mvx

import (
	"fmt"

	mvxAbi "github.com/multiversx/mx-sdk-abi-go/abi"
)

func cloneAbiArgument(abiArgument AbiArgument) (AbiArgument, error) {
	switch convertedAbiArgument := abiArgument.(type) {
	case *mvxAbi.OptionalValue:
		return cloneOptionalValue(convertedAbiArgument)
	case *mvxAbi.VariadicValues:
		return cloneVariadicValues(convertedAbiArgument)
	case *mvxAbi.MultiValue:
		return cloneMultiValue(convertedAbiArgument)
	default:
		return cloneSingleValue(abiArgument)
	}
}

func cloneSingleValue(abiArgument AbiArgument) (mvxAbi.SingleValue, error) {
	switch convertedAbiArgument := abiArgument.(type) {
	case *mvxAbi.U8Value:
		return &mvxAbi.U8Value{}, nil
	case *mvxAbi.U16Value:
		return &mvxAbi.U16Value{}, nil
	case *mvxAbi.U32Value:
		return &mvxAbi.U32Value{}, nil
	case *mvxAbi.U64Value:
		return &mvxAbi.U64Value{}, nil
	case *mvxAbi.BigUIntValue:
		return &mvxAbi.BigUIntValue{}, nil
	case *mvxAbi.I8Value:
		return &mvxAbi.I8Value{}, nil
	case *mvxAbi.I16Value:
		return &mvxAbi.I16Value{}, nil
	case *mvxAbi.I32Value:
		return &mvxAbi.I32Value{}, nil
	case *mvxAbi.I64Value:
		return &mvxAbi.I64Value{}, nil
	case *mvxAbi.BigIntValue:
		return &mvxAbi.BigIntValue{}, nil
	case *mvxAbi.BoolValue:
		return &mvxAbi.BoolValue{}, nil
	case *mvxAbi.BytesValue:
		return &mvxAbi.BytesValue{}, nil
	case *mvxAbi.AddressValue:
		return &mvxAbi.AddressValue{}, nil
	case *mvxAbi.StringValue:
		return &mvxAbi.StringValue{}, nil
	case *mvxAbi.ListValue:
		return cloneListValue(convertedAbiArgument)
	case *mvxAbi.ArrayValue:
		return cloneArrayValue(convertedAbiArgument)
	case *mvxAbi.StructValue:
		return cloneStructValue(convertedAbiArgument)
	case *mvxAbi.OptionValue:
		return cloneOptionValue(convertedAbiArgument)
	default:
		return nil, fmt.Errorf("%w: %T", ErrUnhandledAbiArgumentForClone, convertedAbiArgument)
	}
}

func cloneListValue(listValue *mvxAbi.ListValue) (*mvxAbi.ListValue, error) {
	return &mvxAbi.ListValue{ItemCreator: listValue.ItemCreator}, nil
}

func cloneArrayValue(arrayValue *mvxAbi.ArrayValue) (*mvxAbi.ArrayValue, error) {
	return &mvxAbi.ArrayValue{Length: arrayValue.Length, ItemCreator: arrayValue.ItemCreator}, nil
}

func cloneStructValue(structValue *mvxAbi.StructValue) (*mvxAbi.StructValue, error) {
	clonedFields := make([]mvxAbi.Field, len(structValue.Fields))
	for position, field := range structValue.Fields {
		cloned, err := cloneSingleValue(field.Value)
		if err != nil {
			return nil, err
		}
		clonedFields[position] = mvxAbi.Field{Name: field.Name, Value: cloned}
	}
	return &mvxAbi.StructValue{Fields: clonedFields}, nil
}

func cloneOptionValue(optionValue *mvxAbi.OptionValue) (*mvxAbi.OptionValue, error) {
	cloned, err := cloneSingleValue(optionValue.Value)
	if err != nil {
		return nil, err
	}
	return &mvxAbi.OptionValue{Value: cloned}, nil
}

func cloneOptionalValue(optionalValue *mvxAbi.OptionalValue) (*mvxAbi.OptionalValue, error) {
	cloned, err := cloneAbiArgument(optionalValue.Value)
	if err != nil {
		return nil, err
	}
	return &mvxAbi.OptionalValue{Value: cloned}, nil
}

func cloneVariadicValues(variadicValues *mvxAbi.VariadicValues) (*mvxAbi.VariadicValues, error) {
	return &mvxAbi.VariadicValues{ItemCreator: variadicValues.ItemCreator}, nil
}

func cloneMultiValue(multiValue *mvxAbi.MultiValue) (*mvxAbi.MultiValue, error) {
	clonedItems := make(AbiArguments, len(multiValue.Items))
	for position, item := range multiValue.Items {
		cloned, err := cloneAbiArgument(item)
		if err != nil {
			return nil, err
		}
		clonedItems[position] = cloned
	}
	return &mvxAbi.MultiValue{Items: clonedItems}, nil
}
