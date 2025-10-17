package mvxToEth

import (
	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/common"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/eth"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/mvx"
)

func MvxAndEthAbiFromMvxSignature(signature string) (mvx.AbiArguments, ethAbi.Arguments, error) {
	mvxArgs, err := mvx.ParseMultiversXSignature(signature)
	if err != nil {
		return nil, nil, err
	}

	multiversXAbi, err := mvx.BuildMultiversXAbi(mvxArgs)
	if err != nil {
		return nil, nil, err
	}

	ethArgs, err := MultiversXToEthereumArguments(mvxArgs)
	if err != nil {
		return nil, nil, err
	}

	ethereumAbi, err := eth.BuildEthereumAbi(ethArgs)
	if err != nil {
		return nil, nil, err
	}

	return multiversXAbi, ethereumAbi, nil
}

func MultiversXToEthereumEncoding(context *common.EncodingContext, multiversXAbi mvx.AbiArguments, ethereumAbi ethAbi.Arguments, inputData [][]byte) ([][]byte, error) {
	err := mvx.Serializer.DeserializeParts(inputData, multiversXAbi)
	if err != nil {
		return nil, err
	}

	mvxDecodedData, err := DetachValuesFromMultiversXAbi(context, multiversXAbi, ethereumAbi)
	if err != nil {
		return nil, err
	}

	ethEncodedData, err := ethereumAbi.Pack(mvxDecodedData...)
	if err != nil {
		return nil, err
	}

	return [][]byte{ethEncodedData}, nil
}
