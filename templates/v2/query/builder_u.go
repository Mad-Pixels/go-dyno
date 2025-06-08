package query

// QueryBuilderUniversalTemplate добавляет универсальные методы к QueryBuilder
const QueryBuilderUniversalTemplate = `
// With adds a key condition using the universal operator system
func (qb *QueryBuilder) With(field string, op OperatorType, values ...interface{}) *QueryBuilder {
	if !ValidateValues(op, values) {
		return qb
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
		return qb
	}

	// Проверяем, что это ключевое поле
	isKeyField := field == TableSchema.HashKey || field == TableSchema.RangeKey
	if !isKeyField {
		return qb
	}

	// Строим key condition
	keyCond, err := BuildKeyConditionExpression(field, op, values)
	if err != nil {
		return qb
	}

	qb.KeyConditions[field] = keyCond
	qb.UsedKeys[field] = true

	// Для простых равенств также сохраняем значение в Attributes
	if op == EQ && len(values) == 1 {
		qb.Attributes[field] = values[0]
	}

	return qb
}

// Filter adds a filter condition using the universal operator system
func (qb *QueryBuilder) Filter(field string, op OperatorType, values ...interface{}) *QueryBuilder {
	if !ValidateValues(op, values) {
		return qb
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
		return qb
	}

	// Строим filter condition
	filterCond, err := BuildConditionExpression(field, op, values)
	if err != nil {
		return qb
	}

	qb.FilterConditions = append(qb.FilterConditions, filterCond)
	qb.UsedKeys[field] = true

	// Для простых равенств также сохраняем значение в Attributes
	if op == EQ && len(values) == 1 {
		qb.Attributes[field] = values[0]
	}

	return qb
}

// WithEQ is a convenience method for equality conditions
func (qb *QueryBuilder) WithEQ(field string, value interface{}) *QueryBuilder {
	return qb.With(field, EQ, value)
}

// WithBetween is a convenience method for range conditions
func (qb *QueryBuilder) WithBetween(field string, start, end interface{}) *QueryBuilder {
	return qb.With(field, BETWEEN, start, end)
}

// WithGT is a convenience method for greater than conditions
func (qb *QueryBuilder) WithGT(field string, value interface{}) *QueryBuilder {
	return qb.With(field, GT, value)
}

// WithLT is a convenience method for less than conditions  
func (qb *QueryBuilder) WithLT(field string, value interface{}) *QueryBuilder {
	return qb.With(field, LT, value)
}

// FilterEQ is a convenience method for equality filters
func (qb *QueryBuilder) FilterEQ(field string, value interface{}) *QueryBuilder {
	return qb.Filter(field, EQ, value)
}

// FilterContains is a convenience method for contains filters
func (qb *QueryBuilder) FilterContains(field string, value interface{}) *QueryBuilder {
	return qb.Filter(field, CONTAINS, value)
}

// FilterBeginsWith is a convenience method for begins_with filters
func (qb *QueryBuilder) FilterBeginsWith(field string, value interface{}) *QueryBuilder {
	return qb.Filter(field, BEGINS_WITH, value)
}

// FilterBetween is a convenience method for range filters
func (qb *QueryBuilder) FilterBetween(field string, start, end interface{}) *QueryBuilder {
	return qb.Filter(field, BETWEEN, start, end)
}

// FilterGT is a convenience method for greater than filters
func (qb *QueryBuilder) FilterGT(field string, value interface{}) *QueryBuilder {
	return qb.Filter(field, GT, value)
}

// FilterLT is a convenience method for less than filters
func (qb *QueryBuilder) FilterLT(field string, value interface{}) *QueryBuilder {
	return qb.Filter(field, LT, value)
}
`
