package common

import vmcommon "github.com/multiversx/mx-chain-vm-common-go"

type Argument struct {
	Type      string
	Arguments Arguments
}

type Arguments []*Argument

type EncodingContext struct {
	Accounts vmcommon.AccountsAdapter
}

func BuildEncodingContext(accounts vmcommon.AccountsAdapter) *EncodingContext {
	return &EncodingContext{Accounts: accounts}
}
