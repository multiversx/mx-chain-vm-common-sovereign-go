package eth

import (
	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	convertCommon "github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/common"
)

func BuildEthereumAbi(arguments convertCommon.Arguments) (ethAbi.Arguments, error) {
	abiArguments := make([]ethAbi.Argument, len(arguments))
	for position, argument := range arguments {
		abiType, err := ethAbi.NewType(argument.Type, argument.Type, buildMarshalingArguments(argument.Arguments))
		if err != nil {
			return nil, err
		}
		abiArguments[position] = ethAbi.Argument{Type: abiType}
	}
	return abiArguments, nil
}

func buildMarshalingArguments(arguments convertCommon.Arguments) []ethAbi.ArgumentMarshaling {
	marshalingArguments := make([]ethAbi.ArgumentMarshaling, len(arguments))
	for position, argument := range arguments {
		marshalingArguments[position] = ethAbi.ArgumentMarshaling{
			Name:         convertCommon.BuildArgName(position),
			Type:         argument.Type,
			InternalType: argument.Type,
			Components:   buildMarshalingArguments(argument.Arguments),
		}
	}
	return marshalingArguments
}
