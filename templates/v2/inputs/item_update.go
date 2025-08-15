package inputs

// UpdateInputsTemplate provides input builders for DynamoDB update operations
const UpdateInputsTemplate = `
// UpdateItemInput creates an UpdateItemInput from a complete SchemaItem.
// Automatically extracts the key and updates all non-key attributes.
// Use when you want to update an entire item with new values.
func UpdateItemInput(item SchemaItem) (*dynamodb.UpdateItemInput, error) {
    key, err := KeyInput(item)
    if err != nil {
        return nil, fmt.Errorf("failed to create key from item for update: %v", err)
    }
    allAttributes, err := marshalItemToMap(item)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal item for update: %v", err)
    }
    updates := extractNonKeyAttributes(allAttributes)
    if len(updates) == 0 {
        return nil, fmt.Errorf("no non-key attributes to update")
    }
    updateExpression, attrNames, attrValues := buildUpdateExpression(updates)
   
    return &dynamodb.UpdateItemInput{
        TableName:                 aws.String(TableSchema.TableName),
        Key:                       key,
        UpdateExpression:          aws.String(updateExpression),
        ExpressionAttributeNames:  attrNames,
        ExpressionAttributeValues: attrValues,
    }, nil
}

// UpdateItemInputFromRaw creates an UpdateItemInput from raw key values and update map.
// More efficient for partial updates when you only want to modify specific attributes.
// Use when you know exactly which fields to update without loading the full item.
func UpdateItemInputFromRaw(hashKeyValue any, rangeKeyValue any, updates map[string]any) (*dynamodb.UpdateItemInput, error) {
    if err := validateKeyInputs(hashKeyValue, rangeKeyValue); err != nil {
        return nil, err
    }
    if err := validateUpdatesMap(updates); err != nil {
        return nil, err
    }
    key, err := KeyInputFromRaw(hashKeyValue, rangeKeyValue)
    if err != nil {
        return nil, fmt.Errorf("failed to create key for update: %v", err)
    }
    marshaledUpdates, err := marshalUpdatesWithSchema(updates)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal updates: %v", err)
    }
    updateExpression, attrNames, attrValues := buildUpdateExpression(marshaledUpdates)
   
    return &dynamodb.UpdateItemInput{
        TableName:                 aws.String(TableSchema.TableName),
        Key:                       key,
        UpdateExpression:          aws.String(updateExpression),
        ExpressionAttributeNames:  attrNames,
        ExpressionAttributeValues: attrValues,
    }, nil
}

// UpdateItemInputWithCondition creates a conditional UpdateItemInput.
// Updates the item only if the condition expression evaluates to true.
func UpdateItemInputWithCondition(
    hashKeyValue any, 
    rangeKeyValue any, 
    updates map[string]any, 
    conditionExpression string, 
    conditionAttributeNames map[string]string, 
    conditionAttributeValues map[string]types.AttributeValue,
) (*dynamodb.UpdateItemInput, error) {
    if err := validateKeyInputs(hashKeyValue, rangeKeyValue); err != nil {
        return nil, err
    }
    if err := validateUpdatesMap(updates); err != nil {
        return nil, err
    }
    if err := validateConditionExpression(conditionExpression); err != nil {
        return nil, err
    }
    updateInput, err := UpdateItemInputFromRaw(hashKeyValue, rangeKeyValue, updates)
    if err != nil {
        return nil, err
    }
    updateInput.ConditionExpression = aws.String(conditionExpression)
   
    updateInput.ExpressionAttributeNames, updateInput.ExpressionAttributeValues = mergeExpressionAttributes(
        updateInput.ExpressionAttributeNames,
        updateInput.ExpressionAttributeValues,
        conditionAttributeNames,
        conditionAttributeValues,
    )
    return updateInput, nil
}

// UpdateItemInputWithExpression creates an UpdateItemInput using DynamoDB expression builders.
// Provides maximum flexibility for complex update operations (SET, ADD, REMOVE, DELETE).
// Use for advanced scenarios like atomic increments, list operations, or complex conditions.
// Example: 
//   updateExpr := expression.Set(expression.Name("counter"), expression.Name("counter").Plus(expression.Value(1)))
//   condExpr := expression.Name("version").Equal(expression.Value(currentVersion))
//   input, err := UpdateItemInputWithExpression("user123", nil, updateExpr, &condExpr)
func UpdateItemInputWithExpression(hashKeyValue any, rangeKeyValue any, updateBuilder expression.UpdateBuilder, conditionBuilder *expression.ConditionBuilder) (*dynamodb.UpdateItemInput, error) {
    if err := validateKeyInputs(hashKeyValue, rangeKeyValue); err != nil {
        return nil, err
    }
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
