package scan

// ScanBuilderUniversalTemplate ...
const ScanBuilderUniversalTemplate = `
// Filter adds a filter condition using the universal operator system
func (sb *ScanBuilder) Filter(field string, op OperatorType, values ...any) *ScanBuilder {
    if !ValidateValues(op, values) {
        return sb
    }

    fieldInfo, exists := TableSchema.FieldsMap[field]
    if !exists {
        return sb
    }

    if !ValidateOperator(fieldInfo.DynamoType, op) {
        return sb
    }

    filterCond, err := BuildConditionExpression(field, op, values)
    if err != nil {
        return sb
    }

    sb.FilterConditions = append(sb.FilterConditions, filterCond)
    sb.UsedKeys[field] = true

    if op == EQ && len(values) == 1 {
        sb.Attributes[field] = values[0]
    }

    return sb
}

// FilterEQ is a convenience method for equality filters
func (sb *ScanBuilder) FilterEQ(field string, value any) *ScanBuilder {
    return sb.Filter(field, EQ, value)
}

// FilterContains is a convenience method for contains filters
func (sb *ScanBuilder) FilterContains(field string, value any) *ScanBuilder {
    return sb.Filter(field, CONTAINS, value)
}

// FilterBeginsWith is a convenience method for begins_with filters
func (sb *ScanBuilder) FilterBeginsWith(field string, value any) *ScanBuilder {
    return sb.Filter(field, BEGINS_WITH, value)
}

// FilterBetween is a convenience method for range filters
func (sb *ScanBuilder) FilterBetween(field string, start, end any) *ScanBuilder {
    return sb.Filter(field, BETWEEN, start, end)
}

// FilterGT is a convenience method for greater than filters
func (sb *ScanBuilder) FilterGT(field string, value any) *ScanBuilder {
    return sb.Filter(field, GT, value)
}

// FilterLT is a convenience method for less than filters
func (sb *ScanBuilder) FilterLT(field string, value any) *ScanBuilder {
    return sb.Filter(field, LT, value)
}

// FilterGTE is a convenience method for greater than or equal filters
func (sb *ScanBuilder) FilterGTE(field string, value any) *ScanBuilder {
    return sb.Filter(field, GTE, value)
}

// FilterLTE is a convenience method for less than or equal filters
func (sb *ScanBuilder) FilterLTE(field string, value any) *ScanBuilder {
    return sb.Filter(field, LTE, value)
}

// FilterExists is a convenience method for attribute exists filters
func (sb *ScanBuilder) FilterExists(field string) *ScanBuilder {
    return sb.Filter(field, EXISTS)
}

// FilterNotExists is a convenience method for attribute not exists filters
func (sb *ScanBuilder) FilterNotExists(field string) *ScanBuilder {
    return sb.Filter(field, NOT_EXISTS)
}
`
