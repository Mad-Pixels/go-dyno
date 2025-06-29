package query

// QueryBuilderUtilsTemplate provides utility functions for composite key handling
const QueryBuilderUtilsTemplate = `
// hasAllKeys checks if all non-constant parts of a composite key are available.
// Used to determine if a composite key can be fully constructed from current conditions.
// Constants are always available, variables must be present in UsedKeys.
func (qb *QueryBuilder) hasAllKeys(parts []CompositeKeyPart) bool {
	for _, part := range parts {
		if !part.IsConstant && !qb.UsedKeys[part.Value] {
			return false
		}
	}
	return true
}

// buildCompositeKeyCondition creates a key condition for composite keys.
// Combines multiple key parts into a single equality condition using "#" separator.
// Used internally by the index selection algorithm for complex key structures.
func (qb *QueryBuilder) buildCompositeKeyCondition(parts []CompositeKeyPart) expression.KeyConditionBuilder {
	compositeKeyName := qb.getCompositeKeyName(parts)
	compositeValue := qb.buildCompositeKeyValue(parts)
	return expression.Key(compositeKeyName).Equal(expression.Value(compositeValue))
}

// getCompositeKeyName generates the attribute name for a composite key.
// For single parts, returns the part name directly.
// For multiple parts, joins them with "#" separator for DynamoDB storage.
// Example: ["user", "tenant"] -> "user#tenant"
func (qb *QueryBuilder) getCompositeKeyName(parts []CompositeKeyPart) string {
	switch len(parts) {
	case 0:
		return ""
	case 1:
		return parts[0].Value
	default:
		names := make([]string, len(parts))
		for i, part := range parts {
			names[i] = part.Value
		}
		return strings.Join(names, "#")
	}
}

// buildCompositeKeyValue constructs the actual value for a composite key.
// Combines constant values and variable values from query attributes.
// Uses "#" separator to create a single string value for DynamoDB.
// Example: constant "USER" + variable "123" -> "USER#123"
func (qb *QueryBuilder) buildCompositeKeyValue(parts []CompositeKeyPart) string {
	if len(parts) == 0 {
		return ""
	}
	values := make([]string, len(parts))
	for i, part := range parts {
		if part.IsConstant {
			values[i] = part.Value
		} else {
			values[i] = qb.formatAttributeValue(qb.Attributes[part.Value])
		}
	}
	return strings.Join(values, "#")
}

// formatAttributeValue converts any Go value to its string representation for composite keys.
// Provides optimized fast paths for common types (string, bool) and proper handling
// of complex types through AWS SDK marshaling. Ensures consistent string formatting
// for reliable composite key construction.
func (qb *QueryBuilder) formatAttributeValue(value interface{}) string {
	if value == nil {
		return ""
	}
	
	switch v := value.(type) {
	case string:
		return v
	case bool:
		if v {
			return "true"
		}
		return "false"
	}
	
	av, err := attributevalue.Marshal(value)
	if err != nil {
		return fmt.Sprintf("%v", value)
	}
	
	switch typed := av.(type) {
	case *types.AttributeValueMemberS:
		return typed.Value
	case *types.AttributeValueMemberN:
		return typed.Value
	case *types.AttributeValueMemberBOOL:
		if typed.Value {
			return "true"
		}
		return "false"
	case *types.AttributeValueMemberSS:
		return strings.Join(typed.Value, ",")
	case *types.AttributeValueMemberNS:
		return strings.Join(typed.Value, ",")
	default:
		return fmt.Sprintf("%v", value)
	}
}
`
