package ethToMvx

import (
	"fmt"
	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/multiversx/mx-chain-core-go/core"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/common"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/eth"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/mvx"
	mvxAbi "github.com/multiversx/mx-sdk-abi-go/abi"
	"math/big"
	"reflect"
)

func AttachValuesToMultiversXAbi(context *common.EncodingContext, multiversXAbi mvx.AbiArguments, ethereumAbi ethAbi.Arguments, values []interface{}) error {
	if len(multiversXAbi) != len(values) {
		return ErrInvalidValuesSizeForAttach
	}
	for position, value := range values {
		err := attachValueToAbiArgument(context, multiversXAbi[position], &ethereumAbi[position].Type, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func attachValueToAbiArgument(context *common.EncodingContext, mvxAbiArgument mvx.AbiArgument, ethArgument *ethAbi.Type, value interface{}) (err error) {
	defer func() {
		if recover() != nil {
			err = fmt.Errorf("%w for argument %T and value %v", ErrAttachFailed, mvxAbiArgument, value)
		}
	}()
	switch convertedMvxAbiArgument := mvxAbiArgument.(type) {
	case *mvxAbi.U8Value, *mvxAbi.U16Value, *mvxAbi.U32Value, *mvxAbi.U64Value, *mvxAbi.BigUIntValue:
		attachUint(convertedMvxAbiArgument, ethArgument, value)
	case *mvxAbi.I8Value, *mvxAbi.I16Value, *mvxAbi.I32Value, *mvxAbi.I64Value, *mvxAbi.BigIntValue:
		attachInt(convertedMvxAbiArgument, ethArgument, value)
	case *mvxAbi.BoolValue:
		convertedMvxAbiArgument.Value = value.(bool)
	case *mvxAbi.BytesValue:
		convertedMvxAbiArgument.Value = value.([]byte)
	case *mvxAbi.AddressValue:
		err = attachAddressValue(context, convertedMvxAbiArgument, value)
	case *mvxAbi.StringValue:
		convertedMvxAbiArgument.Value = value.(string)
	case *mvxAbi.ListValue:
		err = attachList(context, convertedMvxAbiArgument, ethArgument, value)
	case *mvxAbi.ArrayValue:
		err = attachArray(context, convertedMvxAbiArgument, ethArgument, value)
	case *mvxAbi.StructValue:
		err = attachStruct(context, mvx.SpreadFields(convertedMvxAbiArgument.Fields), ethArgument.TupleElems, value)
	case *mvxAbi.OptionValue:
		err = attachOption(context, convertedMvxAbiArgument, ethArgument, value)
	case *mvxAbi.OptionalValue:
		err = attachOptional(context, convertedMvxAbiArgument, ethArgument, value)
	case *mvxAbi.VariadicValues:
		err = attachVariadicValues(context, convertedMvxAbiArgument, ethArgument, value)
	case *mvxAbi.MultiValue:
		err = attachStruct(context, convertedMvxAbiArgument.Items, ethArgument.TupleElems, value)
	default:
		err = fmt.Errorf("%w: %T", ErrUnhandledAbiArgumentForAttach, convertedMvxAbiArgument)
	}
	return
}

func attachList(context *common.EncodingContext, listValue *mvxAbi.ListValue, ethArgument *ethAbi.Type, value interface{}) error {
	items, err := attachIterable(context, listValue.ItemCreator, ethArgument.Elem, value)
	if err != nil {
		return err
	}
	listValue.Items = items
	return nil
}

func attachArray(context *common.EncodingContext, arrayValue *mvxAbi.ArrayValue, ethArgument *ethAbi.Type, value interface{}) error {
	items, err := attachIterable(context, arrayValue.ItemCreator, eth.ExtractElem(ethArgument), value)
	if err != nil {
		return err
	}
	if len(items) != int(arrayValue.Size) {
		return ErrInvalidValueSizeForArrayAttach
	}
	arrayValue.Items = items
	return nil
}

func attachOption(context *common.EncodingContext, optionValue *mvxAbi.OptionValue, ethArgument *ethAbi.Type, value interface{}) error {
	return attachConditional(context, optionValue.Value, ethArgument.TupleElems, value, func() { optionValue.Value = nil })
}

func attachOptional(context *common.EncodingContext, optionalValue *mvxAbi.OptionalValue, ethArgument *ethAbi.Type, value interface{}) error {
	return attachConditional(context, optionalValue.Value, ethArgument.TupleElems, value, func() { optionalValue.Value = nil })
}

func attachVariadicValues(context *common.EncodingContext, variadicValues *mvxAbi.VariadicValues, ethArgument *ethAbi.Type, value interface{}) error {
	items, err := attachIterable(context, variadicValues.ItemCreator, ethArgument.Elem, value)
	if err != nil {
		return err
	}
	variadicValues.Items = items
	return nil
}

func attachIterable[T mvx.AbiArgument](context *common.EncodingContext, itemCreator func() T, ethElem *ethAbi.Type, toAttach interface{}) ([]T, error) {
	values := reflect.ValueOf(toAttach)
	mvxAbiArguments := make([]T, values.Len())
	for position := 0; position < values.Len(); position++ {
		mvxAbiArgument := itemCreator()
		err := attachValueToAbiArgument(context, mvxAbiArgument, ethElem, values.Index(position).Interface())
		if err != nil {
			return nil, err
		}
		mvxAbiArguments[position] = mvxAbiArgument
	}
	return mvxAbiArguments, nil
}

func attachStruct(context *common.EncodingContext, mvxAbiArguments mvx.AbiArguments, ethElems []*ethAbi.Type, toAttach interface{}) error {
	values := reflect.ValueOf(toAttach)
	for position, mvxAbiArgument := range mvxAbiArguments {
		err := attachValueToAbiArgument(context, mvxAbiArgument, ethElems[position], values.Field(position).Interface())
		if err != nil {
			return err
		}
	}
	return nil
}

func attachConditional(context *common.EncodingContext, mvxAbiArgument mvx.AbiArgument, ethElems []*ethAbi.Type, value interface{}, onAbsent func()) error {
	marker := mvx.BuildOptionMarker()
	err := attachStruct(context, mvx.AbiArguments{marker, mvxAbiArgument}, ethElems, value)
	if err != nil {
		return err
	}
	if !marker.Value {
		onAbsent()
	}
	return nil
}

func attachAddressValue(context *common.EncodingContext, addressValue *mvxAbi.AddressValue, value interface{}) error {
	ethereumAddress := value.(ethCommon.Address).Bytes()
	addressResponse, err := context.Accounts.RequestAddress(&vmcommon.AddressRequest{
		SourceAddress:       ethereumAddress,
		SourceIdentifier:    core.ETHAddressIdentifier,
		RequestedIdentifier: core.MVXAddressIdentifier,
		SaveOnGenerate:      true,
	})
	if err != nil {
		return err
	}
	addressValue.Value = addressResponse.RequestedAddress
	return nil
}

func attachUint(mvxAbiArgument mvx.AbiArgument, ethArgument *ethAbi.Type, value interface{}) {
	switch convertedMvxAbiArgument := mvxAbiArgument.(type) {
	case *mvxAbi.U8Value:
		if ethArgument.Size != eth.Size8 {
			convertedMvxAbiArgument.Value = uint8(asUint64(value))
		} else {
			convertedMvxAbiArgument.Value = value.(uint8)
		}
	case *mvxAbi.U16Value:
		if ethArgument.Size != eth.Size16 {
			convertedMvxAbiArgument.Value = uint16(asUint64(value))
		} else {
			convertedMvxAbiArgument.Value = value.(uint16)
		}
	case *mvxAbi.U32Value:
		if ethArgument.Size != eth.Size32 {
			value = uint32(asUint64(value))
		} else {
			convertedMvxAbiArgument.Value = value.(uint32)
		}
	case *mvxAbi.U64Value:
		if ethArgument.Size != eth.Size64 {
			convertedMvxAbiArgument.Value = asUint64(value)
		} else {
			convertedMvxAbiArgument.Value = value.(uint64)
		}
	case *mvxAbi.BigUIntValue:
		convertedMvxAbiArgument.Value = value.(*big.Int)
	}
}

func attachInt(mvxAbiArgument mvx.AbiArgument, ethArgument *ethAbi.Type, value interface{}) {
	switch convertedMvxAbiArgument := mvxAbiArgument.(type) {
	case *mvxAbi.I8Value:
		if ethArgument.Size != eth.Size8 {
			convertedMvxAbiArgument.Value = int8(asInt64(value))
		} else {
			convertedMvxAbiArgument.Value = value.(int8)
		}
	case *mvxAbi.I16Value:
		if ethArgument.Size != eth.Size16 {
			convertedMvxAbiArgument.Value = int16(asInt64(value))
		} else {
			convertedMvxAbiArgument.Value = value.(int16)
		}
	case *mvxAbi.I32Value:
		if ethArgument.Size != eth.Size32 {
			value = int32(asInt64(value))
		} else {
			convertedMvxAbiArgument.Value = value.(int32)
		}
	case *mvxAbi.I64Value:
		if ethArgument.Size != eth.Size64 {
			value = asInt64(value)
		} else {
			convertedMvxAbiArgument.Value = value.(int64)
		}
	case *mvxAbi.BigIntValue:
		convertedMvxAbiArgument.Value = value.(*big.Int)
	}
}

func asUint64(value interface{}) uint64 {
	return value.(*big.Int).Uint64()
}

func asInt64(value interface{}) int64 {
	return value.(*big.Int).Int64()
}
