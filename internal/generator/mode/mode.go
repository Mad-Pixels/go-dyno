// Package mode defines the generation modes for code generation.
//
// It provides a type-safe enum for controlling what code is generated,
// allowing users to choose between generating all code (ALL) or just
// the minimum required code (MIN).
package mode

import (
	"strings"

	"github.com/Mad-Pixels/go-dyno/internal/logger"
	"github.com/Mad-Pixels/go-dyno/internal/utils/conv"
)

// Mode represents the code generation mode.
type Mode string

const (
	// ALL generates complete code with all features (default).
	ALL Mode = "ALL"

	// MIN generates minimal code with only essential functionality.
	MIN Mode = "MIN"
)

// Valid modes for validation.
var validModes = map[Mode]bool{
	ALL: true,
	MIN: true,
}

// String returns the string representation of the Mode.
func (m Mode) String() string {
	return string(m)
}

// IsValid checks if the mode is a valid generation mode.
func (m Mode) IsValid() bool {
	return validModes[m]
}

// ParseMode parses a string into a Mode type with case-insensitive matching.
// Returns the parsed mode and an error if the string is not a valid mode.
func ParseMode(s string) (Mode, error) {
	mode := Mode(strings.ToUpper(strings.TrimSpace(s)))
	if !mode.IsValid() {
		return "", logger.NewFailure("invalid generation mode", nil).
			With("mode", s).
			With("available", GetAvailableModes())
	}
	return mode, nil
}

// MustParseMode parses a string into a Mode type and panics on error.
// Should only be used in tests or when the input is guaranteed to be valid.
func MustParseMode(s string) Mode {
	mode, err := ParseMode(s)
	if err != nil {
		panic(err)
	}
	return mode
}

// GetDefault returns the default generation mode.
func GetDefault() Mode {
	return ALL
}

// GetAvailableModes returns a slice of all valid modes sorted alphabetically.
func GetAvailableModes() []string {
	stringModes := make(map[string]bool, len(validModes))
	for mode := range validModes {
		stringModes[string(mode)] = true
	}
	return conv.AvailableKeys(stringModes)
}
