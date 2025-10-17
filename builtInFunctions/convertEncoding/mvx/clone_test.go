package mvx

import (
	mvxAbi "github.com/multiversx/mx-sdk-abi-go/abi"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestCloneAbi(t *testing.T) {
	args := AbiArguments{
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
		&mvxAbi.ListValue{},
		&mvxAbi.ArrayValue{},
		&mvxAbi.StructValue{Fields: []mvxAbi.Field{{Value: &mvxAbi.U8Value{}}}},
		&mvxAbi.OptionValue{Value: &mvxAbi.U8Value{}},
		&mvxAbi.OptionalValue{Value: &mvxAbi.U8Value{}},
		&mvxAbi.VariadicValues{},
		&mvxAbi.MultiValue{Items: AbiArguments{&mvxAbi.U8Value{}}},
	}
	for _, arg := range args {
		cloned, err := cloneAbiArgument(arg)
		require.NoError(t, err)
		require.NotEqual(t, reflect.ValueOf(arg).Pointer(), reflect.ValueOf(cloned).Pointer())
	}
}
