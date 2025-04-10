package mvxToEth

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

func DetachValuesFromMultiversXAbi(context *common.EncodingContext, multiversXAbi mvx.AbiArguments, ethereumAbi ethAbi.Arguments) ([]interface{}, error) {
	values := make([]interface{}, len(multiversXAbi))
	for position, mvxAbiArgument := range multiversXAbi {
		value, err := detachValueFromAbiArgument(context, mvxAbiArgument, &ethereumAbi[position].Type)
		if err != nil {
			return nil, err
		}
		values[position] = value
	}
	return values, nil
}

func detachValueFromAbiArgument(context *common.EncodingContext, mvxAbiArgument mvx.AbiArgument, ethArgument *ethAbi.Type) (value interface{}, err error) {
	defer func() {
		if recover() != nil {
			err = fmt.Errorf("%w for argument %T", ErrDetachFailed, mvxAbiArgument)
		}
	}()
	switch convertedMvxAbiArgument := mvxAbiArgument.(type) {
	case *mvxAbi.U8Value, *mvxAbi.U16Value, *mvxAbi.U32Value, *mvxAbi.U64Value, *mvxAbi.BigUIntValue:
		value = detachUint(convertedMvxAbiArgument, ethArgument)
	case *mvxAbi.I8Value, *mvxAbi.I16Value, *mvxAbi.I32Value, *mvxAbi.I64Value, *mvxAbi.BigIntValue:
		value = detachInt(convertedMvxAbiArgument, ethArgument)
	case *mvxAbi.BoolValue:
		value = convertedMvxAbiArgument.Value
	case *mvxAbi.BytesValue:
		value = convertedMvxAbiArgument.Value
	case *mvxAbi.AddressValue:
		value, err = detachAddressValue(context, convertedMvxAbiArgument)
	case *mvxAbi.StringValue:
		value = convertedMvxAbiArgument.Value
	case *mvxAbi.ListValue:
		value, err = detachSlice(context, ethArgument.GetType(), convertedMvxAbiArgument.Items, ethArgument.Elem)
	case *mvxAbi.ArrayValue:
		value, err = detachArray(context, ethArgument.GetType(), convertedMvxAbiArgument.Items, eth.ExtractElem(ethArgument))
	case *mvxAbi.StructValue:
		value, err = detachStruct(context, ethArgument.GetType(), mvx.SpreadFields(convertedMvxAbiArgument.Fields), ethArgument.TupleElems)
	case *mvxAbi.OptionValue:
		value, err = detachConditional(context, ethArgument.GetType(), convertedMvxAbiArgument.Value, ethArgument.TupleElems)
	case *mvxAbi.OptionalValue:
		value, err = detachConditional(context, ethArgument.GetType(), convertedMvxAbiArgument.Value, ethArgument.TupleElems)
	case *mvxAbi.VariadicValues:
		value, err = detachSlice(context, ethArgument.GetType(), convertedMvxAbiArgument.Items, ethArgument.Elem)
	case *mvxAbi.MultiValue:
		value, err = detachStruct(context, ethArgument.GetType(), convertedMvxAbiArgument.Items, ethArgument.TupleElems)
	default:
		err = fmt.Errorf("%w: %T", ErrUnhandledAbiArgumentForDetach, convertedMvxAbiArgument)
	}
	return
}

func detachSlice[T mvx.AbiArgument](context *common.EncodingContext, sliceType reflect.Type, abiArguments []T, ethElem *ethAbi.Type) (interface{}, error) {
	size := len(abiArguments)
	return detachIterable(context, reflect.MakeSlice(sliceType, size, size), abiArguments, ethElem)
}

func detachArray[T mvx.AbiArgument](context *common.EncodingContext, arrayType reflect.Type, abiArguments []T, ethElem *ethAbi.Type) (interface{}, error) {
	if len(abiArguments) != arrayType.Len() {
		return nil, ErrInvalidValueSizeForArrayDetach
	}
	return detachIterable(context, reflect.New(arrayType).Elem(), abiArguments, ethElem)
}

func detachIterable[T mvx.AbiArgument](context *common.EncodingContext, iterable reflect.Value, mvxAbiArguments []T, ethElem *ethAbi.Type) (interface{}, error) {
	for position, mvxAbiArgument := range mvxAbiArguments {
		detachedValue, err := detachValueFromAbiArgument(context, mvxAbiArgument, ethElem)
		if err != nil {
			return nil, err
		}
		iterable.Index(position).Set(reflect.ValueOf(detachedValue))
	}
	return iterable.Interface(), nil
}

func detachStruct(context *common.EncodingContext, structType reflect.Type, mvxAbiArguments mvx.AbiArguments, ethElems []*ethAbi.Type) (interface{}, error) {
	reflectStruct := reflect.New(structType).Elem()
	for position, mvxAbiArgument := range mvxAbiArguments {
		if mvxAbiArgument != nil {
			detachedValue, err := detachValueFromAbiArgument(context, mvxAbiArgument, ethElems[position])
			if err != nil {
				return nil, err
			}
			reflectStruct.Field(position).Set(reflect.ValueOf(detachedValue))
		}
	}
	return reflectStruct.Interface(), nil
}

func detachConditional(context *common.EncodingContext, conditionalType reflect.Type, mvxAbiArgument mvx.AbiArgument, ethElems []*ethAbi.Type) (interface{}, error) {
	marker := mvx.BuildOptionMarker()
	marker.Value = mvxAbiArgument != nil
	return detachStruct(context, conditionalType, mvx.AbiArguments{marker, mvxAbiArgument}, ethElems)
}

func detachAddressValue(context *common.EncodingContext, addressValue *mvxAbi.AddressValue) (interface{}, error) {
	multiversXAddress := addressValue.Value
	addressResponse, err := context.Accounts.RequestAddress(&vmcommon.AddressRequest{
		SourceAddress:       multiversXAddress,
		SourceIdentifier:    core.MVXAddressIdentifier,
		RequestedIdentifier: core.ETHAddressIdentifier,
		SaveOnGenerate:      true,
	})
	if err != nil {
		return ethCommon.Address{}, err
	}
	return ethCommon.BytesToAddress(addressResponse.RequestedAddress), nil
}

func detachUint(mvxAbiArgument mvx.AbiArgument, ethArgument *ethAbi.Type) (value interface{}) {
	switch convertedMvxAbiArgument := mvxAbiArgument.(type) {
	case *mvxAbi.U8Value:
		if ethArgument.Size != eth.Size8 {
			value = new(big.Int).SetUint64(uint64(convertedMvxAbiArgument.Value))
		} else {
			value = convertedMvxAbiArgument.Value
		}
	case *mvxAbi.U16Value:
		if ethArgument.Size != eth.Size16 {
			value = new(big.Int).SetUint64(uint64(convertedMvxAbiArgument.Value))
		} else {
			value = convertedMvxAbiArgument.Value
		}
	case *mvxAbi.U32Value:
		if ethArgument.Size != eth.Size32 {
			value = new(big.Int).SetUint64(uint64(convertedMvxAbiArgument.Value))
		} else {
			value = convertedMvxAbiArgument.Value
		}
	case *mvxAbi.U64Value:
		if ethArgument.Size != eth.Size64 {
			value = new(big.Int).SetUint64(convertedMvxAbiArgument.Value)
		} else {
			value = convertedMvxAbiArgument.Value
		}
	case *mvxAbi.BigUIntValue:
		value = convertedMvxAbiArgument.Value
	}
	return
}

func detachInt(mvxAbiArgument mvx.AbiArgument, ethArgument *ethAbi.Type) (value interface{}) {
	switch convertedMvxAbiArgument := mvxAbiArgument.(type) {
	case *mvxAbi.I8Value:
		if ethArgument.Size != eth.Size8 {
			value = new(big.Int).SetInt64(int64(convertedMvxAbiArgument.Value))
		} else {
			value = convertedMvxAbiArgument.Value
		}
	case *mvxAbi.I16Value:
		if ethArgument.Size != eth.Size16 {
			value = new(big.Int).SetInt64(int64(convertedMvxAbiArgument.Value))
		} else {
			value = convertedMvxAbiArgument.Value
		}
	case *mvxAbi.I32Value:
		if ethArgument.Size != eth.Size32 {
			value = new(big.Int).SetInt64(int64(convertedMvxAbiArgument.Value))
		} else {
			value = convertedMvxAbiArgument.Value
		}
	case *mvxAbi.I64Value:
		if ethArgument.Size != eth.Size64 {
			value = new(big.Int).SetInt64(convertedMvxAbiArgument.Value)
		} else {
			value = convertedMvxAbiArgument.Value
		}
	case *mvxAbi.BigIntValue:
		value = convertedMvxAbiArgument.Value
	}
	return
}
