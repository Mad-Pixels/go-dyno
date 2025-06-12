package conv

import "strings"

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

// ToUpperInlineCase converts a string to uppercase without underscores.
//
// Examples:
//
//	ToUpperInlineCase("user_id")   → "USERID"
//	ToUpperInlineCase("Api_Token") → "APITOKEN"
func ToUpperInlineCase(s string) string {
	res := strings.ReplaceAll(ToSafeName(s), "_", "")
	return strings.ToUpper(res)
}
