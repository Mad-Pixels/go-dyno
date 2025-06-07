package v2

// ScanBuilderTemplate generates a fluent, type-safe scan builder for DynamoDB operations.
// This template creates a comprehensive ScanBuilder with the following capabilities:
// - Fluent API with method chaining for intuitive scan construction
// - Type-safe filter methods for all table attributes with proper Go types
// - Support for range filters (Between, GreaterThan, LessThan) for numeric attributes
// - String operations (Contains, BeginsWith) for text attributes
// - Parallel scanning for better performance on large tables
// - Projection expressions to limit returned attributes
// - Pagination and limiting capabilities
// - Direct integration with AWS SDK v2 DynamoDB client
const ScanBuilderTemplate = `
// ScanBuilder provides a fluent, type-safe interface for constructing DynamoDB scans.
// Unlike queries, scans examine every item in a table or index and return all matching items.
// Use scans when you need to examine all items or when query patterns don't match your access requirements.
//
// The builder supports:
// - Primary table scans and Global Secondary Index (GSI) scans
// - Filter expressions for all attributes
// - Parallel scanning for better performance on large tables
// - Result limiting and pagination with exclusive start keys
// - Projection expressions to limit returned attributes
//
// Usage pattern:
//   scan := NewScanBuilder().
//       WithIndex("StatusIndex").             // Optional: scan specific index
//       FilterStatus("active").               // Filter conditions
//       FilterCreatedAtGreaterThan(timestamp). // Range filters
//       WithParallelScan(4, 0).              // Optional: parallel scanning
//       Limit(100).                          // Result limit
//       Execute(ctx, dynamoClient)           // Execute scan
//
// Performance Notes:
// - Scans are less efficient than queries for large tables
// - Consider using parallel scans for better throughput
// - Use filters to reduce network traffic, not item examination
// - Scans consume read capacity for all examined items, not just returned items
type ScanBuilder struct {
    // IndexName stores the selected index name for the scan (empty for table scan)
    IndexName           string
    
    // FilterConditions contains all filter expressions for attributes
    FilterConditions    []expression.ConditionBuilder
    
    // UsedKeys tracks which attributes have been provided for filter building
    UsedKeys            map[string]bool
    
    // Attributes stores the actual values for all filter parameters
    Attributes          map[string]interface{}
    
    // LimitValue specifies the maximum number of items to return
    LimitValue          *int
    
    // ExclusiveStartKey enables pagination by specifying where to start the next scan
    ExclusiveStartKey   map[string]types.AttributeValue
    
    // ProjectionAttributes specifies which attributes to return (empty = all attributes)
    ProjectionAttributes []string
    
    // ParallelScanConfig enables parallel scanning for better performance
    ParallelScanConfig  *ParallelScanConfig
}

// ParallelScanConfig contains configuration for parallel scanning
type ParallelScanConfig struct {
    // TotalSegments is the total number of segments to divide the scan into
    TotalSegments int
    
    // Segment is the current segment number (0-based)
    Segment int
}

// NewScanBuilder creates a new ScanBuilder instance with initialized internal maps.
// All internal collections are pre-allocated to avoid nil pointer issues during scan building.
//
// Returns a ready-to-use ScanBuilder with fluent API methods available.
//
// Example:
//   sb := NewScanBuilder()
//   items, err := sb.FilterStatus("active").FilterAge(25).Execute(ctx, client)
func NewScanBuilder() *ScanBuilder {
    return &ScanBuilder{
        FilterConditions: make([]expression.ConditionBuilder, 0),
        UsedKeys:         make(map[string]bool),
        Attributes:       make(map[string]interface{}),
    }
}

// WithIndex specifies which index to scan instead of the main table.
// If not specified, scans the primary table.
//
// Parameters:
//   - indexName: Name of the Global Secondary Index to scan
//
// Example:
//   scan.WithIndex("StatusIndex") // Scan the StatusIndex instead of main table
//
// Returns the ScanBuilder for method chaining.
func (sb *ScanBuilder) WithIndex(indexName string) *ScanBuilder {
    sb.IndexName = indexName
    return sb
}

{{range .AllAttributes}}
// Filter{{ToSafeName .Name | ToUpperCamelCase}} adds a filter condition for "{{.Name}}" attribute.
// This method adds FilterExpression condition to DynamoDB Scan operation.
//
// DynamoDB filter attribute: "{{.Name}}" (type: {{.Type}})
// Go parameter type: {{ToGolangBaseType .}}
//
// Note: Filters are applied after items are read, so they don't reduce read capacity consumption.
// Use queries with key conditions when possible for better performance.
//
// Returns the ScanBuilder for method chaining.
func (sb *ScanBuilder) Filter{{ToSafeName .Name | ToUpperCamelCase}}({{ToSafeName .Name | ToLowerCamelCase}} {{ToGolangBaseType .}}) *ScanBuilder {
    condition := expression.Name("{{.Name}}").Equal(expression.Value({{ToSafeName .Name | ToLowerCamelCase}}))
    sb.FilterConditions = append(sb.FilterConditions, condition)
    sb.Attributes["{{.Name}}"] = {{ToSafeName .Name | ToLowerCamelCase}}
    sb.UsedKeys["{{.Name}}"] = true
    return sb
}
{{end}}

{{range .AllAttributes}}
{{if IsNumericAttr .}}
// Filter{{ToSafeName .Name | ToUpperCamelCase}}Between creates a range filter condition for the "{{.Name}}" attribute.
// This method adds a BETWEEN condition to the scan filter expression.
//
// DynamoDB condition: {{.Name}} BETWEEN start AND end (inclusive)
// Attribute type: {{.Type}} (numeric)
//
// Parameters:
//   - start: Lower bound of the range (inclusive)
//   - end: Upper bound of the range (inclusive)
//
// Example:
//   scan.Filter{{ToSafeName .Name | ToUpperCamelCase}}Between(100, 500) // {{.Name}} between 100 and 500
//
// Returns the ScanBuilder for method chaining.
func (sb *ScanBuilder) Filter{{ToSafeName .Name | ToUpperCamelCase}}Between(start, end {{ToGolangBaseType .}}) *ScanBuilder {
    condition := expression.Name("{{.Name}}").Between(expression.Value(start), expression.Value(end))
    sb.FilterConditions = append(sb.FilterConditions, condition)
    return sb
}

// Filter{{ToSafeName .Name | ToUpperCamelCase}}GreaterThan creates a "greater than" filter condition for the "{{.Name}}" attribute.
// This method adds a > condition to the scan filter expression.
//
// DynamoDB condition: {{.Name}} > value
// Attribute type: {{.Type}} (numeric)
//
// Parameters:
//   - value: The threshold value (exclusive lower bound)
//
// Example:
//   scan.Filter{{ToSafeName .Name | ToUpperCamelCase}}GreaterThan(1000) // {{.Name}} > 1000
//
// Returns the ScanBuilder for method chaining.
func (sb *ScanBuilder) Filter{{ToSafeName .Name | ToUpperCamelCase}}GreaterThan(value {{ToGolangBaseType .}}) *ScanBuilder {
    condition := expression.Name("{{.Name}}").GreaterThan(expression.Value(value))
    sb.FilterConditions = append(sb.FilterConditions, condition)
    return sb
}

// Filter{{ToSafeName .Name | ToUpperCamelCase}}LessThan creates a "less than" filter condition for the "{{.Name}}" attribute.
// This method adds a < condition to the scan filter expression.
//
// DynamoDB condition: {{.Name}} < value
// Attribute type: {{.Type}} (numeric)
//
// Parameters:
//   - value: The threshold value (exclusive upper bound)
//
// Example:
//   scan.Filter{{ToSafeName .Name | ToUpperCamelCase}}LessThan(500) // {{.Name}} < 500
//
// Returns the ScanBuilder for method chaining.
func (sb *ScanBuilder) Filter{{ToSafeName .Name | ToUpperCamelCase}}LessThan(value {{ToGolangBaseType .}}) *ScanBuilder {
    condition := expression.Name("{{.Name}}").LessThan(expression.Value(value))
    sb.FilterConditions = append(sb.FilterConditions, condition)
    return sb
}

// Filter{{ToSafeName .Name | ToUpperCamelCase}}GreaterThanOrEqual creates a ">=" filter condition for the "{{.Name}}" attribute.
// This method adds a >= condition to the scan filter expression.
//
// DynamoDB condition: {{.Name}} >= value
// Attribute type: {{.Type}} (numeric)
//
// Example:
//   scan.Filter{{ToSafeName .Name | ToUpperCamelCase}}GreaterThanOrEqual(100) // {{.Name}} >= 100
//
// Returns the ScanBuilder for method chaining.
func (sb *ScanBuilder) Filter{{ToSafeName .Name | ToUpperCamelCase}}GreaterThanOrEqual(value {{ToGolangBaseType .}}) *ScanBuilder {
    condition := expression.Name("{{.Name}}").GreaterThanEqual(expression.Value(value))
    sb.FilterConditions = append(sb.FilterConditions, condition)
    return sb
}

// Filter{{ToSafeName .Name | ToUpperCamelCase}}LessThanOrEqual creates a "<=" filter condition for the "{{.Name}}" attribute.
// This method adds a <= condition to the scan filter expression.
//
// DynamoDB condition: {{.Name}} <= value
// Attribute type: {{.Type}} (numeric)
//
// Example:
//   scan.Filter{{ToSafeName .Name | ToUpperCamelCase}}LessThanOrEqual(1000) // {{.Name}} <= 1000
//
// Returns the ScanBuilder for method chaining.
func (sb *ScanBuilder) Filter{{ToSafeName .Name | ToUpperCamelCase}}LessThanOrEqual(value {{ToGolangBaseType .}}) *ScanBuilder {
    condition := expression.Name("{{.Name}}").LessThanEqual(expression.Value(value))
    sb.FilterConditions = append(sb.FilterConditions, condition)
    return sb
}
{{end}}
{{end}}

{{range .AllAttributes}}
{{if eq (ToGolangBaseType .) "string"}}
// Filter{{ToSafeName .Name | ToUpperCamelCase}}Contains creates a contains filter condition for the "{{.Name}}" attribute.
// This method adds a contains() condition to the scan filter expression.
//
// DynamoDB condition: contains({{.Name}}, value)
// Attribute type: {{.Type}} (string)
//
// Parameters:
//   - value: The substring to search for
//
// Example:
//   scan.Filter{{ToSafeName .Name | ToUpperCamelCase}}Contains("gmail") // {{.Name}} contains "gmail"
//
// Returns the ScanBuilder for method chaining.
func (sb *ScanBuilder) Filter{{ToSafeName .Name | ToUpperCamelCase}}Contains(value {{ToGolangBaseType .}}) *ScanBuilder {
    condition := expression.Name("{{.Name}}").Contains(value)
    sb.FilterConditions = append(sb.FilterConditions, condition)
    return sb
}

// Filter{{ToSafeName .Name | ToUpperCamelCase}}BeginsWith creates a begins_with filter condition for the "{{.Name}}" attribute.
// This method adds a begins_with() condition to the scan filter expression.
//
// DynamoDB condition: begins_with({{.Name}}, value)
// Attribute type: {{.Type}} (string)
//
// Parameters:
//   - value: The prefix to search for
//
// Example:
//   scan.Filter{{ToSafeName .Name | ToUpperCamelCase}}BeginsWith("user_") // {{.Name}} begins with "user_"
//
// Returns the ScanBuilder for method chaining.
func (sb *ScanBuilder) Filter{{ToSafeName .Name | ToUpperCamelCase}}BeginsWith(value {{ToGolangBaseType .}}) *ScanBuilder {
    condition := expression.Name("{{.Name}}").BeginsWith(value)
    sb.FilterConditions = append(sb.FilterConditions, condition)
    return sb
}
{{end}}
{{end}}

// Limit restricts the maximum number of items returned by the scan.
// This is applied before any filter expressions, so the actual returned count
// may be lower if filters are applied.
//
// DynamoDB behavior: Sets the Limit parameter in ScanInput
// Range: 1 to 1MB of data (DynamoDB limitation)
//
// Parameters:
//   - limit: Maximum number of items to return
//
// Example:
//   scan.Limit(100) // Return at most 100 items
//
// Returns the ScanBuilder for method chaining.
func (sb *ScanBuilder) Limit(limit int) *ScanBuilder {
    sb.LimitValue = &limit
    return sb
}

// StartFrom enables pagination by specifying the exclusive start key for the scan.
// Use the LastEvaluatedKey from a previous scan response to continue pagination.
//
// DynamoDB behavior: Sets ExclusiveStartKey in ScanInput
// Pagination: Essential for handling large result sets
//
// Parameters:
//   - lastEvaluatedKey: The LastEvaluatedKey from the previous scan response
//
// Example:
//   // First scan
//   result1, _ := scan.Limit(100).Execute(ctx, client)
//   
//   // Continue pagination
//   result2, _ := scan.StartFrom(result1.LastEvaluatedKey).Execute(ctx, client)
//
// Returns the ScanBuilder for method chaining.
func (sb *ScanBuilder) StartFrom(lastEvaluatedKey map[string]types.AttributeValue) *ScanBuilder {
    sb.ExclusiveStartKey = lastEvaluatedKey
    return sb
}

// WithProjection specifies which attributes to return in the scan results.
// This can reduce network traffic and response size when you don't need all attributes.
//
// DynamoDB behavior: Sets ProjectionExpression in ScanInput
// Benefits: Reduced network traffic, faster responses, lower costs
//
// Parameters:
//   - attributes: List of attribute names to include in results
//
// Example:
//   scan.WithProjection([]string{"id", "name", "email"}) // Return only these attributes
//
// Returns the ScanBuilder for method chaining.
func (sb *ScanBuilder) WithProjection(attributes []string) *ScanBuilder {
    sb.ProjectionAttributes = attributes
    return sb
}

// WithParallelScan enables parallel scanning for better performance on large tables.
// Divides the table into segments and scans one segment. Use multiple ScanBuilder instances
// with different segment numbers to scan in parallel.
//
// DynamoDB behavior: Sets Segment and TotalSegments in ScanInput
// Performance: Can significantly improve throughput on large tables
//
// Parameters:
//   - totalSegments: Total number of segments to divide the table into (typically 2-4x your read capacity)
//   - segment: Current segment number (0-based, must be < totalSegments)
//
// Example:
//   // Scan segment 0 of 4 segments
//   scan1 := NewScanBuilder().WithParallelScan(4, 0)
//   
//   // Scan segment 1 of 4 segments  
//   scan2 := NewScanBuilder().WithParallelScan(4, 1)
//   
//   // Run both scans concurrently
//
// Returns the ScanBuilder for method chaining.
func (sb *ScanBuilder) WithParallelScan(totalSegments, segment int) *ScanBuilder {
    sb.ParallelScanConfig = &ParallelScanConfig{
        TotalSegments: totalSegments,
        Segment:       segment,
    }
    return sb
}

// BuildScan converts the ScanBuilder state into a complete DynamoDB ScanInput.
// This method handles the final transformation from high-level scan description
// to AWS SDK-compatible request structure.
//
// Process:
// 1. Combine all filter conditions with AND logic
// 2. Construct AWS SDK expression with filter conditions and projections
// 3. Configure ScanInput with all necessary parameters
// 4. Apply limiting, pagination, and parallel scan settings
//
// The generated ScanInput is ready for immediate execution with DynamoDB client.
//
// Returns:
//   - *dynamodb.ScanInput: Complete scan ready for DynamoDB execution
//   - error: Any errors in scan construction
//
// Example generated ScanInput:
//   &dynamodb.ScanInput{
//       TableName: "MyTable",
//       IndexName: "StatusIndex", 
//       FilterExpression: "#status = :s AND #created > :t",
//       ProjectionExpression: "#id, #name, #email",
//       Limit: 100,
//   }
func (sb *ScanBuilder) BuildScan() (*dynamodb.ScanInput, error) {
    input := &dynamodb.ScanInput{
        TableName: aws.String(TableName),
    }
    
    // Set index name if specified
    if sb.IndexName != "" {
        input.IndexName = aws.String(sb.IndexName)
    }
    
    // Build expression if we have filters or projections
    var exprBuilder expression.Builder
    hasExpression := false
    
    // Combine all filter conditions with AND logic
    if len(sb.FilterConditions) > 0 {
        combinedFilter := sb.FilterConditions[0]
        for _, condition := range sb.FilterConditions[1:] {
            combinedFilter = combinedFilter.And(condition)
        }
        exprBuilder = exprBuilder.WithFilter(combinedFilter)
        hasExpression = true
    }
    
    // Add projection expression if specified
    if len(sb.ProjectionAttributes) > 0 {
        var projectionBuilder expression.ProjectionBuilder
        for i, attr := range sb.ProjectionAttributes {
            if i == 0 {
                projectionBuilder = expression.NamesList(expression.Name(attr))
            } else {
                projectionBuilder = projectionBuilder.AddNames(expression.Name(attr))
            }
        }
        exprBuilder = exprBuilder.WithProjection(projectionBuilder)
        hasExpression = true
    }
    
    // Build and apply expression if we have any
    if hasExpression {
        expr, err := exprBuilder.Build()
        if err != nil {
            return nil, fmt.Errorf("failed to build scan expression: %v", err)
        }
        
        if len(sb.FilterConditions) > 0 {
            input.FilterExpression = expr.Filter()
        }
        
        if len(sb.ProjectionAttributes) > 0 {
            input.ProjectionExpression = expr.Projection()
        }
        
        if expr.Names() != nil {
            input.ExpressionAttributeNames = expr.Names()
        }
        
        if expr.Values() != nil {
            input.ExpressionAttributeValues = expr.Values()
        }
    }
    
    // Apply result limiting
    if sb.LimitValue != nil {
        input.Limit = aws.Int32(int32(*sb.LimitValue))
    }
    
    // Enable pagination if start key provided
    if sb.ExclusiveStartKey != nil {
        input.ExclusiveStartKey = sb.ExclusiveStartKey
    }
    
    // Configure parallel scanning if specified
    if sb.ParallelScanConfig != nil {
        input.Segment = aws.Int32(int32(sb.ParallelScanConfig.Segment))
        input.TotalSegments = aws.Int32(int32(sb.ParallelScanConfig.TotalSegments))
    }
    
    return input, nil
}

// Execute performs the complete scan lifecycle: build, execute, and unmarshal results.
// This is the primary method for end-users, providing a seamless experience from
// scan building to typed results.
//
// Execution flow:
// 1. Build optimized DynamoDB ScanInput
// 2. Execute scan against DynamoDB using provided client
// 3. Unmarshal raw DynamoDB items into strongly-typed SchemaItem structs
// 4. Return typed results with comprehensive error handling
//
// The method handles all AWS SDK complexity internally, providing a clean interface
// that returns application-ready data structures.
//
// Parameters:
//   - ctx: Request context for timeout/cancellation control
//   - client: AWS DynamoDB client for scan execution
//
// Returns:
//   - []SchemaItem: Strongly-typed scan results
//   - error: Any errors in scan building, execution, or unmarshaling
//
// Example usage:
//   items, err := NewScanBuilder().
//       FilterStatus("active").
//       FilterCreatedAtGreaterThan(timestamp).
//       WithProjection([]string{"id", "name", "email"}).
//       Limit(100).
//       Execute(ctx, dynamoClient)
//
//   if err != nil {
//       return fmt.Errorf("scan failed: %w", err)
//   }
//
//   for _, item := range items {
//       fmt.Printf("Found item: %s\n", item.Name)
//   }
func (sb *ScanBuilder) Execute(ctx context.Context, client *dynamodb.Client) ([]SchemaItem, error) {
    // Build the complete DynamoDB scan
    input, err := sb.BuildScan()
    if err != nil {
        return nil, err
    }
    
    // Execute scan against DynamoDB
    result, err := client.Scan(ctx, input)
    if err != nil {
        return nil, fmt.Errorf("failed to execute scan: %v", err)
    }
    
    // Unmarshal DynamoDB items into strongly-typed structs
    var items []SchemaItem
    err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal scan result: %v", err)
    }
    
    return items, nil
}
`
