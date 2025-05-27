package v2

// QueryBuilderUtilsTemplate
const QueryBuilderUtilsTemplate = `
func (qb *QueryBuilder) hasAllKeys(parts []CompositeKeyPart) bool {
    for _, part := range parts {
        if !part.IsConstant && !qb.UsedKeys[part.Value] {
            return false 
        }
    }
    return true
}

// buildCompositeKeyCondition 
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
            default:
                
                builder.WriteString(fmt.Sprintf("%v", v))
            }
        }
    }
    
    compositeKeyName := qb.getCompositeKeyName(parts)
    return expression.Key(compositeKeyName).Equal(expression.Value(builder.String()))
}

func (qb *QueryBuilder) getCompositeKeyName(parts []CompositeKeyPart) string {
    if len(parts) == 0 {
        return ""
    }
    
    if len(parts) == 1 {
        return parts[0].Value
    }
    
    if len(parts) <= 3 {
        names := make([]string, 0, len(parts))
        for _, part := range parts {
            names = append(names, part.Value)
        }
        return strings.Join(names, "#")
    }
    
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


func (qb *QueryBuilder) buildCompositeKeyValue(parts []CompositeKeyPart) string {
    if len(parts) == 0 {
        return ""
    }
    
    if len(parts) == 1 {
        if parts[0].IsConstant {
            return parts[0].Value
        }
        return qb.formatAttributeValue(qb.Attributes[parts[0].Value])
    }
    
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

// formatAttributeValue 
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
            return "1" 
        }
        return "0"
    default:
        return fmt.Sprintf("%v", value)
    }
}
`
