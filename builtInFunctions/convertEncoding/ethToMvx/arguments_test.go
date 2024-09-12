package ethToMvx

import (
	convertCommon "github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/common"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/eth"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEthereumToMultiversXArguments(t *testing.T) {
	args, err := eth.ParseEthereumSignature(convertCommon.EthComplexSignature)
	require.NoError(t, err)
	require.NotEmpty(t, args)

	abi, err := eth.BuildEthereumAbi(args)
	require.NoError(t, err)
	require.NotEmpty(t, abi)

	convertedArgs, err := EthereumToMultiversXArguments(abi)
	require.NoError(t, err)
	require.NotEmpty(t, convertedArgs)
}
