package builtInFunctions

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

const (
	esdtIdentifierSeparator  = "-"
	esdtRandomSequenceLength = 6
)

// TopicTokenData groups data that will end up in Topics section of LogEntry
type TopicTokenData struct {
	TokenID []byte
	Nonce   uint64
	Value   *big.Int
}

func addESDTEntryForTransferInVMOutput(
	vmInput *vmcommon.ContractCallInput,
	vmOutput *vmcommon.VMOutput,
	identifier []byte,
	destination []byte,
	topicTokenData []*TopicTokenData,
) {

	topicTokenBytes := make([][]byte, 0)
	for _, tokenData := range topicTokenData {
		nonceBig := big.NewInt(0).SetUint64(tokenData.Nonce)
		topicTokenBytes = append(topicTokenBytes, tokenData.TokenID, nonceBig.Bytes(), tokenData.Value.Bytes())
	}
	topicTokenBytes = append(topicTokenBytes, destination)

	logEntry := &vmcommon.LogEntry{
		Identifier: identifier,
		Address:    vmInput.CallerAddr,
		Topics:     topicTokenBytes,
		Data:       vmcommon.FormatLogDataForCall("", vmInput.Function, vmInput.Arguments),
	}

	if vmOutput.Logs == nil {
		vmOutput.Logs = make([]*vmcommon.LogEntry, 0, 1)
	}

	vmOutput.Logs = append(vmOutput.Logs, logEntry)
}

func addESDTEntryInVMOutput(vmOutput *vmcommon.VMOutput, identifier []byte, tokenID []byte, nonce uint64, value *big.Int, args ...[]byte) {
	entry := newEntryForESDT(identifier, tokenID, nonce, value, args...)

	if vmOutput.Logs == nil {
		vmOutput.Logs = make([]*vmcommon.LogEntry, 0, 1)
	}

	vmOutput.Logs = append(vmOutput.Logs, entry)
}

func newEntryForESDT(identifier, tokenID []byte, nonce uint64, value *big.Int, args ...[]byte) *vmcommon.LogEntry {
	nonceBig := big.NewInt(0).SetUint64(nonce)

	logEntry := &vmcommon.LogEntry{
		Identifier: identifier,
		Topics:     [][]byte{tokenID, nonceBig.Bytes(), value.Bytes()},
	}

	if len(args) > 0 {
		logEntry.Address = args[0]
	}

	if len(args) > 1 {
		logEntry.Topics = append(logEntry.Topics, args[1:]...)
	}

	return logEntry
}

func extractTokenIdentifierAndNonceESDTWipe(esdtPrefix []byte, args []byte) ([]byte, uint64) {
	argsSplit := bytes.Split(args, []byte(esdtIdentifierSeparator))
	invalidLenForArgsWithPrefix := len(argsSplit) < 3 && len(esdtPrefix) != 0
	if len(argsSplit) < 2 || invalidLenForArgsWithPrefix {
		return args, 0
	}

	var tokenId []byte
	var randSeqNonce []byte
	if len(esdtPrefix) == 0 {
		tokenId = argsSplit[0]
		randSeqNonce = argsSplit[1]
	} else {
		tokenId = append(append(argsSplit[0], esdtIdentifierSeparator...), argsSplit[1]...)
		randSeqNonce = argsSplit[2]
	}

	if len(randSeqNonce) <= esdtRandomSequenceLength {
		return args, 0
	}

	identifier := []byte(fmt.Sprintf("%s-%s", tokenId, randSeqNonce[:esdtRandomSequenceLength]))
	nonce := big.NewInt(0).SetBytes(randSeqNonce[esdtRandomSequenceLength:])

	return identifier, nonce.Uint64()
}

func boolToSlice(b bool) []byte {
	return []byte(strconv.FormatBool(b))
}
