package v2

// QueryBuilderBuildTemplate generates the core query building and execution logic for QueryBuilder.
// This template creates intelligent query optimization methods that:
// - Automatically select the most efficient DynamoDB index based on provided attributes
// - Handle both simple and composite key conditions with smart prioritization
// - Build optimized filter expressions for non-key attributes
// - Generate AWS SDK v2 compatible QueryInput with all necessary parameters
// - Execute queries with proper error handling and result unmarshaling
// - Provide seamless integration between fluent API and DynamoDB operations
const QueryBuilderBuildTemplate = `
// Build analyzes the provided query parameters and constructs the optimal DynamoDB query.
// This method implements intelligent index selection using a multi-criteria algorithm:
//
// 1. Index Selection Priority:
//    - PreferredSortKey match (if specified by user)
//    - Composite key complexity (more specific indexes first)
//    - Hash key availability and compatibility
//
// 2. Fallback Strategy:
//    - Secondary indexes (GSI/LSI) are prioritized
//    - Primary table query as last resort
//    - Error if no suitable access pattern found
//
// The method performs comprehensive validation and optimization to ensure
// efficient DynamoDB operations with minimal consumed capacity.
//
// Returns:
//   - indexName: Selected index name ("" for primary table)
//   - keyCondition: DynamoDB key condition expression (hash + optional range)
//   - filterCondition: Filter expression for non-key attributes (nil if none)
//   - exclusiveStartKey: Pagination token for result continuation
//   - error: Query construction error (if no suitable index found)
//
// Example internal flow:
//   1. User calls: query.WithUserId("123").WithStatus("active").WithCreatedGreaterThan(timestamp)
//   2. Build() finds best index for userId + status combination
//   3. Sets created > timestamp as range condition
//   4. Returns optimized query components
func (qb *QueryBuilder) Build() (string, expression.KeyConditionBuilder, *expression.ConditionBuilder, map[string]types.AttributeValue, error) {
    var filterCond *expression.ConditionBuilder

    // Create a copy of indexes for sorting without modifying the original schema
    sortedIndexes := make([]SecondaryIndex, len(TableSchema.SecondaryIndexes))
    copy(sortedIndexes, TableSchema.SecondaryIndexes)
    
    // Smart index prioritization algorithm
    sort.Slice(sortedIndexes, func(i, j int) bool {
        // Priority 1: User-specified preferred sort key
        if qb.PreferredSortKey != "" {
            iMatches := sortedIndexes[i].RangeKey == qb.PreferredSortKey
            jMatches := sortedIndexes[j].RangeKey == qb.PreferredSortKey
            
            if iMatches && !jMatches {
                return true  // Prefer index with matching sort key
            }
            if !iMatches && jMatches {
                return false // Deprioritize index without matching sort key
            }
        }
        
        // Priority 2: Composite key complexity (more specific = higher priority)
        iParts := qb.calculateIndexParts(sortedIndexes[i])
        jParts := qb.calculateIndexParts(sortedIndexes[j])
        
        return iParts > jParts // More complex composite keys are more specific
    })

    // Try each index in priority order
    for _, idx := range sortedIndexes {
        // Check if hash key requirements are satisfied
        hashKeyCondition, hashKeyMatch := qb.buildHashKeyCondition(idx)
        if !hashKeyMatch {
            continue // Skip index if hash key requirements not met
        }

        // Check if range key requirements are satisfied (optional)
        rangeKeyCondition, rangeKeyMatch := qb.buildRangeKeyCondition(idx)
        if !rangeKeyMatch {
            continue // Skip index if range key requirements not met
        }

        // Combine hash and range key conditions
        keyCondition := *hashKeyCondition
        if rangeKeyCondition != nil {
            keyCondition = keyCondition.And(*rangeKeyCondition)
        }

        // Build filter conditions for non-key attributes
        filterCond = qb.buildFilterCondition(idx)

        // Return the first viable index (highest priority due to sorting)
        return idx.Name, keyCondition, filterCond, qb.ExclusiveStartKey, nil
    }

    // Fallback: Try primary table if hash key is available
    if qb.UsedKeys[TableSchema.HashKey] {
        indexName := "" // Empty string indicates primary table query
        keyCondition := expression.Key(TableSchema.HashKey).Equal(expression.Value(qb.Attributes[TableSchema.HashKey]))

        // Add range key condition if available
        if TableSchema.RangeKey != "" && qb.UsedKeys[TableSchema.RangeKey] {
            if cond, exists := qb.KeyConditions[TableSchema.RangeKey]; exists {
                // Use pre-built range condition (Between, GreaterThan, etc.)
                keyCondition = keyCondition.And(cond)
            } else {
                // Use simple equality for range key
                keyCondition = keyCondition.And(expression.Key(TableSchema.RangeKey).Equal(expression.Value(qb.Attributes[TableSchema.RangeKey])))
            }
        }

        // Build filter conditions for all non-key attributes
        var filterConditions []expression.ConditionBuilder
        for attrName, value := range qb.Attributes {
            if attrName != TableSchema.HashKey && attrName != TableSchema.RangeKey {
                filterConditions = append(filterConditions, expression.Name(attrName).Equal(expression.Value(value)))
            }
        }

        // Combine multiple filter conditions with AND logic
        if len(filterConditions) > 0 {
            combinedFilter := filterConditions[0]
            for _, cond := range filterConditions[1:] {
                combinedFilter = combinedFilter.And(cond)
            }
            filterCond = &combinedFilter
        }

        return indexName, keyCondition, filterCond, qb.ExclusiveStartKey, nil
    }

    // No suitable access pattern found
    return "", expression.KeyConditionBuilder{}, nil, nil, fmt.Errorf("no suitable index found for the provided keys")
}

// calculateIndexParts computes the complexity score for an index based on its composite key structure.
// More complex indexes (with more composite key parts) are prioritized as they provide more specific access patterns.
//
// Composite key complexity scoring:
// - Simple keys (user_id): 0 parts
// - Two-part composite (user_id#status): 2 parts  
// - Three-part composite (level#category#status): 3 parts
//
// Higher scores indicate more specific access patterns and better query performance.
//
// Parameters:
//   - idx: SecondaryIndex to analyze
//
// Returns the total number of composite key parts (hash + range).
func (qb *QueryBuilder) calculateIndexParts(idx SecondaryIndex) int {
    parts := 0
    if idx.HashKeyParts != nil {
        parts += len(idx.HashKeyParts)
    }
    if idx.RangeKeyParts != nil {
        parts += len(idx.RangeKeyParts)
    }
    return parts
}

// buildHashKeyCondition constructs the hash key condition for a specific index.
// Handles both simple and composite hash keys with proper validation.
//
// For composite keys: Validates that all required attribute values are available
// For simple keys: Checks that the hash key attribute is set
//
// Returns:
//   - KeyConditionBuilder: The constructed hash key condition
//   - bool: true if condition was successfully built, false if requirements not met
//
// Example composite key logic:
//   HashKey: "user_id#tenant_id", Parts: [{user_id, false}, {tenant_id, false}]
//   Required: qb.Attributes["user_id"] and qb.Attributes["tenant_id"] must exist
//   Result: expression.Key("user_id#tenant_id").Equal(expression.Value("123#org456"))
func (qb *QueryBuilder) buildHashKeyCondition(idx SecondaryIndex) (*expression.KeyConditionBuilder, bool) {
    if idx.HashKeyParts != nil {
        // Composite hash key: check all parts are available
        if qb.hasAllKeys(idx.HashKeyParts) {
            cond := qb.buildCompositeKeyCondition(idx.HashKeyParts)
            return &cond, true
        }
    } else if idx.HashKey != "" && qb.UsedKeys[idx.HashKey] {
        // Simple hash key: direct equality condition
        cond := expression.Key(idx.HashKey).Equal(expression.Value(qb.Attributes[idx.HashKey]))
        return &cond, true
    }
    return nil, false
}

// buildRangeKeyCondition constructs the range key condition for a specific index.
// Supports both simple and composite range keys with flexible condition types.
//
// Range key handling:
// - If no range key: Returns (nil, true) - valid but no range condition
// - If range key exists but not provided: Returns (nil, true) - optional range key
// - If range key provided: Builds appropriate condition (equality, range, etc.)
//
// Range condition types supported:
// - Equality: WithXxxRangeKey(value)
// - Range: WithXxxBetween(start, end), WithXxxGreaterThan(value), WithXxxLessThan(value)
//
// Returns:
//   - KeyConditionBuilder: The constructed range key condition (nil if no range key)
//   - bool: true if the index is compatible, false if requirements not met
func (qb *QueryBuilder) buildRangeKeyCondition(idx SecondaryIndex) (*expression.KeyConditionBuilder, bool) {
    if idx.RangeKeyParts != nil {
        // Composite range key: check all parts are available
        if qb.hasAllKeys(idx.RangeKeyParts) {
            cond := qb.buildCompositeKeyCondition(idx.RangeKeyParts)
            return &cond, true
        }
    } else if idx.RangeKey != "" {
        if qb.UsedKeys[idx.RangeKey] {
            // Check for pre-built range conditions (Between, GreaterThan, LessThan)
            if cond, exists := qb.KeyConditions[idx.RangeKey]; exists {
                return &cond, true
            } else {
                // No pre-built condition, but range key is available for equality
                return nil, true
            }
        } else {
            // Range key exists but not provided - still a valid index for hash-only queries
            return nil, true
        }
    } else {
        // No range key for this index - valid
        return nil, true
    }
    return nil, false
}

// buildFilterCondition creates filter expressions for non-key attributes.
// Filter conditions are applied after the key conditions and don't benefit from index optimization.
//
// Filter logic:
// 1. Identify attributes that are not part of the selected index's key structure
// 2. Create equality conditions for each non-key attribute
// 3. Combine multiple conditions with AND logic
//
// Performance note: Filter expressions consume additional read capacity and may
// require scanning more items than the final result set.
//
// Parameters:
//   - idx: The selected index to determine which attributes are key vs non-key
//
// Returns filter condition builder (nil if no filter conditions needed).
func (qb *QueryBuilder) buildFilterCondition(idx SecondaryIndex) *expression.ConditionBuilder {
    var filterConditions []expression.ConditionBuilder

    // Check each attribute to see if it needs to be filtered
    for attrName, value := range qb.Attributes {
        if qb.isPartOfIndexKey(attrName, idx) {
            continue // Skip key attributes - they're handled in key conditions
        }
        // Add equality filter for non-key attributes
        filterConditions = append(filterConditions, expression.Name(attrName).Equal(expression.Value(value)))
    }

    if len(filterConditions) == 0 {
        return nil // No filter conditions needed
    }

    // Combine multiple filter conditions with AND logic
    combinedFilter := filterConditions[0]
    for _, cond := range filterConditions[1:] {
        combinedFilter = combinedFilter.And(cond)
    }
    return &combinedFilter
}

// isPartOfIndexKey determines if an attribute is part of the selected index's key structure.
// This helps distinguish between key conditions (efficient) and filter conditions (less efficient).
//
// Key attribute identification:
// - Hash key parts (for composite hash keys)
// - Simple hash key
// - Range key parts (for composite range keys)  
// - Simple range key
//
// Parameters:
//   - attrName: Name of the attribute to check
//   - idx: Index to check against
//
// Returns true if the attribute is part of the index key structure.
func (qb *QueryBuilder) isPartOfIndexKey(attrName string, idx SecondaryIndex) bool {
    // Check if attribute is part of composite hash key
    if idx.HashKeyParts != nil {
        for _, part := range idx.HashKeyParts {
            if !part.IsConstant && part.Value == attrName {
                return true
            }
        }
    } else if attrName == idx.HashKey {
        // Simple hash key match
        return true
    }
    
    // Check if attribute is part of composite range key
    if idx.RangeKeyParts != nil {
        for _, part := range idx.RangeKeyParts {
            if !part.IsConstant && part.Value == attrName {
                return true
            }
        }
    } else if attrName == idx.RangeKey {
        // Simple range key match
        return true
    }
    
    return false
}

// BuildQuery converts the QueryBuilder state into a complete DynamoDB QueryInput.
// This method handles the final transformation from high-level query description
// to AWS SDK-compatible request structure.
//
// Process:
// 1. Call Build() to determine optimal index and conditions
// 2. Construct AWS SDK expression with key conditions and filters
// 3. Configure QueryInput with all necessary parameters
// 4. Apply sorting, limiting, and pagination settings
//
// The generated QueryInput is ready for immediate execution with DynamoDB client.
//
// Returns:
//   - *dynamodb.QueryInput: Complete query ready for DynamoDB execution
//   - error: Any errors in query construction
//
// Example generated QueryInput:
//   &dynamodb.QueryInput{
//       TableName: "MyTable",
//       IndexName: "UserStatusIndex", 
//       KeyConditionExpression: "user_id = :u AND #status = :s",
//       FilterExpression: "#created > :t",
//       ScanIndexForward: false,
//       Limit: 10,
//   }
func (qb *QueryBuilder) BuildQuery() (*dynamodb.QueryInput, error) {
    // Get optimized query components
    indexName, keyCond, filterCond, exclusiveStartKey, err := qb.Build()
    if err != nil {
        return nil, err
    }

    // Build AWS SDK expression with key conditions
    exprBuilder := expression.NewBuilder().WithKeyCondition(keyCond)
    if filterCond != nil {
        exprBuilder = exprBuilder.WithFilter(*filterCond)
    }

    // Compile expression into AWS SDK format
    expr, err := exprBuilder.Build()
    if err != nil {
        return nil, fmt.Errorf("failed to build expression: %v", err)
    }

    // Construct complete QueryInput
    input := &dynamodb.QueryInput{
        TableName:                 aws.String(TableName),
        KeyConditionExpression:    expr.KeyCondition(),
        ExpressionAttributeNames:  expr.Names(),
        ExpressionAttributeValues: expr.Values(),
        ScanIndexForward:          aws.Bool(!qb.SortDescending), // Convert to DynamoDB format
    }

    // Set index name if querying a secondary index
    if indexName != "" {
        input.IndexName = aws.String(indexName)
    }

    // Add filter expression if present
    if filterCond != nil {
        input.FilterExpression = expr.Filter()
    }

    // Apply result limiting
    if qb.LimitValue != nil {
        input.Limit = aws.Int32(int32(*qb.LimitValue))
    }

    // Enable pagination if start key provided
    if exclusiveStartKey != nil {
        input.ExclusiveStartKey = exclusiveStartKey
    }

    return input, nil
}

// Execute performs the complete query lifecycle: build, execute, and unmarshal results.
// This is the primary method for end-users, providing a seamless experience from
// query building to typed results.
//
// Execution flow:
// 1. Build optimized DynamoDB QueryInput
// 2. Execute query against DynamoDB using provided client
// 3. Unmarshal raw DynamoDB items into strongly-typed SchemaItem structs
// 4. Return typed results with comprehensive error handling
//
// The method handles all AWS SDK complexity internally, providing a clean interface
// that returns application-ready data structures.
//
// Parameters:
//   - ctx: Request context for timeout/cancellation control
//   - client: AWS DynamoDB client for query execution
//
// Returns:
//   - []SchemaItem: Strongly-typed query results
//   - error: Any errors in query building, execution, or unmarshaling
//
// Example usage:
//   items, err := NewQueryBuilder().
//       WithUserId("user123").
//       WithCreatedGreaterThan(timestamp).
//       OrderByDesc().
//       Limit(50).
//       Execute(ctx, dynamoClient)
//
//   if err != nil {
//       return fmt.Errorf("query failed: %w", err)
//   }
//
//   for _, item := range items {
//       fmt.Printf("Found item: %s\n", item.Name)
//   }
func (qb *QueryBuilder) Execute(ctx context.Context, client *dynamodb.Client) ([]SchemaItem, error) {
    // Build the complete DynamoDB query
    input, err := qb.BuildQuery()
    if err != nil {
        return nil, err
    }

    // Execute query against DynamoDB
    result, err := client.Query(ctx, input)
    if err != nil {
        return nil, fmt.Errorf("failed to execute query: %v", err)
    }

    // Unmarshal DynamoDB items into strongly-typed structs
    var items []SchemaItem
    err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal result: %v", err)
    }

    return items, nil
}`
