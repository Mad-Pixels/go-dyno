package inputs

// DeleteInputsTemplate ...
const DeleteInputsTemplate = `
// DeleteItemInput ...
func DeleteItemInput(item SchemaItem) (*dynamodb.DeleteItemInput, error) {
    key, err := KeyInput(item)
    if err != nil {
        return nil, fmt.Errorf("failed to create key from item for delete: %v", err)
    }
   
    return &dynamodb.DeleteItemInput{
        TableName: aws.String(TableSchema.TableName),
        Key:       key,
    }, nil
}

// DeleteItemInputFromRaw ...
func DeleteItemInputFromRaw(hashKeyValue any, rangeKeyValue any) (*dynamodb.DeleteItemInput, error) {
    // All validations at the beginning
    if err := validateKeyInputs(hashKeyValue, rangeKeyValue); err != nil {
        return nil, err
    }
    
    // Pure business logic after validation
    key, err := KeyInputFromRaw(hashKeyValue, rangeKeyValue)
    if err != nil {
        return nil, fmt.Errorf("failed to create key for delete: %v", err)
    }
   
    return &dynamodb.DeleteItemInput{
        TableName: aws.String(TableSchema.TableName),
        Key:       key,
    }, nil
}

// DeleteItemInputWithCondition ...
func DeleteItemInputWithCondition(hashKeyValue any, rangeKeyValue any, conditionExpression string, expressionAttributeNames map[string]string, expressionAttributeValues map[string]types.AttributeValue) (*dynamodb.DeleteItemInput, error) {
    // All validations at the beginning
    if err := validateKeyInputs(hashKeyValue, rangeKeyValue); err != nil {
        return nil, err
    }
    if err := validateConditionExpression(conditionExpression); err != nil {
        return nil, err
    }
    
    // Pure business logic after validation
    key, err := KeyInputFromRaw(hashKeyValue, rangeKeyValue)
    if err != nil {
        return nil, fmt.Errorf("failed to create key for conditional delete: %v", err)
    }
   
    input := &dynamodb.DeleteItemInput{
        TableName:           aws.String(TableSchema.TableName),
        Key:                 key,
        ConditionExpression: aws.String(conditionExpression),
    }
   
    if expressionAttributeNames != nil {
        input.ExpressionAttributeNames = expressionAttributeNames
    }
   
    if expressionAttributeValues != nil {
        input.ExpressionAttributeValues = expressionAttributeValues
    }
   
    return input, nil
}

// BatchDeleteItemsInput ...
func BatchDeleteItemsInput(keys []map[string]types.AttributeValue) (*dynamodb.BatchWriteItemInput, error) {
    // All validations at the beginning
    if err := validateBatchSize(len(keys), "delete"); err != nil {
        return nil, err
    }
   
    if len(keys) == 0 {
        return &dynamodb.BatchWriteItemInput{}, nil
    }
   
    // Pure business logic after validation
    writeRequests := make([]types.WriteRequest, 0, len(keys))
    for _, key := range keys {
        writeRequests = append(writeRequests, types.WriteRequest{
            DeleteRequest: &types.DeleteRequest{
                Key: key,
            },
        })
    }
   
    return &dynamodb.BatchWriteItemInput{
        RequestItems: map[string][]types.WriteRequest{
            TableSchema.TableName: writeRequests,
        },
    }, nil
}

// BatchDeleteItemsInputFromRaw ...
func BatchDeleteItemsInputFromRaw(items []SchemaItem) (*dynamodb.BatchWriteItemInput, error) {
    // All validations at the beginning
    if err := validateBatchSize(len(items), "delete"); err != nil {
        return nil, err
    }
   
    if len(items) == 0 {
        return &dynamodb.BatchWriteItemInput{}, nil
    }
   
    // Pure business logic after validation
    keys := make([]map[string]types.AttributeValue, 0, len(items))
    for _, item := range items {
        key, err := KeyInput(item)
        if err != nil {
            return nil, fmt.Errorf("failed to create key from item: %v", err)
        }
        keys = append(keys, key)
    }
   
    return BatchDeleteItemsInput(keys)
}
`
