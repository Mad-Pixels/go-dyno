package v2

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
    var compositeKeyValue string
    for i, part := range parts {
        var valueStr string
        if part.IsConstant {
            valueStr = part.Value
        } else {
            value := qb.Attributes[part.Value]
            valueStr = fmt.Sprintf("%v", value)
        }
        if i > 0 {
            compositeKeyValue += "#"
        }
        compositeKeyValue += valueStr
    }
    compositeKeyName := qb.getCompositeKeyName(parts)
    return expression.Key(compositeKeyName).Equal(expression.Value(compositeKeyValue))
}

func (qb *QueryBuilder) getCompositeKeyName(parts []CompositeKeyPart) string {
    var names []string
    for _, part := range parts {
        names = append(names, part.Value)
    }
    return strings.Join(names, "#")
}`
