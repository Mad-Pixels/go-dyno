package inputs

// UpdateInputsTemplate ...
const UpdateInputsTemplate = `
// UpdateItemInput ...
func UpdateItemInput(item SchemaItem) (*dynamodb.UpdateItemInput, error) {
    key, err := KeyInput(item)
    if err != nil {
        return nil, fmt.Errorf("failed to create key from item for update: %v", err)
    }
   
    allAttributes, err := attributevalue.MarshalMap(item)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal item for update: %v", err)
    }
   
    updates := make(map[string]types.AttributeValue)
    for attrName, attrValue := range allAttributes {
        if attrName != TableSchema.HashKey && attrName != TableSchema.RangeKey {
            updates[attrName] = attrValue
        }
    }
   
    if len(updates) == 0 {
        return nil, fmt.Errorf("no non-key attributes to update")
    }
   
    var updateExpressionParts []string
    expressionAttributeNames := make(map[string]string)
    expressionAttributeValues := make(map[string]types.AttributeValue)
   
    i := 0
    for attrName, attrValue := range updates {
        nameKey := fmt.Sprintf("#attr%d", i)
        valueKey := fmt.Sprintf(":val%d", i)
       
        updateExpressionParts = append(updateExpressionParts, fmt.Sprintf("%s = %s", nameKey, valueKey))
        expressionAttributeNames[nameKey] = attrName
        expressionAttributeValues[valueKey] = attrValue
        i++
    }
   
    updateExpression := "SET " + strings.Join(updateExpressionParts, ", ")
   
    return &dynamodb.UpdateItemInput{
        TableName:                 aws.String(TableSchema.TableName),
        Key:                       key,
        UpdateExpression:          aws.String(updateExpression),
        ExpressionAttributeNames:  expressionAttributeNames,
        ExpressionAttributeValues: expressionAttributeValues,
    }, nil
}

// UpdateItemInputFromRaw ...
func UpdateItemInputFromRaw(hashKeyValue interface{}, rangeKeyValue interface{}, updates map[string]interface{}) (*dynamodb.UpdateItemInput, error) {
    key, err := KeyInputFromRaw(hashKeyValue, rangeKeyValue)
    if err != nil {
        return nil, fmt.Errorf("failed to create key for update: %v", err)
    }
   
    if len(updates) == 0 {
        return nil, fmt.Errorf("no updates provided")
    }
   
    var updateExpressionParts []string
    expressionAttributeNames := make(map[string]string)
    expressionAttributeValues := make(map[string]types.AttributeValue)
   
    i := 0
    for attrName, value := range updates {
        nameKey := fmt.Sprintf("#attr%d", i)
        valueKey := fmt.Sprintf(":val%d", i)
       
        updateExpressionParts = append(updateExpressionParts, fmt.Sprintf("%s = %s", nameKey, valueKey))
        expressionAttributeNames[nameKey] = attrName
       
        av, err := attributevalue.Marshal(value)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal update value for %s: %v", attrName, err)
        }
        expressionAttributeValues[valueKey] = av
        i++
    }
   
    updateExpression := "SET " + strings.Join(updateExpressionParts, ", ")
   
    return &dynamodb.UpdateItemInput{
        TableName:                 aws.String(TableSchema.TableName),
        Key:                       key,
        UpdateExpression:          aws.String(updateExpression),
        ExpressionAttributeNames:  expressionAttributeNames,
        ExpressionAttributeValues: expressionAttributeValues,
    }, nil
}

// UpdateItemInputWithCondition ...
func UpdateItemInputWithCondition(hashKeyValue interface{}, rangeKeyValue interface{}, updates map[string]interface{}, conditionExpression string, conditionAttributeNames map[string]string, conditionAttributeValues map[string]types.AttributeValue) (*dynamodb.UpdateItemInput, error) {
    updateInput, err := UpdateItemInputFromRaw(hashKeyValue, rangeKeyValue, updates)
    if err != nil {
        return nil, err
    }
   
    updateInput.ConditionExpression = aws.String(conditionExpression)
   
    if conditionAttributeNames != nil {
        for key, value := range conditionAttributeNames {
            updateInput.ExpressionAttributeNames[key] = value
        }
    }
   
    if conditionAttributeValues != nil {
        for key, value := range conditionAttributeValues {
            updateInput.ExpressionAttributeValues[key] = value
        }
    }
   
    return updateInput, nil
}

// UpdateItemInputWithExpression ...
func UpdateItemInputWithExpression(hashKeyValue interface{}, rangeKeyValue interface{}, updateBuilder expression.UpdateBuilder, conditionBuilder *expression.ConditionBuilder) (*dynamodb.UpdateItemInput, error) {
    key, err := KeyInputFromRaw(hashKeyValue, rangeKeyValue)
    if err != nil {
        return nil, fmt.Errorf("failed to create key for expression update: %v", err)
    }
   
    var expr expression.Expression
    if conditionBuilder != nil {
        expr, err = expression.NewBuilder().
            WithUpdate(updateBuilder).
            WithCondition(*conditionBuilder).
            Build()
    } else {
        expr, err = expression.NewBuilder().
            WithUpdate(updateBuilder).
            Build()
    }
   
    if err != nil {
        return nil, fmt.Errorf("failed to build update expression: %v", err)
    }
   
    input := &dynamodb.UpdateItemInput{
        TableName:                 aws.String(TableSchema.TableName),
        Key:                       key,
        UpdateExpression:          expr.Update(),
        ExpressionAttributeNames:  expr.Names(),
        ExpressionAttributeValues: expr.Values(),
    }
   
    if conditionBuilder != nil {
        input.ConditionExpression = expr.Condition()
    }
   
    return input, nil
}`
