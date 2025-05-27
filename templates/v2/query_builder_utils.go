package v2

// QueryBuilderUtilsTemplate generates optimized utility functions for QueryBuilder composite key operations.
const QueryBuilderUtilsTemplate = `
// hasAllKeys validates that all non-constant parts of a composite key have corresponding values.
// Used internally by QueryBuilder to determine if a composite key can be built before
// attempting to create DynamoDB query conditions.
//
// Parameters:
//   - parts: Slice of CompositeKeyPart defining the key structure
//
// Returns true only if all dynamic parts have values in qb.UsedKeys, false otherwise.
//
// Example:
//
//	parts := []CompositeKeyPart{
//		{IsConstant: false, Value: "user_id"},    // needs value
//		{IsConstant: true,  Value: "active"},     // constant, always available
//		{IsConstant: false, Value: "tenant_id"},  // needs value
//	}
//	// Returns true only if both "user_id" and "tenant_id" are in qb.UsedKeys
func (qb *QueryBuilder) hasAllKeys(parts []CompositeKeyPart) bool {
	for _, part := range parts {
		if !part.IsConstant && !qb.UsedKeys[part.Value] {
			return false
		}
	}
	return true
}

// buildCompositeKeyCondition constructs a DynamoDB KeyConditionBuilder for composite keys.
// Optimized for performance with pre-allocated string builder and type-specific value formatting.
//
// The function builds the composite key value by joining parts with "#" separator:
// - Constant parts: use literal string values
// - Dynamic parts: format attribute values using type-specific conversion
//
// Example:
//
//	parts := []CompositeKeyPart{
//		{IsConstant: false, Value: "user_id"},   // qb.Attributes["user_id"] = "123"
//		{IsConstant: true,  Value: "active"},    // literal "active"
//		{IsConstant: false, Value: "year"},      // qb.Attributes["year"] = 2024
//	}
//	// Builds condition: expression.Key("user_id#active#year").Equal(expression.Value("123#active#2024"))
func (qb *QueryBuilder) buildCompositeKeyCondition(parts []CompositeKeyPart) expression.KeyConditionBuilder {
	// Pre-allocate builder capacity: average 12 characters per part (name + separator + value)
	estimatedSize := len(parts) * 12
	var builder strings.Builder
	builder.Grow(estimatedSize)

	for i, part := range parts {
		if i > 0 {
			builder.WriteByte('#')
		}

		if part.IsConstant {
			builder.WriteString(part.Value)
		} else {
			value := qb.Attributes[part.Value]
			switch v := value.(type) {
			case string:
				builder.WriteString(v)
			case int:
				builder.WriteString(strconv.Itoa(v))
			case int64:
				builder.WriteString(strconv.FormatInt(v, 10))
			case bool:
				if v {
					builder.WriteString("true")
				} else {
					builder.WriteString("false")
				}
			default:
				builder.WriteString(fmt.Sprintf("%v", v))
			}
		}
	}
	compositeKeyName := qb.getCompositeKeyName(parts)
	return expression.Key(compositeKeyName).Equal(expression.Value(builder.String()))
}

// getCompositeKeyName generates the attribute name for a composite key from its parts.
//
// Parameters:
//   - parts: Slice of CompositeKeyPart defining the key structure
//
// Returns the composite key attribute name (e.g., "user_id#status#year").
//
// Example:
//
//	parts := []CompositeKeyPart{
//		{Value: "user_id"},
//		{Value: "status"},
//		{Value: "created_year"}
//	}
//	// Returns: "user_id#status#created_year"
func (qb *QueryBuilder) getCompositeKeyName(parts []CompositeKeyPart) string {
	switch len(parts) {
	case 0:
		return ""
	case 1:
		return parts[0].Value
	case 2, 3:
		names := make([]string, 0, len(parts))
		for _, part := range parts {
			names = append(names, part.Value)
		}
		return strings.Join(names, "#")
	default:
		estimatedSize := len(parts) * 10
		var builder strings.Builder
		builder.Grow(estimatedSize)

		for i, part := range parts {
			if i > 0 {
				builder.WriteByte('#')
			}
			builder.WriteString(part.Value)
		}
		return builder.String()
	}
}

// buildCompositeKeyValue constructs the actual value string for a composite key.
// Similar to buildCompositeKeyCondition but focuses on value generation rather than condition building.
//
// Used when setting composite key values in QueryBuilder attributes before building conditions.
// The generated value is what gets stored in DynamoDB as the composite key.
//
// Performance features:
// - Pre-allocated string builder with capacity estimation
// - Single-part optimization (direct value return)
// - Small keys optimization for 2-3 parts
// - Delegates to formatAttributeValue for consistent type handling
//
// Example:
//
//	parts := []CompositeKeyPart{
//		{IsConstant: false, Value: "user_id"},   // qb.Attributes["user_id"] = "user123"
//		{IsConstant: true,  Value: "active"},    // literal "active"
//		{IsConstant: false, Value: "is_public"}, // qb.Attributes["is_public"] = 1
//	}
//	// Returns: "user123#active#1"
func (qb *QueryBuilder) buildCompositeKeyValue(parts []CompositeKeyPart) string {
	switch len(parts) {
	case 0:
		return ""

	case 1:
		if parts[0].IsConstant {
			return parts[0].Value
		}
		return qb.formatAttributeValue(qb.Attributes[parts[0].Value])

	case 2:
		var part1, part2 string

		if parts[0].IsConstant {
			part1 = parts[0].Value
		} else {
			part1 = qb.formatAttributeValue(qb.Attributes[parts[0].Value])
		}
		if parts[1].IsConstant {
			part2 = parts[1].Value
		} else {
			part2 = qb.formatAttributeValue(qb.Attributes[parts[1].Value])
		}
		return part1 + "#" + part2

	case 3:
		var part1, part2, part3 string

		if parts[0].IsConstant {
			part1 = parts[0].Value
		} else {
			part1 = qb.formatAttributeValue(qb.Attributes[parts[0].Value])
		}
		if parts[1].IsConstant {
			part2 = parts[1].Value
		} else {
			part2 = qb.formatAttributeValue(qb.Attributes[parts[1].Value])
		}
		if parts[2].IsConstant {
			part3 = parts[2].Value
		} else {
			part3 = qb.formatAttributeValue(qb.Attributes[parts[2].Value])
		}
		return part1 + "#" + part2 + "#" + part3
	default:
		estimatedSize := len(parts) * 12
		var builder strings.Builder
		builder.Grow(estimatedSize)

		for i, part := range parts {
			if i > 0 {
				builder.WriteByte('#')
			}

			if part.IsConstant {
				builder.WriteString(part.Value)
			} else {
				value := qb.Attributes[part.Value]
				builder.WriteString(qb.formatAttributeValue(value))
			}
		}
		return builder.String()
	}
}

// formatAttributeValue converts Go values to their string representation for DynamoDB composite keys.
// Provides consistent, DynamoDB-compatible formatting across all QueryBuilder operations.
//
// Type conversion rules:
// - string: Pass through unchanged
// - int/int64: Convert to decimal string representation
// - bool: Convert to "1" (true) or "0" (false) for DynamoDB compatibility
// - other types: Use fmt.Sprintf as fallback (slower but comprehensive)
//
// The bool → "1"/"0" conversion aligns with DynamoDB's numeric boolean storage pattern
// commonly used in applications (e.g., is_public: 1 for true, 0 for false).
//
// Parameters:
//   - value: The Go value to format for DynamoDB storage
//
// Returns string representation suitable for composite key construction.
//
// Example:
//
//	formatAttributeValue("user123")    // → "user123"
//	formatAttributeValue(42)           // → "42"
//	formatAttributeValue(true)         // → "1"
//	formatAttributeValue(false)        // → "0"
func (qb *QueryBuilder) formatAttributeValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case bool:
		if v {
			return "1" // DynamoDB-compatible boolean: true = 1
		}
		return "0" // DynamoDB-compatible boolean: false = 0
	default:
		return fmt.Sprintf("%v", value)
	}
}
`
