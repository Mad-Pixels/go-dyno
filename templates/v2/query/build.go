package query

// QueryBuilderBuildTemplate provides intelligent query building with automatic index selection
const QueryBuilderBuildTemplate = `
// Build analyzes the query conditions and selects the optimal index for execution.
// Implements smart index selection algorithm considering:
// - Preferred sort key hints from user
// - Number of composite key parts matched
// - Index efficiency for the given query pattern
// Returns index name, key conditions, filter conditions, pagination key, and any errors.
func (qb *QueryBuilder) Build() (string, expression.KeyConditionBuilder, *expression.ConditionBuilder, map[string]types.AttributeValue, error) {
    var filterCond *expression.ConditionBuilder
    sortedIndexes := make([]SecondaryIndex, len(TableSchema.SecondaryIndexes))
    copy(sortedIndexes, TableSchema.SecondaryIndexes)
    
    sort.Slice(sortedIndexes, func(i, j int) bool {
        if qb.PreferredSortKey != "" {
            iMatches := sortedIndexes[i].RangeKey == qb.PreferredSortKey
            jMatches := sortedIndexes[j].RangeKey == qb.PreferredSortKey
            
            if iMatches && !jMatches {
                return true
            }
            if !iMatches && jMatches {
                return false
            }
        }
        iParts := qb.calculateIndexParts(sortedIndexes[i])
        jParts := qb.calculateIndexParts(sortedIndexes[j])
        return iParts > jParts
    })

    for _, idx := range sortedIndexes {
        hashKeyCondition, hashKeyMatch := qb.buildHashKeyCondition(idx)
        if !hashKeyMatch {
            continue
        }
        rangeKeyCondition, rangeKeyMatch := qb.buildRangeKeyCondition(idx)
        if !rangeKeyMatch {
            continue
        }
        keyCondition := *hashKeyCondition
        if rangeKeyCondition != nil {
            keyCondition = keyCondition.And(*rangeKeyCondition)
        }
        filterCond = qb.buildFilterCondition(idx)
        return idx.Name, keyCondition, filterCond, qb.ExclusiveStartKey, nil
    }

    if qb.UsedKeys[TableSchema.HashKey] {
        indexName := ""
        keyCondition := expression.Key(TableSchema.HashKey).Equal(expression.Value(qb.Attributes[TableSchema.HashKey]))

        if TableSchema.RangeKey != "" && qb.UsedKeys[TableSchema.RangeKey] {
            if cond, exists := qb.KeyConditions[TableSchema.RangeKey]; exists {
                keyCondition = keyCondition.And(cond)
            } else {
                keyCondition = keyCondition.And(expression.Key(TableSchema.RangeKey).Equal(expression.Value(qb.Attributes[TableSchema.RangeKey])))
            }
        }
        var filterConditions []expression.ConditionBuilder
        filterConditions = append(filterConditions, qb.FilterConditions...)
        for attrName, value := range qb.Attributes {
            if attrName != TableSchema.HashKey && attrName != TableSchema.RangeKey {
                filterConditions = append(filterConditions, expression.Name(attrName).Equal(expression.Value(value)))
            }
        }
        if len(filterConditions) > 0 {
            combinedFilter := filterConditions[0]
            for _, cond := range filterConditions[1:] {
                combinedFilter = combinedFilter.And(cond)
            }
            filterCond = &combinedFilter
        }
        return indexName, keyCondition, filterCond, qb.ExclusiveStartKey, nil
    }
    return "", expression.KeyConditionBuilder{}, nil, nil, fmt.Errorf("no suitable index found for the provided keys")
}

// calculateIndexParts counts the number of composite key parts in an index.
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

// buildHashKeyCondition creates the hash key condition for a given index.
func (qb *QueryBuilder) buildHashKeyCondition(idx SecondaryIndex) (*expression.KeyConditionBuilder, bool) {
    if idx.HashKeyParts != nil {
        if qb.hasAllKeys(idx.HashKeyParts) {
            cond := qb.buildCompositeKeyCondition(idx.HashKeyParts)
            return &cond, true
        }
    } else if idx.HashKey != "" && qb.UsedKeys[idx.HashKey] {
        cond := expression.Key(idx.HashKey).Equal(expression.Value(qb.Attributes[idx.HashKey]))
        return &cond, true
    }
    return nil, false
}

// buildRangeKeyCondition creates the range key condition for a given index.
func (qb *QueryBuilder) buildRangeKeyCondition(idx SecondaryIndex) (*expression.KeyConditionBuilder, bool) {
    if idx.RangeKeyParts != nil {
        if qb.hasAllKeys(idx.RangeKeyParts) {
            cond := qb.buildCompositeKeyCondition(idx.RangeKeyParts)
            return &cond, true
        }
    } else if idx.RangeKey != "" {
        if qb.UsedKeys[idx.RangeKey] {
            if cond, exists := qb.KeyConditions[idx.RangeKey]; exists {
                return &cond, true
            } else {
                cond := expression.Key(idx.RangeKey).Equal(expression.Value(qb.Attributes[idx.RangeKey]))
                return &cond, true
            }
        } else {
            return nil, true
        }
    } else {
        return nil, true
    }
    return nil, false
}

// buildFilterCondition creates filter conditions for attributes not part of the index keys.
func (qb *QueryBuilder) buildFilterCondition(idx SecondaryIndex) *expression.ConditionBuilder {
    var filterConditions []expression.ConditionBuilder
    
    filterConditions = append(filterConditions, qb.FilterConditions...)
    for attrName, value := range qb.Attributes {
        if qb.isPartOfIndexKey(attrName, idx) {
            continue
        }
        filterConditions = append(filterConditions, expression.Name(attrName).Equal(expression.Value(value)))
    }
    if len(filterConditions) == 0 {
        return nil
    }
    combinedFilter := filterConditions[0]
    for _, cond := range filterConditions[1:] {
        combinedFilter = combinedFilter.And(cond)
    }
    return &combinedFilter
}

// isPartOfIndexKey checks if an attribute is part of the index's key structure.
func (qb *QueryBuilder) isPartOfIndexKey(attrName string, idx SecondaryIndex) bool {
    if idx.HashKeyParts != nil {
        for _, part := range idx.HashKeyParts {
            if !part.IsConstant && part.Value == attrName {
                return true
            }
        }
    } else if attrName == idx.HashKey {
        return true
    }
    if idx.RangeKeyParts != nil {
        for _, part := range idx.RangeKeyParts {
            if !part.IsConstant && part.Value == attrName {
                return true
            }
        }
    } else if attrName == idx.RangeKey {
        return true
    }
    return false
}

// BuildQuery constructs the final DynamoDB QueryInput with all expressions and parameters.
// Combines key conditions, filter conditions, pagination, and sorting options.
func (qb *QueryBuilder) BuildQuery() (*dynamodb.QueryInput, error) {
    indexName, keyCond, filterCond, exclusiveStartKey, err := qb.Build()
    if err != nil {
        return nil, err
    }
    exprBuilder := expression.NewBuilder().WithKeyCondition(keyCond)
    if filterCond != nil {
        exprBuilder = exprBuilder.WithFilter(*filterCond)
    }
    expr, err := exprBuilder.Build()
    if err != nil {
        return nil, fmt.Errorf("failed to build expression: %v", err)
    }
    input := &dynamodb.QueryInput{
        TableName:                 aws.String(TableName),
        KeyConditionExpression:    expr.KeyCondition(),
        ExpressionAttributeNames:  expr.Names(),
        ExpressionAttributeValues: expr.Values(),
        ScanIndexForward:          aws.Bool(!qb.SortDescending),
    }
    if indexName != "" {
        input.IndexName = aws.String(indexName)
    }
    if filterCond != nil {
        input.FilterExpression = expr.Filter()
    }
    if qb.LimitValue != nil {
        input.Limit = aws.Int32(int32(*qb.LimitValue))
    }
    if exclusiveStartKey != nil {
        input.ExclusiveStartKey = exclusiveStartKey
    }
    return input, nil
}

// Execute runs the query against DynamoDB and returns strongly-typed results.
// Handles the complete query lifecycle: build input, execute, unmarshal results.
func (qb *QueryBuilder) Execute(ctx context.Context, client *dynamodb.Client) ([]SchemaItem, error) {
    input, err := qb.BuildQuery()
    if err != nil {
        return nil, err
    }
    result, err := client.Query(ctx, input)
    if err != nil {
        return nil, fmt.Errorf("failed to execute query: %v", err)
    }
    var items []SchemaItem
    err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal result: %v", err)
    }
    return items, nil
}
`
