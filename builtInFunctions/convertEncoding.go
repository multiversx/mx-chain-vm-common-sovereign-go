package builtInFunctions

import (
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding"
	"sync"

	"github.com/multiversx/mx-chain-vm-common-go"
)

const (
	minArgsCount      = 2
	signaturePosition = 0
	dataStartPosition = 1
)

type ConvertEncoding struct {
	baseAlwaysActiveHandler
	baseOperationCost vmcommon.BaseOperationCost
	function          string
	mutExecution      sync.RWMutex
	accounts          vmcommon.AccountsAdapter
}

// NewConvertEncodingFunc converts the arguments using the specified configuration for the encoding
func NewConvertEncodingFunc(
	baseOperationCost vmcommon.BaseOperationCost,
	accounts vmcommon.AccountsAdapter,
	function string,
) (*ConvertEncoding, error) {
	if check.IfNil(accounts) {
		return nil, ErrNilAccountsAdapter
	}

	return &ConvertEncoding{
		baseOperationCost: baseOperationCost,
		function:          function,
		mutExecution:      sync.RWMutex{},
		accounts:          accounts,
	}, nil
}

// SetNewGasConfig is called whenever gas cost is changed
func (ce *ConvertEncoding) SetNewGasConfig(gasCost *vmcommon.GasCost) {
	if gasCost == nil {
		return
	}

	ce.mutExecution.Lock()
	ce.baseOperationCost = gasCost.BaseOperationCost
	ce.mutExecution.Unlock()
}

// ProcessBuiltinFunction resolves the convert encoding function call
func (ce *ConvertEncoding) ProcessBuiltinFunction(
	_, _ vmcommon.UserAccountHandler,
	vmInput *vmcommon.ContractCallInput,
) (*vmcommon.VMOutput, error) {
	ce.mutExecution.RLock()
	defer ce.mutExecution.RUnlock()

	err := doCommonValidation(vmInput)
	if err != nil {
		return nil, err
	}

	signature := string(vmInput.Arguments[signaturePosition])
	inputData := vmInput.Arguments[dataStartPosition:]

	gasToUse := ce.calculateGasToUse(inputData)
	if vmInput.GasProvided < gasToUse {
		return nil, ErrNotEnoughGas
	}

	encodingHandler := convertEncoding.NewHandler(ce.accounts)
	outputData, err := ce.convertData(encodingHandler, signature, inputData)
	if err != nil {
		return nil, err
	}

	return &vmcommon.VMOutput{
		ReturnCode:   vmcommon.Ok,
		GasRemaining: vmInput.GasProvided - gasToUse,
		ReturnData:   outputData,
	}, nil
}

func (ce *ConvertEncoding) convertData(encodingHandler *convertEncoding.Handler, signature string, inputData [][]byte) ([][]byte, error) {
	switch ce.function {
	case core.BuiltInFunctionEthereumToMultiversXEncodingWithMultiversXSignature:
		return encodingHandler.EthToMvxEncodingWithMvxSignature(signature, inputData)
	case core.BuiltInFunctionEthereumToMultiversXEncodingWithEthereumSignature:
		return encodingHandler.EthToMvxEncodingWithEthSignature(signature, inputData)
	case core.BuiltInFunctionMultiversXToEthereumEncodingWithMultiversXSignature:
		return encodingHandler.MvxToEthEncodingWithMvxSignature(signature, inputData)
	case core.BuiltInFunctionMultiversXToEthereumEncodingWithEthereumSignature:
		return encodingHandler.MvxToEthEncodingWithEthSignature(signature, inputData)
	default:
		return nil, ErrInvalidArguments
	}
}

func (ce *ConvertEncoding) calculateGasToUse(inputData [][]byte) uint64 {
	totalLength := uint64(0)
	for _, arg := range inputData {
		totalLength += uint64(len(arg))
	}
	return totalLength * ce.baseOperationCost.CompilePerByte
}

func doCommonValidation(vmInput *vmcommon.ContractCallInput) error {
	if vmInput == nil {
		return ErrNilVmInput
	}
	if len(vmInput.Arguments) < minArgsCount {
		return ErrInvalidArguments
	}
	if vmInput.CallValue == nil {
		return ErrNilValue
	}
	if vmInput.CallValue.Cmp(zero) != 0 {
		return ErrBuiltInFunctionCalledWithValue
	}

	return nil
}

// IsInterfaceNil returns true if underlying object in nil
func (ce *ConvertEncoding) IsInterfaceNil() bool {
	return ce == nil
}
