package conv

// TrimLeftN returns a substring of the input string `s` with the first `start` characters removed.
//
// If `start` is greater than or equal to the length of the string, it returns an empty string.
// This function is useful in templates and code generation where controlled slicing is needed.
//
// Examples:
//
//	TrimLeftN("[]int", 2)       → "int"
//	TrimLeftN("##Hello", 2)     → "Hello"
//	TrimLeftN("GoLang", 0)      → "GoLang"
//	TrimLeftN("short", 10)      → ""
func TrimLeftN(s string, start int) string {
	if start < 0 {
		return s
	}
	if start >= len(s) {
		return ""
	}
	return s[start:]
}

// IsFloatType checks whether the provided Go type name is a floating-point type.
//
// This is commonly used in code generation scenarios to determine how numeric types
// should be handled (e.g., when converting to DynamoDB number sets).
//
// Examples:
//
//	IsFloatType("float32") → true
//	IsFloatType("float64") → true
//	IsFloatType("int64")   → false
//	IsFloatType("double")  → false
func IsFloatType(goType string) bool {
	return goType == "float32" || goType == "float64"
}
