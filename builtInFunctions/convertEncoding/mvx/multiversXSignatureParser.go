package mvx

import (
	convertCommon "github.com/multiversx/mx-chain-vm-common-go/builtInFunctions/convertEncoding/common"
	"strings"
)

var delimiters = []string{BeginType, EndType, convertCommon.Comma}

func ParseMultiversXSignature(signature string) (convertCommon.Arguments, error) {
	arguments, remainder, err := parseMultiversXExpression(signature)
	if remainder != "" {
		return convertCommon.Arguments{}, convertCommon.ErrExpectedBlankRemainder
	}
	return arguments, err
}

func parseMultiversXExpression(expression string) (convertCommon.Arguments, string, error) {
	arguments := convertCommon.Arguments{}
	expression = strings.TrimSpace(expression)

	if expression == "" {
		return nil, "", convertCommon.ErrBlankExpression
	}

	for len(expression) > 0 {
		var err error
		var token string

		token, expression, err = extractMvxToken(expression)
		if err != nil {
			return nil, "", err
		}

		var children convertCommon.Arguments
		if convertCommon.StartsWith(expression, BeginType) {
			children, expression, err = extractTypes(expression[1:])
		}
		if err != nil {
			return nil, "", err
		}

		argument := &convertCommon.Argument{Type: token, Arguments: children}
		arguments = append(arguments, argument)

		if convertCommon.StartsWith(expression, EndType) {
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

func extractTypes(expression string) (convertCommon.Arguments, string, error) {
	children, expression, err := parseMultiversXExpression(expression)
	if err != nil {
		return nil, "", err
	}
	if !convertCommon.StartsWith(expression, EndType) {
		return nil, "", ErrExpectedTypeEnd
	}
	return children, expression[1:], nil
}

func extractMvxToken(expression string) (string, string, error) {
	return convertCommon.ExtractToken(expression, delimiters)
}

func ExtractMvxArraySize(arrayType string) (int, error) {
	return convertCommon.ExtractNumericSuffix(arrayType, Array)
}
