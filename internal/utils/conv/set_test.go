package conv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAvailableKeys_BasicMap(t *testing.T) {
	input := map[string]bool{"A": true, "C": true, "B": true}
	result := AvailableKeys(input)
	expected := []string{"A", "B", "C"}
	assert.Equal(t, expected, result)
}

func TestAvailableKeys_EmptyMap(t *testing.T) {
	input := map[string]bool{}
	result := AvailableKeys(input)
	expected := []string{}
	assert.Equal(t, expected, result)
}

func TestAvailableKeys_SingleKey(t *testing.T) {
	input := map[string]bool{"test": true}
	result := AvailableKeys(input)
	expected := []string{"test"}
	assert.Equal(t, expected, result)
}

func TestAvailableKeys_MixedBoolValues(t *testing.T) {
	input := map[string]bool{"A": true, "B": false, "C": true}
	result := AvailableKeys(input)
	expected := []string{"A", "B", "C"}
	assert.Equal(t, expected, result)
}
