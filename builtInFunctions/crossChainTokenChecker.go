package builtInFunctions

import (
	"bytes"
	"fmt"

	"github.com/multiversx/mx-chain-core-go/data/esdt"
)

type crossChainTokenChecker struct {
	selfESDTPrefix       []byte
	whiteListedAddresses map[string]struct{}
}

// NewCrossChainTokenChecker creates a new cross chain token checker
func NewCrossChainTokenChecker(selfESDTPrefix []byte, whiteListedAddresses map[string]struct{}) (*crossChainTokenChecker, error) {
	ctc := &crossChainTokenChecker{
		selfESDTPrefix:       selfESDTPrefix,
		whiteListedAddresses: whiteListedAddresses,
	}

	if !areChainParamsCompatible(selfESDTPrefix, whiteListedAddresses) {
		return nil, ErrInvalidCrossChainConfig
	}

	if len(selfESDTPrefix) == 0 {
		return ctc, nil
	}

	if !esdt.IsValidTokenPrefix(string(selfESDTPrefix)) {
		return nil, fmt.Errorf("%w: %s", ErrInvalidTokenPrefix, selfESDTPrefix)
	}

	return ctc, nil
}

func areChainParamsCompatible(selfESDTPrefix []byte, whiteListedAddresses map[string]struct{}) bool {
	validConfigMainChain := len(selfESDTPrefix) == 0 && len(whiteListedAddresses) != 0
	validConfigSingleShardChain := len(selfESDTPrefix) != 0 && len(whiteListedAddresses) == 0

	return validConfigMainChain || validConfigSingleShardChain
}

// IsCrossChainOperation checks if the provided token comes from another chain/sovereign shard
func (ctc *crossChainTokenChecker) IsCrossChainOperation(tokenID []byte) bool {
	tokenPrefix, hasPrefix := esdt.IsValidPrefixedToken(string(tokenID))
	// no prefix or malformed token in main chain operation
	if !hasPrefix && len(ctc.selfESDTPrefix) == 0 {
		return false
	}

	return !bytes.Equal([]byte(tokenPrefix), ctc.selfESDTPrefix)
}

// IsAllowedToMint checks whether an address is allowed to mint/create/add quantity a token
func (ctc *crossChainTokenChecker) IsAllowedToMint(address []byte, tokenID []byte) bool {
	return ctc.isWhiteListed(address) && ctc.IsCrossChainOperation(tokenID)
}

func (ctc *crossChainTokenChecker) isWhiteListed(address []byte) bool {
	_, found := ctc.whiteListedAddresses[string(address)]
	return found
}

// IsInterfaceNil checks if the underlying pointer is nil
func (ctc *crossChainTokenChecker) IsInterfaceNil() bool {
	return ctc == nil
}
