package builtInFunctions

// BaseAlwaysActiveHandler defines the base always active handler
type BaseAlwaysActiveHandler struct {
}

// IsActive returns true as this built-in function is always active
func (b BaseAlwaysActiveHandler) IsActive() bool {
	return trueHandler()
}

// IsInterfaceNil always returns false
func (b BaseAlwaysActiveHandler) IsInterfaceNil() bool {
	return false
}

type baseActiveHandler struct {
	activeHandler func() bool
}

// IsActive returns true if function is active
func (b *baseActiveHandler) IsActive() bool {
	return b.activeHandler()
}

// IsInterfaceNil returns true if there is no value under the interface
func (b *baseActiveHandler) IsInterfaceNil() bool {
	return b == nil
}
