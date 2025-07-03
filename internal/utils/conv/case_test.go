package conv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToUpperCamelCase_BasicSnakeCase(t *testing.T) {
	result := ToUpperCamelCase("user_id")
	assert.Equal(t, "UserId", result)
}

func TestToUpperCamelCase_BasicKebabCase(t *testing.T) {
	result := ToUpperCamelCase("user-name")
	assert.Equal(t, "UserName", result)
}

func TestToUpperCamelCase_StartsWithDigit(t *testing.T) {
	result := ToUpperCamelCase("1type")
	assert.Equal(t, "X1type", result)
}

func TestToUpperCamelCase_WithSpecialChars(t *testing.T) {
	result := ToUpperCamelCase("full#access")
	assert.Equal(t, "FullAccess", result)
}

func TestToUpperCamelCase_MultipleSpecialChars(t *testing.T) {
	result := ToUpperCamelCase("!@#special-case")
	assert.Equal(t, "SpecialCase", result)
}

func TestToUpperCamelCase_EmptyString(t *testing.T) {
	result := ToUpperCamelCase("")
	assert.Equal(t, "Xxx", result)
}

func TestToUpperCamelCase_SingleChar(t *testing.T) {
	result := ToUpperCamelCase("a")
	assert.Equal(t, "A", result)
}

func TestToUpperCamelCase_AlreadyUpperCamel(t *testing.T) {
	result := ToUpperCamelCase("UserName")
	assert.Equal(t, "UserName", result)
}

func TestToUpperCamelCase_OnlySpecialChars(t *testing.T) {
	result := ToUpperCamelCase("!@#")
	assert.Equal(t, "Xxx", result)
}

func TestToLowerCamelCase_BasicSnakeCase(t *testing.T) {
	result := ToLowerCamelCase("user_id")
	assert.Equal(t, "userId", result)
}

func TestToLowerCamelCase_UpperCaseInput(t *testing.T) {
	result := ToLowerCamelCase("Type")
	assert.Equal(t, "xType", result)
}

func TestToLowerCamelCase_StartsWithDigit(t *testing.T) {
	result := ToLowerCamelCase("1invalid")
	assert.Equal(t, "x1invalid", result)
}

func TestToLowerCamelCase_EmptyString(t *testing.T) {
	result := ToLowerCamelCase("")
	assert.Equal(t, "xxx", result)
}

func TestToLowerCamelCase_SingleChar(t *testing.T) {
	result := ToLowerCamelCase("A")
	assert.Equal(t, "a", result)
}

func TestToLowerCamelCase_AlreadyLowerCamel(t *testing.T) {
	result := ToLowerCamelCase("userName")
	assert.Equal(t, "userName", result)
}

func TestToLowerCamelCase_WithKebabCase(t *testing.T) {
	result := ToLowerCamelCase("user-name")
	assert.Equal(t, "userName", result)
}

func TestToLowerCamelCase_OnlySpecialChars(t *testing.T) {
	result := ToLowerCamelCase("!@#")
	assert.Equal(t, "xxx", result)
}

func TestToLowerInlineCase_BasicSnakeCase(t *testing.T) {
	result := ToLowerInlineCase("user_id")
	assert.Equal(t, "userid", result)
}

func TestToLowerInlineCase_LongSnakeCase(t *testing.T) {
	result := ToLowerInlineCase("snake_case_value")
	assert.Equal(t, "snakecasevalue", result)
}

func TestToLowerInlineCase_MixedCase(t *testing.T) {
	result := ToLowerInlineCase("User_Name")
	assert.Equal(t, "username", result)
}

func TestToLowerInlineCase_NoUnderscores(t *testing.T) {
	result := ToLowerInlineCase("username")
	assert.Equal(t, "username", result)
}

func TestToLowerInlineCase_EmptyString(t *testing.T) {
	result := ToLowerInlineCase("")
	assert.Equal(t, "xxx", result)
}

func TestToLowerInlineCase_WithSpecialChars(t *testing.T) {
	result := ToLowerInlineCase("user#name")
	assert.Equal(t, "username", result)
}

func TestToLowerInlineCase_OnlySpecialChars(t *testing.T) {
	result := ToLowerInlineCase("!@#")
	assert.Equal(t, "xxx", result)
}

func TestToUpperInlineCase_BasicSnakeCase(t *testing.T) {
	result := ToUpperInlineCase("user_id")
	assert.Equal(t, "USERID", result)
}

func TestToUpperInlineCase_MixedCase(t *testing.T) {
	result := ToUpperInlineCase("Api_Token")
	assert.Equal(t, "APITOKEN", result)
}

func TestToUpperInlineCase_NoUnderscores(t *testing.T) {
	result := ToUpperInlineCase("apitoken")
	assert.Equal(t, "APITOKEN", result)
}

func TestToUpperInlineCase_EmptyString(t *testing.T) {
	result := ToUpperInlineCase("")
	assert.Equal(t, "XXX", result)
}

func TestToUpperInlineCase_SingleChar(t *testing.T) {
	result := ToUpperInlineCase("a")
	assert.Equal(t, "A", result)
}

func TestToUpperInlineCase_WithSpecialChars(t *testing.T) {
	result := ToUpperInlineCase("user#name")
	assert.Equal(t, "USERNAME", result)
}

func TestToUpperInlineCase_MultipleUnderscores(t *testing.T) {
	result := ToUpperInlineCase("user__name___id")
	assert.Equal(t, "USERNAMEID", result)
}

func TestToUpperInlineCase_OnlySpecialChars(t *testing.T) {
	result := ToUpperInlineCase("!@#")
	assert.Equal(t, "XXX", result)
}
