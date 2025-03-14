package builtInFunctions

import (
	"math/big"
	"sync"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-vm-common-go"
)

type esdtNFTupdate struct {
	baseActiveHandler
	vmcommon.BlockchainDataProvider
	keyPrefix             []byte
	esdtStorageHandler    vmcommon.ESDTNFTStorageHandler
	globalSettingsHandler vmcommon.ESDTGlobalSettingsHandler
	rolesHandler          vmcommon.ESDTRoleHandler
	gasConfig             vmcommon.BaseOperationCost
	funcGasCost           uint64
	marshaller            marshal.Marshalizer
	enableEpochsHandler   vmcommon.EnableEpochsHandler
	mutExecution          sync.RWMutex
}

// NewESDTNFTUpdateAttributesFunc returns the esdt NFT update attribute built-in function component
func NewESDTNFTUpdateAttributesFunc(
	funcGasCost uint64,
	gasConfig vmcommon.BaseOperationCost,
	esdtStorageHandler vmcommon.ESDTNFTStorageHandler,
	globalSettingsHandler vmcommon.ESDTGlobalSettingsHandler,
	rolesHandler vmcommon.ESDTRoleHandler,
	enableEpochsHandler vmcommon.EnableEpochsHandler,
	marshaller marshal.Marshalizer,
) (*esdtNFTupdate, error) {
	if check.IfNil(esdtStorageHandler) {
		return nil, ErrNilESDTNFTStorageHandler
	}
	if check.IfNil(globalSettingsHandler) {
		return nil, ErrNilGlobalSettingsHandler
	}
	if check.IfNil(rolesHandler) {
		return nil, ErrNilRolesHandler
	}
	if check.IfNil(enableEpochsHandler) {
		return nil, ErrNilEnableEpochsHandler
	}
	if check.IfNil(marshaller) {
		return nil, ErrNilMarshalizer
	}

	e := &esdtNFTupdate{
		keyPrefix:              []byte(baseESDTKeyPrefix),
		esdtStorageHandler:     esdtStorageHandler,
		funcGasCost:            funcGasCost,
		mutExecution:           sync.RWMutex{},
		globalSettingsHandler:  globalSettingsHandler,
		gasConfig:              gasConfig,
		rolesHandler:           rolesHandler,
		BlockchainDataProvider: NewBlockchainDataProvider(),
		marshaller:             marshaller,
		enableEpochsHandler:    enableEpochsHandler,
	}

	e.baseActiveHandler.activeHandler = func() bool {
		return enableEpochsHandler.IsFlagEnabled(ESDTNFTImprovementV1Flag)
	}

	return e, nil
}

// SetNewGasConfig is called whenever gas cost is changed
func (e *esdtNFTupdate) SetNewGasConfig(gasCost *vmcommon.GasCost) {
	if gasCost == nil {
		return
	}

	e.mutExecution.Lock()
	e.funcGasCost = gasCost.BuiltInCost.ESDTNFTUpdateAttributes
	e.gasConfig = gasCost.BaseOperationCost
	e.mutExecution.Unlock()
}

// ProcessBuiltinFunction resolves ESDT NFT update attributes function call
// Requires 3 arguments:
// arg0 - token identifier
// arg1 - nonce
// arg2 - new attributes
func (e *esdtNFTupdate) ProcessBuiltinFunction(
	acntSnd, _ vmcommon.UserAccountHandler,
	vmInput *vmcommon.ContractCallInput,
) (*vmcommon.VMOutput, error) {
	e.mutExecution.RLock()
	defer e.mutExecution.RUnlock()

	err := checkESDTNFTCreateBurnAddInput(acntSnd, vmInput, e.funcGasCost)
	if err != nil {
		return nil, err
	}
	if len(vmInput.Arguments) != 3 {
		return nil, ErrInvalidArguments
	}

	err = e.rolesHandler.CheckAllowedToExecute(acntSnd, vmInput.Arguments[0], []byte(core.ESDTRoleNFTUpdateAttributes))
	if err != nil {
		return nil, err
	}

	gasCostForStore := uint64(len(vmInput.Arguments[2])) * e.gasConfig.StorePerByte
	if vmInput.GasProvided < e.funcGasCost+gasCostForStore {
		return nil, ErrNotEnoughGas
	}

	esdtTokenKey := append(e.keyPrefix, vmInput.Arguments[0]...)
	nonce := big.NewInt(0).SetBytes(vmInput.Arguments[1]).Uint64()
	if nonce == 0 {
		return nil, ErrNFTDoesNotHaveMetadata
	}
	esdtData, err := e.esdtStorageHandler.GetESDTNFTTokenOnSender(acntSnd, esdtTokenKey, nonce)
	if err != nil {
		return nil, err
	}

	metaDataVersion, _, err := getMetaDataVersion(esdtData, e.enableEpochsHandler, e.marshaller)
	if err != nil {
		return nil, err
	}

	esdtData.TokenMetaData.Attributes = vmInput.Arguments[2]
	metaDataVersion.Attributes = e.CurrentRound()

	err = changeEsdtVersion(esdtData, metaDataVersion, e.enableEpochsHandler, e.marshaller)
	if err != nil {
		return nil, err
	}

	properties := vmcommon.NftSaveArgs{
		MustUpdateAllFields:         true,
		IsReturnWithError:           vmInput.ReturnCallAfterError,
		KeepMetaDataOnZeroLiquidity: true,
	}
	_, err = e.esdtStorageHandler.SaveESDTNFTToken(acntSnd.AddressBytes(), acntSnd, esdtTokenKey, nonce, esdtData, properties)
	if err != nil {
		return nil, err
	}

	vmOutput := &vmcommon.VMOutput{
		ReturnCode:   vmcommon.Ok,
		GasRemaining: vmInput.GasProvided - e.funcGasCost - gasCostForStore,
	}

	addESDTEntryInVMOutput(vmOutput, []byte(core.BuiltInFunctionESDTNFTUpdateAttributes), vmInput.Arguments[0], nonce, big.NewInt(0), vmInput.CallerAddr, vmInput.Arguments[2])

	return vmOutput, nil
}

// IsInterfaceNil returns true if underlying object in nil
func (e *esdtNFTupdate) IsInterfaceNil() bool {
	return e == nil
}
