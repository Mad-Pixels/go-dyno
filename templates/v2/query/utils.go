package query

// QueryBuilderUtilsTemplate with simplified formatAttributeValue
const QueryBuilderUtilsTemplate = `
func (qb *QueryBuilder) hasAllKeys(parts []CompositeKeyPart) bool {
	for _, part := range parts {
		if !part.IsConstant && !qb.UsedKeys[part.Value] {
			return false
		}
	}
	return true
}

func (qb *QueryBuilder) buildCompositeKeyCondition(parts []CompositeKeyPart) expression.KeyConditionBuilder {
	compositeKeyName := qb.getCompositeKeyName(parts)
	compositeValue := qb.buildCompositeKeyValue(parts)
	return expression.Key(compositeKeyName).Equal(expression.Value(compositeValue))
}

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

func (qb *QueryBuilder) formatAttributeValue(value interface{}) string {
	if value == nil {
		return ""
	}

	// Fast path for common simple types
	switch v := value.(type) {
	case string:
		return v
	case bool:
		if v {
			return "true"
		}
		return "false"
	}

	// Use AWS SDK marshaling for complex types and numbers
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
		// Fallback to string representation
		return fmt.Sprintf("%v", value)
	}
}
`
