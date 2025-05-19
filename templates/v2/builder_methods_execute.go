package templates

var QueryBuilderExecuteMethodsTemplate = `
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
}
`
