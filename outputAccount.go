package vmcommon

import (
	"math/big"
)

var _ OutputAccountHandler = (*OutputAccount)(nil)

// GetAddress returns the account address
func (oa *OutputAccount) GetAddress() []byte {
	return oa.Address
}

// GetNonce returns the account nonce
func (oa *OutputAccount) GetNonce() uint64 {
	return oa.Nonce
}

// GetBalance returns the account balance
func (oa *OutputAccount) GetBalance() *big.Int {
	return oa.Balance
}

// GetStorageUpdates returns the storage updates map
func (oa *OutputAccount) GetStorageUpdates() map[string]*StorageUpdate {
	return oa.StorageUpdates
}

// GetCode returns the smart contract code
func (oa *OutputAccount) GetCode() []byte {
	return oa.Code
}

// GetCodeMetadata returns the metadata of the smart contract code
func (oa *OutputAccount) GetCodeMetadata() []byte {
	return oa.CodeMetadata
}

// GetCodeDeployerAddress returns the deployer's address
func (oa *OutputAccount) GetCodeDeployerAddress() []byte {
	return oa.CodeDeployerAddress
}

// GetBalanceDelta returns the balance delta
func (oa *OutputAccount) GetBalanceDelta() *big.Int {
	return oa.BalanceDelta
}

// GetOutputTransfers returns the list of output transfers
func (oa *OutputAccount) GetOutputTransfers() []OutputTransfer {
	return oa.OutputTransfers
}

// GetGasUsed returns the gas used
func (oa *OutputAccount) GetGasUsed() uint64 {
	return oa.GasUsed
}

// GetBytesAddedToStorage returns the bytes added to storage
func (oa *OutputAccount) GetBytesAddedToStorage() uint64 {
	return oa.BytesAddedToStorage
}

// GetBytesDeletedFromStorage returns the bytes deleted from storage
func (oa *OutputAccount) GetBytesDeletedFromStorage() uint64 {
	return oa.BytesDeletedFromStorage
}

// GetBytesConsumedByTxAsNetworking returns the bytes consumed by networking
func (oa *OutputAccount) GetBytesConsumedByTxAsNetworking() uint64 {
	return oa.BytesConsumedByTxAsNetworking
}

// SetNonce sets the account nonce
func (oa *OutputAccount) SetNonce(nonce uint64) {
	oa.Nonce = nonce
}

// SetAddress sets the account address
func (oa *OutputAccount) SetAddress(address []byte) {
	oa.Address = address
}

// SetBalance sets the account balance
func (oa *OutputAccount) SetBalance(balance *big.Int) {
	oa.Balance = balance
}

// SetStorageUpdates sets the storage updates
func (oa *OutputAccount) SetStorageUpdates(updates map[string]*StorageUpdate) {
	oa.StorageUpdates = updates
}

// SetCode sets the smart contract code
func (oa *OutputAccount) SetCode(code []byte) {
	oa.Code = code
}

// SetCodeMetadata sets the code metadata
func (oa *OutputAccount) SetCodeMetadata(metadata []byte) {
	oa.CodeMetadata = metadata
}

// SetCodeDeployerAddress sets the deployer address
func (oa *OutputAccount) SetCodeDeployerAddress(address []byte) {
	oa.CodeDeployerAddress = address
}

// SetBalanceDelta sets the balance delta
func (oa *OutputAccount) SetBalanceDelta(delta *big.Int) {
	oa.BalanceDelta = delta
}

// SetOutputTransfers sets the output transfers
func (oa *OutputAccount) SetOutputTransfers(transfers []OutputTransfer) {
	oa.OutputTransfers = transfers
}

// SetGasUsed sets the gas used
func (oa *OutputAccount) SetGasUsed(gasUsed uint64) {
	oa.GasUsed = gasUsed
}

// SetBytesAddedToStorage sets the bytes added to storage
func (oa *OutputAccount) SetBytesAddedToStorage(bytes uint64) {
	oa.BytesAddedToStorage = bytes
}

// SetBytesDeletedFromStorage sets the bytes deleted from storage
func (oa *OutputAccount) SetBytesDeletedFromStorage(bytes uint64) {
	oa.BytesDeletedFromStorage = bytes
}

// SetBytesConsumedByTxAsNetworking sets the bytes consumed by networking
func (oa *OutputAccount) SetBytesConsumedByTxAsNetworking(bytes uint64) {
	oa.BytesConsumedByTxAsNetworking = bytes
}
