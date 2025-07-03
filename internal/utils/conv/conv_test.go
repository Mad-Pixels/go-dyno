package conv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToCamelCase_BasicSnakeCase(t *testing.T) {
	result := toCamelCase("user_id")
	assert.Equal(t, "UserId", result)
}

func TestToCamelCase_BasicKebabCase(t *testing.T) {
	result := toCamelCase("user-name")
	assert.Equal(t, "UserName", result)
}

func TestToCamelCase_HashSeparator(t *testing.T) {
	result := toCamelCase("user#name")
	assert.Equal(t, "UserName", result)
}

func TestToCamelCase_MixedSeparators(t *testing.T) {
	result := toCamelCase("user_name-id#value")
	assert.Equal(t, "UserNameIdValue", result)
}

func TestToCamelCase_MultipleSeparatorsInRow(t *testing.T) {
	result := toCamelCase("user__name")
	assert.Equal(t, "UserName", result)
}

func TestToCamelCase_StartingWithSeparator(t *testing.T) {
	result := toCamelCase("_user_name")
	assert.Equal(t, "UserName", result)
}

func TestToCamelCase_EndingWithSeparator(t *testing.T) {
	result := toCamelCase("user_name_")
	assert.Equal(t, "UserName", result)
}

func TestToCamelCase_EmptyString(t *testing.T) {
	result := toCamelCase("")
	assert.Equal(t, "", result)
}

func TestToCamelCase_SingleChar(t *testing.T) {
	result := toCamelCase("a")
	assert.Equal(t, "A", result)
}

func TestToCamelCase_SingleSeparator(t *testing.T) {
	result := toCamelCase("_")
	assert.Equal(t, "", result)
}

func TestToCamelCase_OnlySeparators(t *testing.T) {
	result := toCamelCase("___")
	assert.Equal(t, "", result)
}

func TestToCamelCase_NoSeparators(t *testing.T) {
	result := toCamelCase("username")
	assert.Equal(t, "Username", result)
}

func TestToCamelCase_AlreadyCamelCase(t *testing.T) {
	result := toCamelCase("UserName")
	assert.Equal(t, "UserName", result)
}

func TestToCamelCase_WithNumbers(t *testing.T) {
	result := toCamelCase("user_id_123")
	assert.Equal(t, "UserId123", result)
}

func TestToCamelCase_WithSpecialChars(t *testing.T) {
	result := toCamelCase("user@name")
	assert.Equal(t, "User@name", result)
}

func TestToCamelCase_UnicodeChars(t *testing.T) {
	result := toCamelCase("пользователь_имя")
	assert.Equal(t, "ПользовательИмя", result)
}

func TestToCamelCase_MixedCase(t *testing.T) {
	result := toCamelCase("User_NAME")
	assert.Equal(t, "UserNAME", result)
}

func TestReservedWords_ControlFlow(t *testing.T) {
	assert.True(t, reservedWords["break"])
	assert.True(t, reservedWords["continue"])
	assert.True(t, reservedWords["return"])
	assert.True(t, reservedWords["fallthrough"])
	assert.True(t, reservedWords["goto"])
}

func TestReservedWords_Conditionals(t *testing.T) {
	assert.True(t, reservedWords["if"])
	assert.True(t, reservedWords["else"])
	assert.True(t, reservedWords["for"])
	assert.True(t, reservedWords["range"])
	assert.True(t, reservedWords["switch"])
	assert.True(t, reservedWords["case"])
	assert.True(t, reservedWords["default"])
	assert.True(t, reservedWords["select"])
}

func TestReservedWords_Declarations(t *testing.T) {
	assert.True(t, reservedWords["var"])
	assert.True(t, reservedWords["const"])
	assert.True(t, reservedWords["type"])
	assert.True(t, reservedWords["struct"])
	assert.True(t, reservedWords["interface"])
	assert.True(t, reservedWords["map"])
	assert.True(t, reservedWords["chan"])
	assert.True(t, reservedWords["func"])
	assert.True(t, reservedWords["package"])
	assert.True(t, reservedWords["import"])
	assert.True(t, reservedWords["defer"])
}

func TestReservedWords_SpecialKeywords(t *testing.T) {
	assert.True(t, reservedWords["any"])
	assert.True(t, reservedWords["go"])
}

func TestReservedWords_NonReserved(t *testing.T) {
	assert.False(t, reservedWords["user"])
	assert.False(t, reservedWords["name"])
	assert.False(t, reservedWords["value"])
	assert.False(t, reservedWords["string"])
	assert.False(t, reservedWords["int"])
	assert.False(t, reservedWords["bool"])
}

func TestReservedWords_CaseSensitive(t *testing.T) {
	assert.True(t, reservedWords["type"])
	assert.False(t, reservedWords["Type"])
	assert.False(t, reservedWords["TYPE"])
}

func TestReservedWords_EmptyAndSpaces(t *testing.T) {
	assert.False(t, reservedWords[""])
	assert.False(t, reservedWords[" type "])
	assert.False(t, reservedWords["type "])
	assert.False(t, reservedWords[" type"])
}
