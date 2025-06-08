package query

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
			case []string:
				builder.WriteString(strings.Join(v, ","))
			case []int:
				strs := make([]string, len(v))
				for i, num := range v {
					strs[i] = strconv.Itoa(num)
				}
				builder.WriteString(strings.Join(strs, ","))
			default:
				builder.WriteString(fmt.Sprintf("%v", v))
			}
		}
	}
	compositeKeyName := qb.getCompositeKeyName(parts)
	return expression.Key(compositeKeyName).Equal(expression.Value(builder.String()))
}

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

func (qb *QueryBuilder) formatAttributeValue(value interface{}) string {
   switch v := value.(type) {
   case string:
   	return v
   case int, int8, int16, int32, int64:
   	return fmt.Sprintf("%d", v)
   case uint, uint8, uint16, uint32, uint64:
   	return fmt.Sprintf("%d", v)
   case float32, float64:
   	return fmt.Sprintf("%g", v)
   case bool:
   	if v {
   		return "true"
   	}
   	return "false"
   case []string:
   	return strings.Join(v, ",")
   case []int:
   	strs := make([]string, len(v))
   	for i, num := range v {
   		strs[i] = strconv.Itoa(num)
   	}
   	return strings.Join(strs, ",")
   case []int8:
   	strs := make([]string, len(v))
   	for i, num := range v {
   		strs[i] = fmt.Sprintf("%d", num)
   	}
   	return strings.Join(strs, ",")
   case []int16:
   	strs := make([]string, len(v))
   	for i, num := range v {
   		strs[i] = fmt.Sprintf("%d", num)
   	}
   	return strings.Join(strs, ",")
   case []int32:
   	strs := make([]string, len(v))
   	for i, num := range v {
   		strs[i] = fmt.Sprintf("%d", num)
   	}
   	return strings.Join(strs, ",")
   case []int64:
   	strs := make([]string, len(v))
   	for i, num := range v {
   		strs[i] = fmt.Sprintf("%d", num)
   	}
   	return strings.Join(strs, ",")
   case []uint:
   	strs := make([]string, len(v))
   	for i, num := range v {
   		strs[i] = fmt.Sprintf("%d", num)
   	}
   	return strings.Join(strs, ",")
   case []uint8:
   	strs := make([]string, len(v))
   	for i, num := range v {
   		strs[i] = fmt.Sprintf("%d", num)
   	}
   	return strings.Join(strs, ",")
   case []uint16:
   	strs := make([]string, len(v))
   	for i, num := range v {
   		strs[i] = fmt.Sprintf("%d", num)
   	}
   	return strings.Join(strs, ",")
   case []uint32:
   	strs := make([]string, len(v))
   	for i, num := range v {
   		strs[i] = fmt.Sprintf("%d", num)
   	}
   	return strings.Join(strs, ",")
   case []uint64:
   	strs := make([]string, len(v))
   	for i, num := range v {
   		strs[i] = fmt.Sprintf("%d", num)
   	}
   	return strings.Join(strs, ",")
   case []float32:
   	strs := make([]string, len(v))
   	for i, num := range v {
   		strs[i] = fmt.Sprintf("%g", num)
   	}
   	return strings.Join(strs, ",")
   case []float64:
   	strs := make([]string, len(v))
   	for i, num := range v {
   		strs[i] = fmt.Sprintf("%g", num)
   	}
   	return strings.Join(strs, ",")
   default:
   	return fmt.Sprintf("%v", value)
   }
}
`
