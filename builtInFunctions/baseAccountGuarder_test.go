package builtInFunctions

import (
	"math/big"
	"strings"
	"testing"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	mockvm "github.com/multiversx/mx-chain-vm-common-go/mock"
	"github.com/stretchr/testify/require"
)

var marshallerMock = &mockvm.MarshalizerMock{}

func TestNewBaseAccountGuarder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args        func() BaseAccountGuarderArgs
		expectedErr error
	}{
		{
			args: func() BaseAccountGuarderArgs {
				args := createBaseAccountGuarderArgs()
				args.Marshaller = nil
				return args
			},
			expectedErr: ErrNilMarshalizer,
		},
		{
			args: func() BaseAccountGuarderArgs {
				args := createBaseAccountGuarderArgs()
				args.EnableEpochsHandler = nil
				return args
			},
			expectedErr: ErrNilEnableEpochsHandler,
		},
		{
			args: func() BaseAccountGuarderArgs {
				args := createBaseAccountGuarderArgs()
				args.GuardedAccountHandler = nil
				return args
			},
			expectedErr: ErrNilGuardedAccountHandler,
		},
		{
			args: func() BaseAccountGuarderArgs {
				return createBaseAccountGuarderArgs()
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		instance, err := newBaseAccountGuarder(test.args())
		if test.expectedErr != nil {
			require.Nil(t, instance)
			require.Equal(t, test.expectedErr, err)
		} else {
			require.NotNil(t, instance)
			require.Nil(t, err)
		}
	}
}

func TestBaseAccountGuarder_CheckArgs(t *testing.T) {
	t.Parallel()

	address := generateRandomByteArray(pubKeyLen)
	account := mockvm.NewUserAccount(address)

	guardianAddress := generateRandomByteArray(pubKeyLen)
	vmInput := getDefaultVmInput([][]byte{guardianAddress})
	vmInput.CallerAddr = address
	vmInput.RecipientAddr = account.Address

	tests := []struct {
		vmInput       func() *vmcommon.ContractCallInput
		senderAccount vmcommon.UserAccountHandler
		expectedErr   error
		noOfArgs      uint32
	}{
		{
			vmInput: func() *vmcommon.ContractCallInput {
				return vmInput
			},
			senderAccount: nil,
			expectedErr:   ErrNilUserAccount,
			noOfArgs:      1,
		},
		{
			vmInput: func() *vmcommon.ContractCallInput {
				return nil
			},
			senderAccount: account,
			expectedErr:   ErrNilVmInput,
			noOfArgs:      1,
		},
		{
			vmInput: func() *vmcommon.ContractCallInput {
				return vmInput
			},
			senderAccount: mockvm.NewUserAccount([]byte("userAddress2")),
			expectedErr:   ErrOperationNotPermitted,
			noOfArgs:      1,
		},
		{
			vmInput: func() *vmcommon.ContractCallInput {
				input := *vmInput
				input.CallerAddr = []byte("userAddress2")
				input.RecipientAddr = mockvm.NewUserAccount([]byte("userAddress2")).Address
				return &input
			},
			senderAccount: account,
			expectedErr:   ErrOperationNotPermitted,
			noOfArgs:      1,
		},
		{
			vmInput: func() *vmcommon.ContractCallInput {
				input := *vmInput
				input.CallerAddr = []byte("userAddress2")
				return &input
			},
			senderAccount: account,
			expectedErr:   ErrOperationNotPermitted,
			noOfArgs:      1,
		},
		{
			vmInput: func() *vmcommon.ContractCallInput {
				input := *vmInput
				input.CallValue = nil
				return &input
			},
			senderAccount: account,
			expectedErr:   ErrNilValue,
			noOfArgs:      1,
		},
		{
			vmInput: func() *vmcommon.ContractCallInput {
				input := *vmInput
				input.CallValue = big.NewInt(1)
				return &input
			},
			senderAccount: account,
			expectedErr:   ErrBuiltInFunctionCalledWithValue,
			noOfArgs:      1,
		},
		{
			vmInput: func() *vmcommon.ContractCallInput {
				input := *vmInput
				input.Arguments = [][]byte{guardianAddress, guardianAddress}
				return &input
			},
			senderAccount: account,
			expectedErr:   ErrInvalidNumberOfArguments,
			noOfArgs:      1,
		},
		{
			vmInput: func() *vmcommon.ContractCallInput {
				input := *vmInput
				input.GasProvided = 0
				return &input
			},
			senderAccount: account,
			expectedErr:   ErrNotEnoughGas,
			noOfArgs:      1,
		},
		{
			vmInput: func() *vmcommon.ContractCallInput {
				return vmInput
			},
			senderAccount: account,
			expectedErr:   nil,
			noOfArgs:      1,
		},
	}

	args := createBaseAccountGuarderArgs()
	baseAccGuarder, _ := newBaseAccountGuarder(args)

	for _, test := range tests {
		err := baseAccGuarder.checkBaseAccountGuarderArgs(test.senderAccount, test.vmInput(), test.noOfArgs)
		if test.expectedErr != nil {
			require.Error(t, err)
			require.True(t, strings.Contains(err.Error(), test.expectedErr.Error()))
		} else {
			require.Nil(t, err)
		}
	}
}

func createBaseAccountGuarderArgs() BaseAccountGuarderArgs {
	return BaseAccountGuarderArgs{
		Marshaller: marshallerMock,
		EnableEpochsHandler: &mockvm.EnableEpochsHandlerStub{
			IsSetGuardianEnabledField: false,
		},
		FuncGasCost:           100000,
		GuardedAccountHandler: &mockvm.GuardedAccountHandlerStub{},
	}
}
