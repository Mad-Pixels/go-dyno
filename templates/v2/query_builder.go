package v2

// QueryBuilderTemplate generates a fluent, type-safe query builder for DynamoDB operations.
// This template creates a comprehensive QueryBuilder with the following capabilities:
// - Fluent API with method chaining for intuitive query construction
// - Type-safe methods for all table attributes with proper Go types
// - Automatic composite key handling and value generation
// - Smart index selection based on provided query parameters
// - Range query support (Between, GreaterThan, LessThan) for numeric attributes
// - Sorting, pagination, and filtering capabilities
// - Direct integration with AWS SDK v2 DynamoDB client
const QueryBuilderTemplate = `
// QueryBuilder provides a fluent, type-safe interface for constructing DynamoDB queries.
// It automatically selects the most efficient index based on provided attributes and
// handles complex composite key operations transparently.
//
// The builder supports:
// - Primary table queries and Global Secondary Index (GSI) queries
// - Composite key construction for complex access patterns
// - Range conditions (Between, GreaterThan, LessThan) on sort keys
// - Query result ordering (ascending/descending)
// - Result limiting and pagination with exclusive start keys
// - Filter expressions for non-key attributes
//
// Usage pattern:
//   query := NewQueryBuilder().
//       WithUserId("user123").             // Set hash key
//       WithCreatedGreaterThan(timestamp). // Range condition
//       OrderByDesc().                     // Sort order
//       Limit(10).                         // Result limit
//       Execute(ctx, dynamoClient)         // Execute query
type QueryBuilder struct {
    // IndexName stores the selected index name for the query
    IndexName           string
    
    // KeyConditions holds DynamoDB key condition expressions for hash and range keys
    KeyConditions       map[string]expression.KeyConditionBuilder
    
    // FilterConditions contains non-key attribute filter conditions
    FilterConditions    []expression.ConditionBuilder
    
    // UsedKeys tracks which attributes have been provided for query building
    UsedKeys            map[string]bool
    
    // Attributes stores the actual values for all query parameters
    Attributes          map[string]interface{}
    
    // SortDescending controls the sort order for query results
    SortDescending      bool
    
    // LimitValue specifies the maximum number of items to return
    LimitValue          *int
    
    // ExclusiveStartKey enables pagination by specifying where to start the next query
    ExclusiveStartKey   map[string]types.AttributeValue
    
    // PreferredSortKey allows manual index selection when multiple indexes are available
    PreferredSortKey    string
}

// NewQueryBuilder creates a new QueryBuilder instance with initialized internal maps.
// All internal collections are pre-allocated to avoid nil pointer issues during query building.
//
// Returns a ready-to-use QueryBuilder with fluent API methods available.
//
// Example:
//   qb := NewQueryBuilder()
//   items, err := qb.WithUserId("123").WithStatus("active").Execute(ctx, client)
func NewQueryBuilder() *QueryBuilder {
    return &QueryBuilder{
        KeyConditions:   make(map[string]expression.KeyConditionBuilder),
        UsedKeys:        make(map[string]bool),
        Attributes:      make(map[string]interface{}),
    }
}

{{range .Attributes}}
// With{{ToSafeName .Name | ToUpperCamelCase}} sets the key condition for "{{.Name}}" attribute.
// This method sets KeyConditionExpression for DynamoDB Query operation.
//
// DynamoDB key attribute: "{{.Name}}" (type: {{.Type}})
// Go parameter type: {{ToGolangBaseType .}}
//
// Returns the QueryBuilder for method chaining.
func (qb *QueryBuilder) With{{ToSafeName .Name | ToUpperCamelCase}}({{ToSafeName .Name | ToLowerCamelCase}} {{ToGolangBaseType .}}) *QueryBuilder {
    qb.Attributes["{{.Name}}"] = {{ToSafeName .Name | ToLowerCamelCase}}
    qb.UsedKeys["{{.Name}}"] = true
    return qb
}
{{end}}

{{range .CommonAttributes}}
// Filter{{ToSafeName .Name | ToUpperCamelCase}} sets the filter condition for "{{.Name}}" attribute.
// This method adds FilterExpression condition to DynamoDB Query operation.
//
// DynamoDB filter attribute: "{{.Name}}" (type: {{.Type}})
// Go parameter type: {{ToGolangBaseType .}}
//
// Returns the QueryBuilder for method chaining.
func (qb *QueryBuilder) Filter{{ToSafeName .Name | ToUpperCamelCase}}({{ToSafeName .Name | ToLowerCamelCase}} {{ToGolangBaseType .}}) *QueryBuilder {
    qb.Attributes["{{.Name}}"] = {{ToSafeName .Name | ToLowerCamelCase}}
    qb.UsedKeys["{{.Name}}"] = true
    return qb
}
{{end}}

{{range .SecondaryIndexes}}
{{if gt (len .HashKeyParts) 0}}
{{- $hasNonConstant := false -}}
{{- range .HashKeyParts -}}{{- if not .IsConstant -}}{{- $hasNonConstant = true -}}{{- end -}}{{- end -}}
{{- if $hasNonConstant}}
// With{{ToUpperCamelCase .Name}}HashKey sets the composite hash key for the "{{.Name}}" index.
// This method constructs a composite key from multiple attribute values for efficient querying.
//
// Composite key structure: "{{.HashKey}}" = {{range $i, $part := .HashKeyParts}}{{if $i}} + "#" + {{end}}{{if $part.IsConstant}}"{{$part.Value}}"{{else}}{{$part.Value}}{{end}}{{end}}
// Index: {{.Name}} ({{.ProjectionType}} projection)
//
// Parameters:
{{range .HashKeyParts}}{{if not .IsConstant}}//   - {{.Value | ToLowerCamelCase}} {{ToGolangAttrType .Value $.AllAttributes}}: Value for "{{.Value}}" part of composite key
{{end}}{{end}}//
// The method automatically:
// 1. Stores individual attribute values for filtering
// 2. Constructs the composite key value (e.g., "value1#value2")
// 3. Creates the appropriate DynamoDB key condition
//
// Example:
//   query.With{{ToUpperCamelCase .Name}}HashKey({{range $i, $part := .HashKeyParts}}{{if not $part.IsConstant}}{{if $i}}, {{end}}"example{{$i}}"{{end}}{{end}}) // Creates composite key
//
// Returns the QueryBuilder for method chaining.
func (qb *QueryBuilder) With{{ToUpperCamelCase .Name}}HashKey({{range $i, $part := .HashKeyParts}}{{if not $part.IsConstant}}{{if $i}}, {{end}}{{$part.Value | ToLowerCamelCase}} {{ToGolangAttrType $part.Value $.AllAttributes}}{{end}}{{end}}) *QueryBuilder {
    // Store individual attribute values for potential filtering
    {{range .HashKeyParts}}{{if not .IsConstant}}
    qb.Attributes["{{.Value}}"] = {{.Value | ToLowerCamelCase}}
    qb.UsedKeys["{{.Value}}"] = true
    {{end}}{{end}}
    
    // Build composite key value from parts
    compositeValue := qb.buildCompositeKeyValue([]CompositeKeyPart{
        {{range .HashKeyParts}}
        {{if .IsConstant}}
        {IsConstant: true, Value: "{{.Value}}"},  // Constant part: "{{.Value}}"
        {{else}}
        {IsConstant: false, Value: "{{.Value}}"}, // Dynamic part: {{.Value}}
        {{end}}
        {{end}}
    })
    
    // Set composite key value and condition
    qb.Attributes["{{.HashKey}}"] = compositeValue
    qb.UsedKeys["{{.HashKey}}"] = true
    qb.KeyConditions["{{.HashKey}}"] = expression.Key("{{.HashKey}}").Equal(expression.Value(compositeValue))
    return qb
}
{{end}}
{{else if .HashKey}}
// With{{ToUpperCamelCase .Name}}HashKey sets the hash key for the "{{.Name}}" index.
// This method enables querying using the {{.Name}} index with the specified hash key value.
//
// Index: {{.Name}} ({{.ProjectionType}} projection)
// Hash key: "{{.HashKey}}" ({{ToGolangAttrType .HashKey $.AllAttributes}})
{{if .RangeKey}}// Range key: "{{.RangeKey}}" (optional for this method){{end}}
//
// Example:
//   query.With{{ToUpperCamelCase .Name}}HashKey(hashValue) // Query {{.Name}} index
//
// Returns the QueryBuilder for method chaining.
func (qb *QueryBuilder) With{{ToUpperCamelCase .Name}}HashKey({{.HashKey | ToLowerCamelCase}} {{ToGolangAttrType .HashKey $.AllAttributes}}) *QueryBuilder {
    qb.Attributes["{{.HashKey}}"] = {{.HashKey | ToLowerCamelCase}}
    qb.UsedKeys["{{.HashKey}}"] = true
    qb.KeyConditions["{{.HashKey}}"] = expression.Key("{{.HashKey}}").Equal(expression.Value({{.HashKey | ToLowerCamelCase}}))
    return qb
}
{{end}}
{{end}}

// WithPreferredSortKey manually specifies which sort key to prefer when multiple indexes are available.
// This method provides fine-grained control over index selection for query optimization.
//
// When multiple indexes can satisfy a query, the QueryBuilder automatically selects the most
// efficient one. Use this method to override the automatic selection if needed.
//
// Parameters:
//   - key: The name of the preferred sort key attribute
//
// Example:
//   query.WithPreferredSortKey("created_at") // Prefer indexes with created_at as sort key
//
// Returns the QueryBuilder for method chaining.
func (qb *QueryBuilder) WithPreferredSortKey(key string) *QueryBuilder {
    qb.PreferredSortKey = key
    return qb
}

{{range .SecondaryIndexes}}
{{if gt (len .RangeKeyParts) 0}}
{{- $hasNonConstant := false -}}
{{- range .RangeKeyParts -}}{{- if not .IsConstant -}}{{- $hasNonConstant = true -}}{{- end -}}{{- end -}}
{{- if $hasNonConstant}}
// With{{ToUpperCamelCase .Name}}RangeKey sets the composite range key for the "{{.Name}}" index.
// This method constructs a composite range key from multiple attribute values.
//
// Composite range key: "{{.RangeKey}}" = {{range $i, $part := .RangeKeyParts}}{{if $i}} + "#" + {{end}}{{if $part.IsConstant}}"{{$part.Value}}"{{else}}{{$part.Value}}{{end}}{{end}}
// Index: {{.Name}} ({{.ProjectionType}} projection)
//
// Note: Range key conditions should be set separately using WithXxxBetween, WithXxxGreaterThan, etc.
// This method only sets the range key value for exact matching.
//
// Parameters:
{{range .RangeKeyParts}}{{if not .IsConstant}}//   - {{.Value | ToLowerCamelCase}} {{ToGolangAttrType .Value $.AllAttributes}}: Value for "{{.Value}}" part of composite range key
{{end}}{{end}}//
// Example:
//   query.With{{ToUpperCamelCase .Name}}RangeKey({{range $i, $part := .RangeKeyParts}}{{if not $part.IsConstant}}{{if $i}}, {{end}}"value{{$i}}"{{end}}{{end}}) // Sets composite range key
//
// Returns the QueryBuilder for method chaining.
func (qb *QueryBuilder) With{{ToUpperCamelCase .Name}}RangeKey({{range $i, $part := .RangeKeyParts}}{{if not $part.IsConstant}}{{if $i}}, {{end}}{{$part.Value | ToLowerCamelCase}} {{ToGolangAttrType $part.Value $.AllAttributes}}{{end}}{{end}}) *QueryBuilder {
    // Store individual attribute values for potential filtering
    {{range .RangeKeyParts}}{{if not .IsConstant}}
    qb.Attributes["{{.Value}}"] = {{.Value | ToLowerCamelCase}}
    qb.UsedKeys["{{.Value}}"] = true
    {{end}}{{end}}
    return qb
}
{{end}}
{{else if .RangeKey}}
// With{{ToUpperCamelCase .Name}}RangeKey sets the range key for the "{{.Name}}" index.
// This method enables exact matching on the range key for the {{.Name}} index.
//
// Index: {{.Name}} ({{.ProjectionType}} projection)
// Range key: "{{.RangeKey}}" ({{ToGolangAttrType .RangeKey $.AllAttributes}})
//
// For range conditions (>, <, BETWEEN), use the dedicated range methods:
// - With{{ToSafeName .RangeKey | ToUpperCamelCase}}GreaterThan()
// - With{{ToSafeName .RangeKey | ToUpperCamelCase}}LessThan()
// - With{{ToSafeName .RangeKey | ToUpperCamelCase}}Between()
//
// Example:
//   query.With{{ToUpperCamelCase .Name}}RangeKey(rangeValue) // Exact range key match
//
// Returns the QueryBuilder for method chaining.
func (qb *QueryBuilder) With{{ToUpperCamelCase .Name}}RangeKey({{.RangeKey | ToLowerCamelCase}} {{ToGolangAttrType .RangeKey $.AllAttributes}}) *QueryBuilder {
    qb.Attributes["{{.RangeKey}}"] = {{.RangeKey | ToLowerCamelCase}}
    qb.UsedKeys["{{.RangeKey}}"] = true
    qb.KeyConditions["{{.RangeKey}}"] = expression.Key("{{.RangeKey}}").Equal(expression.Value({{.RangeKey | ToLowerCamelCase}}))
    return qb
}
{{end}}
{{end}}

{{range .AllAttributes}}
{{if IsNumericAttr .}}
// With{{ToUpperCamelCase .Name}}Between creates a range condition for the "{{.Name}}" attribute.
// This method is particularly useful for sort keys in queries where you need to find items
// within a specific numeric range.
//
// DynamoDB condition: {{.Name}} BETWEEN start AND end (inclusive)
// Attribute type: {{.Type}} (numeric)
//
// Parameters:
//   - start: Lower bound of the range (inclusive)
//   - end: Upper bound of the range (inclusive)
//
// Example:
//   query.With{{ToUpperCamelCase .Name}}Between(100, 500) // {{.Name}} between 100 and 500
//
// Returns the QueryBuilder for method chaining.
func (qb *QueryBuilder) With{{ToUpperCamelCase .Name}}Between(start, end {{ToGolangBaseType .}}) *QueryBuilder {
    qb.KeyConditions["{{.Name}}"] = expression.Key("{{.Name}}").Between(expression.Value(start), expression.Value(end))
    qb.UsedKeys["{{.Name}}"] = true
    return qb
}

// With{{ToUpperCamelCase .Name}}GreaterThan creates a "greater than" condition for the "{{.Name}}" attribute.
// This method is useful for sort key queries where you need items after a specific value.
//
// DynamoDB condition: {{.Name}} > value
// Attribute type: {{.Type}} (numeric)
//
// Parameters:
//   - value: The threshold value (exclusive lower bound)
//
// Example:
//   query.With{{ToUpperCamelCase .Name}}GreaterThan(1000) // {{.Name}} > 1000
//
// Returns the QueryBuilder for method chaining.
func (qb *QueryBuilder) With{{ToUpperCamelCase .Name}}GreaterThan(value {{ToGolangBaseType .}}) *QueryBuilder {
    qb.KeyConditions["{{.Name}}"] = expression.Key("{{.Name}}").GreaterThan(expression.Value(value))
    qb.UsedKeys["{{.Name}}"] = true
    return qb
}

// With{{ToUpperCamelCase .Name}}LessThan creates a "less than" condition for the "{{.Name}}" attribute.
// This method is useful for sort key queries where you need items before a specific value.
//
// DynamoDB condition: {{.Name}} < value
// Attribute type: {{.Type}} (numeric)
//
// Parameters:
//   - value: The threshold value (exclusive upper bound)
//
// Example:
//   query.With{{ToUpperCamelCase .Name}}LessThan(500) // {{.Name}} < 500
//
// Returns the QueryBuilder for method chaining.
func (qb *QueryBuilder) With{{ToUpperCamelCase .Name}}LessThan(value {{ToGolangBaseType .}}) *QueryBuilder {
    qb.KeyConditions["{{.Name}}"] = expression.Key("{{.Name}}").LessThan(expression.Value(value))
    qb.UsedKeys["{{.Name}}"] = true
    return qb
}
{{end}}
{{end}}

// OrderByDesc sets the query to return results in descending order.
// This affects the sort order of the range key in DynamoDB queries.
//
// DynamoDB behavior: Sets ScanIndexForward=false in the query
// Default: Ascending order (ScanIndexForward=true)
//
// Example:
//   query.OrderByDesc() // Most recent items first (if sorting by timestamp)
//
// Returns the QueryBuilder for method chaining.
func (qb *QueryBuilder) OrderByDesc() *QueryBuilder {
    qb.SortDescending = true
    return qb
}

// OrderByAsc sets the query to return results in ascending order.
// This is the default behavior, but can be used to explicitly override a previous OrderByDesc().
//
// DynamoDB behavior: Sets ScanIndexForward=true in the query
// Default: This is the default sort order
//
// Example:
//   query.OrderByAsc() // Oldest items first (if sorting by timestamp)
//
// Returns the QueryBuilder for method chaining.
func (qb *QueryBuilder) OrderByAsc() *QueryBuilder {
    qb.SortDescending = false
    return qb
}

// Limit restricts the maximum number of items returned by the query.
// This is applied before any filter expressions, so the actual returned count
// may be lower if filters are applied.
//
// DynamoDB behavior: Sets the Limit parameter in QueryInput
// Range: 1 to 1MB of data (DynamoDB limitation)
//
// Parameters:
//   - limit: Maximum number of items to return
//
// Example:
//   query.Limit(50) // Return at most 50 items
//
// Returns the QueryBuilder for method chaining.
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
    qb.LimitValue = &limit
    return qb
}

// StartFrom enables pagination by specifying the exclusive start key for the query.
// Use the LastEvaluatedKey from a previous query response to continue pagination.
//
// DynamoDB behavior: Sets ExclusiveStartKey in QueryInput
// Pagination: Essential for handling large result sets
//
// Parameters:
//   - lastEvaluatedKey: The LastEvaluatedKey from the previous query response
//
// Example:
//   // First query
//   result1, _ := query.Limit(10).Execute(ctx, client)
//   
//   // Continue pagination
//   result2, _ := query.StartFrom(result1.LastEvaluatedKey).Execute(ctx, client)
//
// Returns the QueryBuilder for method chaining.
func (qb *QueryBuilder) StartFrom(lastEvaluatedKey map[string]types.AttributeValue) *QueryBuilder {
    qb.ExclusiveStartKey = lastEvaluatedKey
    return qb
}` + QueryBuilderBuildTemplate + QueryBuilderUtilsTemplate
