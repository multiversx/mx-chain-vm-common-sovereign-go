package builtInFunctions

import (
	"math/big"
	"sync"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/marshal"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

type esdtModifyCreator struct {
	baseActiveHandler
	vmcommon.BlockchainDataProvider
	globalSettingsHandler vmcommon.GlobalMetadataHandler
	storageHandler        vmcommon.ESDTNFTStorageHandler
	rolesHandler          vmcommon.ESDTRoleHandler
	accounts              vmcommon.AccountsAdapter
	enableEpochsHandler   vmcommon.EnableEpochsHandler
	funcGasCost           uint64
	marshaller            marshal.Marshalizer
	mutExecution          sync.RWMutex
}

// NewESDTModifyCreatorFunc returns the esdt modify creator built-in function component
func NewESDTModifyCreatorFunc(
	funcGasCost uint64,
	accounts vmcommon.AccountsAdapter,
	globalSettingsHandler vmcommon.GlobalMetadataHandler,
	storageHandler vmcommon.ESDTNFTStorageHandler,
	rolesHandler vmcommon.ESDTRoleHandler,
	enableEpochsHandler vmcommon.EnableEpochsHandler,
	marshaller marshal.Marshalizer,
) (*esdtModifyCreator, error) {
	if check.IfNil(accounts) {
		return nil, ErrNilAccountsAdapter
	}
	if check.IfNil(globalSettingsHandler) {
		return nil, ErrNilGlobalSettingsHandler
	}
	if check.IfNil(enableEpochsHandler) {
		return nil, ErrNilEnableEpochsHandler
	}
	if check.IfNil(storageHandler) {
		return nil, ErrNilESDTNFTStorageHandler
	}
	if check.IfNil(rolesHandler) {
		return nil, ErrNilRolesHandler
	}
	if check.IfNil(marshaller) {
		return nil, ErrNilMarshalizer
	}

	e := &esdtModifyCreator{
		accounts:               accounts,
		globalSettingsHandler:  globalSettingsHandler,
		storageHandler:         storageHandler,
		rolesHandler:           rolesHandler,
		funcGasCost:            funcGasCost,
		enableEpochsHandler:    enableEpochsHandler,
		mutExecution:           sync.RWMutex{},
		BlockchainDataProvider: NewBlockchainDataProvider(),
		marshaller:             marshaller,
	}

	e.baseActiveHandler.activeHandler = func() bool {
		return enableEpochsHandler.IsFlagEnabled(DynamicEsdtFlag)
	}

	return e, nil
}

// ProcessBuiltinFunction saves the token type in the system account
func (e *esdtModifyCreator) ProcessBuiltinFunction(acntSnd, _ vmcommon.UserAccountHandler, vmInput *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	err := checkUpdateArguments(vmInput, acntSnd, e.baseActiveHandler, 2, e.rolesHandler, core.ESDTRoleModifyCreator)
	if err != nil {
		return nil, err
	}

	e.mutExecution.RLock()
	funcGasCost := e.funcGasCost
	e.mutExecution.RUnlock()

	if vmInput.GasProvided < funcGasCost {
		return nil, ErrNotEnoughGas
	}

	esdtInfo, err := getEsdtInfo(vmInput, acntSnd, e.storageHandler, e.globalSettingsHandler)
	if err != nil {
		return nil, err
	}

	metaDataVersion, _, err := getMetaDataVersion(esdtInfo.esdtData, e.enableEpochsHandler, e.marshaller)
	if err != nil {
		return nil, err
	}

	esdtInfo.esdtData.TokenMetaData.Creator = vmInput.CallerAddr
	metaDataVersion.Creator = e.CurrentRound()

	err = changeEsdtVersion(esdtInfo.esdtData, metaDataVersion, e.enableEpochsHandler, e.marshaller)
	if err != nil {
		return nil, err
	}

	err = saveESDTMetaDataInfo(esdtInfo, e.storageHandler, acntSnd, vmInput.ReturnCallAfterError)
	if err != nil {
		return nil, err
	}

	vmOutput := &vmcommon.VMOutput{
		ReturnCode:   vmcommon.Ok,
		GasRemaining: vmInput.GasProvided - funcGasCost,
	}

	addESDTEntryInVMOutput(vmOutput, []byte(core.ESDTModifyCreator), vmInput.Arguments[tokenIDIndex], esdtInfo.esdtData.TokenMetaData.Nonce, big.NewInt(0), [][]byte{vmInput.CallerAddr}...)

	return vmOutput, nil
}

// SetNewGasConfig is called whenever gas cost is changed
func (e *esdtModifyCreator) SetNewGasConfig(gasCost *vmcommon.GasCost) {
	if gasCost == nil {
		return
	}

	e.mutExecution.Lock()
	e.funcGasCost = gasCost.BuiltInCost.ESDTModifyCreator
	e.mutExecution.Unlock()
}

// IsInterfaceNil returns true if there is no value under the interface
func (e *esdtModifyCreator) IsInterfaceNil() bool {
	return e == nil
}
