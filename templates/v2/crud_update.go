package v2

// CrudUpdateTemplate ...
const CrudUpdateTemplate = `
// UpdateItem creates an UpdateItemInput using an existing SchemaItem.
// Updates all non-key attributes from the provided item.
// This is the primary method for updating items using complete objects.
//
// Example usage:
//   item := SchemaItem{
//       Id: "user123", Created: 1640995200,
//       Name: "Updated Name", Age: 30,
//   }
//   updateInput, err := UpdateItem(item)
//   if err != nil {
//       return err
//   }
//   _, err = dynamoClient.UpdateItem(ctx, updateInput)
func UpdateItem(item SchemaItem) (*dynamodb.UpdateItemInput, error) {
    key, err := CreateKey(item)
    if err != nil {
        return nil, fmt.Errorf("failed to create key from item for update: %v", err)
    }
    
    // Marshal the entire item to get all attributes
    allAttributes, err := attributevalue.MarshalMap(item)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal item for update: %v", err)
    }
    
    // Remove key attributes from update
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

// UpdateItemFromRaw creates an UpdateItemInput for DynamoDB update operation using raw values.
// Uses SET action to update specific attributes while preserving others.
// Use this when you have individual key values and update map rather than a complete SchemaItem.
//
// Example usage:
//   updates := map[string]interface{}{
//       "name": "Updated Name",
//       "age": 30,
//       "is_active": true,
//   }
//   updateInput, err := UpdateItemFromRaw("user123", 1640995200, updates)
//   if err != nil {
//       return err
//   }
//   _, err = dynamoClient.UpdateItem(ctx, updateInput)
func UpdateItemFromRaw(hashKeyValue interface{}, rangeKeyValue interface{}, updates map[string]interface{}) (*dynamodb.UpdateItemInput, error) {
    key, err := CreateKeyFromRaw(hashKeyValue, rangeKeyValue)
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

// UpdateItemWithCondition creates an UpdateItemInput with a condition expression
// Useful for conditional updates (e.g., update only if version matches, optimistic locking)
//
// Example usage:
//   updates := map[string]interface{}{
//       "name": "New Name",
//       "version": 2,
//   }
//   updateInput, err := UpdateItemWithCondition(
//       "user123", 1640995200,
//       updates,
//       "#version = :currentVersion",
//       map[string]string{"#version": "version"},
//       map[string]types.AttributeValue{":currentVersion": &types.AttributeValueMemberN{Value: "1"}},
//   )
//   if err != nil {
//       return err
//   }
//   _, err = dynamoClient.UpdateItem(ctx, updateInput)
func UpdateItemWithCondition(hashKeyValue interface{}, rangeKeyValue interface{}, updates map[string]interface{}, conditionExpression string, conditionAttributeNames map[string]string, conditionAttributeValues map[string]types.AttributeValue) (*dynamodb.UpdateItemInput, error) {
    updateInput, err := UpdateItemFromRaw(hashKeyValue, rangeKeyValue, updates)
    if err != nil {
        return nil, err
    }
    
    updateInput.ConditionExpression = aws.String(conditionExpression)
    
    // Merge condition attribute names with update attribute names
    if conditionAttributeNames != nil {
        for key, value := range conditionAttributeNames {
            updateInput.ExpressionAttributeNames[key] = value
        }
    }
    
    // Merge condition attribute values with update attribute values
    if conditionAttributeValues != nil {
        for key, value := range conditionAttributeValues {
            updateInput.ExpressionAttributeValues[key] = value
        }
    }
    
    return updateInput, nil
}

// UpdateItemWithExpression creates an UpdateItemInput using DynamoDB expression builder
// Provides full control over update expressions with type safety
//
// Example usage:
//   import "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
//   
//   update := expression.Set(expression.Name("name"), expression.Value("New Name")).
//             Add(expression.Name("views"), expression.Value(1))
//   condition := expression.Equal(expression.Name("version"), expression.Value(1))
//   
//   updateInput, err := UpdateItemWithExpression(
//       "user123", 1640995200,
//       update, &condition,
//   )
//   if err != nil {
//       return err
//   }
//   _, err = dynamoClient.UpdateItem(ctx, updateInput)
func UpdateItemWithExpression(hashKeyValue interface{}, rangeKeyValue interface{}, updateBuilder expression.UpdateBuilder, conditionBuilder *expression.ConditionBuilder) (*dynamodb.UpdateItemInput, error) {
    key, err := CreateKeyFromRaw(hashKeyValue, rangeKeyValue)
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
