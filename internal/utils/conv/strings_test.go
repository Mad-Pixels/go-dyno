package conv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimLeftN_BasicUsage(t *testing.T) {
	result := TrimLeftN("[]int", 2)
	assert.Equal(t, "int", result)
}

func TestTrimLeftN_WithSymbols(t *testing.T) {
	result := TrimLeftN("##Hello", 2)
	assert.Equal(t, "Hello", result)
}

func TestTrimLeftN_ZeroStart(t *testing.T) {
	result := TrimLeftN("GoLang", 0)
	assert.Equal(t, "GoLang", result)
}

func TestTrimLeftN_EmptyString(t *testing.T) {
	result := TrimLeftN("", 5)
	assert.Equal(t, "", result)
}

func TestTrimLeftN_NegativeStart(t *testing.T) {
	result := TrimLeftN("test", -1)
	assert.Equal(t, "test", result)
}

func TestIsFloatType_Float32(t *testing.T) {
	result := IsFloatType("float32")
	assert.True(t, result)
}

func TestIsFloatType_Float64(t *testing.T) {
	result := IsFloatType("float64")
	assert.True(t, result)
}

func TestIsFloatType_Double(t *testing.T) {
	result := IsFloatType("double")
	assert.False(t, result)
}

func TestIsFloatType_EmptyString(t *testing.T) {
	result := IsFloatType("")
	assert.False(t, result)
}

func TestIsFloatType_CaseSensitive(t *testing.T) {
	result := IsFloatType("Float32")
	assert.False(t, result)
}

func TestIsFloatType_WithSpaces(t *testing.T) {
	result := IsFloatType(" float32 ")
	assert.False(t, result)
}
