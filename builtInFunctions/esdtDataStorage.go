package builtInFunctions

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/esdt"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/parsers"
)

const existsOnShard = byte(1)

type queryOptions struct {
	isCustomSystemAccountSet bool
	customSystemAccount      vmcommon.UserAccountHandler
}

func defaultQueryOptions() queryOptions {
	return queryOptions{}
}

type esdtDataStorage struct {
	accounts                      vmcommon.AccountsAdapter
	globalSettingsHandler         vmcommon.GlobalMetadataHandler
	marshaller                    vmcommon.Marshalizer
	keyPrefix                     []byte
	shardCoordinator              vmcommon.Coordinator
	txDataParser                  vmcommon.CallArgsParser
	enableEpochsHandler           vmcommon.EnableEpochsHandler
	crossChainTokenCheckerHandler CrossChainTokenCheckerHandler
}

// ArgsNewESDTDataStorage defines the argument list for new esdt data storage handler
type ArgsNewESDTDataStorage struct {
	Accounts                      vmcommon.AccountsAdapter
	GlobalSettingsHandler         vmcommon.GlobalMetadataHandler
	Marshalizer                   vmcommon.Marshalizer
	EnableEpochsHandler           vmcommon.EnableEpochsHandler
	ShardCoordinator              vmcommon.Coordinator
	CrossChainTokenCheckerHandler CrossChainTokenCheckerHandler
}

// NewESDTDataStorage creates a new esdt data storage handler
func NewESDTDataStorage(args ArgsNewESDTDataStorage) (*esdtDataStorage, error) {
	if check.IfNil(args.Accounts) {
		return nil, ErrNilAccountsAdapter
	}
	if check.IfNil(args.GlobalSettingsHandler) {
		return nil, ErrNilGlobalSettingsHandler
	}
	if check.IfNil(args.Marshalizer) {
		return nil, ErrNilMarshalizer
	}
	if check.IfNil(args.EnableEpochsHandler) {
		return nil, ErrNilEnableEpochsHandler
	}
	if check.IfNil(args.ShardCoordinator) {
		return nil, ErrNilShardCoordinator
	}
	if check.IfNil(args.CrossChainTokenCheckerHandler) {
		return nil, ErrNilCrossChainTokenChecker
	}

	e := &esdtDataStorage{
		accounts:                      args.Accounts,
		globalSettingsHandler:         args.GlobalSettingsHandler,
		marshaller:                    args.Marshalizer,
		keyPrefix:                     []byte(baseESDTKeyPrefix),
		shardCoordinator:              args.ShardCoordinator,
		txDataParser:                  parsers.NewCallArgsParser(),
		enableEpochsHandler:           args.EnableEpochsHandler,
		crossChainTokenCheckerHandler: args.CrossChainTokenCheckerHandler,
	}

	return e, nil
}

// GetESDTNFTTokenOnSender gets the nft token on sender account
func (e *esdtDataStorage) GetESDTNFTTokenOnSender(
	accnt vmcommon.UserAccountHandler,
	esdtTokenKey []byte,
	nonce uint64,
) (*esdt.ESDigitalToken, error) {
	esdtData, isNew, err := e.GetESDTNFTTokenOnDestination(accnt, esdtTokenKey, nonce)
	if err != nil {
		return nil, err
	}
	if isNew {
		return nil, ErrNewNFTDataOnSenderAddress
	}

	return esdtData, nil
}

// GetESDTNFTTokenOnDestination gets the nft token on destination account
func (e *esdtDataStorage) GetESDTNFTTokenOnDestination(
	accnt vmcommon.UserAccountHandler,
	esdtTokenKey []byte,
	nonce uint64,
) (*esdt.ESDigitalToken, bool, error) {
	return e.getESDTNFTTokenOnDestinationWithAccountsAdapterOptions(accnt, esdtTokenKey, nonce, defaultQueryOptions())
}

// GetESDTNFTTokenOnDestinationWithCustomSystemAccount gets the nft token on destination account by using a custom system account
func (e *esdtDataStorage) GetESDTNFTTokenOnDestinationWithCustomSystemAccount(
	accnt vmcommon.UserAccountHandler,
	esdtTokenKey []byte,
	nonce uint64,
	customSystemAccount vmcommon.UserAccountHandler,
) (*esdt.ESDigitalToken, bool, error) {
	if check.IfNil(customSystemAccount) {
		return nil, false, ErrNilUserAccount
	}

	queryOpts := queryOptions{
		isCustomSystemAccountSet: true,
		customSystemAccount:      customSystemAccount,
	}

	return e.getESDTNFTTokenOnDestinationWithAccountsAdapterOptions(accnt, esdtTokenKey, nonce, queryOpts)
}

// GetMetaDataFromSystemAccount gets the metadata from the system account
func (e *esdtDataStorage) GetMetaDataFromSystemAccount(esdtTokenKey []byte, nonce uint64) (*esdt.ESDigitalToken, error) {
	esdtNFTTokenKey := computeESDTNFTTokenKey(esdtTokenKey, nonce)

	esdtData, _, err := e.getESDTDigitalTokenDataFromSystemAccount(esdtNFTTokenKey, defaultQueryOptions())
	if err != nil {
		return nil, err
	}

	return esdtData, nil
}

func (e *esdtDataStorage) getESDTNFTTokenOnDestinationWithAccountsAdapterOptions(
	accnt vmcommon.UserAccountHandler,
	esdtTokenKey []byte,
	nonce uint64,
	options queryOptions,
) (*esdt.ESDigitalToken, bool, error) {
	esdtNFTTokenKey := computeESDTNFTTokenKey(esdtTokenKey, nonce)
	esdtData := &esdt.ESDigitalToken{
		Value: big.NewInt(0),
		Type:  uint32(core.Fungible),
	}
	marshaledData, _, err := accnt.AccountDataHandler().RetrieveValue(esdtNFTTokenKey)
	if core.IsGetNodeFromDBError(err) {
		return nil, false, err
	}
	if err != nil || len(marshaledData) == 0 {
		return esdtData, true, nil
	}

	err = e.marshaller.Unmarshal(esdtData, marshaledData)
	if err != nil {
		return nil, false, err
	}

	if !e.enableEpochsHandler.IsFlagEnabled(SaveToSystemAccountFlag) || nonce == 0 || metaDataOnUserAccount(esdtData.Type) {
		return esdtData, false, nil
	}

	retrievedEsdtData, _, err := e.getESDTDigitalTokenDataFromSystemAccount(esdtNFTTokenKey, options)
	if err != nil {
		return nil, false, err
	}
	if retrievedEsdtData != nil && retrievedEsdtData.TokenMetaData != nil {
		esdtData.TokenMetaData = retrievedEsdtData.TokenMetaData
		if core.IsDynamicESDT(esdtData.Type) {
			esdtData.Reserved = retrievedEsdtData.Reserved
		}
	}

	return esdtData, false, nil
}

func (e *esdtDataStorage) getESDTDigitalTokenDataFromSystemAccount(
	tokenKey []byte,
	options queryOptions,
) (*esdt.ESDigitalToken, vmcommon.UserAccountHandler, error) {
	systemAcc, err := e.getSystemAccount(options)
	if err != nil {
		return nil, nil, err
	}

	marshaledData, _, err := systemAcc.AccountDataHandler().RetrieveValue(tokenKey)
	if core.IsGetNodeFromDBError(err) {
		return nil, systemAcc, err
	}
	if err != nil || len(marshaledData) == 0 {
		return nil, systemAcc, nil
	}

	esdtData := &esdt.ESDigitalToken{}
	err = e.marshaller.Unmarshal(esdtData, marshaledData)
	if err != nil {
		return nil, nil, err
	}

	return esdtData, systemAcc, nil
}

func (e *esdtDataStorage) getESDTMetaDataFromSystemAccount(
	tokenKey []byte,
	options queryOptions,
) (*esdt.MetaData, error) {
	esdtData, _, err := e.getESDTDigitalTokenDataFromSystemAccount(tokenKey, options)
	if err != nil {
		return nil, err
	}
	if esdtData == nil {
		return nil, nil
	}

	return esdtData.TokenMetaData, nil
}

// CheckCollectionIsFrozenForAccount returns
func (e *esdtDataStorage) checkCollectionIsFrozenForAccount(
	accnt vmcommon.UserAccountHandler,
	esdtTokenKey []byte,
	nonce uint64,
	isReturnWithError bool,
) error {
	if !e.enableEpochsHandler.IsFlagEnabled(CheckFrozenCollectionFlag) {
		return nil
	}
	if nonce == 0 || isReturnWithError {
		return nil
	}

	esdtData := &esdt.ESDigitalToken{
		Value: big.NewInt(0),
		Type:  uint32(core.Fungible),
	}
	marshaledData, _, err := accnt.AccountDataHandler().RetrieveValue(esdtTokenKey)
	if core.IsGetNodeFromDBError(err) {
		return err
	}
	if err != nil || len(marshaledData) == 0 {
		return nil
	}

	err = e.marshaller.Unmarshal(esdtData, marshaledData)
	if err != nil {
		return err
	}

	esdtUserMetaData := ESDTUserMetadataFromBytes(esdtData.Properties)
	if esdtUserMetaData.Frozen {
		return ErrESDTIsFrozenForAccount
	}

	return nil
}

func (e *esdtDataStorage) checkFrozenPauseProperties(
	acnt vmcommon.UserAccountHandler,
	esdtTokenKey []byte,
	nonce uint64,
	esdtData *esdt.ESDigitalToken,
	isReturnWithError bool,
) error {
	err := checkFrozeAndPause(acnt.AddressBytes(), esdtTokenKey, esdtData, e.globalSettingsHandler, isReturnWithError)
	if err != nil {
		return err
	}

	esdtNFTTokenKey := computeESDTNFTTokenKey(esdtTokenKey, nonce)
	err = checkFrozeAndPause(acnt.AddressBytes(), esdtNFTTokenKey, esdtData, e.globalSettingsHandler, isReturnWithError)
	if err != nil {
		return err
	}

	return e.checkCollectionIsFrozenForAccount(acnt, esdtTokenKey, nonce, isReturnWithError)
}

// AddToLiquiditySystemAcc will increase/decrease the liquidity for ESDT Tokens on the metadata
func (e *esdtDataStorage) AddToLiquiditySystemAcc(
	esdtTokenKey []byte,
	tokenType uint32,
	nonce uint64,
	transferValue *big.Int,
	keepMetadataOnZeroLiquidity bool,
) error {
	isSaveToSystemAccountFlagEnabled := e.enableEpochsHandler.IsFlagEnabled(SaveToSystemAccountFlag)
	isSendAlwaysFlagEnabled := e.enableEpochsHandler.IsFlagEnabled(SendAlwaysFlag)
	if !isSaveToSystemAccountFlagEnabled || !isSendAlwaysFlagEnabled || nonce == 0 || metaDataOnUserAccount(tokenType) {
		return nil
	}

	esdtNFTTokenKey := computeESDTNFTTokenKey(esdtTokenKey, nonce)
	esdtData, systemAcc, err := e.getESDTDigitalTokenDataFromSystemAccount(esdtNFTTokenKey, defaultQueryOptions())
	if err != nil {
		return err
	}

	if esdtData == nil {
		return ErrNilESDTData
	}

	// old style metaData - nothing to do
	if len(esdtData.Reserved) == 0 {
		return nil
	}

	if e.enableEpochsHandler.IsFlagEnabled(FixOldTokenLiquidityFlag) {
		// old tokens which were transferred intra shard before the activation of this flag
		if esdtData.Value.Cmp(zero) == 0 && transferValue.Cmp(zero) < 0 {
			if !wasMetaDataUpdated(esdtData.Reserved) {
				esdtData.Reserved = nil
			}
			return e.marshalAndSaveData(systemAcc, esdtData, esdtNFTTokenKey)
		}
	}

	esdtData.Value.Add(esdtData.Value, transferValue)
	if esdtData.Value.Cmp(zero) < 0 {
		return ErrInvalidLiquidityForESDT
	}

	if esdtData.Value.Cmp(zero) == 0 && !keepMetadataOnZeroLiquidity {
		err = systemAcc.AccountDataHandler().SaveKeyValue(esdtNFTTokenKey, nil)
		if err != nil {
			return err
		}

		return e.accounts.SaveAccount(systemAcc)
	}

	return e.marshalAndSaveData(systemAcc, esdtData, esdtNFTTokenKey)
}

func (e *esdtDataStorage) shouldSaveMetadataInSystemAccount(esdtDataType uint32) bool {
	if !e.enableEpochsHandler.IsFlagEnabled(SaveToSystemAccountFlag) {
		return false
	}

	return !metaDataOnUserAccount(esdtDataType)
}

func metaDataOnUserAccount(esdtDataType uint32) bool {
	return esdtDataType == uint32(core.NonFungibleV2) ||
		esdtDataType == uint32(core.DynamicNFT) ||
		esdtDataType == uint32(core.Fungible)
}

// SaveMetaDataToSystemAccount saves the metadata to the system account
func (e *esdtDataStorage) SaveMetaDataToSystemAccount(tokenKey []byte, nonce uint64, esdtData *esdt.ESDigitalToken) error {
	systemAcc, err := e.loadSystemAccount()
	if err != nil {
		return err
	}

	esdtNFTTokenKey := computeESDTNFTTokenKey(tokenKey, nonce)
	return e.marshalAndSaveData(systemAcc, esdtData, esdtNFTTokenKey)
}

// SaveESDTNFTToken saves the nft token to the account and system account
func (e *esdtDataStorage) SaveESDTNFTToken(
	senderAddress []byte,
	acnt vmcommon.UserAccountHandler,
	esdtTokenKey []byte,
	nonce uint64,
	esdtData *esdt.ESDigitalToken,
	properties vmcommon.NftSaveArgs,
) ([]byte, error) {
	err := e.checkFrozenPauseProperties(acnt, esdtTokenKey, nonce, esdtData, properties.IsReturnWithError)
	if err != nil {
		return nil, err
	}

	esdtIdentifier, err := getEsdtIdentifierWithoutBaseKeyPrefix(esdtTokenKey)
	if err != nil {
		return nil, err
	}

	if !e.crossChainTokenCheckerHandler.IsCrossChainOperation(esdtIdentifier) {
		tokenTypeChanged, newTokenType, err := e.checkTokenTypeChanged(esdtTokenKey, esdtData)
		if err != nil {
			return nil, err
		}
		if tokenTypeChanged {
			esdtData.Type = newTokenType
			if newTokenType == uint32(core.NonFungibleV2) {
				err = e.removeMetaDataFromSystemAccount(esdtTokenKey, nonce)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	esdtNFTTokenKey := computeESDTNFTTokenKey(esdtTokenKey, nonce)
	senderShardID := e.shardCoordinator.ComputeId(senderAddress)
	if e.shouldSaveMetadataInSystemAccount(esdtData.Type) {
		err = e.saveESDTMetaDataToSystemAccount(acnt, senderShardID, esdtNFTTokenKey, nonce, esdtData, properties.MustUpdateAllFields)
		if err != nil {
			return nil, err
		}
	}

	if esdtData.Value.Cmp(zero) <= 0 && !properties.KeepMetaDataOnZeroLiquidity {
		return nil, acnt.AccountDataHandler().SaveKeyValue(esdtNFTTokenKey, nil)
	}

	if !e.enableEpochsHandler.IsFlagEnabled(SaveToSystemAccountFlag) || metaDataOnUserAccount(esdtData.Type) {
		marshaledData, errMarshal := e.marshaller.Marshal(esdtData)
		if errMarshal != nil {
			return nil, errMarshal
		}

		return marshaledData, acnt.AccountDataHandler().SaveKeyValue(esdtNFTTokenKey, marshaledData)
	}

	esdtDataOnAccount := &esdt.ESDigitalToken{
		Type:       esdtData.Type,
		Value:      esdtData.Value,
		Properties: esdtData.Properties,
	}
	marshaledData, err := e.marshaller.Marshal(esdtDataOnAccount)
	if err != nil {
		return nil, err
	}

	return marshaledData, acnt.AccountDataHandler().SaveKeyValue(esdtNFTTokenKey, marshaledData)
}

func getEsdtIdentifierWithoutBaseKeyPrefix(esdtTokenKey []byte) ([]byte, error) {
	if len(esdtTokenKey) < len(baseESDTKeyPrefix) {
		return nil, fmt.Errorf("invalid esdt token key %s", string(esdtTokenKey))
	}

	return esdtTokenKey[len(baseESDTKeyPrefix):], nil
}

func (e *esdtDataStorage) checkTokenTypeChanged(esdtTokenKey []byte, esdtData *esdt.ESDigitalToken) (bool, uint32, error) {
	if !e.enableEpochsHandler.IsFlagEnabled(DynamicEsdtFlag) {
		return false, 0, nil
	}

	if esdtData.Type == uint32(core.NonFungibleV2) ||
		esdtData.Type == uint32(core.DynamicNFT) ||
		esdtData.Type == uint32(core.DynamicSFT) ||
		esdtData.Type == uint32(core.DynamicMeta) ||
		esdtData.Type == uint32(core.Fungible) {

		return false, 0, nil
	}

	newTokenType, err := e.globalSettingsHandler.GetTokenType(esdtTokenKey)
	if err != nil {
		return false, 0, err
	}

	if esdtData.Type == newTokenType {
		return false, 0, nil
	}

	return true, newTokenType, nil
}

func (e *esdtDataStorage) removeMetaDataFromSystemAccount(esdtTokenKey []byte, nonce uint64) error {
	systemAcc, err := getSystemAccount(e.accounts)
	if err != nil {
		return err
	}

	esdtNFTTokenKey := computeESDTNFTTokenKey(esdtTokenKey, nonce)
	err = systemAcc.AccountDataHandler().SaveKeyValue(esdtNFTTokenKey, nil)
	if err != nil {
		return err
	}

	return e.accounts.SaveAccount(systemAcc)
}

func (e *esdtDataStorage) saveESDTMetaDataToSystemAccount(
	userAcc vmcommon.UserAccountHandler,
	senderShardID uint32,
	esdtNFTTokenKey []byte,
	nonce uint64,
	esdtData *esdt.ESDigitalToken,
	mustUpdateAllFields bool,
) error {
	if nonce == 0 {
		return nil
	}
	if esdtData.TokenMetaData == nil {
		return nil
	}

	systemAcc, err := e.getSystemAccount(defaultQueryOptions())
	if err != nil {
		return err
	}

	currentSaveData, _, err := systemAcc.AccountDataHandler().RetrieveValue(esdtNFTTokenKey)
	if core.IsGetNodeFromDBError(err) {
		return err
	}
	shouldUpdateCurrentSavedData, err := e.saveMetadataIfRequired(esdtNFTTokenKey, systemAcc, currentSaveData, esdtData)
	if err != nil {
		return err
	}

	if !mustUpdateAllFields && len(currentSaveData) > 0 && !shouldUpdateCurrentSavedData {
		return nil
	}

	esdtDataOnSystemAcc := &esdt.ESDigitalToken{
		Type:          esdtData.Type,
		Value:         big.NewInt(0),
		TokenMetaData: esdtData.TokenMetaData,
		Properties:    make([]byte, e.shardCoordinator.NumberOfShards()),
	}
	isSendAlwaysFlagEnabled := e.enableEpochsHandler.IsFlagEnabled(SendAlwaysFlag)
	if len(currentSaveData) == 0 && isSendAlwaysFlagEnabled {
		esdtDataOnSystemAcc.Properties = nil
		esdtDataOnSystemAcc.Reserved = []byte{1}

		err = e.setReservedToNilForOldToken(esdtDataOnSystemAcc, userAcc, esdtNFTTokenKey)
		if err != nil {
			return err
		}
	}
	if core.IsDynamicESDT(esdtData.Type) && len(esdtData.Reserved) != 0 {
		esdtDataOnSystemAcc.Reserved = esdtData.Reserved
	}

	if !isSendAlwaysFlagEnabled {
		selfID := e.shardCoordinator.SelfId()
		if selfID != core.MetachainShardId {
			esdtDataOnSystemAcc.Properties[selfID] = existsOnShard
		}
		if senderShardID != core.MetachainShardId {
			esdtDataOnSystemAcc.Properties[senderShardID] = existsOnShard
		}
	}

	return e.marshalAndSaveData(systemAcc, esdtDataOnSystemAcc, esdtNFTTokenKey)
}

func (e *esdtDataStorage) saveMetadataIfRequired(
	esdtNFTTokenKey []byte,
	systemAcc vmcommon.UserAccountHandler,
	currentSaveData []byte,
	esdtData *esdt.ESDigitalToken,
) (bool, error) {
	if !e.enableEpochsHandler.IsFlagEnabled(AlwaysSaveTokenMetaDataFlag) {
		return false, nil
	}
	if !e.enableEpochsHandler.IsFlagEnabled(SendAlwaysFlag) {
		// do not re-write the metadata if it is not sent, as it will cause data loss
		return false, nil
	}
	if len(currentSaveData) == 0 {
		// optimization: do not try to write here the token metadata, it will be written automatically by the next step
		return false, nil
	}

	esdtDataOnSystemAcc := &esdt.ESDigitalToken{}
	err := e.marshaller.Unmarshal(esdtDataOnSystemAcc, currentSaveData)
	if err != nil {
		return false, err
	}

	if wasMetaDataUpdated(esdtData.Reserved) {
		return true, nil
	}
	if len(esdtDataOnSystemAcc.Reserved) > 0 {
		return false, nil
	}

	esdtDataOnSystemAcc.TokenMetaData = esdtData.TokenMetaData
	return false, e.marshalAndSaveData(systemAcc, esdtDataOnSystemAcc, esdtNFTTokenKey)
}

func wasMetaDataUpdated(version []byte) bool {
	if version == nil || bytes.Equal(version, []byte{0}) || bytes.Equal(version, []byte{1}) {
		return false
	}

	return true
}

func (e *esdtDataStorage) setReservedToNilForOldToken(
	esdtDataOnSystemAcc *esdt.ESDigitalToken,
	userAcc vmcommon.UserAccountHandler,
	esdtNFTTokenKey []byte,
) error {
	if !e.enableEpochsHandler.IsFlagEnabled(FixOldTokenLiquidityFlag) {
		return nil
	}

	if check.IfNil(userAcc) {
		return ErrNilUserAccount
	}
	dataOnUserAcc, _, errNotCritical := userAcc.AccountDataHandler().RetrieveValue(esdtNFTTokenKey)
	if core.IsGetNodeFromDBError(errNotCritical) {
		return errNotCritical
	}
	shouldIgnoreToken := errNotCritical != nil || len(dataOnUserAcc) == 0
	if shouldIgnoreToken {
		return nil
	}

	esdtDataOnUserAcc := &esdt.ESDigitalToken{}
	err := e.marshaller.Unmarshal(esdtDataOnUserAcc, dataOnUserAcc)
	if err != nil {
		return err
	}

	// tokens which were last moved before flagOptimizeNFTStore keep the esdt metaData on the user account
	// these are not compatible with the new liquidity model,so we set the reserved field to nil
	if esdtDataOnUserAcc.TokenMetaData != nil {
		esdtDataOnSystemAcc.Reserved = nil
	}

	return nil
}

func (e *esdtDataStorage) marshalAndSaveData(
	systemAcc vmcommon.UserAccountHandler,
	esdtData *esdt.ESDigitalToken,
	esdtNFTTokenKey []byte,
) error {
	marshaledData, err := e.marshaller.Marshal(esdtData)
	if err != nil {
		return err
	}

	err = systemAcc.AccountDataHandler().SaveKeyValue(esdtNFTTokenKey, marshaledData)
	if err != nil {
		return err
	}

	return e.accounts.SaveAccount(systemAcc)
}

func (e *esdtDataStorage) getSystemAccount(options queryOptions) (vmcommon.UserAccountHandler, error) {
	if options.isCustomSystemAccountSet && !check.IfNil(options.customSystemAccount) {
		return options.customSystemAccount, nil
	}

	return e.loadSystemAccount()
}

func (e *esdtDataStorage) loadSystemAccount() (vmcommon.UserAccountHandler, error) {
	systemSCAccount, err := e.accounts.LoadAccount(vmcommon.SystemAccountAddress)
	if err != nil {
		return nil, err
	}

	userAcc, ok := systemSCAccount.(vmcommon.UserAccountHandler)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	return userAcc, nil
}

//TODO: merge properties in case of shard merge

// WasAlreadySentToDestinationShardAndUpdateState checks whether NFT metadata was sent to destination shard or not
// and saves the destination shard as sent
func (e *esdtDataStorage) WasAlreadySentToDestinationShardAndUpdateState(
	tickerID []byte,
	nonce uint64,
	dstAddress []byte,
) (bool, error) {
	if !e.enableEpochsHandler.IsFlagEnabled(SaveToSystemAccountFlag) {
		return false, nil
	}

	if nonce == 0 {
		return true, nil
	}
	dstShardID := e.shardCoordinator.ComputeId(dstAddress)
	if dstShardID == e.shardCoordinator.SelfId() {
		return true, nil
	}

	if e.enableEpochsHandler.IsFlagEnabled(SendAlwaysFlag) {
		return false, nil
	}

	if dstShardID == core.MetachainShardId {
		return true, nil
	}
	esdtTokenKey := append(e.keyPrefix, tickerID...)
	esdtNFTTokenKey := computeESDTNFTTokenKey(esdtTokenKey, nonce)

	esdtData, systemAcc, err := e.getESDTDigitalTokenDataFromSystemAccount(esdtNFTTokenKey, defaultQueryOptions())
	if err != nil {
		return false, err
	}
	if esdtData == nil {
		return false, nil
	}

	if uint32(len(esdtData.Properties)) < e.shardCoordinator.NumberOfShards() {
		newSlice := make([]byte, e.shardCoordinator.NumberOfShards())
		copy(newSlice, esdtData.Properties)
		esdtData.Properties = newSlice
	}

	if esdtData.Properties[dstShardID] > 0 {
		return true, nil
	}

	esdtData.Properties[dstShardID] = existsOnShard
	return false, e.marshalAndSaveData(systemAcc, esdtData, esdtNFTTokenKey)
}

// SaveNFTMetaData this saves the NFT metadata to the system account even if there was an error in processing
func (e *esdtDataStorage) SaveNFTMetaData(
	tx data.TransactionHandler,
) error {
	if !e.enableEpochsHandler.IsFlagEnabled(SaveToSystemAccountFlag) {
		return nil
	}
	if e.enableEpochsHandler.IsFlagEnabled(SendAlwaysFlag) {
		return nil
	}
	if check.IfNil(tx) {
		return ErrNilTransactionHandler
	}

	sndShardID := e.shardCoordinator.ComputeId(tx.GetSndAddr())
	dstShardID := e.shardCoordinator.ComputeId(tx.GetRcvAddr())
	isCrossShardTxAtDest := sndShardID != dstShardID && e.shardCoordinator.SelfId() == dstShardID
	if !isCrossShardTxAtDest {
		return nil
	}

	function, arguments, err := e.txDataParser.ParseData(string(tx.GetData()))
	if err != nil {
		return nil
	}
	if len(arguments) < 4 {
		return nil
	}

	switch function {
	case core.BuiltInFunctionESDTNFTTransfer:
		return e.addMetaDataToProperAccountFromNFTTransfer(tx.GetSndAddr(), sndShardID, arguments)
	case core.BuiltInFunctionMultiESDTNFTTransfer:
		return e.addMetaDataToProperAccountFromMultiTransfer(tx.GetSndAddr(), sndShardID, arguments)
	default:
		return nil
	}
}

func (e *esdtDataStorage) addMetaDataToProperAccountFromNFTTransfer(
	senderAddr []byte,
	sndShardID uint32,
	arguments [][]byte,
) error {
	if !bytes.Equal(arguments[3], zeroByteArray) {
		esdtTransferData := &esdt.ESDigitalToken{}
		err := e.marshaller.Unmarshal(esdtTransferData, arguments[3])
		if err != nil {
			return err
		}
		esdtTokenKey := append(e.keyPrefix, arguments[0]...)
		nonce := big.NewInt(0).SetBytes(arguments[1]).Uint64()
		esdtNFTTokenKey := computeESDTNFTTokenKey(esdtTokenKey, nonce)

		if metaDataOnUserAccount(esdtTransferData.Type) {
			return e.saveMetaDataToUserAccount(senderAddr, esdtNFTTokenKey, esdtTransferData)
		}

		return e.saveESDTMetaDataToSystemAccount(nil, sndShardID, esdtNFTTokenKey, nonce, esdtTransferData, true)
	}
	return nil
}

func (e *esdtDataStorage) saveMetaDataToUserAccount(address []byte, esdtNFTTokenKey []byte, esdtTransferData *esdt.ESDigitalToken) error {
	sndAcc, err := e.getAccount(address)
	if err != nil {
		return err
	}

	marshaledData, errMarshal := e.marshaller.Marshal(esdtTransferData)
	if errMarshal != nil {
		return errMarshal
	}

	return sndAcc.AccountDataHandler().SaveKeyValue(esdtNFTTokenKey, marshaledData)
}

func (e *esdtDataStorage) getAccount(address []byte) (vmcommon.UserAccountHandler, error) {
	acc, err := e.accounts.LoadAccount(address)
	if err != nil {
		return nil, err
	}

	userAcc, ok := acc.(vmcommon.UserAccountHandler)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	return userAcc, nil
}

func (e *esdtDataStorage) addMetaDataToProperAccountFromMultiTransfer(
	senderAddr []byte,
	sndShardID uint32,
	arguments [][]byte,
) error {
	numOfTransfers := big.NewInt(0).SetBytes(arguments[0]).Uint64()
	if numOfTransfers == 0 {
		return fmt.Errorf("%w, 0 tokens to transfer", ErrInvalidArguments)
	}
	minNumOfArguments := numOfTransfers*argumentsPerTransfer + 1
	if uint64(len(arguments)) < minNumOfArguments {
		return fmt.Errorf("%w, invalid number of arguments", ErrInvalidArguments)
	}

	startIndex := uint64(1)
	for i := uint64(0); i < numOfTransfers; i++ {
		tokenStartIndex := startIndex + i*argumentsPerTransfer
		tokenID := arguments[tokenStartIndex]
		nonce := big.NewInt(0).SetBytes(arguments[tokenStartIndex+1]).Uint64()

		if nonce > 0 && len(arguments[tokenStartIndex+2]) > vmcommon.MaxLengthForValueToOptTransfer {
			esdtTransferData := &esdt.ESDigitalToken{}
			marshaledNFTTransfer := arguments[tokenStartIndex+2]
			err := e.marshaller.Unmarshal(esdtTransferData, marshaledNFTTransfer)
			if err != nil {
				return fmt.Errorf("%w for token %s", err, string(tokenID))
			}

			esdtTokenKey := append(e.keyPrefix, tokenID...)
			esdtNFTTokenKey := computeESDTNFTTokenKey(esdtTokenKey, nonce)

			if metaDataOnUserAccount(esdtTransferData.Type) {
				err = e.saveMetaDataToUserAccount(senderAddr, esdtNFTTokenKey, esdtTransferData)
				if err != nil {
					return err
				}

				continue
			}

			err = e.saveESDTMetaDataToSystemAccount(nil, sndShardID, esdtNFTTokenKey, nonce, esdtTransferData, true)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// IsInterfaceNil returns true if underlying object in nil
func (e *esdtDataStorage) IsInterfaceNil() bool {
	return e == nil
}
