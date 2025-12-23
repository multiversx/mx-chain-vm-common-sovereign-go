package vmcommon

var _ ContractCreateInputHandler = (*ContractCreateInput)(nil)

// GetVMInput returns the embedded VMInput
func (c *ContractCreateInput) GetVMInput() *VMInput {
	return &c.VMInput
}

// GetContractCode returns the contract code
func (c *ContractCreateInput) GetContractCode() []byte {
	return c.ContractCode
}

// GetContractCodeMetadata returns the contract code metadata
func (c *ContractCreateInput) GetContractCodeMetadata() []byte {
	return c.ContractCodeMetadata
}

var _ ContractCallInputHandler = (*ContractCallInput)(nil)

// SetVMInput sets the embedded VMInput
func (c *ContractCallInput) SetVMInput(vmInput VMInput) {
	c.VMInput = vmInput
}

// GetVMInput returns the embedded VMInput
func (c *ContractCallInput) GetVMInput() *VMInput {
	return &c.VMInput
}

// GetRecipientAddr returns the recipient address
func (c *ContractCallInput) GetRecipientAddr() []byte {
	return c.RecipientAddr
}

// GetFunction returns the smart contract function name
func (c *ContractCallInput) GetFunction() string {
	return c.Function
}

// GetAllowInitFunction returns whether calling the init function is allowed
func (c *ContractCallInput) GetAllowInitFunction() bool {
	return c.AllowInitFunction
}

// SetRecipientAddr sets the recipient address
func (c *ContractCallInput) SetRecipientAddr(address []byte) {
	c.RecipientAddr = address
}
