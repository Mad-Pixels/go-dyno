package v2

// CrudDeleteTemplate ...
const CrudDeleteTemplate = `
// DeleteItem creates a DeleteItemInput using an existing SchemaItem.
// This is the primary method for deleting items using complete objects.
//
// Example usage:
//   item := SchemaItem{Id: "user123", Created: 1640995200}
//   deleteInput, err := DeleteItem(item)
//   if err != nil {
//       return err
//   }
//   _, err = dynamoClient.DeleteItem(ctx, deleteInput)
func DeleteItem(item SchemaItem) (*dynamodb.DeleteItemInput, error) {
    key, err := CreateKey(item)
    if err != nil {
        return nil, fmt.Errorf("failed to create key from item for delete: %v", err)
    }
    
    return &dynamodb.DeleteItemInput{
        TableName: aws.String(TableSchema.TableName),
        Key:       key,
    }, nil
}

// DeleteItemFromRaw creates a DeleteItemInput for DynamoDB delete operation using raw key values.
// Use this when you have individual key values rather than a complete SchemaItem.
//
// Example usage:
//   deleteInput, err := DeleteItemFromRaw("user123", 1640995200)
//   if err != nil {
//       return err
//   }
//   _, err = dynamoClient.DeleteItem(ctx, deleteInput)
func DeleteItemFromRaw(hashKeyValue interface{}, rangeKeyValue interface{}) (*dynamodb.DeleteItemInput, error) {
    key, err := CreateKeyFromRaw(hashKeyValue, rangeKeyValue)
    if err != nil {
        return nil, fmt.Errorf("failed to create key for delete: %v", err)
    }
    
    return &dynamodb.DeleteItemInput{
        TableName: aws.String(TableSchema.TableName),
        Key:       key,
    }, nil
}

// DeleteItemWithCondition creates a DeleteItemInput with a condition expression
// Useful for conditional deletes (e.g., delete only if version matches)
//
// Example usage:
//   deleteInput, err := DeleteItemWithCondition(
//       "user123", 1640995200,
//       "#version = :v",
//       map[string]string{"#version": "version"},
//       map[string]types.AttributeValue{":v": &types.AttributeValueMemberN{Value: "1"}},
//   )
//   if err != nil {
//       return err
//   }
//   _, err = dynamoClient.DeleteItem(ctx, deleteInput)
func DeleteItemWithCondition(hashKeyValue interface{}, rangeKeyValue interface{}, conditionExpression string, expressionAttributeNames map[string]string, expressionAttributeValues map[string]types.AttributeValue) (*dynamodb.DeleteItemInput, error) {
    key, err := CreateKeyFromRaw(hashKeyValue, rangeKeyValue)
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

// BatchDeleteItems creates BatchWriteItemInput for deleting multiple items
// Takes slice of key pairs (hash, range) and creates batch delete request
// Maximum 25 items per batch (DynamoDB limitation)
//
// Example usage:
//   keys := []map[string]types.AttributeValue{
//       {"id": &types.AttributeValueMemberS{Value: "user1"}, "created": &types.AttributeValueMemberN{Value: "123"}},
//       {"id": &types.AttributeValueMemberS{Value: "user2"}, "created": &types.AttributeValueMemberN{Value: "456"}},
//   }
//   batchInput, err := BatchDeleteItems(keys)
//   if err != nil {
//       return err
//   }
//   _, err = dynamoClient.BatchWriteItem(ctx, batchInput)
func BatchDeleteItems(keys []map[string]types.AttributeValue) (*dynamodb.BatchWriteItemInput, error) {
    if len(keys) == 0 {
        return &dynamodb.BatchWriteItemInput{}, nil
    }
    
    if len(keys) > 25 {
        return nil, fmt.Errorf("batch delete supports maximum 25 items, got %d", len(keys))
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

// BatchDeleteItemsFromItems creates BatchWriteItemInput for deleting items by SchemaItem slice
// Extracts keys from items and creates batch delete request
// Maximum 25 items per batch (DynamoDB limitation)
//
// Example usage:
//   items := []SchemaItem{
//       {Id: "user1", Created: 123, Name: "John"},
//       {Id: "user2", Created: 456, Name: "Jane"},
//   }
//   batchInput, err := BatchDeleteItemsFromItems(items)
//   if err != nil {
//       return err
//   }
//   _, err = dynamoClient.BatchWriteItem(ctx, batchInput)
func BatchDeleteItemsFromItems(items []SchemaItem) (*dynamodb.BatchWriteItemInput, error) {
    if len(items) == 0 {
        return &dynamodb.BatchWriteItemInput{}, nil
    }
    
    keys := make([]map[string]types.AttributeValue, 0, len(items))
    for _, item := range items {
        key, err := CreateKey(item)
        if err != nil {
            return nil, fmt.Errorf("failed to create key from item: %v", err)
        }
        keys = append(keys, key)
    }
    
    return BatchDeleteItems(keys)
}
`
