package builtInFunctions

import (
	"errors"
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/esdt"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewESDTModifyRoyaltiesFunc(t *testing.T) {
	t.Parallel()

	t.Run("nil accounts adapter", func(t *testing.T) {
		t.Parallel()

		e, err := NewESDTModifyRoyaltiesFunc(0, nil, nil, nil, nil, nil)
		assert.Nil(t, e)
		assert.Equal(t, ErrNilAccountsAdapter, err)
	})
	t.Run("nil global settings handler", func(t *testing.T) {
		t.Parallel()

		e, err := NewESDTModifyRoyaltiesFunc(0, &mock.AccountsStub{}, nil, nil, nil, nil)
		assert.Nil(t, e)
		assert.Equal(t, ErrNilGlobalSettingsHandler, err)
	})
	t.Run("nil enable epochs handler", func(t *testing.T) {
		t.Parallel()

		e, err := NewESDTModifyRoyaltiesFunc(0, &mock.AccountsStub{}, &mock.GlobalSettingsHandlerStub{}, nil, nil, nil)
		assert.Nil(t, e)
		assert.Equal(t, ErrNilEnableEpochsHandler, err)
	})
	t.Run("nil storage handler", func(t *testing.T) {
		t.Parallel()

		e, err := NewESDTModifyRoyaltiesFunc(0, &mock.AccountsStub{}, &mock.GlobalSettingsHandlerStub{}, nil, nil, &mock.EnableEpochsHandlerStub{})
		assert.Nil(t, e)
		assert.Equal(t, ErrNilESDTNFTStorageHandler, err)
	})
	t.Run("nil roles handler", func(t *testing.T) {
		t.Parallel()

		e, err := NewESDTModifyRoyaltiesFunc(0, &mock.AccountsStub{}, &mock.GlobalSettingsHandlerStub{}, &mock.ESDTNFTStorageHandlerStub{}, nil, &mock.EnableEpochsHandlerStub{})
		assert.Nil(t, e)
		assert.Equal(t, ErrNilRolesHandler, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		funcGasCost := uint64(10)
		e, err := NewESDTModifyRoyaltiesFunc(funcGasCost, &mock.AccountsStub{}, &mock.GlobalSettingsHandlerStub{}, &mock.ESDTNFTStorageHandlerStub{}, &mock.ESDTRoleHandlerStub{}, &mock.EnableEpochsHandlerStub{})
		assert.NotNil(t, e)
		assert.Nil(t, err)
		assert.Equal(t, funcGasCost, e.funcGasCost)
	})
}

func TestESDTModifyRoyalties_ProcessBuiltinFunction(t *testing.T) {
	t.Parallel()

	t.Run("nil vmInput", func(t *testing.T) {
		t.Parallel()

		e, _ := NewESDTModifyRoyaltiesFunc(0, &mock.AccountsStub{}, &mock.GlobalSettingsHandlerStub{}, &mock.ESDTNFTStorageHandlerStub{}, &mock.ESDTRoleHandlerStub{}, &mock.EnableEpochsHandlerStub{})
		vmOutput, err := e.ProcessBuiltinFunction(nil, nil, nil)
		assert.Nil(t, vmOutput)
		assert.Equal(t, ErrNilVmInput, err)
	})
	t.Run("nil CallValue", func(t *testing.T) {
		t.Parallel()

		e, _ := NewESDTModifyRoyaltiesFunc(0, &mock.AccountsStub{}, &mock.GlobalSettingsHandlerStub{}, &mock.ESDTNFTStorageHandlerStub{}, &mock.ESDTRoleHandlerStub{}, &mock.EnableEpochsHandlerStub{})
		vmInput := &vmcommon.ContractCallInput{
			VMInput: vmcommon.VMInput{
				CallValue: nil,
			},
		}
		vmOutput, err := e.ProcessBuiltinFunction(nil, nil, vmInput)
		assert.Nil(t, vmOutput)
		assert.Equal(t, ErrNilValue, err)
	})
	t.Run("call value not zero", func(t *testing.T) {
		t.Parallel()

		e, _ := NewESDTModifyRoyaltiesFunc(0, &mock.AccountsStub{}, &mock.GlobalSettingsHandlerStub{}, &mock.ESDTNFTStorageHandlerStub{}, &mock.ESDTRoleHandlerStub{}, &mock.EnableEpochsHandlerStub{})
		vmInput := &vmcommon.ContractCallInput{
			VMInput: vmcommon.VMInput{
				CallValue: big.NewInt(10),
			},
		}
		vmOutput, err := e.ProcessBuiltinFunction(nil, nil, vmInput)
		assert.Nil(t, vmOutput)
		assert.Equal(t, ErrBuiltInFunctionCalledWithValue, err)
	})
	t.Run("recipient address is not caller address", func(t *testing.T) {
		t.Parallel()

		e, _ := NewESDTModifyRoyaltiesFunc(0, &mock.AccountsStub{}, &mock.GlobalSettingsHandlerStub{}, &mock.ESDTNFTStorageHandlerStub{}, &mock.ESDTRoleHandlerStub{}, &mock.EnableEpochsHandlerStub{})
		vmInput := &vmcommon.ContractCallInput{
			VMInput: vmcommon.VMInput{
				CallValue:  big.NewInt(0),
				CallerAddr: []byte("caller"),
			},
			RecipientAddr: []byte("recipient"),
		}
		vmOutput, err := e.ProcessBuiltinFunction(nil, nil, vmInput)
		assert.Nil(t, vmOutput)
		assert.Equal(t, ErrInvalidRcvAddr, err)
	})
	t.Run("nil sender account", func(t *testing.T) {
		t.Parallel()

		e, _ := NewESDTModifyRoyaltiesFunc(0, &mock.AccountsStub{}, &mock.GlobalSettingsHandlerStub{}, &mock.ESDTNFTStorageHandlerStub{}, &mock.ESDTRoleHandlerStub{}, &mock.EnableEpochsHandlerStub{})
		vmInput := &vmcommon.ContractCallInput{
			VMInput: vmcommon.VMInput{
				CallValue:  big.NewInt(0),
				CallerAddr: []byte("caller"),
			},
			RecipientAddr: []byte("caller"),
		}
		vmOutput, err := e.ProcessBuiltinFunction(nil, nil, vmInput)
		assert.Nil(t, vmOutput)
		assert.Equal(t, ErrNilUserAccount, err)
	})
	t.Run("built-in function is not active", func(t *testing.T) {
		t.Parallel()

		enableEpochsHandler := &mock.EnableEpochsHandlerStub{
			IsFlagEnabledCalled: func(flag core.EnableEpochFlag) bool {
				return false
			},
		}
		e, _ := NewESDTModifyRoyaltiesFunc(0, &mock.AccountsStub{}, &mock.GlobalSettingsHandlerStub{}, &mock.ESDTNFTStorageHandlerStub{}, &mock.ESDTRoleHandlerStub{}, enableEpochsHandler)
		vmInput := &vmcommon.ContractCallInput{
			VMInput: vmcommon.VMInput{
				CallValue:  big.NewInt(0),
				CallerAddr: []byte("caller"),
			},
			RecipientAddr: []byte("caller"),
		}
		vmOutput, err := e.ProcessBuiltinFunction(mock.NewUserAccount([]byte("addr")), nil, vmInput)
		assert.Nil(t, vmOutput)
		assert.Equal(t, ErrBuiltInFunctionIsNotActive, err)
	})
	t.Run("invalid number of arguments", func(t *testing.T) {
		t.Parallel()

		enableEpochsHandler := &mock.EnableEpochsHandlerStub{
			IsFlagEnabledCalled: func(flag core.EnableEpochFlag) bool {
				return true
			},
		}
		e, _ := NewESDTModifyRoyaltiesFunc(0, &mock.AccountsStub{}, &mock.GlobalSettingsHandlerStub{}, &mock.ESDTNFTStorageHandlerStub{}, &mock.ESDTRoleHandlerStub{}, enableEpochsHandler)
		vmInput := &vmcommon.ContractCallInput{
			VMInput: vmcommon.VMInput{
				CallValue:  big.NewInt(0),
				CallerAddr: []byte("caller"),
				Arguments:  [][]byte{},
			},
			RecipientAddr: []byte("caller"),
		}
		vmOutput, err := e.ProcessBuiltinFunction(mock.NewUserAccount([]byte("addr")), nil, vmInput)
		assert.Nil(t, vmOutput)
		assert.Equal(t, ErrInvalidNumberOfArguments, err)
	})
	t.Run("check allowed to execute failed", func(t *testing.T) {
		t.Parallel()

		allowedToExecuteCalled := false
		expectedErr := errors.New("expected error")
		enableEpochsHandler := &mock.EnableEpochsHandlerStub{
			IsFlagEnabledCalled: func(flag core.EnableEpochFlag) bool {
				return true
			},
		}
		rolesHandler := &mock.ESDTRoleHandlerStub{
			CheckAllowedToExecuteCalled: func(account vmcommon.UserAccountHandler, tokenID []byte, role []byte) error {
				allowedToExecuteCalled = true
				return expectedErr
			},
		}
		e, _ := NewESDTModifyRoyaltiesFunc(0, &mock.AccountsStub{}, &mock.GlobalSettingsHandlerStub{}, &mock.ESDTNFTStorageHandlerStub{}, rolesHandler, enableEpochsHandler)
		vmInput := &vmcommon.ContractCallInput{
			VMInput: vmcommon.VMInput{
				CallValue:  big.NewInt(0),
				CallerAddr: []byte("caller"),
				Arguments:  [][]byte{[]byte("tokenID"), {}, {}, {}, {}, {}, {}},
			},
			RecipientAddr: []byte("caller"),
		}
		vmOutput, err := e.ProcessBuiltinFunction(mock.NewUserAccount([]byte("addr")), nil, vmInput)
		assert.Nil(t, vmOutput)
		assert.Equal(t, expectedErr, err)
		assert.True(t, allowedToExecuteCalled)
	})
	t.Run("only changes the royalties", func(t *testing.T) {
		t.Parallel()

		getESDTNFTTokenOnSenderCalled := false
		saveESDTNFTTokenCalled := false
		tokenId := []byte("tokenID")
		esdtTokenKey := append([]byte(baseESDTKeyPrefix), tokenId...)
		nonce := uint64(15)

		enableEpochsHandler := &mock.EnableEpochsHandlerStub{
			IsFlagEnabledCalled: func(flag core.EnableEpochFlag) bool {
				return true
			},
		}
		globalSettingsHandler := &mock.GlobalSettingsHandlerStub{
			GetTokenTypeCalled: func(key []byte) (uint32, error) {
				assert.Equal(t, esdtTokenKey, key)
				return uint32(core.DynamicNFT), nil
			},
		}
		accounts := &mock.AccountsStub{}
		oldMetaData := &esdt.MetaData{
			Nonce:      nonce,
			Name:       []byte("name"),
			Creator:    []byte("creator"),
			Royalties:  10,
			Hash:       []byte("hash"),
			URIs:       [][]byte{[]byte("uri")},
			Attributes: []byte("attributes"),
		}
		storageHandler := &mock.ESDTNFTStorageHandlerStub{
			GetESDTNFTTokenOnSenderCalled: func(acnt vmcommon.UserAccountHandler, esdtTokenKey []byte, n uint64) (*esdt.ESDigitalToken, error) {
				getESDTNFTTokenOnSenderCalled = true
				return &esdt.ESDigitalToken{
					Value:         big.NewInt(1),
					TokenMetaData: oldMetaData,
				}, nil
			},
			SaveESDTNFTTokenCalled: func(senderAddress []byte, acnt vmcommon.UserAccountHandler, tokenKey []byte, n uint64, esdtData *esdt.ESDigitalToken, mustUpdateAllFields bool, isReturnWithError bool) ([]byte, error) {
				assert.Equal(t, esdtTokenKey, tokenKey)
				assert.Equal(t, nonce, n)
				assert.Equal(t, oldMetaData.Name, esdtData.TokenMetaData.Name)
				assert.Equal(t, oldMetaData.URIs, esdtData.TokenMetaData.URIs)
				assert.Equal(t, uint32(50), esdtData.TokenMetaData.Royalties)
				assert.Equal(t, oldMetaData.Hash, esdtData.TokenMetaData.Hash)
				assert.Equal(t, oldMetaData.Attributes, esdtData.TokenMetaData.Attributes)
				assert.Equal(t, oldMetaData.Creator, esdtData.TokenMetaData.Creator)
				saveESDTNFTTokenCalled = true
				return nil, nil
			},
		}
		e, _ := NewESDTModifyRoyaltiesFunc(101, accounts, globalSettingsHandler, storageHandler, &mock.ESDTRoleHandlerStub{}, enableEpochsHandler)

		vmInput := &vmcommon.ContractCallInput{
			VMInput: vmcommon.VMInput{
				CallValue:   big.NewInt(0),
				CallerAddr:  []byte("caller"),
				GasProvided: 1000,
				Arguments:   [][]byte{tokenId, {15}, {50}},
			},
			RecipientAddr: []byte("caller"),
		}

		vmOutput, err := e.ProcessBuiltinFunction(mock.NewUserAccount([]byte("addr")), nil, vmInput)
		assert.Nil(t, err)
		assert.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
		assert.Equal(t, uint64(899), vmOutput.GasRemaining)
		assert.True(t, saveESDTNFTTokenCalled)
		assert.True(t, getESDTNFTTokenOnSenderCalled)
	})
}

func TestESDTModifyRoyalties_SetNewGasConfig(t *testing.T) {
	t.Parallel()

	e, _ := NewESDTModifyRoyaltiesFunc(0, &mock.AccountsStub{}, &mock.GlobalSettingsHandlerStub{}, &mock.ESDTNFTStorageHandlerStub{}, &mock.ESDTRoleHandlerStub{}, &mock.EnableEpochsHandlerStub{})

	newGasCost := &vmcommon.GasCost{
		BuiltInCost: vmcommon.BuiltInCost{
			ESDTModifyRoyalties: 10,
		},
	}
	e.SetNewGasConfig(newGasCost)

	assert.Equal(t, newGasCost.BuiltInCost.ESDTModifyRoyalties, e.funcGasCost)
}
