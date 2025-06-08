package scan

const ScanBuilderBuildTemplate = `
// BuildScan ...
func (sb *ScanBuilder) BuildScan() (*dynamodb.ScanInput, error) {
    input := &dynamodb.ScanInput{
        TableName: aws.String(TableName),
    }
    
    if sb.IndexName != "" {
        input.IndexName = aws.String(sb.IndexName)
    }
    
    var exprBuilder expression.Builder
    hasExpression := false
    
    if len(sb.FilterConditions) > 0 {
        combinedFilter := sb.FilterConditions[0]
        for _, condition := range sb.FilterConditions[1:] {
            combinedFilter = combinedFilter.And(condition)
        }
        exprBuilder = exprBuilder.WithFilter(combinedFilter)
        hasExpression = true
    }
    
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
    
    if sb.LimitValue != nil {
        input.Limit = aws.Int32(int32(*sb.LimitValue))
    }
    
    if sb.ExclusiveStartKey != nil {
        input.ExclusiveStartKey = sb.ExclusiveStartKey
    }
    
    if sb.ParallelScanConfig != nil {
        input.Segment = aws.Int32(int32(sb.ParallelScanConfig.Segment))
        input.TotalSegments = aws.Int32(int32(sb.ParallelScanConfig.TotalSegments))
    }
    
    return input, nil
}

// Execute ...
func (sb *ScanBuilder) Execute(ctx context.Context, client *dynamodb.Client) ([]SchemaItem, error) {
    input, err := sb.BuildScan()
    if err != nil {
        return nil, err
    }
    
    result, err := client.Scan(ctx, input)
    if err != nil {
        return nil, fmt.Errorf("failed to execute scan: %v", err)
    }
    
    var items []SchemaItem
    err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal scan result: %v", err)
    }
    
    return items, nil
}
`
