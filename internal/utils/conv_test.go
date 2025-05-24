package utils

import (
	"testing"

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
