package query

// QueryBuilderTemplate provides the main QueryBuilder struct and basic methods
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
func NewQueryBuilder() *QueryBuilder {
    return &QueryBuilder{
        FilterMixin:       NewFilterMixin(),
        PaginationMixin:   NewPaginationMixin(),
        KeyConditionMixin: NewKeyConditionMixin(),
    }
}

// Limit sets the maximum number of items and returns QueryBuilder for method chaining.
// Controls the number of items returned in a single request.
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
    qb.PaginationMixin.Limit(limit)
    return qb
}

// StartFrom sets the exclusive start key and returns QueryBuilder for method chaining.
// Use LastEvaluatedKey from previous response for pagination.
func (qb *QueryBuilder) StartFrom(lastEvaluatedKey map[string]types.AttributeValue) *QueryBuilder {
    qb.PaginationMixin.StartFrom(lastEvaluatedKey)
    return qb
}

// OrderByDesc sets descending sort order and returns QueryBuilder for method chaining.
// Only affects sort key ordering, not filter results.
func (qb *QueryBuilder) OrderByDesc() *QueryBuilder {
    qb.KeyConditionMixin.OrderByDesc()
    return qb
}

// OrderByAsc sets ascending sort order and returns QueryBuilder for method chaining.
// This is the default sort order.
func (qb *QueryBuilder) OrderByAsc() *QueryBuilder {
    qb.KeyConditionMixin.OrderByAsc()
    return qb
}

// WithPreferredSortKey sets the preferred sort key and returns QueryBuilder for method chaining.
// Hints the index selection algorithm when multiple indexes could satisfy the query.
func (qb *QueryBuilder) WithPreferredSortKey(key string) *QueryBuilder {
    qb.KeyConditionMixin.WithPreferredSortKey(key)
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
    for i, part := range nonConstantParts {
        if i < len(values) {
            qb.Attributes[part.Value] = values[i]
            qb.UsedKeys[part.Value] = true
        }
    }
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
                Name:             index.Name,
                Type:             getIndexType(index),
                HashKey:          index.HashKey,
                RangeKey:         index.RangeKey,
                IsHashComposite:  len(index.HashKeyParts) > 0,
                IsRangeComposite: len(index.RangeKeyParts) > 0,
                HashKeyParts:     countNonConstantParts(index.HashKeyParts),
                RangeKeyParts:    countNonConstantParts(index.RangeKeyParts),
                ProjectionType:   index.ProjectionType,
            }
        }
    }
    return nil
}

// IndexInfo provides metadata about a table index.
type IndexInfo struct {
    Name             string
    Type             string
    HashKey          string
    RangeKey         string
    IsHashComposite  bool
    IsRangeComposite bool
    HashKeyParts     int
    RangeKeyParts    int
    ProjectionType   string
}

// getIndexType returns human-readable index type.
func getIndexType(index SecondaryIndex) string {
    if index.HashKey == "" {
        return "LSI"
    }
    return "GSI"
}

// countNonConstantParts counts non-constant parts in composite key.
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
