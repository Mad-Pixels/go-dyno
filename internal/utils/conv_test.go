package utils

import (
	"testing"

	"github.com/Mad-Pixels/go-dyno/internal/schema/common"
	"github.com/stretchr/testify/assert"
)

func TestToSafeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "   -Привет-123_foo-",
			expected: "x123_foo",
		},
		{
			input:    "$$$",
			expected: "xxx",
		},
		{
			input:    "9lives",
			expected: "x9lives",
		},
		{
			input:    "Hello-World",
			expected: "Hello_World",
		},
		{
			input:    "_foo_bar_",
			expected: "foo_bar",
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, ToSafeName(tt.input), "input: %q", tt.input)
	}
}

func TestToUpperCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "   -Привет-123_foo-",
			expected: "X123Foo",
		},
		{
			input:    "$$$",
			expected: "Xxx",
		},
		{
			input:    "9lives",
			expected: "X9lives",
		},
		{
			input:    "Hello-World",
			expected: "HelloWorld",
		},
		{
			input:    "_foo_bar_",
			expected: "FooBar",
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, ToUpperCamelCase(tt.input), "input: %q", tt.input)
	}
}

func TestToLowerCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "   -Привет-123_foo-",
			expected: "x123Foo",
		},
		{
			input:    "$$$",
			expected: "xxx",
		},
		{
			input:    "9lives",
			expected: "x9lives",
		},
		{
			input:    "Hello-World",
			expected: "helloWorld",
		},
		{
			input:    "_foo_bar_",
			expected: "fooBar",
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, ToLowerCamelCase(tt.input), "input: %q", tt.input)
	}
}

func TestToGolangBaseType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"S", "string"},
		{"N", "int"},
		{"B", "bool"},
		{"SS", "[]string"},
		{"NS", "[]int"},
		{"UNKNOWN", "any"},
		{"", "any"},
	}

	for _, tt := range tests {
		result := ToGolangBaseType(tt.input)
		assert.Equal(t, tt.expected, result, "input: %q", tt.input)
	}
}

func TestToGolangZeroType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"S", `""`},
		{"N", "0"},
		{"B", "false"},
		{"SS", "nil"},
		{"NS", "nil"},
		{"UNKNOWN", "nil"},
		{"", "nil"},
	}

	for _, tt := range tests {
		result := ToGolangZeroType(tt.input)
		assert.Equal(t, tt.expected, result, "input: %q", tt.input)
	}
}

func TestToGolangAttrType(t *testing.T) {
	attrs := []common.Attribute{
		{Name: "id", Type: "S"},
		{Name: "count", Type: "N"},
		{Name: "is_active", Type: "B"},
		{Name: "tags", Type: "SS"},
		{Name: "scores", Type: "NS"},
	}

	tests := []struct {
		attrName string
		expected string
	}{
		{"id", "string"},
		{"count", "int"},
		{"is_active", "bool"},
		{"tags", "[]string"},
		{"scores", "[]int"},
		{"missing", "any"},
	}

	for _, tt := range tests {
		result := ToGolangAttrType(tt.attrName, attrs)
		assert.Equal(t, tt.expected, result, "attrName: %q", tt.attrName)
	}
}
