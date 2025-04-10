package convertEncoding

import (
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/common"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/ethToMvx"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/mvxToEth"
)

type Handler struct {
	context *common.EncodingContext
}

func NewHandler(accounts vmcommon.AccountsAdapter) *Handler {
	return &Handler{context: common.BuildEncodingContext(accounts)}
}

func (h *Handler) EthToMvxEncodingWithMvxSignature(signature string, inputData [][]byte) ([][]byte, error) {
	multiversXAbi, ethereumAbi, err := mvxToEth.MvxAndEthAbiFromMvxSignature(signature)
	if err != nil {
		return nil, err
	}

	return ethToMvx.EthereumToMultiversXEncoding(h.context, multiversXAbi, ethereumAbi, inputData)
}

func (h *Handler) EthToMvxEncodingWithEthSignature(signature string, inputData [][]byte) ([][]byte, error) {
	multiversXAbi, ethereumAbi, err := ethToMvx.MvxAndEthAbiFromEthSignature(signature)
	if err != nil {
		return nil, err
	}

	return ethToMvx.EthereumToMultiversXEncoding(h.context, multiversXAbi, ethereumAbi, inputData)
}

func (h *Handler) MvxToEthEncodingWithMvxSignature(signature string, inputData [][]byte) ([][]byte, error) {
	multiversXAbi, ethereumAbi, err := mvxToEth.MvxAndEthAbiFromMvxSignature(signature)
	if err != nil {
		return nil, err
	}

	return mvxToEth.MultiversXToEthereumEncoding(h.context, multiversXAbi, ethereumAbi, inputData)
}

func (h *Handler) MvxToEthEncodingWithEthSignature(signature string, inputData [][]byte) ([][]byte, error) {
	multiversXAbi, ethereumAbi, err := ethToMvx.MvxAndEthAbiFromEthSignature(signature)
	if err != nil {
		return nil, err
	}

	return mvxToEth.MultiversXToEthereumEncoding(h.context, multiversXAbi, ethereumAbi, inputData)
}
