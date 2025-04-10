package mvxToEth

import (
	"bytes"
	convertCommon "github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/common"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/eth"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/mvx"
	mvxAbi "github.com/multiversx/mx-sdk-abi-go/abi"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestDetachValuesFromMultiversXAbi(t *testing.T) {
	mvxArgs := mvx.AbiArguments{
		&mvxAbi.U8Value{Value: uint8(1)},
		&mvxAbi.U16Value{Value: uint16(2)},
		&mvxAbi.U32Value{Value: uint32(3)},
		&mvxAbi.U64Value{Value: uint64(4)},
		&mvxAbi.BigUIntValue{Value: big.NewInt(5)},
		&mvxAbi.I8Value{Value: int8(6)},
		&mvxAbi.I16Value{Value: int16(7)},
		&mvxAbi.I32Value{Value: int32(8)},
		&mvxAbi.I64Value{Value: int64(9)},
		&mvxAbi.BigIntValue{Value: big.NewInt(10)},
		&mvxAbi.BoolValue{Value: true},
		&mvxAbi.BytesValue{Value: []byte{0x11, 0x12}},
		&mvxAbi.AddressValue{Value: bytes.Repeat([]byte{0x13, 0x14}, 16)},
		&mvxAbi.StringValue{Value: convertCommon.BuildArgName(15)},
		&mvxAbi.ListValue{Items: []mvxAbi.SingleValue{
			&mvxAbi.StringValue{Value: convertCommon.BuildArgName(16)},
			&mvxAbi.StringValue{Value: convertCommon.BuildArgName(17)},
		}},
		&mvxAbi.ArrayValue{Items: []mvxAbi.SingleValue{
			&mvxAbi.StringValue{Value: convertCommon.BuildArgName(18)},
			&mvxAbi.StringValue{Value: convertCommon.BuildArgName(19)},
		}},
		&mvxAbi.StructValue{Fields: []mvxAbi.Field{
			{Value: &mvxAbi.U8Value{Value: uint8(20)}},
		}},
		&mvxAbi.OptionValue{Value: &mvxAbi.U8Value{Value: uint8(21)}},
		&mvxAbi.OptionValue{Value: &mvxAbi.ListValue{Items: []mvxAbi.SingleValue{&mvxAbi.U32Value{Value: uint32(22)}}}},
		&mvxAbi.OptionalValue{Value: nil},
		&mvxAbi.VariadicValues{Items: mvx.AbiArguments{
			&mvxAbi.StringValue{Value: convertCommon.BuildArgName(23)},
			&mvxAbi.StringValue{Value: convertCommon.BuildArgName(24)},
			&mvxAbi.StringValue{Value: convertCommon.BuildArgName(25)},
		}},
		&mvxAbi.MultiValue{Items: mvx.AbiArguments{
			&mvxAbi.U8Value{Value: uint8(26)},
			&mvxAbi.StringValue{Value: convertCommon.BuildArgName(27)},
		}},
	}
	ethArgs := convertCommon.Arguments{
		&convertCommon.Argument{Type: eth.Uint8},
		&convertCommon.Argument{Type: eth.Uint16},
		&convertCommon.Argument{Type: eth.Uint32},
		&convertCommon.Argument{Type: eth.Uint64},
		&convertCommon.Argument{Type: eth.Uint256},
		&convertCommon.Argument{Type: eth.Int8},
		&convertCommon.Argument{Type: eth.Int16},
		&convertCommon.Argument{Type: eth.Int32},
		&convertCommon.Argument{Type: eth.Int64},
		&convertCommon.Argument{Type: eth.Int256},
		&convertCommon.Argument{Type: eth.Bool},
		&convertCommon.Argument{Type: eth.Bytes},
		&convertCommon.Argument{Type: eth.Address},
		&convertCommon.Argument{Type: eth.String},
		&convertCommon.Argument{Type: eth.String + eth.BeginArray + eth.EndArray},
		&convertCommon.Argument{Type: eth.String + eth.BeginArray + "2" + eth.EndArray},
		&convertCommon.Argument{Type: eth.Tuple, Arguments: convertCommon.Arguments{{Type: eth.Uint8}}},
		&convertCommon.Argument{Type: eth.Tuple, Arguments: convertCommon.Arguments{{Type: eth.Bool}, {Type: eth.Uint8}}},
		&convertCommon.Argument{Type: eth.Tuple, Arguments: convertCommon.Arguments{{Type: eth.Bool}, {Type: eth.Uint32 + eth.BeginArray + eth.EndArray}}},
		&convertCommon.Argument{Type: eth.Tuple, Arguments: convertCommon.Arguments{{Type: eth.Bool}, {Type: eth.Uint32 + eth.BeginArray + eth.EndArray}}},
		&convertCommon.Argument{Type: eth.String + eth.BeginArray + eth.EndArray},
		&convertCommon.Argument{Type: eth.Tuple, Arguments: convertCommon.Arguments{{Type: eth.Uint8}, {Type: eth.String}}},
	}
	ethAbi, err := eth.BuildEthereumAbi(ethArgs)
	require.NoError(t, err)
	values, err := DetachValuesFromMultiversXAbi(convertCommon.BuildTestEncodingContext(), mvxArgs, ethAbi)
	require.NoError(t, err)
	require.Equal(t, len(values), len(mvxArgs))
}
