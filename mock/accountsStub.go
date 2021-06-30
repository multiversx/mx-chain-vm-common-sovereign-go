package mock

import (
	"context"
	"errors"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// AccountsStub -
type AccountsStub struct {
	GetExistingAccountCalled func(address []byte) (vmcommon.AccountHandler, error)
	LoadAccountCalled        func(address []byte) (vmcommon.AccountHandler, error)
	SaveAccountCalled        func(account vmcommon.AccountHandler) error
	RemoveAccountCalled      func(address []byte) error
	CommitCalled             func() ([]byte, error)
	JournalLenCalled         func() int
	RevertToSnapshotCalled   func(snapshot int) error
	RootHashCalled           func() ([]byte, error)
	RecreateTrieCalled       func(rootHash []byte) error
	SnapshotStateCalled      func(rootHash []byte)
	SetStateCheckpointCalled func(rootHash []byte)
	IsPruningEnabledCalled   func() bool
	GetNumCheckpointsCalled  func() uint32
	GetCodeCalled            func([]byte) []byte
}

// GetCode -
func (as *AccountsStub) GetCode(codeHash []byte) []byte {
	if as.GetCodeCalled != nil {
		return as.GetCodeCalled(codeHash)
	}
	return nil
}

// LoadAccount -
func (as *AccountsStub) LoadAccount(address []byte) (vmcommon.AccountHandler, error) {
	if as.LoadAccountCalled != nil {
		return as.LoadAccountCalled(address)
	}
	return nil, errNotImplemented
}

// SaveAccount -
func (as *AccountsStub) SaveAccount(account vmcommon.AccountHandler) error {
	if as.SaveAccountCalled != nil {
		return as.SaveAccountCalled(account)
	}
	return nil
}

var errNotImplemented = errors.New("not implemented")

// Commit -
func (as *AccountsStub) Commit() ([]byte, error) {
	if as.CommitCalled != nil {
		return as.CommitCalled()
	}

	return nil, errNotImplemented
}

// GetExistingAccount -
func (as *AccountsStub) GetExistingAccount(address []byte) (vmcommon.AccountHandler, error) {
	if as.GetExistingAccountCalled != nil {
		return as.GetExistingAccountCalled(address)
	}

	return nil, errNotImplemented
}

// JournalLen -
func (as *AccountsStub) JournalLen() int {
	if as.JournalLenCalled != nil {
		return as.JournalLenCalled()
	}

	return 0
}

// RemoveAccount -
func (as *AccountsStub) RemoveAccount(address []byte) error {
	if as.RemoveAccountCalled != nil {
		return as.RemoveAccountCalled(address)
	}

	return errNotImplemented
}

// RevertToSnapshot -
func (as *AccountsStub) RevertToSnapshot(snapshot int) error {
	if as.RevertToSnapshotCalled != nil {
		return as.RevertToSnapshotCalled(snapshot)
	}

	return errNotImplemented
}

// RootHash -
func (as *AccountsStub) RootHash() ([]byte, error) {
	if as.RootHashCalled != nil {
		return as.RootHashCalled()
	}

	return nil, errNotImplemented
}

// RecreateTrie -
func (as *AccountsStub) RecreateTrie(rootHash []byte) error {
	if as.RecreateTrieCalled != nil {
		return as.RecreateTrieCalled(rootHash)
	}

	return errNotImplemented
}

// SnapshotState -
func (as *AccountsStub) SnapshotState(rootHash []byte, _ context.Context) {
	if as.SnapshotStateCalled != nil {
		as.SnapshotStateCalled(rootHash)
	}
}

// SetStateCheckpoint -
func (as *AccountsStub) SetStateCheckpoint(rootHash []byte, _ context.Context) {
	if as.SetStateCheckpointCalled != nil {
		as.SetStateCheckpointCalled(rootHash)
	}
}

// IsPruningEnabled -
func (as *AccountsStub) IsPruningEnabled() bool {
	if as.IsPruningEnabledCalled != nil {
		return as.IsPruningEnabledCalled()
	}

	return false
}

// GetNumCheckpoints -
func (as *AccountsStub) GetNumCheckpoints() uint32 {
	if as.GetNumCheckpointsCalled != nil {
		return as.GetNumCheckpointsCalled()
	}

	return 0
}

// IsInterfaceNil returns true if there is no value under the interface
func (as *AccountsStub) IsInterfaceNil() bool {
	return as == nil
}
