package common

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestBuildArgName(t *testing.T) {
	require.Equal(t, BuildArgName(1), ArgName+"1")
	require.Equal(t, BuildArgName(5), ArgName+"5")
}

func TestExtractNumericSuffix(t *testing.T) {
	suffix := 10
	extractedSuffix, err := ExtractNumericSuffix(BuildArgName(suffix), ArgName)
	require.NoError(t, err)
	require.Equal(t, suffix, extractedSuffix)

	_, err = ExtractNumericSuffix(ArgName, ArgName)
	require.Error(t, err)
}

func TestExtractNestedArgument(t *testing.T) {
	expectedNestedComponent := &Argument{Type: ArgName}
	extractedNestedComponent, err := ExtractNestedArgument(&Argument{Arguments: Arguments{expectedNestedComponent}})
	require.NoError(t, err)
	require.Equal(t, expectedNestedComponent.Type, extractedNestedComponent.Type)
	require.Equal(t, reflect.ValueOf(expectedNestedComponent).Pointer(), reflect.ValueOf(extractedNestedComponent).Pointer())

	_, err = ExtractNestedArgument(&Argument{})
	require.Equal(t, err, ErrExpectedOneNestedArgument)
}

func TestExtractToken(t *testing.T) {
	firstPart := BuildArgName(1)
	secondPart := Comma + BuildArgName(2)
	first, second, err := ExtractToken(firstPart+secondPart, []string{Comma})
	require.NoError(t, err)
	require.Equal(t, first, firstPart)
	require.Equal(t, second, secondPart)

	first, second, err = ExtractToken(firstPart, []string{Comma})
	require.NoError(t, err)
	require.Equal(t, first, firstPart)
	require.Equal(t, second, "")

	_, _, err = ExtractToken("", []string{Comma})
	require.Equal(t, err, ErrBlankTokenExpression)

	_, _, err = ExtractToken(secondPart, []string{Comma})
	require.Equal(t, err, ErrExpressionStartsWithDelimiter)
}

func TestStartsWith(t *testing.T) {
	require.True(t, StartsWith(Comma, Comma))
	require.False(t, StartsWith("", Comma))
	require.False(t, StartsWith(ArgName, Comma))
}

func TestStartsWithDelimiter(t *testing.T) {
	require.True(t, startsWithDelimiter(Comma, []string{Comma}))
	require.False(t, startsWithDelimiter("", []string{}))
	require.False(t, startsWithDelimiter(ArgName, []string{Comma}))
}

func TestIsDelimiter(t *testing.T) {
	require.True(t, isDelimiter(Comma, []string{Comma}))
	require.False(t, isDelimiter(Comma, []string{}))
	require.False(t, isDelimiter(PartsSeparator, []string{Comma}))
}
