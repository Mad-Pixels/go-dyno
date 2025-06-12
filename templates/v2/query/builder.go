package query

// QueryBuilderTemplate provides the main QueryBuilder struct with mixin composition
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

{{range .SecondaryIndexes}}
{{if gt (len .HashKeyParts) 0}}
{{- $hasNonConstant := false -}}
{{- range .HashKeyParts -}}{{- if not .IsConstant -}}{{- $hasNonConstant = true -}}{{- end -}}{{- end -}}
{{- if $hasNonConstant}}
// WithIndexHashKey sets composite hash key for the index.
// Automatically builds the composite key from the provided components.
// Example: query.WithIndexHashKey(value1, value2)
func (qb *QueryBuilder) With{{ToUpperCamelCase .Name}}HashKey({{range $i, $part := .HashKeyParts}}{{if not $part.IsConstant}}{{if $i}}, {{end}}{{$part.Value | ToLowerCamelCase}} {{ToGolangAttrType $part.Value $.AllAttributes}}{{end}}{{end}}) *QueryBuilder {
    {{range .HashKeyParts}}{{if not .IsConstant}}
    qb.Attributes["{{.Value}}"] = {{.Value | ToLowerCamelCase}}
    qb.UsedKeys["{{.Value}}"] = true
    {{end}}{{end}}
    
    compositeValue := qb.buildCompositeKeyValue([]CompositeKeyPart{
        {{range .HashKeyParts}}
        {{if .IsConstant}}
        {IsConstant: true, Value: "{{.Value}}"},
        {{else}}
        {IsConstant: false, Value: "{{.Value}}"},
        {{end}}
        {{end}}
    })
    
    qb.Attributes["{{.HashKey}}"] = compositeValue
    qb.UsedKeys["{{.HashKey}}"] = true
    qb.KeyConditions["{{.HashKey}}"] = expression.Key("{{.HashKey}}").Equal(expression.Value(compositeValue))
    return qb
}
{{end}}
{{else if .HashKey}}
// WithIndexHashKey sets hash key for the index.
// Provides type-safe access to the index partition key.
// Example: query.WithIndexHashKey(value)
func (qb *QueryBuilder) With{{ToUpperCamelCase .Name}}HashKey({{.HashKey | ToLowerCamelCase}} {{ToGolangAttrType .HashKey $.AllAttributes}}) *QueryBuilder {
    qb.Attributes["{{.HashKey}}"] = {{.HashKey | ToLowerCamelCase}}
    qb.UsedKeys["{{.HashKey}}"] = true
    qb.KeyConditions["{{.HashKey}}"] = expression.Key("{{.HashKey}}").Equal(expression.Value({{.HashKey | ToLowerCamelCase}}))
    return qb
}
{{end}}
{{end}}

{{range .SecondaryIndexes}}
{{if gt (len .RangeKeyParts) 0}}
{{- $hasNonConstant := false -}}
{{- range .RangeKeyParts -}}{{- if not .IsConstant -}}{{- $hasNonConstant = true -}}{{- end -}}{{- end -}}
{{- if $hasNonConstant}}
// WithIndexRangeKey sets composite range key for the index.
// Automatically builds the composite key from the provided components.
// Example: query.WithIndexRangeKey(value1, value2)
func (qb *QueryBuilder) With{{ToUpperCamelCase .Name}}RangeKey({{range $i, $part := .RangeKeyParts}}{{if not $part.IsConstant}}{{if $i}}, {{end}}{{$part.Value | ToLowerCamelCase}} {{ToGolangAttrType $part.Value $.AllAttributes}}{{end}}{{end}}) *QueryBuilder {
    {{range .RangeKeyParts}}{{if not .IsConstant}}
    qb.Attributes["{{.Value}}"] = {{.Value | ToLowerCamelCase}}
    qb.UsedKeys["{{.Value}}"] = true
    {{end}}{{end}}
    return qb
}
{{end}}
{{else if .RangeKey}}
// WithIndexRangeKey sets range key for the index.
// Provides type-safe access to the index sort key.
// Example: query.WithIndexRangeKey(value)
func (qb *QueryBuilder) With{{ToUpperCamelCase .Name}}RangeKey({{.RangeKey | ToLowerCamelCase}} {{ToGolangAttrType .RangeKey $.AllAttributes}}) *QueryBuilder {
    qb.Attributes["{{.RangeKey}}"] = {{.RangeKey | ToLowerCamelCase}}
    qb.UsedKeys["{{.RangeKey}}"] = true
    qb.KeyConditions["{{.RangeKey}}"] = expression.Key("{{.RangeKey}}").Equal(expression.Value({{.RangeKey | ToLowerCamelCase}}))
    return qb
}
{{end}}
{{end}}
`
