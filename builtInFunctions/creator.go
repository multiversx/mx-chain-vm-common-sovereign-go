package builtInFunctions

import (
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/mitchellh/mapstructure"
)

var _ vmcommon.BuiltInFunctionFactory = (*builtInFuncCreator)(nil)

// ArgsCreateBuiltInFunctionContainer defines the input arguments to create built in functions container
type ArgsCreateBuiltInFunctionContainer struct {
	GasMap                              map[string]map[string]uint64
	MapDNSAddresses                     map[string]struct{}
	EnableUserNameChange                bool
	Marshalizer                         vmcommon.Marshalizer
	Accounts                            vmcommon.AccountsAdapter
	ShardCoordinator                    vmcommon.Coordinator
	EpochNotifier                       vmcommon.EpochNotifier
	ESDTNFTImprovementV1ActivationEpoch uint32
	ESDTTransferRoleEnableEpoch         uint32
	GlobalMintBurnDisableEpoch          uint32
	ESDTTransferToMetaEnableEpoch       uint32
	NFTCreateMultiShardEnableEpoch      uint32
	SaveNFTToSystemAccountEnableEpoch   uint32
	CheckCorrectTokenIDEnableEpoch      uint32
	SendESDTMetadataAlwaysEnableEpoch   uint32
	CheckFunctionArgumentEnableEpoch    uint32
	FixAsyncCallbackCheckEnableEpoch    uint32
	FixOldTokenLiquidityEnableEpoch     uint32
	MaxNumOfAddressesForTransferRole    uint32
	ConfigAddress                       []byte
}

type builtInFuncCreator struct {
	mapDNSAddresses                     map[string]struct{}
	enableUserNameChange                bool
	marshaller                          vmcommon.Marshalizer
	accounts                            vmcommon.AccountsAdapter
	builtInFunctions                    vmcommon.BuiltInFunctionContainer
	gasConfig                           *vmcommon.GasCost
	shardCoordinator                    vmcommon.Coordinator
	epochNotifier                       vmcommon.EpochNotifier
	esdtStorageHandler                  vmcommon.ESDTNFTStorageHandler
	esdtGlobalSettingsHandler           vmcommon.ESDTGlobalSettingsHandler
	esdtNFTImprovementV1ActivationEpoch uint32
	esdtTransferRoleEnableEpoch         uint32
	globalMintBurnDisableEpoch          uint32
	esdtTransferToMetaEnableEpoch       uint32
	nftCreateMultiShardEnableEpoch      uint32
	saveNFTToSystemAccountEnableEpoch   uint32
	checkCorrectTokenIDEnableEpoch      uint32
	sendESDTMetadataAlwaysEnableEpoch   uint32
	checkFunctionArgumentEnableEpoch    uint32
	fixAsnycCallbackCheckEnableEpoch    uint32
	fixOldTokenLiquidityEnableEpoch     uint32
	maxNumOfAddressesForTransferRole    uint32
	configAddress                       []byte
}

// NewBuiltInFunctionsCreator creates a component which will instantiate the built in functions contracts
func NewBuiltInFunctionsCreator(args ArgsCreateBuiltInFunctionContainer) (*builtInFuncCreator, error) {
	if check.IfNil(args.Marshalizer) {
		return nil, ErrNilMarshalizer
	}
	if check.IfNil(args.Accounts) {
		return nil, ErrNilAccountsAdapter
	}
	if args.MapDNSAddresses == nil {
		return nil, ErrNilDnsAddresses
	}
	if check.IfNil(args.ShardCoordinator) {
		return nil, ErrNilShardCoordinator
	}
	if check.IfNil(args.EpochNotifier) {
		return nil, ErrNilEpochHandler
	}

	b := &builtInFuncCreator{
		mapDNSAddresses:                     args.MapDNSAddresses,
		enableUserNameChange:                args.EnableUserNameChange,
		marshaller:                          args.Marshalizer,
		accounts:                            args.Accounts,
		shardCoordinator:                    args.ShardCoordinator,
		epochNotifier:                       args.EpochNotifier,
		esdtNFTImprovementV1ActivationEpoch: args.ESDTNFTImprovementV1ActivationEpoch,
		esdtTransferRoleEnableEpoch:         args.ESDTTransferRoleEnableEpoch,
		globalMintBurnDisableEpoch:          args.GlobalMintBurnDisableEpoch,
		esdtTransferToMetaEnableEpoch:       args.ESDTTransferToMetaEnableEpoch,
		nftCreateMultiShardEnableEpoch:      args.NFTCreateMultiShardEnableEpoch,
		saveNFTToSystemAccountEnableEpoch:   args.SaveNFTToSystemAccountEnableEpoch,
		checkCorrectTokenIDEnableEpoch:      args.CheckCorrectTokenIDEnableEpoch,
		sendESDTMetadataAlwaysEnableEpoch:   args.SendESDTMetadataAlwaysEnableEpoch,
		checkFunctionArgumentEnableEpoch:    args.CheckFunctionArgumentEnableEpoch,
		maxNumOfAddressesForTransferRole:    args.MaxNumOfAddressesForTransferRole,
		configAddress:                       args.ConfigAddress,
	}

	var err error
	b.gasConfig, err = createGasConfig(args.GasMap)
	if err != nil {
		return nil, err
	}
	b.builtInFunctions = NewBuiltInFunctionContainer()

	return b, nil
}

// GasScheduleChange is called when gas schedule is changed, thus all contracts must be updated
func (b *builtInFuncCreator) GasScheduleChange(gasSchedule map[string]map[string]uint64) {
	newGasConfig, err := createGasConfig(gasSchedule)
	if err != nil {
		return
	}

	b.gasConfig = newGasConfig
	for key := range b.builtInFunctions.Keys() {
		builtInFunc, errGet := b.builtInFunctions.Get(key)
		if errGet != nil {
			return
		}

		builtInFunc.SetNewGasConfig(b.gasConfig)
	}
}

// NFTStorageHandler will return the esdt storage handler from the built in functions factory
func (b *builtInFuncCreator) NFTStorageHandler() vmcommon.SimpleESDTNFTStorageHandler {
	return b.esdtStorageHandler
}

// ESDTGlobalSettingsHandler will return the esdt global settings handler from the built in functions factory
func (b *builtInFuncCreator) ESDTGlobalSettingsHandler() vmcommon.ESDTGlobalSettingsHandler {
	return b.esdtGlobalSettingsHandler
}

// BuiltInFunctionContainer will return the built in function container
func (b *builtInFuncCreator) BuiltInFunctionContainer() vmcommon.BuiltInFunctionContainer {
	return b.builtInFunctions
}

// CreateBuiltInFunctionContainer will create the list of built-in functions
func (b *builtInFuncCreator) CreateBuiltInFunctionContainer() error {

	b.builtInFunctions = NewBuiltInFunctionContainer()
	var newFunc vmcommon.BuiltinFunction
	newFunc = NewClaimDeveloperRewardsFunc(b.gasConfig.BuiltInCost.ClaimDeveloperRewards)
	err := b.builtInFunctions.Add(core.BuiltInFunctionClaimDeveloperRewards, newFunc)
	if err != nil {
		return err
	}

	newFunc = NewChangeOwnerAddressFunc(b.gasConfig.BuiltInCost.ChangeOwnerAddress)
	err = b.builtInFunctions.Add(core.BuiltInFunctionChangeOwnerAddress, newFunc)
	if err != nil {
		return err
	}

	newFunc, err = NewSaveUserNameFunc(b.gasConfig.BuiltInCost.SaveUserName, b.mapDNSAddresses, b.enableUserNameChange)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(core.BuiltInFunctionSetUserName, newFunc)
	if err != nil {
		return err
	}

	newFunc, err = NewSaveKeyValueStorageFunc(b.gasConfig.BaseOperationCost, b.gasConfig.BuiltInCost.SaveKeyValue)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(core.BuiltInFunctionSaveKeyValue, newFunc)
	if err != nil {
		return err
	}

	globalSettingsFunc, err := NewESDTGlobalSettingsFunc(b.accounts, b.marshaller, true, core.BuiltInFunctionESDTPause, 0, b.epochNotifier)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(core.BuiltInFunctionESDTPause, globalSettingsFunc)
	if err != nil {
		return err
	}
	b.esdtGlobalSettingsHandler = globalSettingsFunc

	setRoleFunc, err := NewESDTRolesFunc(b.marshaller, true)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(core.BuiltInFunctionSetESDTRole, setRoleFunc)
	if err != nil {
		return err
	}

	newFunc, err = NewESDTTransferFunc(
		b.gasConfig.BuiltInCost.ESDTTransfer,
		b.marshaller,
		globalSettingsFunc,
		b.shardCoordinator,
		setRoleFunc,
		b.esdtTransferToMetaEnableEpoch,
		b.checkCorrectTokenIDEnableEpoch,
		b.epochNotifier,
	)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(core.BuiltInFunctionESDTTransfer, newFunc)
	if err != nil {
		return err
	}

	newFunc, err = NewESDTBurnFunc(b.gasConfig.BuiltInCost.ESDTBurn, b.marshaller, globalSettingsFunc, b.globalMintBurnDisableEpoch, b.epochNotifier)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(core.BuiltInFunctionESDTBurn, newFunc)
	if err != nil {
		return err
	}

	newFunc, err = NewESDTFreezeWipeFunc(b.marshaller, true, false)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(core.BuiltInFunctionESDTFreeze, newFunc)
	if err != nil {
		return err
	}

	newFunc, err = NewESDTFreezeWipeFunc(b.marshaller, false, false)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(core.BuiltInFunctionESDTUnFreeze, newFunc)
	if err != nil {
		return err
	}

	newFunc, err = NewESDTFreezeWipeFunc(b.marshaller, false, true)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(core.BuiltInFunctionESDTWipe, newFunc)
	if err != nil {
		return err
	}

	newFunc, err = NewESDTGlobalSettingsFunc(b.accounts, b.marshaller, false, core.BuiltInFunctionESDTUnPause, 0, b.epochNotifier)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(core.BuiltInFunctionESDTUnPause, newFunc)
	if err != nil {
		return err
	}

	newFunc, err = NewESDTRolesFunc(b.marshaller, false)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(core.BuiltInFunctionUnSetESDTRole, newFunc)
	if err != nil {
		return err
	}

	newFunc, err = NewESDTLocalBurnFunc(b.gasConfig.BuiltInCost.ESDTLocalBurn, b.marshaller, globalSettingsFunc, setRoleFunc)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(core.BuiltInFunctionESDTLocalBurn, newFunc)
	if err != nil {
		return err
	}

	newFunc, err = NewESDTLocalMintFunc(b.gasConfig.BuiltInCost.ESDTLocalMint, b.marshaller, globalSettingsFunc, setRoleFunc)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(core.BuiltInFunctionESDTLocalMint, newFunc)
	if err != nil {
		return err
	}

	args := ArgsNewESDTDataStorage{
		Accounts:                        b.accounts,
		GlobalSettingsHandler:           globalSettingsFunc,
		Marshalizer:                     b.marshaller,
		SaveToSystemEnableEpoch:         b.saveNFTToSystemAccountEnableEpoch,
		EpochNotifier:                   b.epochNotifier,
		ShardCoordinator:                b.shardCoordinator,
		SendAlwaysEnableEpoch:           b.sendESDTMetadataAlwaysEnableEpoch,
		FixOldTokenLiquidityEnableEpoch: b.fixOldTokenLiquidityEnableEpoch,
	}
	b.esdtStorageHandler, err = NewESDTDataStorage(args)
	if err != nil {
		return err
	}

	newFunc, err = NewESDTNFTAddQuantityFunc(b.gasConfig.BuiltInCost.ESDTNFTAddQuantity, b.esdtStorageHandler, globalSettingsFunc, setRoleFunc, b.saveNFTToSystemAccountEnableEpoch, b.epochNotifier)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(core.BuiltInFunctionESDTNFTAddQuantity, newFunc)
	if err != nil {
		return err
	}

	newFunc, err = NewESDTNFTBurnFunc(b.gasConfig.BuiltInCost.ESDTNFTBurn, b.esdtStorageHandler, globalSettingsFunc, setRoleFunc)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(core.BuiltInFunctionESDTNFTBurn, newFunc)
	if err != nil {
		return err
	}

	newFunc, err = NewESDTNFTCreateFunc(b.gasConfig.BuiltInCost.ESDTNFTCreate, b.gasConfig.BaseOperationCost, b.marshaller, globalSettingsFunc, setRoleFunc, b.esdtStorageHandler, b.accounts, b.saveNFTToSystemAccountEnableEpoch, b.epochNotifier)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(core.BuiltInFunctionESDTNFTCreate, newFunc)
	if err != nil {
		return err
	}

	newFunc, err = NewESDTNFTTransferFunc(
		b.gasConfig.BuiltInCost.ESDTNFTTransfer,
		b.marshaller,
		globalSettingsFunc,
		b.accounts,
		b.shardCoordinator,
		b.gasConfig.BaseOperationCost,
		setRoleFunc,
		b.esdtTransferToMetaEnableEpoch,
		b.saveNFTToSystemAccountEnableEpoch,
		b.checkCorrectTokenIDEnableEpoch,
		b.esdtStorageHandler,
		b.epochNotifier,
	)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(core.BuiltInFunctionESDTNFTTransfer, newFunc)
	if err != nil {
		return err
	}

	newFunc, err = NewESDTNFTCreateRoleTransfer(b.marshaller, b.accounts, b.shardCoordinator)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(core.BuiltInFunctionESDTNFTCreateRoleTransfer, newFunc)
	if err != nil {
		return err
	}

	newFunc, err = NewESDTNFTUpdateAttributesFunc(b.gasConfig.BuiltInCost.ESDTNFTUpdateAttributes, b.gasConfig.BaseOperationCost, b.esdtStorageHandler, globalSettingsFunc, setRoleFunc, b.esdtNFTImprovementV1ActivationEpoch, b.epochNotifier)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(core.BuiltInFunctionESDTNFTUpdateAttributes, newFunc)
	if err != nil {
		return err
	}

	newFunc, err = NewESDTNFTAddUriFunc(b.gasConfig.BuiltInCost.ESDTNFTAddURI, b.gasConfig.BaseOperationCost, b.esdtStorageHandler, globalSettingsFunc, setRoleFunc, b.esdtNFTImprovementV1ActivationEpoch, b.epochNotifier)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(core.BuiltInFunctionESDTNFTAddURI, newFunc)
	if err != nil {
		return err
	}

	newFunc, err = NewESDTNFTMultiTransferFunc(
		b.gasConfig.BuiltInCost.ESDTNFTMultiTransfer,
		b.marshaller,
		globalSettingsFunc,
		b.accounts,
		b.shardCoordinator,
		b.gasConfig.BaseOperationCost,
		b.esdtNFTImprovementV1ActivationEpoch,
		b.epochNotifier,
		setRoleFunc,
		b.esdtTransferToMetaEnableEpoch,
		b.checkCorrectTokenIDEnableEpoch,
		b.esdtStorageHandler,
	)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(core.BuiltInFunctionMultiESDTNFTTransfer, newFunc)
	if err != nil {
		return err
	}

	newFunc, err = NewESDTGlobalSettingsFunc(b.accounts, b.marshaller, true, core.BuiltInFunctionESDTSetLimitedTransfer, b.esdtTransferRoleEnableEpoch, b.epochNotifier)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(core.BuiltInFunctionESDTSetLimitedTransfer, newFunc)
	if err != nil {
		return err
	}

	newFunc, err = NewESDTGlobalSettingsFunc(b.accounts, b.marshaller, false, core.BuiltInFunctionESDTUnSetLimitedTransfer, b.esdtTransferRoleEnableEpoch, b.epochNotifier)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(core.BuiltInFunctionESDTUnSetLimitedTransfer, newFunc)
	if err != nil {
		return err
	}

	argsNewDeleteFunc := ArgsNewESDTDeleteMetadata{
		FuncGasCost:     b.gasConfig.BuiltInCost.ESDTNFTBurn,
		Marshalizer:     b.marshaller,
		Accounts:        b.accounts,
		ActivationEpoch: b.sendESDTMetadataAlwaysEnableEpoch,
		EpochNotifier:   b.epochNotifier,
		AllowedAddress:  b.configAddress,
		Delete:          true,
	}
	newFunc, err = NewESDTDeleteMetadataFunc(argsNewDeleteFunc)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(vmcommon.ESDTDeleteMetadata, newFunc)
	if err != nil {
		return err
	}

	argsNewDeleteFunc.Delete = false
	newFunc, err = NewESDTDeleteMetadataFunc(argsNewDeleteFunc)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(vmcommon.ESDTAddMetadata, newFunc)
	if err != nil {
		return err
	}

	newFunc, err = NewESDTGlobalSettingsFunc(b.accounts, b.marshaller, true, vmcommon.BuiltInFunctionESDTSetBurnRoleForAll, b.sendESDTMetadataAlwaysEnableEpoch, b.epochNotifier)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(vmcommon.BuiltInFunctionESDTSetBurnRoleForAll, newFunc)
	if err != nil {
		return err
	}

	newFunc, err = NewESDTGlobalSettingsFunc(b.accounts, b.marshaller, false, vmcommon.BuiltInFunctionESDTUnSetBurnRoleForAll, b.sendESDTMetadataAlwaysEnableEpoch, b.epochNotifier)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(vmcommon.BuiltInFunctionESDTUnSetBurnRoleForAll, newFunc)
	if err != nil {
		return err
	}

	newFunc, err = NewESDTTransferRoleAddressFunc(b.accounts, b.marshaller, b.sendESDTMetadataAlwaysEnableEpoch, b.epochNotifier, b.maxNumOfAddressesForTransferRole, false)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(vmcommon.BuiltInFunctionESDTTransferRoleDeleteAddress, newFunc)
	if err != nil {
		return err
	}

	newFunc, err = NewESDTTransferRoleAddressFunc(b.accounts, b.marshaller, b.sendESDTMetadataAlwaysEnableEpoch, b.epochNotifier, b.maxNumOfAddressesForTransferRole, true)
	if err != nil {
		return err
	}
	err = b.builtInFunctions.Add(vmcommon.BuiltInFunctionESDTTransferRoleAddAddress, newFunc)
	if err != nil {
		return err
	}

	return nil
}

func createGasConfig(gasMap map[string]map[string]uint64) (*vmcommon.GasCost, error) {
	baseOps := &vmcommon.BaseOperationCost{}
	err := mapstructure.Decode(gasMap[core.BaseOperationCostString], baseOps)
	if err != nil {
		return nil, err
	}

	err = check.ForZeroUintFields(*baseOps)
	if err != nil {
		return nil, err
	}

	builtInOps := &vmcommon.BuiltInCost{}
	err = mapstructure.Decode(gasMap[core.BuiltInCostString], builtInOps)
	if err != nil {
		return nil, err
	}

	err = check.ForZeroUintFields(*builtInOps)
	if err != nil {
		return nil, err
	}

	gasCost := vmcommon.GasCost{
		BaseOperationCost: *baseOps,
		BuiltInCost:       *builtInOps,
	}

	return &gasCost, nil
}

// SetPayableHandler sets the payableCheck interface to the needed functions
func (b *builtInFuncCreator) SetPayableHandler(payableHandler vmcommon.PayableHandler) error {
	payableChecker, err := NewPayableCheckFunc(
		payableHandler,
		b.checkFunctionArgumentEnableEpoch,
		b.fixAsnycCallbackCheckEnableEpoch,
		b.epochNotifier,
	)
	if err != nil {
		return err
	}

	listOfTransferFunc := []string{
		core.BuiltInFunctionMultiESDTNFTTransfer,
		core.BuiltInFunctionESDTNFTTransfer,
		core.BuiltInFunctionESDTTransfer}

	for _, transferFunc := range listOfTransferFunc {
		builtInFunc, err := b.builtInFunctions.Get(transferFunc)
		if err != nil {
			return err
		}

		esdtTransferFunc, ok := builtInFunc.(vmcommon.AcceptPayableChecker)
		if !ok {
			return ErrWrongTypeAssertion
		}

		err = esdtTransferFunc.SetPayableChecker(payableChecker)
		if err != nil {
			return err
		}
	}

	return nil
}

// IsInterfaceNil returns true if underlying object is nil
func (b *builtInFuncCreator) IsInterfaceNil() bool {
	return b == nil
}
