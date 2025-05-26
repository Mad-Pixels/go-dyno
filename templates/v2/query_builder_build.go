package v2

// QueryBuilderBuildTemplate ...
const QueryBuilderBuildTemplate = `
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
    
    iParts := 0
    if sortedIndexes[i].HashKeyParts != nil {
        iParts += len(sortedIndexes[i].HashKeyParts)
    }
    if sortedIndexes[i].RangeKeyParts != nil {
        iParts += len(sortedIndexes[i].RangeKeyParts)
    }
    
    jParts := 0
    if sortedIndexes[j].HashKeyParts != nil {
        jParts += len(sortedIndexes[j].HashKeyParts)
    }
    if sortedIndexes[j].RangeKeyParts != nil {
        jParts += len(sortedIndexes[j].RangeKeyParts)
    }
    
    return iParts > jParts
})

    for _, idx := range sortedIndexes {
        var hashKeyCondition, rangeKeyCondition *expression.KeyConditionBuilder
        var hashKeyMatch, rangeKeyMatch bool

        if idx.HashKeyParts != nil {
            if qb.hasAllKeys(idx.HashKeyParts) {
                cond := qb.buildCompositeKeyCondition(idx.HashKeyParts)
                hashKeyCondition = &cond
                hashKeyMatch = true
            }
        } else if idx.HashKey != "" && qb.UsedKeys[idx.HashKey] {
            cond := expression.Key(idx.HashKey).Equal(expression.Value(qb.Attributes[idx.HashKey]))
            hashKeyCondition = &cond
            hashKeyMatch = true
        }

        if !hashKeyMatch {
            continue // Этот индекс не подходит
        }

if idx.RangeKeyParts != nil {
    if qb.hasAllKeys(idx.RangeKeyParts) {
        cond := qb.buildCompositeKeyCondition(idx.RangeKeyParts)
        rangeKeyCondition = &cond
        rangeKeyMatch = true
    }
} else if idx.RangeKey != "" {
    if qb.UsedKeys[idx.RangeKey] {
        if cond, exists := qb.KeyConditions[idx.RangeKey]; exists {
            rangeKeyCondition = &cond
            rangeKeyMatch = true
        } else {
            rangeKeyMatch = true
        }
    } else {
        rangeKeyMatch = true
    }
} else {
    rangeKeyMatch = true
}

        if !rangeKeyMatch {
            continue
        }

        keyCondition := *hashKeyCondition
        if rangeKeyCondition != nil {
            keyCondition = keyCondition.And(*rangeKeyCondition)
        }

        for attrName, value := range qb.Attributes {
            isPartOfHashKey := false
            isPartOfRangeKey := false
            
            if idx.HashKeyParts != nil {
                for _, part := range idx.HashKeyParts {
                    if !part.IsConstant && part.Value == attrName {
                        isPartOfHashKey = true
                        break
                    }
                }
            } else if attrName == idx.HashKey {
                isPartOfHashKey = true
            }
            
            if idx.RangeKeyParts != nil {
                for _, part := range idx.RangeKeyParts {
                    if !part.IsConstant && part.Value == attrName {
                        isPartOfRangeKey = true
                        break
                    }
                }
            } else if attrName == idx.RangeKey {
                isPartOfRangeKey = true
            }
            
            // Если атрибут не является частью ключа, добавляем его в фильтр
            if !isPartOfHashKey && !isPartOfRangeKey {
                cond := expression.Name(attrName).Equal(expression.Value(value))
                qb.FilterConditions = append(qb.FilterConditions, cond)
            }
        }

        if len(qb.FilterConditions) > 0 {
            combinedFilter := qb.FilterConditions[0]
            for _, cond := range qb.FilterConditions[1:] {
                combinedFilter = combinedFilter.And(cond)
            }
            filterCond = &combinedFilter
        }

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

        for attrName, value := range qb.Attributes {
            if attrName != TableSchema.HashKey && attrName != TableSchema.RangeKey {
                cond := expression.Name(attrName).Equal(expression.Value(value))
                qb.FilterConditions = append(qb.FilterConditions, cond)
            }
        }

        if len(qb.FilterConditions) > 0 {
            combinedFilter := qb.FilterConditions[0]
            for _, cond := range qb.FilterConditions[1:] {
                combinedFilter = combinedFilter.And(cond)
            }
            filterCond = &combinedFilter
        }

        return indexName, keyCondition, filterCond, qb.ExclusiveStartKey, nil
    }

    return "", expression.KeyConditionBuilder{}, nil, nil, fmt.Errorf("no suitable index found for the provided keys")
}

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
}`
