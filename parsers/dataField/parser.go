package datafield

import (
	"encoding/json"
	"math/big"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/elrond-vm-common/parsers"
)

const (
	operationTransfer                 = `transfer`
	minArgumentsQuantityOperationESDT = 2
	minArgumentsQuantityOperationNFT  = 3
	numArgsRelayedV2                  = 4
)

type operationDataFieldParser struct {
	argsParser         vmcommon.CallArgsParser
	shardCoordinator   vmcommon.Coordinator
	esdtTransferParser vmcommon.ESDTTransferParser
}

// NewOperationDataFieldParser will return a new instance of operationDataFieldParser
func NewOperationDataFieldParser(args *ArgsOperationDataFieldParser) (*operationDataFieldParser, error) {
	if check.IfNil(args.ShardCoordinator) {
		return nil, core.ErrNilShardCoordinator
	}
	if check.IfNil(args.Marshalizer) {
		return nil, core.ErrNilMarshalizer
	}

	argsParser := parsers.NewCallArgsParser()
	esdtTransferParser, err := parsers.NewESDTTransferParser(args.Marshalizer)
	if err != nil {
		return nil, err
	}

	return &operationDataFieldParser{
		argsParser:         argsParser,
		shardCoordinator:   args.ShardCoordinator,
		esdtTransferParser: esdtTransferParser,
	}, nil
}

// Parse will parse the provided data field
func (odp *operationDataFieldParser) Parse(dataField []byte, sender, receiver []byte) *ResponseParseData {
	return odp.parse(dataField, sender, receiver, false)
}

func (odp *operationDataFieldParser) parse(dataField []byte, sender, receiver []byte, ignoreRelayed bool) *ResponseParseData {
	responseParse := &ResponseParseData{
		Operation: operationTransfer,
	}

	function, args, err := odp.argsParser.ParseData(string(dataField))
	if err != nil {
		return responseParse
	}

	switch function {
	case core.BuiltInFunctionESDTTransfer:
		return odp.parseESDTTransfer(args, sender, receiver)
	case core.BuiltInFunctionESDTNFTTransfer:
		return odp.parseESDTNFTTransfer(args, sender, receiver)
	case core.BuiltInFunctionMultiESDTNFTTransfer:
		return odp.parseMultiESDTNFTTransfer(args, sender, receiver)
	case core.BuiltInFunctionESDTLocalBurn, core.BuiltInFunctionESDTLocalMint:
		return parseQuantityOperationESDT(args, function)
	case core.BuiltInFunctionESDTWipe, core.BuiltInFunctionESDTFreeze, core.BuiltInFunctionESDTUnFreeze:
		return parseBlockingOperationESDT(args, function)
	case core.BuiltInFunctionESDTNFTCreate, core.BuiltInFunctionESDTNFTBurn, core.BuiltInFunctionESDTNFTAddQuantity:
		return parseQuantityOperationNFT(args, function)
	case core.RelayedTransaction, core.RelayedTransactionV2:
		if ignoreRelayed {
			return &ResponseParseData{
				IsRelayed: true,
			}
		}
		return odp.parseRelayed(function, args, receiver)
	}

	if function != "" && core.IsSmartContractAddress(receiver) {
		responseParse.Function = function
	}

	return responseParse
}

func (odp *operationDataFieldParser) parseRelayed(function string, args [][]byte, receiver []byte) *ResponseParseData {
	if len(args) == 0 {
		return &ResponseParseData{
			IsRelayed: true,
		}
	}

	tx, ok := extractInnerTx(function, args, receiver)
	if !ok {
		return &ResponseParseData{
			IsRelayed: true,
		}
	}

	res := odp.parse(tx.Data, tx.SndAddr, tx.RcvAddr, true)
	if res.IsRelayed {
		return &ResponseParseData{
			IsRelayed: true,
		}
	}

	receivers := [][]byte{tx.RcvAddr}
	receiversShardID := []uint32{odp.shardCoordinator.ComputeId(tx.RcvAddr)}
	if res.Operation == core.BuiltInFunctionMultiESDTNFTTransfer || res.Operation == core.BuiltInFunctionESDTNFTTransfer {
		receivers = res.Receivers
		receiversShardID = res.ReceiversShardID
	}

	return &ResponseParseData{
		Operation:        res.Operation,
		Function:         res.Function,
		ESDTValues:       res.ESDTValues,
		Tokens:           res.Tokens,
		Receivers:        receivers,
		ReceiversShardID: receiversShardID,
		IsRelayed:        true,
	}
}

func extractInnerTx(function string, args [][]byte, receiver []byte) (*transaction.Transaction, bool) {
	tx := &transaction.Transaction{}

	if function == core.RelayedTransaction {
		err := json.Unmarshal(args[0], &tx)

		return tx, err == nil
	}

	if len(args) != numArgsRelayedV2 {
		return nil, false
	}

	// sender of the inner tx is the receiver of the relayed tx
	tx.SndAddr = receiver
	tx.RcvAddr = args[0]
	tx.Data = args[2]

	return tx, true
}

func parseBlockingOperationESDT(args [][]byte, funcName string) *ResponseParseData {
	responseData := &ResponseParseData{
		Operation: funcName,
	}

	if len(args) == 0 {
		return responseData
	}

	token, nonce := extractTokenIdentifierAndNonce(args[0])
	tokenStr := string(token)
	if nonce != 0 {
		tokenStr = computeTokenIdentifier(tokenStr, nonce)
	}

	responseData.Tokens = append(responseData.Tokens, tokenStr)
	return responseData
}

func parseQuantityOperationESDT(args [][]byte, funcName string) *ResponseParseData {
	responseData := &ResponseParseData{
		Operation: funcName,
	}

	if len(args) < minArgumentsQuantityOperationESDT {
		return responseData
	}

	responseData.Tokens = append(responseData.Tokens, string(args[0]))
	responseData.ESDTValues = append(responseData.ESDTValues, big.NewInt(0).SetBytes(args[1]).String())

	return responseData
}

func parseQuantityOperationNFT(args [][]byte, funcName string) *ResponseParseData {
	responseData := &ResponseParseData{
		Operation: funcName,
	}

	if len(args) < minArgumentsQuantityOperationNFT {
		return responseData
	}

	nonce := big.NewInt(0).SetBytes(args[1]).Uint64()
	token := computeTokenIdentifier(string(args[0]), nonce)
	responseData.Tokens = append(responseData.Tokens, token)

	value := big.NewInt(0).SetBytes(args[2]).String()
	responseData.ESDTValues = append(responseData.ESDTValues, value)

	return responseData
}
