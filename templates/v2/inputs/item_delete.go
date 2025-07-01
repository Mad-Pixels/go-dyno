package inputs

// DeleteInputsTemplate provides input builders for DynamoDB delete operations
const DeleteInputsTemplate = `
// DeleteItemInput creates a DeleteItemInput from a complete SchemaItem.
// Extracts the primary key from the item for the delete operation.
// Use when you have the full item and want to delete it.
// Example: input, err := DeleteItemInput(userItem)
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

// DeleteItemInputFromRaw creates a DeleteItemInput from raw key values.
// Use when you only have the key values and want to delete the item.
// More efficient than DeleteItemInput when you don't have the full item.
// Example: input, err := DeleteItemInputFromRaw("user123", "session456")
func DeleteItemInputFromRaw(hashKeyValue any, rangeKeyValue any) (*dynamodb.DeleteItemInput, error) {
    if err := validateKeyInputs(hashKeyValue, rangeKeyValue); err != nil {
        return nil, err
    }
    
    key, err := KeyInputFromRaw(hashKeyValue, rangeKeyValue)
    if err != nil {
        return nil, fmt.Errorf("failed to create key for delete: %v", err)
    }
    return &dynamodb.DeleteItemInput{
        TableName: aws.String(TableSchema.TableName),
        Key:       key,
    }, nil
}

// DeleteItemInputWithCondition creates a conditional DeleteItemInput.
// Deletes the item only if the condition expression evaluates to true.
// Prevents accidental deletion and enables optimistic locking patterns.
// Example: DeleteItemInputWithCondition("user123", nil, "attribute_exists(#status)", {"#status": "status"}, nil)
func DeleteItemInputWithCondition(hashKeyValue any, rangeKeyValue any, conditionExpression string, expressionAttributeNames map[string]string, expressionAttributeValues map[string]types.AttributeValue) (*dynamodb.DeleteItemInput, error) {
    if err := validateKeyInputs(hashKeyValue, rangeKeyValue); err != nil {
        return nil, err
    }
    if err := validateConditionExpression(conditionExpression); err != nil {
        return nil, err
    }
    
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

// BatchDeleteItemsInput creates a BatchWriteItemInput for deleting multiple items.
// Takes pre-built key maps and creates delete requests for batch operation.
// Limited to 25 items per batch due to DynamoDB constraints.
// Example: BatchDeleteItemsInput([]map[string]types.AttributeValue{key1, key2})
func BatchDeleteItemsInput(keys []map[string]types.AttributeValue) (*dynamodb.BatchWriteItemInput, error) {
    if err := validateBatchSize(len(keys), "delete"); err != nil {
        return nil, err
    }
    if len(keys) == 0 {
        return &dynamodb.BatchWriteItemInput{}, nil
    }
   
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

// BatchDeleteItemsInputFromRaw creates a BatchWriteItemInput from SchemaItems.
// Extracts keys from each item and creates batch delete requests.
// More convenient than BatchDeleteItemsInput when you have full items.
// Example: BatchDeleteItemsInputFromRaw([]SchemaItem{item1, item2, item3})
func BatchDeleteItemsInputFromRaw(items []SchemaItem) (*dynamodb.BatchWriteItemInput, error) {
    if err := validateBatchSize(len(items), "delete"); err != nil {
        return nil, err
    }
    if len(items) == 0 {
        return &dynamodb.BatchWriteItemInput{}, nil
    }
   
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
