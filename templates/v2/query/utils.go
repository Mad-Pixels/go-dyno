package query

// QueryBuilderUtilsTemplate ...
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
			builder.WriteString(qb.formatAttributeValue(value))
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
   case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
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
   	return qb.formatIntSlice(len(v), func(i int) int64 { return int64(v[i]) })
   case []int8:
   	return qb.formatIntSlice(len(v), func(i int) int64 { return int64(v[i]) })
   case []int16:
   	return qb.formatIntSlice(len(v), func(i int) int64 { return int64(v[i]) })
   case []int32:
   	return qb.formatIntSlice(len(v), func(i int) int64 { return int64(v[i]) })
   case []int64:
   	return qb.formatIntSlice(len(v), func(i int) int64 { return v[i] })
   case []uint:
   	return qb.formatIntSlice(len(v), func(i int) int64 { return int64(v[i]) })
   case []uint8:
   	return qb.formatIntSlice(len(v), func(i int) int64 { return int64(v[i]) })
   case []uint16:
   	return qb.formatIntSlice(len(v), func(i int) int64 { return int64(v[i]) })
   case []uint32:
   	return qb.formatIntSlice(len(v), func(i int) int64 { return int64(v[i]) })
   case []uint64:
   	return qb.formatIntSlice(len(v), func(i int) int64 { return int64(v[i]) })
   
   case []float32:
   	return qb.formatFloatSlice(len(v), func(i int) float64 { return float64(v[i]) })
   case []float64:
   	return qb.formatFloatSlice(len(v), func(i int) float64 { return v[i] })
   
   default:
   	return fmt.Sprintf("%v", value)
   }
}

func (qb *QueryBuilder) formatIntSlice(length int, getValue func(int) int64) string {
	if length == 0 {
		return ""
	}
	
	strs := make([]string, length)
	for i := 0; i < length; i++ {
		strs[i] = strconv.FormatInt(getValue(i), 10)
	}
	return strings.Join(strs, ",")
}

func (qb *QueryBuilder) formatFloatSlice(length int, getValue func(int) float64) string {
	if length == 0 {
		return ""
	}
	
	strs := make([]string, length)
	for i := 0; i < length; i++ {
		strs[i] = strconv.FormatFloat(getValue(i), 'g', -1, 64)
	}
	return strings.Join(strs, ",")
}
`
