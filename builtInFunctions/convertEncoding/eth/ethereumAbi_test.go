package eth

import (
	convertCommon "github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/common"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuildEthereumAbi(t *testing.T) {
	args, err := ParseEthereumSignature(convertCommon.EthComplexSignature)
	require.NoError(t, err)
	require.NotEmpty(t, args)

	abi, err := BuildEthereumAbi(args)
	require.NoError(t, err)
	require.NotEmpty(t, abi)

	args, err = ParseEthereumSignature(convertCommon.BuildArgName(0))
	require.NoError(t, err)
	require.NotEmpty(t, args)

	_, err = BuildEthereumAbi(args)
	require.Error(t, err)
}
