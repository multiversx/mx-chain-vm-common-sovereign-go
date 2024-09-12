package ethToMvx

import (
	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/common"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/eth"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/mvx"
)

const EthereumToMultiversXInputDataExpectedSize = 1

func validateEthereumToMultiversXEncodingInput(inputData [][]byte) error {
	if len(inputData) != EthereumToMultiversXInputDataExpectedSize {
		return ErrEthToMvxExpectedOneArgument
	}
	return nil
}

func MvxAndEthAbiFromEthSignature(signature string) (mvx.AbiArguments, ethAbi.Arguments, error) {
	ethArgs, err := eth.ParseEthereumSignature(signature)
	if err != nil {
		return nil, nil, err
	}

	ethereumAbi, err := eth.BuildEthereumAbi(ethArgs)
	if err != nil {
		return nil, nil, err
	}

	mvxArgs, err := EthereumToMultiversXArguments(ethereumAbi)
	if err != nil {
		return nil, nil, err
	}

	multiversXAbi, err := mvx.BuildMultiversXAbi(mvxArgs)
	if err != nil {
		return nil, nil, err
	}

	return multiversXAbi, ethereumAbi, nil
}

func EthereumToMultiversXEncoding(context *common.EncodingContext, multiversXAbi mvx.AbiArguments, ethereumAbi ethAbi.Arguments, inputData [][]byte) ([][]byte, error) {
	err := validateEthereumToMultiversXEncodingInput(inputData)
	if err != nil {
		return nil, err
	}

	ethDecodedData, err := ethereumAbi.Unpack(inputData[0])
	if err != nil {
		return nil, err
	}

	err = AttachValuesToMultiversXAbi(context, multiversXAbi, ethereumAbi, ethDecodedData)
	if err != nil {
		return nil, err
	}

	return mvx.Serializer.SerializeToParts(multiversXAbi)
}
