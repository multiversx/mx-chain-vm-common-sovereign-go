package builtInFunctions

import (
	"math/big"
	"sync"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-vm-common-go"
)

type esdtNFTAddUri struct {
	baseActiveHandler
	vmcommon.BlockchainDataProvider
	keyPrefix             []byte
	esdtStorageHandler    vmcommon.ESDTNFTStorageHandler
	globalSettingsHandler vmcommon.ESDTGlobalSettingsHandler
	rolesHandler          vmcommon.ESDTRoleHandler
	gasConfig             vmcommon.BaseOperationCost
	enableEpochsHandler   vmcommon.EnableEpochsHandler
	funcGasCost           uint64
	marshaller            marshal.Marshalizer
	mutExecution          sync.RWMutex
}

// NewESDTNFTAddUriFunc returns the esdt NFT add URI built-in function component
func NewESDTNFTAddUriFunc(
	funcGasCost uint64,
	gasConfig vmcommon.BaseOperationCost,
	esdtStorageHandler vmcommon.ESDTNFTStorageHandler,
	globalSettingsHandler vmcommon.ESDTGlobalSettingsHandler,
	rolesHandler vmcommon.ESDTRoleHandler,
	enableEpochsHandler vmcommon.EnableEpochsHandler,
	marshaller marshal.Marshalizer,
) (*esdtNFTAddUri, error) {
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

	e := &esdtNFTAddUri{
		keyPrefix:              []byte(baseESDTKeyPrefix),
		esdtStorageHandler:     esdtStorageHandler,
		funcGasCost:            funcGasCost,
		mutExecution:           sync.RWMutex{},
		globalSettingsHandler:  globalSettingsHandler,
		gasConfig:              gasConfig,
		rolesHandler:           rolesHandler,
		enableEpochsHandler:    enableEpochsHandler,
		BlockchainDataProvider: NewBlockchainDataProvider(),
		marshaller:             marshaller,
	}

	e.baseActiveHandler.activeHandler = func() bool {
		return enableEpochsHandler.IsFlagEnabled(ESDTNFTImprovementV1Flag)
	}

	return e, nil
}

// SetNewGasConfig is called whenever gas cost is changed
func (e *esdtNFTAddUri) SetNewGasConfig(gasCost *vmcommon.GasCost) {
	if gasCost == nil {
		return
	}

	e.mutExecution.Lock()
	e.funcGasCost = gasCost.BuiltInCost.ESDTNFTAddURI
	e.gasConfig = gasCost.BaseOperationCost
	e.mutExecution.Unlock()
}

// ProcessBuiltinFunction resolves ESDT NFT add uris function call
// Requires 3 arguments:
// arg0 - token identifier
// arg1 - nonce
// arg[2:] - uris to add
func (e *esdtNFTAddUri) ProcessBuiltinFunction(
	acntSnd, _ vmcommon.UserAccountHandler,
	vmInput *vmcommon.ContractCallInput,
) (*vmcommon.VMOutput, error) {
	e.mutExecution.RLock()
	defer e.mutExecution.RUnlock()

	err := checkESDTNFTCreateBurnAddInput(acntSnd, vmInput, e.funcGasCost)
	if err != nil {
		return nil, err
	}
	if len(vmInput.Arguments) < 3 {
		return nil, ErrInvalidArguments
	}

	err = e.rolesHandler.CheckAllowedToExecute(acntSnd, vmInput.Arguments[0], []byte(core.ESDTRoleNFTAddURI))
	if err != nil {
		return nil, err
	}

	gasCostForStore := e.getGasCostForURIStore(vmInput)
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

	esdtData.TokenMetaData.URIs = append(esdtData.TokenMetaData.URIs, vmInput.Arguments[2:]...)
	metaDataVersion.URIs = e.CurrentRound()

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

	extraTopics := append([][]byte{vmInput.CallerAddr}, vmInput.Arguments[2:]...)
	addESDTEntryInVMOutput(vmOutput, []byte(core.BuiltInFunctionESDTNFTAddURI), vmInput.Arguments[0], nonce, big.NewInt(0), extraTopics...)

	return vmOutput, nil
}

func (e *esdtNFTAddUri) getGasCostForURIStore(vmInput *vmcommon.ContractCallInput) uint64 {
	lenURIs := 0
	for _, uri := range vmInput.Arguments[2:] {
		lenURIs += len(uri)
	}
	return uint64(lenURIs) * e.gasConfig.StorePerByte
}

// IsInterfaceNil returns true if underlying object in nil
func (e *esdtNFTAddUri) IsInterfaceNil() bool {
	return e == nil
}
