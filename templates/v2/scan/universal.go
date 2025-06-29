package scan

// ScanBuilderUniversalTemplate provides universal operator support for ScanBuilder
const ScanBuilderUniversalTemplate = `
// Filter adds a filter condition using the universal operator system.
// Works with any table attribute and validates operator compatibility with field types.
// Filters are applied after items are read from DynamoDB, affecting the returned results.
// Example: scan.Filter("status", EQ, "active").Filter("tags", CONTAINS, "premium")
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

// FilterEQ is a convenience method for equality filter conditions.
// Most commonly used filter for exact matches.
// Example: scan.FilterEQ("status", "active")
func (sb *ScanBuilder) FilterEQ(field string, value any) *ScanBuilder {
    return sb.Filter(field, EQ, value)
}

// FilterContains is a convenience method for contains filter conditions.
// Works with String attributes (substring search) and Set attributes (membership check).
// Example: scan.FilterContains("description", "urgent") or scan.FilterContains("tags", "vip")
func (sb *ScanBuilder) FilterContains(field string, value any) *ScanBuilder {
    return sb.Filter(field, CONTAINS, value)
}

// FilterBeginsWith is a convenience method for begins_with filter conditions.
// Only works with String attributes for prefix matching.
// Useful for hierarchical data or prefix-based searches.
// Example: scan.FilterBeginsWith("email", "admin@")
func (sb *ScanBuilder) FilterBeginsWith(field string, value any) *ScanBuilder {
    return sb.Filter(field, BEGINS_WITH, value)
}

// FilterBetween is a convenience method for range filter conditions.
// Works with comparable types (strings, numbers, dates) for inclusive range filtering.
// Useful for date ranges, score ranges, or any bounded searches.
// Example: scan.FilterBetween("score", 80, 100)
func (sb *ScanBuilder) FilterBetween(field string, start, end any) *ScanBuilder {
    return sb.Filter(field, BETWEEN, start, end)
}

// FilterGT is a convenience method for greater than filter conditions.
// Works with comparable types for threshold-based filtering.
// Example: scan.FilterGT("last_login", cutoffDate)
func (sb *ScanBuilder) FilterGT(field string, value any) *ScanBuilder {
    return sb.Filter(field, GT, value)
}

// FilterLT is a convenience method for less than filter conditions.
// Works with comparable types for upper bound filtering.
// Example: scan.FilterLT("attempts", maxAttempts)
func (sb *ScanBuilder) FilterLT(field string, value any) *ScanBuilder {
    return sb.Filter(field, LT, value)
}

// FilterGTE is a convenience method for greater than or equal filter conditions.
// Works with comparable types for inclusive lower bound filtering.
// Example: scan.FilterGTE("age", minimumAge)
func (sb *ScanBuilder) FilterGTE(field string, value any) *ScanBuilder {
    return sb.Filter(field, GTE, value)
}

// FilterLTE is a convenience method for less than or equal filter conditions.
// Works with comparable types for inclusive upper bound filtering.
// Example: scan.FilterLTE("file_size", maxFileSize)
func (sb *ScanBuilder) FilterLTE(field string, value any) *ScanBuilder {
    return sb.Filter(field, LTE, value)
}

// FilterExists is a convenience method for attribute exists filter conditions.
// Checks if the specified attribute exists in the item, regardless of its value.
// Useful for filtering items that have optional attributes populated.
// Example: scan.FilterExists("optional_field")
func (sb *ScanBuilder) FilterExists(field string) *ScanBuilder {
    return sb.Filter(field, EXISTS)
}

// FilterNotExists is a convenience method for attribute not exists filter conditions.
// Checks if the specified attribute does not exist in the item.
// Useful for finding items missing certain attributes or filtering incomplete records.
// Example: scan.FilterNotExists("deprecated_field")
func (sb *ScanBuilder) FilterNotExists(field string) *ScanBuilder {
    return sb.Filter(field, NOT_EXISTS)
}
`
