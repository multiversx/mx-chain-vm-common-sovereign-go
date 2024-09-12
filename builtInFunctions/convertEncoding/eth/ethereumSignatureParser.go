package eth

import (
	convertCommon "github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/common"
	"strings"
)

var delimiters = []string{BeginTuple, EndTuple, convertCommon.Comma}

func ParseEthereumSignature(signature string) (convertCommon.Arguments, error) {
	arguments, remainder, err := parseEthereumExpression(signature)
	if remainder != "" {
		return convertCommon.Arguments{}, convertCommon.ErrExpectedBlankRemainder
	}
	return arguments, err
}

func parseEthereumExpression(expression string) (convertCommon.Arguments, string, error) {
	arguments := convertCommon.Arguments{}
	expression = strings.TrimSpace(expression)

	if expression == "" {
		return nil, "", convertCommon.ErrBlankExpression
	}

	for len(expression) > 0 {
		var err error
		var argument *convertCommon.Argument

		if convertCommon.StartsWith(expression, BeginTuple) {
			argument, expression, err = extractTuple(expression[1:])
		} else {
			argument, expression, err = extractSimpleToken(expression)
		}
		if err != nil {
			return nil, "", err
		}

		arguments = append(arguments, argument)

		if convertCommon.StartsWith(expression, EndTuple) {
			return arguments, expression, nil
		}
		if convertCommon.StartsWith(expression, convertCommon.Comma) {
			expression = expression[1:]
			if expression == "" {
				return nil, "", convertCommon.ErrExpectedExpressionAfterComma
			}
		}
	}

	return arguments, "", nil
}

func extractTuple(expression string) (*convertCommon.Argument, string, error) {
	children, expression, err := parseEthereumExpression(expression)
	if err != nil {
		return nil, "", err
	}
	if !convertCommon.StartsWith(expression, EndTuple) {
		return nil, "", ErrExpectedTupleEnd
	}

	var arrayType string
	expression = expression[1:]

	if convertCommon.StartsWith(expression, BeginArray) {
		arrayType, expression, err = extractToken(expression)
		if err != nil {
			return nil, "", err
		}
	}
	return &convertCommon.Argument{Type: Tuple + arrayType, Arguments: children}, expression, nil
}

func extractSimpleToken(expression string) (*convertCommon.Argument, string, error) {
	token, expression, err := extractToken(expression)
	if err != nil {
		return nil, "", err
	}
	return &convertCommon.Argument{Type: token}, expression, nil
}

func extractToken(expression string) (string, string, error) {
	return convertCommon.ExtractToken(expression, delimiters)
}
