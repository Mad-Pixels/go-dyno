package query

// QueryBuilderWithTemplate provides With methods for key conditions
const QueryBuilderWithTemplate = `
// With adds key condition and returns QueryBuilder for method chaining.
// Only works with partition and sort key attributes for efficient querying.
// Example: query.With("user_id", EQ, "123").With("created_at", GT, timestamp)
func (qb *QueryBuilder) With(field string, op OperatorType, values ...any) *QueryBuilder {
    qb.KeyConditionMixin.With(field, op, values...)
    if op == EQ && len(values) == 1 {
        qb.Attributes[field] = values[0]
        qb.UsedKeys[field] = true
    }
    return qb
}

{{if IsALL .Mode}}
// CONVENIENCE METHODS - Only available in ALL mode

// WithEQ adds equality key condition and returns QueryBuilder for method chaining.
// Required for partition keys, commonly used for sort keys.
// Example: query.WithEQ("user_id", "123")
func (qb *QueryBuilder) WithEQ(field string, value any) *QueryBuilder {
    qb.KeyConditionMixin.WithEQ(field, value)
    qb.Attributes[field] = value
    qb.UsedKeys[field] = true
    return qb
}

// WithBetween adds range key condition and returns QueryBuilder for method chaining.
// Only valid for sort keys, not partition keys.
// Example: query.WithBetween("timestamp", startTime, endTime)
func (qb *QueryBuilder) WithBetween(field string, start, end any) *QueryBuilder {
    qb.KeyConditionMixin.WithBetween(field, start, end)
    qb.Attributes[field+"_start"] = start
    qb.Attributes[field+"_end"] = end
    qb.UsedKeys[field] = true
    return qb
}

// WithGT adds greater than key condition and returns QueryBuilder for method chaining.
// Only valid for sort keys in range queries.
// Example: query.WithGT("created_at", yesterday)
func (qb *QueryBuilder) WithGT(field string, value any) *QueryBuilder {
    qb.KeyConditionMixin.WithGT(field, value)
    qb.Attributes[field] = value 
    qb.UsedKeys[field] = true
    return qb
}

// WithGTE adds greater than or equal key condition and returns QueryBuilder for method chaining.
// Only valid for sort keys in range queries.
// Example: query.WithGTE("score", minimumScore)
func (qb *QueryBuilder) WithGTE(field string, value any) *QueryBuilder {
    qb.KeyConditionMixin.WithGTE(field, value)
    qb.Attributes[field] = value
    qb.UsedKeys[field] = true
    return qb
}

// WithLT adds less than key condition and returns QueryBuilder for method chaining.
// Only valid for sort keys in range queries.
// Example: query.WithLT("expiry_date", now)
func (qb *QueryBuilder) WithLT(field string, value any) *QueryBuilder {
    qb.KeyConditionMixin.WithLT(field, value)
    qb.Attributes[field] = value
    qb.UsedKeys[field] = true
    return qb
}

// WithLTE adds less than or equal key condition and returns QueryBuilder for method chaining.
// Only valid for sort keys in range queries.
// Example: query.WithLTE("price", maxBudget)
func (qb *QueryBuilder) WithLTE(field string, value any) *QueryBuilder {
    qb.KeyConditionMixin.WithLTE(field, value)
    qb.Attributes[field] = value
    qb.UsedKeys[field] = true
    return qb
}

// WithIndexHashKey sets hash key for any index by name.
// Automatically handles both simple and composite keys based on schema metadata.
// For composite keys, pass values in the order they appear in the schema.
// Example: query.WithIndexHashKey("user-status-index", "user123")
// Example: query.WithIndexHashKey("tenant-user-index", "tenant1", "user123") // composite
func (qb *QueryBuilder) WithIndexHashKey(indexName string, values ...any) *QueryBuilder {
    index := qb.getIndexByName(indexName)
    if index == nil {
        return qb
    }
    
    if index.HashKeyParts != nil {
        nonConstantParts := qb.getNonConstantParts(index.HashKeyParts)
        if len(values) != len(nonConstantParts) {
            return qb
        }
        qb.setCompositeKey(index.HashKey, index.HashKeyParts, values)
    } else {
        if len(values) != 1 {
            return qb
        }
        qb.Attributes[index.HashKey] = values[0]
        qb.UsedKeys[index.HashKey] = true
        qb.KeyConditions[index.HashKey] = expression.Key(index.HashKey).Equal(expression.Value(values[0]))
    }
    return qb
}

// WithIndexRangeKey sets range key for any index by name.
// Automatically handles both simple and composite keys based on schema metadata.
// For composite keys, pass values in the order they appear in the schema.
// Example: query.WithIndexRangeKey("user-status-index", "active")
// Example: query.WithIndexRangeKey("date-type-index", "2023-01-01", "ORDER") // composite
func (qb *QueryBuilder) WithIndexRangeKey(indexName string, values ...any) *QueryBuilder {
    index := qb.getIndexByName(indexName)
    if index == nil || index.RangeKey == "" {
        return qb
    }
    
    if index.RangeKeyParts != nil {
        nonConstantParts := qb.getNonConstantParts(index.RangeKeyParts)
        if len(values) != len(nonConstantParts) {
            return qb
        }
        qb.setCompositeKey(index.RangeKey, index.RangeKeyParts, values)
    } else {
        if len(values) != 1 {
            return qb
        }
        qb.Attributes[index.RangeKey] = values[0]
        qb.UsedKeys[index.RangeKey] = true
        qb.KeyConditions[index.RangeKey] = expression.Key(index.RangeKey).Equal(expression.Value(values[0]))
    }
    return qb
}

// WithIndexRangeKeyBetween sets range key condition for any index with BETWEEN operator.
// Only works with simple range keys, not composite ones.
// Example: query.WithIndexRangeKeyBetween("date-index", startDate, endDate)
func (qb *QueryBuilder) WithIndexRangeKeyBetween(indexName string, start, end any) *QueryBuilder {
    index := qb.getIndexByName(indexName)
    if index == nil || index.RangeKey == "" || index.RangeKeyParts != nil {
        return qb 
    }
    qb.KeyConditions[index.RangeKey] = expression.Key(index.RangeKey).Between(expression.Value(start), expression.Value(end))
    qb.UsedKeys[index.RangeKey] = true
    qb.Attributes[index.RangeKey+"_start"] = start
    qb.Attributes[index.RangeKey+"_end"] = end
    return qb
}

// WithIndexRangeKeyGT sets range key condition for any index with GT operator.
// Only works with simple range keys, not composite ones.
// Example: query.WithIndexRangeKeyGT("score-index", 100)
func (qb *QueryBuilder) WithIndexRangeKeyGT(indexName string, value any) *QueryBuilder {
    index := qb.getIndexByName(indexName)
    if index == nil || index.RangeKey == "" || index.RangeKeyParts != nil {
        return qb
    }
    qb.KeyConditions[index.RangeKey] = expression.Key(index.RangeKey).GreaterThan(expression.Value(value))
    qb.UsedKeys[index.RangeKey] = true
    qb.Attributes[index.RangeKey] = value
    return qb
}

// WithIndexRangeKeyLT sets range key condition for any index with LT operator.
// Only works with simple range keys, not composite ones.
// Example: query.WithIndexRangeKeyLT("timestamp-index", cutoffTime)
func (qb *QueryBuilder) WithIndexRangeKeyLT(indexName string, value any) *QueryBuilder {
    index := qb.getIndexByName(indexName)
    if index == nil || index.RangeKey == "" || index.RangeKeyParts != nil {
        return qb
    }
    qb.KeyConditions[index.RangeKey] = expression.Key(index.RangeKey).LessThan(expression.Value(value))
    qb.UsedKeys[index.RangeKey] = true
    qb.Attributes[index.RangeKey] = value
    return qb
}

// WithIndexRangeKeyGTE sets range key condition for any index with GTE operator.
// Only works with simple range keys, not composite ones.
// Example: query.WithIndexRangeKeyGTE("score-index", 100)
func (qb *QueryBuilder) WithIndexRangeKeyGTE(indexName string, value any) *QueryBuilder {
    index := qb.getIndexByName(indexName)
    if index == nil || index.RangeKey == "" || index.RangeKeyParts != nil {
        return qb
    }
    qb.KeyConditions[index.RangeKey] = expression.Key(index.RangeKey).GreaterThanEqual(expression.Value(value))
    qb.UsedKeys[index.RangeKey] = true
    qb.Attributes[index.RangeKey] = value
    return qb
}

// WithIndexRangeKeyLTE sets range key condition for any index with LTE operator.
// Only works with simple range keys, not composite ones.
// Example: query.WithIndexRangeKeyLTE("timestamp-index", cutoffTime)
func (qb *QueryBuilder) WithIndexRangeKeyLTE(indexName string, value any) *QueryBuilder {
    index := qb.getIndexByName(indexName)
    if index == nil || index.RangeKey == "" || index.RangeKeyParts != nil {
        return qb
    }
    qb.KeyConditions[index.RangeKey] = expression.Key(index.RangeKey).LessThanEqual(expression.Value(value))
    qb.UsedKeys[index.RangeKey] = true
    qb.Attributes[index.RangeKey] = value
    return qb
}
{{end}}
`
