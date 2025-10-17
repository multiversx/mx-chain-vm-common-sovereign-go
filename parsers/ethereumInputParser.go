package parsers

func ParseEthereumCallInput(input []byte) ([]byte, []byte, error) {
	if len(input) < EVMSelectorSize {
		return nil, nil, ErrUnexpectedInputSize
	}
	return input[:EVMSelectorSize], input[EVMSelectorSize:], nil
}
