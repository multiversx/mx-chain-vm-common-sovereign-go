package mvxToEth

import (
	convertCommon "github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/common"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/mvx"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMultiversXToEthereumArguments(t *testing.T) {
	testMultiversXToEthereumArguments(t, convertCommon.MvxComplexSignature1)
	testMultiversXToEthereumArguments(t, convertCommon.MvxComplexSignature2)
	testMultiversXToEthereumArguments(t, convertCommon.MvxComplexSignature3)
}

func testMultiversXToEthereumArguments(t *testing.T, signature string) {
	args, err := mvx.ParseMultiversXSignature(signature)
	require.NoError(t, err)
	require.NotEmpty(t, args)

	convertedArgs, err := MultiversXToEthereumArguments(args)
	require.NoError(t, err)
	require.NotEmpty(t, convertedArgs)
}
