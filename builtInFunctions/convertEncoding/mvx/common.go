package mvx

import (
	convertCommon "github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/common"
	mvxAbi "github.com/multiversx/mx-sdk-abi-go/abi"
)

type AbiArgument = interface{}

type AbiArguments = []AbiArgument

func SpreadFields(fields []mvxAbi.Field) AbiArguments {
	items := make(AbiArguments, len(fields))
	for position, field := range fields {
		items[position] = field.Value
	}
	return items
}

func BuildOptionMarker() *mvxAbi.BoolValue {
	return &mvxAbi.BoolValue{}
}

const (
	BeginType = "<"
	EndType   = ">"

	U8              = "u8"
	U16             = "u16"
	U32             = "u32"
	U64             = "u64"
	BigUint         = "BigUint"
	I8              = "i8"
	I16             = "i16"
	I32             = "i32"
	I64             = "i64"
	BigInt          = "BigInt"
	Bool            = "bool"
	Bytes           = "bytes"
	Address         = "Address"
	String          = "utf-8 string"
	TokenIdentifier = "TokenIdentifier"
	List            = "List"
	Array           = "array"
	Tuple           = "tuple"
	Option          = "Option"
	Variadic        = "variadic"
	Optional        = "optional"
	Multi           = "multi"
)

var Serializer, _ = mvxAbi.NewSerializer(mvxAbi.ArgsNewSerializer{PartsSeparator: convertCommon.PartsSeparator})
