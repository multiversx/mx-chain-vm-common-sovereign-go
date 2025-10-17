package common

import (
	"strconv"
	"strings"
)

func BuildArgName(position int) string {
	return ArgName + strconv.Itoa(position)
}

func ExtractNumericSuffix(word string, prefix string) (int, error) {
	return strconv.Atoi(strings.ReplaceAll(word, prefix, ""))
}

func ExtractNestedArgument(argument *Argument) (*Argument, error) {
	if len(argument.Arguments) != SingleExpectedComponent {
		return nil, ErrExpectedOneNestedArgument
	}
	return argument.Arguments[0], nil
}

func StartsWith(expression string, character string) bool {
	return len(expression) > 0 && string(expression[0]) == character
}

func ExtractToken(expression string, delimiters []string) (string, string, error) {
	expression = strings.TrimSpace(expression)
	if len(expression) == 0 {
		return "", "", ErrBlankTokenExpression
	}
	if startsWithDelimiter(expression, delimiters) {
		return "", "", ErrExpressionStartsWithDelimiter
	}
	for position, character := range expression {
		if isDelimiter(string(character), delimiters) {
			return expression[:position], expression[position:], nil
		}
	}
	return expression, "", nil
}

func startsWithDelimiter(expression string, delimiters []string) bool {
	return len(expression) > 0 && isDelimiter(string(expression[0]), delimiters)
}

func isDelimiter(character string, delimiters []string) bool {
	for _, delimiter := range delimiters {
		if character == delimiter {
			return true
		}
	}
	return false
}
