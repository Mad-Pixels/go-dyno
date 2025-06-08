package inputs

// UpdateInputsTemplate ...
const UpdateInputsTemplate = `
// UpdateItemInput creates update request from SchemaItem
func UpdateItemInput(item SchemaItem) (*dynamodb.UpdateItemInput, error) {
    key, err := KeyInput(item)
    if err != nil {
        return nil, fmt.Errorf("failed to create key from item for update: %v", err)
    }
   
    // Use helper to marshal item
    allAttributes, err := marshalItemToMap(item)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal item for update: %v", err)
    }
   
    // Use helper to extract non-key attributes
    updates := extractNonKeyAttributes(allAttributes)
    if len(updates) == 0 {
        return nil, fmt.Errorf("no non-key attributes to update")
    }
   
    // Use helper to build expression
    updateExpression, attrNames, attrValues := buildUpdateExpression(updates)
   
    return &dynamodb.UpdateItemInput{
        TableName:                 aws.String(TableSchema.TableName),
        Key:                       key,
        UpdateExpression:          aws.String(updateExpression),
        ExpressionAttributeNames:  attrNames,
        ExpressionAttributeValues: attrValues,
    }, nil
}

// UpdateItemInputFromRaw creates update request from raw values
func UpdateItemInputFromRaw(hashKeyValue interface{}, rangeKeyValue interface{}, updates map[string]interface{}) (*dynamodb.UpdateItemInput, error) {
    // All validations at the beginning
    if err := validateKeyInputs(hashKeyValue, rangeKeyValue); err != nil {
        return nil, err
    }
    if err := validateUpdatesMap(updates); err != nil {
        return nil, err
    }

    // Pure business logic after validation
    key, err := KeyInputFromRaw(hashKeyValue, rangeKeyValue)
    if err != nil {
        return nil, fmt.Errorf("failed to create key for update: %v", err)
    }
   
    // Use helper to marshal raw updates
    marshaledUpdates, err := marshalRawUpdates(updates)
    if err != nil {
        return nil, err
    }
   
    // Use helper to build expression
    updateExpression, attrNames, attrValues := buildUpdateExpression(marshaledUpdates)
   
    return &dynamodb.UpdateItemInput{
        TableName:                 aws.String(TableSchema.TableName),
        Key:                       key,
        UpdateExpression:          aws.String(updateExpression),
        ExpressionAttributeNames:  attrNames,
        ExpressionAttributeValues: attrValues,
    }, nil
}

// UpdateItemInputWithCondition creates conditional update request
func UpdateItemInputWithCondition(hashKeyValue interface{}, rangeKeyValue interface{}, updates map[string]interface{}, conditionExpression string, conditionAttributeNames map[string]string, conditionAttributeValues map[string]types.AttributeValue) (*dynamodb.UpdateItemInput, error) {
    // All validations at the beginning
    if err := validateKeyInputs(hashKeyValue, rangeKeyValue); err != nil {
        return nil, err
    }
    if err := validateUpdatesMap(updates); err != nil {
        return nil, err
    }
    if err := validateConditionExpression(conditionExpression); err != nil {
        return nil, err
    }

    // Pure business logic after validation
    updateInput, err := UpdateItemInputFromRaw(hashKeyValue, rangeKeyValue, updates)
    if err != nil {
        return nil, err
    }
   
    updateInput.ConditionExpression = aws.String(conditionExpression)
   
    // Use helper to merge expression attributes
    updateInput.ExpressionAttributeNames, updateInput.ExpressionAttributeValues = mergeExpressionAttributes(
        updateInput.ExpressionAttributeNames,
        updateInput.ExpressionAttributeValues,
        conditionAttributeNames,
        conditionAttributeValues,
    )
   
    return updateInput, nil
}

// UpdateItemInputWithExpression creates update with expression builder
func UpdateItemInputWithExpression(hashKeyValue interface{}, rangeKeyValue interface{}, updateBuilder expression.UpdateBuilder, conditionBuilder *expression.ConditionBuilder) (*dynamodb.UpdateItemInput, error) {
    // All validations at the beginning
    if err := validateKeyInputs(hashKeyValue, rangeKeyValue); err != nil {
        return nil, err
    }
    
    // Pure business logic after validation
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
}
`
