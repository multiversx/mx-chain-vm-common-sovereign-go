package ethToMvx

import (
	ethCommon "github.com/ethereum/go-ethereum/common"
	convertCommon "github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/common"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/eth"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/mvx"
	mvxAbi "github.com/multiversx/mx-sdk-abi-go/abi"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestAttachValuesToMultiversXAbi(t *testing.T) {
	mvxArgs := mvx.AbiArguments{
		&mvxAbi.U8Value{},
		&mvxAbi.U16Value{},
		&mvxAbi.U32Value{},
		&mvxAbi.U64Value{},
		&mvxAbi.BigUIntValue{},
		&mvxAbi.I8Value{},
		&mvxAbi.I16Value{},
		&mvxAbi.I32Value{},
		&mvxAbi.I64Value{},
		&mvxAbi.BigIntValue{},
		&mvxAbi.BoolValue{},
		&mvxAbi.BytesValue{},
		&mvxAbi.AddressValue{},
		&mvxAbi.StringValue{},
		&mvxAbi.ListValue{ItemCreator: func() mvxAbi.SingleValue { return &mvxAbi.StringValue{} }},
		&mvxAbi.ArrayValue{Size: 2, ItemCreator: func() mvxAbi.SingleValue { return &mvxAbi.StringValue{} }},
		&mvxAbi.StructValue{Fields: []mvxAbi.Field{{Value: &mvxAbi.U8Value{}}}},
		&mvxAbi.OptionValue{Value: &mvxAbi.U8Value{}},
		&mvxAbi.OptionalValue{Value: &mvxAbi.U8Value{}},
		&mvxAbi.VariadicValues{ItemCreator: func() mvx.AbiArgument { return &mvxAbi.StringValue{} }},
		&mvxAbi.MultiValue{Items: mvx.AbiArguments{&mvxAbi.U8Value{}, &mvxAbi.StringValue{}}},
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
		&convertCommon.Argument{Type: eth.Tuple, Arguments: convertCommon.Arguments{{Type: eth.Bool}, {Type: eth.Uint8}}},
		&convertCommon.Argument{Type: eth.String + eth.BeginArray + eth.EndArray},
		&convertCommon.Argument{Type: eth.Tuple, Arguments: convertCommon.Arguments{{Type: eth.Uint8}, {Type: eth.String}}},
	}
	values := mvx.AbiArguments{
		uint8(1),
		uint16(2),
		uint32(3),
		uint64(4),
		big.NewInt(5),
		int8(6),
		int16(7),
		int32(8),
		int64(9),
		big.NewInt(10),
		true,
		[]byte{0x11, 0x12},
		ethCommon.BytesToAddress([]byte{0x13, 0x14}),
		convertCommon.BuildArgName(15),
		[]string{convertCommon.BuildArgName(16), convertCommon.BuildArgName(17)},
		[2]string{convertCommon.BuildArgName(18), convertCommon.BuildArgName(19)},
		struct{ ArgName1 uint8 }{uint8(20)},
		struct {
			ArgName1 bool
			ArgName2 uint8
		}{false, uint8(21)},
		struct {
			ArgName1 bool
			ArgName2 uint8
		}{true, uint8(22)},
		[]string{convertCommon.BuildArgName(23), convertCommon.BuildArgName(24), convertCommon.BuildArgName(25)},
		struct {
			ArgName1 uint8
			ArgName2 string
		}{uint8(26), convertCommon.BuildArgName(27)},
	}
	ethAbi, err := eth.BuildEthereumAbi(ethArgs)
	require.NoError(t, err)
	err = AttachValuesToMultiversXAbi(convertCommon.BuildTestEncodingContext(), mvxArgs, ethAbi, values)
	require.NoError(t, err)
}
