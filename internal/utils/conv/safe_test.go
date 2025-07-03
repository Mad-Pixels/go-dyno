package conv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSafeName_StartsWithDigit(t *testing.T) {
	result := ToSafeName("1test")
	assert.Equal(t, "x1test", result)
}

func TestToSafeName_ReservedKeyword(t *testing.T) {
	result := ToSafeName("type")
	assert.Equal(t, "xtype", result)
}

func TestToSafeName_ValidIdentifier(t *testing.T) {
	result := ToSafeName("hello_world")
	assert.Equal(t, "hello_world", result)
}

func TestToSafeName_SpecialCharsAtStart(t *testing.T) {
	result := ToSafeName("$$$abc")
	assert.Equal(t, "abc", result)
}

func TestToSafeName_EmptyString(t *testing.T) {
	result := ToSafeName("")
	assert.Equal(t, "xxx", result)
}

func TestToSafeName_OnlySpecialChars(t *testing.T) {
	result := ToSafeName("!@#$%")
	assert.Equal(t, "xxx", result)
}

func TestToSafeName_ValidAlphanumeric(t *testing.T) {
	result := ToSafeName("test123")
	assert.Equal(t, "test123", result)
}

func TestToSafeName_MixedValidInvalid(t *testing.T) {
	result := ToSafeName("test@name#123")
	assert.Equal(t, "test_name_123", result)
}

func TestToSafeName_AllReservedWords_Break(t *testing.T) {
	result := ToSafeName("break")
	assert.Equal(t, "xbreak", result)
}

func TestToSafeName_AllReservedWords_Continue(t *testing.T) {
	result := ToSafeName("continue")
	assert.Equal(t, "xcontinue", result)
}

func TestToSafeName_AllReservedWords_Return(t *testing.T) {
	result := ToSafeName("return")
	assert.Equal(t, "xreturn", result)
}

func TestToSafeName_AllReservedWords_Var(t *testing.T) {
	result := ToSafeName("var")
	assert.Equal(t, "xvar", result)
}

func TestToSafeName_AllReservedWords_Func(t *testing.T) {
	result := ToSafeName("func")
	assert.Equal(t, "xfunc", result)
}

func TestToSafeName_ReservedWordMixedCase(t *testing.T) {
	result := ToSafeName("Type")
	assert.Equal(t, "xType", result)
}

func TestToSafeName_ReservedWordUpperCase(t *testing.T) {
	result := ToSafeName("TYPE")
	assert.Equal(t, "xTYPE", result)
}

func TestToSafeName_NonReservedWord(t *testing.T) {
	result := ToSafeName("username")
	assert.Equal(t, "username", result)
}

func TestToSafeName_SingleChar(t *testing.T) {
	result := ToSafeName("a")
	assert.Equal(t, "a", result)
}

func TestToSafeName_SingleDigit(t *testing.T) {
	result := ToSafeName("1")
	assert.Equal(t, "x1", result)
}

func TestToSafeName_SingleSpecialChar(t *testing.T) {
	result := ToSafeName("@")
	assert.Equal(t, "xxx", result)
}

func TestToSafeName_UnicodeChars(t *testing.T) {
	result := ToSafeName("тест")
	assert.Equal(t, "xxx", result)
}

func TestToSafeName_LeadingSpecialChars(t *testing.T) {
	result := ToSafeName("!!!test")
	assert.Equal(t, "test", result)
}

func TestToSafeName_TrailingSpecialChars(t *testing.T) {
	result := ToSafeName("test!!!")
	assert.Equal(t, "test", result)
}

func TestToSafeName_LeadingAndTrailingSpecialChars(t *testing.T) {
	result := ToSafeName("!!!test!!!")
	assert.Equal(t, "test", result)
}

func TestToSafeName_SpecialCharsInMiddle(t *testing.T) {
	result := ToSafeName("hello@world#test")
	assert.Equal(t, "hello_world_test", result)
}

func TestToSafeName_ConsecutiveSpecialChars(t *testing.T) {
	result := ToSafeName("hello@@world")
	assert.Equal(t, "hello__world", result)
}

func TestToSafeName_MixedCaseValid(t *testing.T) {
	result := ToSafeName("HelloWorld")
	assert.Equal(t, "HelloWorld", result)
}

func TestToSafeName_UnderscoreOnly(t *testing.T) {
	result := ToSafeName("_")
	assert.Equal(t, "xxx", result)
}

func TestToSafeName_ValidWithUnderscores(t *testing.T) {
	result := ToSafeName("hello_world_123")
	assert.Equal(t, "hello_world_123", result)
}

func TestToSafeName_NumbersInMiddle(t *testing.T) {
	result := ToSafeName("test123abc")
	assert.Equal(t, "test123abc", result)
}

func TestToSafeName_StartWithDigitAfterTrim(t *testing.T) {
	result := ToSafeName("$$$123abc")
	assert.Equal(t, "x123abc", result)
}

func TestToSafeName_ReservedAfterTrim(t *testing.T) {
	result := ToSafeName("$$$type")
	assert.Equal(t, "xtype", result)
}

func TestToSafeName_SpacesAndTabs(t *testing.T) {
	result := ToSafeName("hello world\ttest")
	assert.Equal(t, "hello_world_test", result)
}

func TestToSafeName_NewlinesAndCarriageReturns(t *testing.T) {
	result := ToSafeName("hello\nworld\rtest")
	assert.Equal(t, "hello_world_test", result)
}
