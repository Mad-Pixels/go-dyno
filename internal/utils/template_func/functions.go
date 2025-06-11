package templatefunc

import (
	"strings"
	"unicode"
)

// ToUpperCamelCase converts a string into UpperCamelCase format,
// ensuring the result is safe for use in Go identifiers.
//
// Examples:
//
//	ToUpperCamelCase("user_id")         → "UserId"
//	ToUpperCamelCase("user-name")       → "UserName"
//	ToUpperCamelCase("1type")           → "X1type"
//	ToUpperCamelCase("full#access")     → "FullAccess"
//	ToUpperCamelCase("!@#special-case") → "XxxSpecialCase"
func ToUpperCamelCase(s string) string {
	res := ToSafeName(toCamelCase(s))
	return strings.ToUpper(res[:1]) + res[1:]
}

// ToLowerCamelCase converts a string into lowerCamelCase format,
// ensuring the result is safe for use in Go identifiers.
//
// Examples:
//
//	ToLowerCamelCase("user_id")  → "userId"
//	ToLowerCamelCase("Type")     → "type"
//	ToLowerCamelCase("1invalid") → "x1invalid"
func ToLowerCamelCase(s string) string {
	res := ToSafeName(toCamelCase(s))
	return strings.ToLower(res[:1]) + res[1:]
}

// ToLowerInlineCase converts a string to lowercase without underscores.
//
// Examples:
//
//	ToLowerInlineCase("user_id")          → "userid"
//	ToLowerInlineCase("snake_case_value") → "snakecasevalue"
func ToLowerInlineCase(s string) string {
	res := strings.ReplaceAll(ToSafeName(s), "_", "")
	return strings.ToLower(res)
}

// ToUpperinlineCase converts a string to uppercase without underscores.
//
// Examples:
//
//	ToUpperinlineCase("user_id")   → "USERID"
//	ToUpperinlineCase("Api_Token") → "APITOKEN"
func ToUpperinlineCase(s string) string {
	res := strings.ReplaceAll(ToSafeName(s), "_", "")
	return strings.ToUpper(res)
}

// ToSafeName sanitizes any string into a Go-safe identifier:
// - replaces invalid characters with underscores
// - prefixes with 'x' if it starts with a number or is a reserved keyword
//
// Examples:
//
//	ToSafeName("1test")       → "x1test"
//	ToSafeName("type")        → "xtype"
//	ToSafeName("hello_world") → "hello_world"
//	ToSafeName("$$$abc")      → "abc"
func ToSafeName(s string) string {
	s = strings.TrimFunc(s, func(r rune) bool {
		return (r < 'A' || r > 'Z') &&
			(r < 'a' || r > 'z') &&
			(r < '0' || r > '9')
	})
	var b strings.Builder
	for _, r := range s {
		switch {
		case (r >= 'A' && r <= 'Z') ||
			(r >= 'a' && r <= 'z') ||
			(r >= '0' && r <= '9'):
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	switch {
	case b.Len() == 0:
		return "xxx"
	case unicode.IsDigit(rune(b.String()[0])) || reservedWords[strings.ToLower(b.String())]:
		return "x" + b.String()
	default:
		return b.String()
	}
}

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
	if start >= len(s) {
		return ""
	}
	return s[start:]
}
