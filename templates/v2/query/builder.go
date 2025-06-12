package query

// QueryBuilderTemplate provides the main QueryBuilder struct with universal index methods
const QueryBuilderTemplate = `
// QueryBuilder provides a fluent interface for building type-safe DynamoDB queries.
// Combines FilterMixin, PaginationMixin, and KeyConditionMixin for comprehensive query building.
// Supports automatic index selection, composite keys, and all DynamoDB query patterns.
type QueryBuilder struct {
    FilterMixin       // Filter conditions for any table attribute
    PaginationMixin   // Limit and pagination support
    KeyConditionMixin // Key conditions for partition and sort keys
    IndexName string  // Optional index name override
}

// NewQueryBuilder creates a new QueryBuilder instance with initialized mixins.
// All mixins are properly initialized for immediate use.
// Example: query := NewQueryBuilder().WithEQ("user_id", "123").FilterEQ("status", "active")
func NewQueryBuilder() *QueryBuilder {
    return &QueryBuilder{
        FilterMixin:       NewFilterMixin(),
        PaginationMixin:   NewPaginationMixin(),
        KeyConditionMixin: NewKeyConditionMixin(),
    }
}

// Filter adds a filter condition and returns QueryBuilder for method chaining.
// Wraps FilterMixin.Filter with fluent interface support.
func (qb *QueryBuilder) Filter(field string, op OperatorType, values ...any) *QueryBuilder {
    qb.FilterMixin.Filter(field, op, values...)
    return qb
}

// FilterEQ adds equality filter and returns QueryBuilder for method chaining.
// Example: query.FilterEQ("status", "active")
func (qb *QueryBuilder) FilterEQ(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterEQ(field, value)
    return qb
}

// FilterContains adds contains filter and returns QueryBuilder for method chaining.
// Works with String attributes (substring) and Set attributes (membership).
// Example: query.FilterContains("tags", "premium")
func (qb *QueryBuilder) FilterContains(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterContains(field, value)
    return qb
}

// FilterNotContains adds not contains filter and returns QueryBuilder for method chaining.
// Opposite of FilterContains for exclusion filtering.
func (qb *QueryBuilder) FilterNotContains(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterNotContains(field, value)
    return qb
}

// FilterBeginsWith adds begins_with filter and returns QueryBuilder for method chaining.
// Only works with String attributes for prefix matching.
// Example: query.FilterBeginsWith("email", "admin@")
func (qb *QueryBuilder) FilterBeginsWith(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterBeginsWith(field, value)
    return qb
}

// FilterBetween adds range filter and returns QueryBuilder for method chaining.
// Works with comparable types for inclusive range filtering.
// Example: query.FilterBetween("score", 80, 100)
func (qb *QueryBuilder) FilterBetween(field string, start, end any) *QueryBuilder {
    qb.FilterMixin.FilterBetween(field, start, end)
    return qb
}

// FilterGT adds greater than filter and returns QueryBuilder for method chaining.
// Example: query.FilterGT("last_login", cutoffDate)
func (qb *QueryBuilder) FilterGT(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterGT(field, value)
    return qb
}

// FilterLT adds less than filter and returns QueryBuilder for method chaining.
// Example: query.FilterLT("attempts", maxAttempts)
func (qb *QueryBuilder) FilterLT(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterLT(field, value)
    return qb
}

// FilterGTE adds greater than or equal filter and returns QueryBuilder for method chaining.
// Example: query.FilterGTE("age", minimumAge)
func (qb *QueryBuilder) FilterGTE(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterGTE(field, value)
    return qb
}

// FilterLTE adds less than or equal filter and returns QueryBuilder for method chaining.
// Example: query.FilterLTE("file_size", maxFileSize)
func (qb *QueryBuilder) FilterLTE(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterLTE(field, value)
    return qb
}

// FilterExists adds attribute exists filter and returns QueryBuilder for method chaining.
// Checks if the specified attribute exists in the item.
// Example: query.FilterExists("optional_field")
func (qb *QueryBuilder) FilterExists(field string) *QueryBuilder {
    qb.FilterMixin.FilterExists(field)
    return qb
}

// FilterNotExists adds attribute not exists filter and returns QueryBuilder for method chaining.
// Checks if the specified attribute does not exist in the item.
func (qb *QueryBuilder) FilterNotExists(field string) *QueryBuilder {
    qb.FilterMixin.FilterNotExists(field)
    return qb
}

// FilterNE adds not equal filter and returns QueryBuilder for method chaining.
// Example: query.FilterNE("status", "deleted")
func (qb *QueryBuilder) FilterNE(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterNE(field, value)
    return qb
}

// FilterIn adds IN filter and returns QueryBuilder for method chaining.
// For scalar values only - use FilterContains for DynamoDB Sets.
// Example: query.FilterIn("category", "books", "electronics", "clothing")
func (qb *QueryBuilder) FilterIn(field string, values ...any) *QueryBuilder {
    qb.FilterMixin.FilterIn(field, values...)
    return qb
}

// FilterNotIn adds NOT_IN filter and returns QueryBuilder for method chaining.
// For scalar values only - use FilterNotContains for DynamoDB Sets.
func (qb *QueryBuilder) FilterNotIn(field string, values ...any) *QueryBuilder {
    qb.FilterMixin.FilterNotIn(field, values...)
    return qb
}

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

// WithPreferredSortKey sets the preferred sort key and returns QueryBuilder for method chaining.
// Hints the index selection algorithm when multiple indexes could satisfy the query.
// Example: query.WithPreferredSortKey("created_at")
func (qb *QueryBuilder) WithPreferredSortKey(key string) *QueryBuilder {
    qb.KeyConditionMixin.WithPreferredSortKey(key)
    return qb
}

// OrderByDesc sets descending sort order and returns QueryBuilder for method chaining.
// Only affects sort key ordering, not filter results.
// Example: query.OrderByDesc() // newest first
func (qb *QueryBuilder) OrderByDesc() *QueryBuilder {
    qb.KeyConditionMixin.OrderByDesc()
    return qb
}

// OrderByAsc sets ascending sort order and returns QueryBuilder for method chaining.
// This is the default sort order.
// Example: query.OrderByAsc() // oldest first
func (qb *QueryBuilder) OrderByAsc() *QueryBuilder {
    qb.KeyConditionMixin.OrderByAsc()
    return qb
}

// Limit sets the maximum number of items and returns QueryBuilder for method chaining.
// Controls the number of items returned in a single request.
// Example: query.Limit(25)
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
    qb.PaginationMixin.Limit(limit)
    return qb
}

// StartFrom sets the exclusive start key and returns QueryBuilder for method chaining.
// Use LastEvaluatedKey from previous response for pagination.
// Example: query.StartFrom(previousResponse.LastEvaluatedKey)
func (qb *QueryBuilder) StartFrom(lastEvaluatedKey map[string]types.AttributeValue) *QueryBuilder {
    qb.PaginationMixin.StartFrom(lastEvaluatedKey)
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
        // Composite key
        nonConstantParts := qb.getNonConstantParts(index.HashKeyParts)
        if len(values) != len(nonConstantParts) {
            return qb // Wrong number of values
        }
        qb.setCompositeKey(index.HashKey, index.HashKeyParts, values)
    } else {
        // Simple key
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
        // Composite key
        nonConstantParts := qb.getNonConstantParts(index.RangeKeyParts)
        if len(values) != len(nonConstantParts) {
            return qb
        }
        qb.setCompositeKey(index.RangeKey, index.RangeKeyParts, values)
    } else {
        // Simple key
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
        return qb // Only works with simple range keys
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

// HELPER METHODS for universal index access

// getIndexByName finds index by name in schema metadata.
func (qb *QueryBuilder) getIndexByName(indexName string) *SecondaryIndex {
    for i := range TableSchema.SecondaryIndexes {
        if TableSchema.SecondaryIndexes[i].Name == indexName {
            return &TableSchema.SecondaryIndexes[i]
        }
    }
    return nil
}

// getNonConstantParts returns only non-constant parts of composite key.
func (qb *QueryBuilder) getNonConstantParts(parts []CompositeKeyPart) []CompositeKeyPart {
    var result []CompositeKeyPart
    for _, part := range parts {
        if !part.IsConstant {
            result = append(result, part)
        }
    }
    return result
}

// setCompositeKey builds and sets composite key from parts and values.
func (qb *QueryBuilder) setCompositeKey(keyName string, parts []CompositeKeyPart, values []any) {
    nonConstantParts := qb.getNonConstantParts(parts)
    
    // Map values to their respective attributes
    for i, part := range nonConstantParts {
        if i < len(values) {
            qb.Attributes[part.Value] = values[i]
            qb.UsedKeys[part.Value] = true
        }
    }
    
    // Build composite value
    compositeValue := qb.buildCompositeKeyValue(parts)
    qb.Attributes[keyName] = compositeValue
    qb.UsedKeys[keyName] = true
    qb.KeyConditions[keyName] = expression.Key(keyName).Equal(expression.Value(compositeValue))
}

// SCHEMA INTROSPECTION METHODS

// GetIndexNames returns all available index names.
func GetIndexNames() []string {
    names := make([]string, len(TableSchema.SecondaryIndexes))
    for i, index := range TableSchema.SecondaryIndexes {
        names[i] = index.Name
    }
    return names
}

// GetIndexInfo returns detailed information about an index.
func GetIndexInfo(indexName string) *IndexInfo {
    for _, index := range TableSchema.SecondaryIndexes {
        if index.Name == indexName {
            return &IndexInfo{
                Name:              index.Name,
                Type:              getIndexType(index),
                HashKey:           index.HashKey,
                RangeKey:          index.RangeKey,
                IsHashComposite:   len(index.HashKeyParts) > 0,
                IsRangeComposite:  len(index.RangeKeyParts) > 0,
                HashKeyParts:      countNonConstantParts(index.HashKeyParts),
                RangeKeyParts:     countNonConstantParts(index.RangeKeyParts),
                ProjectionType:    index.ProjectionType,
            }
        }
    }
    return nil
}

// IndexInfo provides metadata about a table index.
type IndexInfo struct {
    Name              string   // Index name
    Type              string   // "GSI" or "LSI"
    HashKey           string   // Hash key attribute name
    RangeKey          string   // Range key attribute name (empty if none)
    IsHashComposite   bool     // Whether hash key is composite
    IsRangeComposite  bool     // Whether range key is composite
    HashKeyParts      int      // Number of non-constant hash key parts
    RangeKeyParts     int      // Number of non-constant range key parts
    ProjectionType    string   // "ALL", "KEYS_ONLY", or "INCLUDE"
}

func getIndexType(index SecondaryIndex) string {
    // GSI has different hash key than main table, LSI has same hash key
    if index.HashKey != TableSchema.HashKey {
        return "GSI"
    }
    return "LSI"
}

func countNonConstantParts(parts []CompositeKeyPart) int {
    count := 0
    for _, part := range parts {
        if !part.IsConstant {
            count++
        }
    }
    return count
}
`
