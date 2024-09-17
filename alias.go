package vmcommon

import (
	"bytes"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/multiversx/mx-chain-core-go/core"
)

var multiversXBlankAddress = make([]byte, len(core.SystemAccountAddress))

var ethereumBlankAddress = make([]byte, ethCommon.AddressLength)

// AliasSaveRequest is the save request for an alias
type AliasSaveRequest struct {
	MultiversXAddress []byte
	AliasAddress      []byte
	AliasIdentifier   core.AddressIdentifier
}

// AddressRequest is the request for an address
type AddressRequest struct {
	SourceAddress       []byte
	SourceIdentifier    core.AddressIdentifier
	RequestedIdentifier core.AddressIdentifier
	SaveOnGenerate      bool
}

// AddressResponse is the response for an address request
type AddressResponse struct {
	RequestedAddress  []byte
	MultiversXAddress []byte
}

func ValidateAliasSaveRequest(request *AliasSaveRequest) error {
	if request == nil {
		return ErrNilRequest
	}
	switch request.AliasIdentifier {
	case core.InvalidAddressIdentifier, core.MVXAddressIdentifier:
		return core.ErrInvalidAddressIdentifier
	default:
		return nil
	}
}

func ValidateAddressRequest(request *AddressRequest) error {
	if request == nil {
		return ErrNilRequest
	}
	if request.SourceIdentifier == core.InvalidAddressIdentifier {
		return ErrInvalidSourceIdentifier
	}
	if request.RequestedIdentifier == core.InvalidAddressIdentifier {
		return ErrInvalidRequestedIdentifier
	}
	if request.SourceIdentifier == request.RequestedIdentifier {
		return ErrSourceIdentifierMatchesRequestedIdentifier
	}
	return nil
}

func EnhanceAddressRequest(request *AddressRequest) error {
	if len(request.SourceAddress) == 0 {
		blankAddress, err := RequestBlankAddress(request.SourceIdentifier)
		if err != nil {
			return err
		}
		request.SourceAddress = blankAddress
	}
	return nil
}

func IsBlankAddress(address []byte, addressIdentifier core.AddressIdentifier) bool {
	switch addressIdentifier {
	case core.MVXAddressIdentifier:
		return bytes.Equal(address, multiversXBlankAddress)
	case core.ETHAddressIdentifier:
		return bytes.Equal(address, ethereumBlankAddress)
	default:
		return false
	}
}

func RequestBlankAddress(addressIdentifier core.AddressIdentifier) ([]byte, error) {
	switch addressIdentifier {
	case core.MVXAddressIdentifier:
		return multiversXBlankAddress, nil
	case core.ETHAddressIdentifier:
		return ethereumBlankAddress, nil
	default:
		return nil, ErrIdentifierNotHandledForBlankAddress
	}
}
