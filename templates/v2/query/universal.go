package query

// QueryBuilderUniversalTemplate provides universal operator support for QueryBuilder
const QueryBuilderUniversalTemplate = `
// With adds a key condition using the universal operator system.
// Only works with partition and sort key attributes - non-key attributes are ignored.
// Provides type-safe operator validation based on the schema field types.
// Example: query.With("user_id", EQ, "123").With("created_at", GT, timestamp)
func (qb *QueryBuilder) With(field string, op OperatorType, values ...any) *QueryBuilder {
    if !ValidateValues(op, values) {
        return qb
    }

    fieldInfo, exists := TableSchema.FieldsMap[field]
    if !exists {
        return qb
    }

    if !fieldInfo.IsKey {
        return qb
    }

    if !ValidateOperator(fieldInfo.DynamoType, op) {
        return qb
    }

    keyCond, err := BuildKeyConditionExpression(field, op, values)
    if err != nil {
        return qb
    }

    qb.KeyConditions[field] = keyCond
    qb.UsedKeys[field] = true

    if op == EQ && len(values) == 1 {
        qb.Attributes[field] = values[0]
    }

    return qb
}

// Filter adds a filter condition using the universal operator system.
// Works with any table attribute and validates operator compatibility with field types.
// Uses O(1) schema lookup for efficient field validation and type checking.
// Example: query.Filter("status", EQ, "active").Filter("tags", CONTAINS, "premium")
func (qb *QueryBuilder) Filter(field string, op OperatorType, values ...any) *QueryBuilder {
    if !ValidateValues(op, values) {
        return qb
    }

    // O(1) field lookup with pre-computed type information
    fieldInfo, exists := TableSchema.FieldsMap[field]
    if !exists {
        return qb
    }

    // Validate operator compatibility with DynamoDB field type
    if !ValidateOperator(fieldInfo.DynamoType, op) {
        return qb
    }

    // Build type-safe filter condition
    filterCond, err := BuildConditionExpression(field, op, values)
    if err != nil {
        return qb
    }

    qb.FilterConditions = append(qb.FilterConditions, filterCond)
    qb.UsedKeys[field] = true

    // Store simple equality values for index selection optimization
    if op == EQ && len(values) == 1 {
        qb.Attributes[field] = values[0]
    }

    return qb
}

// WithEQ is a convenience method for equality key conditions.
// Required for partition keys, commonly used for sort keys.
// Example: query.WithEQ("user_id", "123")
func (qb *QueryBuilder) WithEQ(field string, value any) *QueryBuilder {
    return qb.With(field, EQ, value)
}

// WithBetween is a convenience method for range key conditions.
// Only valid for sort keys, not partition keys.
// Example: query.WithBetween("timestamp", startTime, endTime)
func (qb *QueryBuilder) WithBetween(field string, start, end any) *QueryBuilder {
    return qb.With(field, BETWEEN, start, end)
}

// WithGT is a convenience method for greater than key conditions.
// Only valid for sort keys in range queries.
// Example: query.WithGT("created_at", yesterday)
func (qb *QueryBuilder) WithGT(field string, value any) *QueryBuilder {
    return qb.With(field, GT, value)
}

// WithGTE is a convenience method for greater than or equal key conditions.
// Only valid for sort keys in range queries.
// Example: query.WithGTE("score", minimumScore)
func (qb *QueryBuilder) WithGTE(field string, value any) *QueryBuilder {
    return qb.With(field, GTE, value)
}

// WithLT is a convenience method for less than key conditions.
// Only valid for sort keys in range queries.
// Example: query.WithLT("expiry_date", now)
func (qb *QueryBuilder) WithLT(field string, value any) *QueryBuilder {
    return qb.With(field, LT, value)
}

// WithLTE is a convenience method for less than or equal key conditions.
// Only valid for sort keys in range queries.
// Example: query.WithLTE("price", maxBudget)
func (qb *QueryBuilder) WithLTE(field string, value any) *QueryBuilder {
    return qb.With(field, LTE, value)
}

// FilterEQ is a convenience method for equality filter conditions.
// Works with any non-key attribute for post-query filtering.
// Example: query.FilterEQ("status", "active")
func (qb *QueryBuilder) FilterEQ(field string, value any) *QueryBuilder {
    return qb.Filter(field, EQ, value)
}

// FilterContains is a convenience method for contains filter conditions.
// Works with String attributes (substring) and Set attributes (membership).
// Example: query.FilterContains("description", "urgent") or query.FilterContains("tags", "vip")
func (qb *QueryBuilder) FilterContains(field string, value any) *QueryBuilder {
    return qb.Filter(field, CONTAINS, value)
}

// FilterBeginsWith is a convenience method for begins_with filter conditions.
// Only works with String attributes for prefix matching.
// Example: query.FilterBeginsWith("email", "admin@")
func (qb *QueryBuilder) FilterBeginsWith(field string, value any) *QueryBuilder {
    return qb.Filter(field, BEGINS_WITH, value)
}

// FilterBetween is a convenience method for range filter conditions.
// Works with comparable types (strings, numbers, dates).
// Example: query.FilterBetween("score", 80, 100)
func (qb *QueryBuilder) FilterBetween(field string, start, end any) *QueryBuilder {
    return qb.Filter(field, BETWEEN, start, end)
}

// FilterGT is a convenience method for greater than filter conditions.
// Works with comparable types for post-query filtering.
// Example: query.FilterGT("last_login", cutoffDate)
func (qb *QueryBuilder) FilterGT(field string, value any) *QueryBuilder {
    return qb.Filter(field, GT, value)
}

// FilterGTE is a convenience method for greater than or equal filter conditions.
// Works with comparable types for inclusive range filtering.
// Example: query.FilterGTE("age", minimumAge)
func (qb *QueryBuilder) FilterGTE(field string, value any) *QueryBuilder {
    return qb.Filter(field, GTE, value)
}

// FilterLT is a convenience method for less than filter conditions.
// Works with comparable types for upper bound filtering.
// Example: query.FilterLT("attempts", maxAttempts)
func (qb *QueryBuilder) FilterLT(field string, value any) *QueryBuilder {
    return qb.Filter(field, LT, value)
}

// FilterLTE is a convenience method for less than or equal filter conditions.
// Works with comparable types for inclusive upper bound filtering.
// Example: query.FilterLTE("file_size", maxFileSize)
func (qb *QueryBuilder) FilterLTE(field string, value any) *QueryBuilder {
    return qb.Filter(field, LTE, value)
}
`
