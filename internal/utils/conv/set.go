package conv

import "slices"

// AvailableKeys returns a sorted list of keys from the given map.
//
// This is typically used to extract and display all supported options
// (e.g. valid types, projection modes, etc.) for error messages or documentation.
//
// Example:
//
//	input := map[string]bool{"A": true, "C": true, "B": true}
//	output := AvailableKeys(input) // → []string{"A", "B", "C"}
func AvailableKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
