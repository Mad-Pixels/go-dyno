package scan

// ScanBuilderUniversalTemplate ...
const ScanBuilderUniversalTemplate = `
// Filter adds a filter condition using the universal operator system
func (sb *ScanBuilder) Filter(field string, op OperatorType, values ...interface{}) *ScanBuilder {
	if !ValidateValues(op, values) {
		return sb
	}

	// Простая проверка, что поле существует в наших атрибутах
	fieldExists := false
	for _, attr := range TableSchema.Attributes {
		if attr.Name == field {
			fieldExists = true
			break
		}
	}
	for _, attr := range TableSchema.CommonAttributes {
		if attr.Name == field {
			fieldExists = true
			break
		}
	}
	
	if !fieldExists {
		return sb
	}

	// Строим filter condition
	filterCond, err := BuildConditionExpression(field, op, values)
	if err != nil {
		return sb
	}

	sb.FilterConditions = append(sb.FilterConditions, filterCond)
	sb.UsedKeys[field] = true

	// Для простых равенств также сохраняем значение в Attributes
	if op == EQ && len(values) == 1 {
		sb.Attributes[field] = values[0]
	}

	return sb
}

// FilterEQ is a convenience method for equality filters
func (sb *ScanBuilder) FilterEQ(field string, value interface{}) *ScanBuilder {
	return sb.Filter(field, EQ, value)
}

// FilterContains is a convenience method for contains filters
func (sb *ScanBuilder) FilterContains(field string, value interface{}) *ScanBuilder {
	return sb.Filter(field, CONTAINS, value)
}

// FilterBeginsWith is a convenience method for begins_with filters
func (sb *ScanBuilder) FilterBeginsWith(field string, value interface{}) *ScanBuilder {
	return sb.Filter(field, BEGINS_WITH, value)
}

// FilterBetween is a convenience method for range filters
func (sb *ScanBuilder) FilterBetween(field string, start, end interface{}) *ScanBuilder {
	return sb.Filter(field, BETWEEN, start, end)
}

// FilterGT is a convenience method for greater than filters
func (sb *ScanBuilder) FilterGT(field string, value interface{}) *ScanBuilder {
	return sb.Filter(field, GT, value)
}

// FilterLT is a convenience method for less than filters
func (sb *ScanBuilder) FilterLT(field string, value interface{}) *ScanBuilder {
	return sb.Filter(field, LT, value)
}

// FilterGTE is a convenience method for greater than or equal filters
func (sb *ScanBuilder) FilterGTE(field string, value interface{}) *ScanBuilder {
	return sb.Filter(field, GTE, value)
}

// FilterLTE is a convenience method for less than or equal filters
func (sb *ScanBuilder) FilterLTE(field string, value interface{}) *ScanBuilder {
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
