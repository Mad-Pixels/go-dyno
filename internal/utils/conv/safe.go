package conv

import (
	"strings"
	"unicode"
)

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
