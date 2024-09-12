package common

import (
	"errors"
)

var ErrExpectedOneNestedArgument = errors.New("expected one nested argument")

var ErrBlankExpression = errors.New("blank expression")

var ErrBlankTokenExpression = errors.New("blank token expression")

var ErrExpectedBlankRemainder = errors.New("expected blank remainder")

var ErrExpectedExpressionAfterComma = errors.New("expected expression after comma")

var ErrExpressionStartsWithDelimiter = errors.New("expression starts with delimiter")
