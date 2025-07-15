package builtInFunctions

import (
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/stretchr/testify/require"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

func TestNewEntryForNFT(t *testing.T) {
	t.Parallel()

	vmOutput := &vmcommon.VMOutput{}
	addESDTEntryInVMOutput(vmOutput, []byte(core.BuiltInFunctionESDTNFTCreate), []byte("my-token"), 5, big.NewInt(1), []byte("caller"), []byte("receiver"))
	require.Equal(t, &vmcommon.LogEntry{
		Identifier: []byte(core.BuiltInFunctionESDTNFTCreate),
		Address:    []byte("caller"),
		Topics:     [][]byte{[]byte("my-token"), big.NewInt(0).SetUint64(5).Bytes(), big.NewInt(1).Bytes(), []byte("receiver")},
		Data:       nil,
	}, vmOutput.Logs[0])
}

func TestExtractTokenIdentifierAndNonceESDTWipe(t *testing.T) {
	t.Parallel()

	prefix := []byte{}
	token := []byte("TOKEN")
	identifier, nonce := extractTokenIdentifierAndNonceESDTWipe(prefix, token)
	require.Equal(t, uint64(0), nonce)
	require.Equal(t, token, identifier)

	prefix = []byte("prf")
	token = []byte("TOKEN")
	identifier, nonce = extractTokenIdentifierAndNonceESDTWipe(prefix, token)
	require.Equal(t, uint64(0), nonce)
	require.Equal(t, token, identifier)

	prefix = []byte("prf")
	token = []byte("TOKEN-1a2b3c")
	identifier, nonce = extractTokenIdentifierAndNonceESDTWipe(prefix, token)
	require.Equal(t, uint64(0), nonce)
	require.Equal(t, token, identifier)

	prefix = []byte{}
	token = []byte("TOKEN-1a2b3c") // no nonce
	identifier, nonce = extractTokenIdentifierAndNonceESDTWipe(prefix, token)
	require.Equal(t, uint64(0), nonce)
	require.Equal(t, token, identifier)

	prefix = []byte{}
	token = []byte("TOKEN-1a2b3c")
	tokenNonce := big.NewInt(1)
	tokenWithNonce := append(token, tokenNonce.Bytes()...)
	identifier, nonce = extractTokenIdentifierAndNonceESDTWipe(prefix, tokenWithNonce)
	require.Equal(t, tokenNonce.Uint64(), nonce)
	require.Equal(t, token, identifier)

	prefix = []byte{}
	token = []byte("prf-TOKEN-a1b2c3")
	identifier, nonce = extractTokenIdentifierAndNonceESDTWipe(prefix, token)
	require.Equal(t, uint64(0), nonce)
	require.Equal(t, token, identifier)

	prefix = []byte("prf")
	token = []byte("prf-TOKEN-a1b2c3") // no nonce
	identifier, nonce = extractTokenIdentifierAndNonceESDTWipe(prefix, token)
	require.Equal(t, uint64(0), nonce)
	require.Equal(t, token, identifier)

	prefix = []byte("prf")
	token = []byte("prf-TOKEN-a1b2c3")
	tokenNonce = big.NewInt(2)
	tokenWithNonce = append(token, tokenNonce.Bytes()...)
	identifier, nonce = extractTokenIdentifierAndNonceESDTWipe(prefix, tokenWithNonce)
	require.Equal(t, tokenNonce.Uint64(), nonce)
	require.Equal(t, token, identifier)

}
