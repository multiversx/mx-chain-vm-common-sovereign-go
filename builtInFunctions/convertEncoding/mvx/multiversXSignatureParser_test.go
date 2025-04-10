package mvx

import (
	convertCommon "github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/common"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseMultiversXSignature(t *testing.T) {
	testParseMultiversXSignature(t, convertCommon.MvxComplexSignature1)
	testParseMultiversXSignature(t, convertCommon.MvxComplexSignature2)
	testParseMultiversXSignature(t, convertCommon.MvxComplexSignature3)

	_, err := ParseMultiversXSignature("")
	require.Equal(t, convertCommon.ErrBlankExpression, err)

	_, err = ParseMultiversXSignature(convertCommon.Comma)
	require.Equal(t, convertCommon.ErrExpressionStartsWithDelimiter, err)

	_, err = ParseMultiversXSignature(Tuple + BeginType)
	require.Equal(t, convertCommon.ErrBlankExpression, err)

	_, err = ParseMultiversXSignature(Tuple + BeginType + Address)
	require.Equal(t, ErrExpectedTypeEnd, err)

	_, err = ParseMultiversXSignature(Address + convertCommon.Comma)
	require.Equal(t, convertCommon.ErrExpectedExpressionAfterComma, err)

	_, err = ParseMultiversXSignature(Address + EndType)
	require.Equal(t, convertCommon.ErrExpectedBlankRemainder, err)
}

func testParseMultiversXSignature(t *testing.T, signature string) {
	args, err := ParseMultiversXSignature(signature)
	require.NoError(t, err)
	validateArguments(t, args)
}

func validateArguments(t *testing.T, args convertCommon.Arguments) {
	require.NotEmpty(t, args)
	for _, arg := range args {
		if arg.Type == Tuple {
			validateArguments(t, arg.Arguments)
		}
	}
}
