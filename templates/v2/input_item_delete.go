package v2

// DeleteItemInputTemplate generates utility functions for creating DynamoDB DeleteItemInput.
// This template creates helper functions for deleting items using complete objects,
// raw values, conditional deletes, and batch operations.
const DeleteItemInputTemplate = `
// DeleteItemInput creates DeleteItemInput from a SchemaItem.
//
// Parameters:
//   - item: SchemaItem containing key values
//
// Returns:
//   - *dynamodb.DeleteItemInput: Ready for DeleteItem operation
//   - error: If key extraction fails
//
// Example:
//   item := SchemaItem{Id: "user123", Created: 1640995200}
//   deleteInput, err := DeleteItemInput(item)
//   _, err = client.DeleteItem(ctx, deleteInput)
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

// DeleteItemInputFromRaw creates DeleteItemInput from raw key values.
//
// Parameters:
//   - hashKeyValue: Hash key value (any type)
//   - rangeKeyValue: Range key value (any type, nil if no range key)
//
// Returns:
//   - *dynamodb.DeleteItemInput: Ready for DeleteItem operation
//   - error: If key creation fails
//
// Example:
//   deleteInput, err := DeleteItemInputFromRaw("user123", 1640995200)
//   _, err = client.DeleteItem(ctx, deleteInput)
func DeleteItemInputFromRaw(hashKeyValue interface{}, rangeKeyValue interface{}) (*dynamodb.DeleteItemInput, error) {
   key, err := KeyInputFromRaw(hashKeyValue, rangeKeyValue)
   if err != nil {
       return nil, fmt.Errorf("failed to create key for delete: %v", err)
   }
   
   return &dynamodb.DeleteItemInput{
       TableName: aws.String(TableSchema.TableName),
       Key:       key,
   }, nil
}

// DeleteItemInputWithCondition creates DeleteItemInput with conditional expression (optimistic locking).
//
// Parameters:
//   - hashKeyValue: Hash key value
//   - rangeKeyValue: Range key value (nil if no range key)
//   - conditionExpression: DynamoDB condition (e.g., "#version = :v")
//   - expressionAttributeNames: Attribute name mappings
//   - expressionAttributeValues: Attribute value mappings
//
// Returns:
//   - *dynamodb.DeleteItemInput: Ready for conditional DeleteItem
//   - error: If key creation fails
//
// Example:
//   deleteInput, err := DeleteItemInputWithCondition(
//       "user123", 1640995200,
//       "#version = :v",
//       map[string]string{"#version": "version"},
//       map[string]types.AttributeValue{":v": &types.AttributeValueMemberN{Value: "1"}},
//   )
//   _, err = client.DeleteItem(ctx, deleteInput)
func DeleteItemInputWithCondition(hashKeyValue interface{}, rangeKeyValue interface{}, conditionExpression string, expressionAttributeNames map[string]string, expressionAttributeValues map[string]types.AttributeValue) (*dynamodb.DeleteItemInput, error) {
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

// BatchDeleteItemsInput creates BatchWriteItemInput for deleting multiple items (max 25).
//
// Parameters:
//   - keys: Slice of DynamoDB keys to delete
//
// Returns:
//   - *dynamodb.BatchWriteItemInput: Ready for BatchWriteItem operation
//   - error: If batch size exceeds limit
//
// Example:
//   keys := []map[string]types.AttributeValue{
//       {"id": &types.AttributeValueMemberS{Value: "user1"}},
//       {"id": &types.AttributeValueMemberS{Value: "user2"}},
//   }
//   batchInput, err := BatchDeleteItemsInput(keys)
//   _, err = client.BatchWriteItem(ctx, batchInput)
func BatchDeleteItemsInput(keys []map[string]types.AttributeValue) (*dynamodb.BatchWriteItemInput, error) {
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

// BatchDeleteItemsInputFromRaw creates BatchWriteItemInput from SchemaItem slice (max 25).
//
// Parameters:
//   - items: Slice of SchemaItems to extract keys from
//
// Returns:
//   - *dynamodb.BatchWriteItemInput: Ready for BatchWriteItem operation
//   - error: If key extraction fails or batch size exceeds limit
//
// Example:
//   items := []SchemaItem{
//       {Id: "user1", Created: 123},
//       {Id: "user2", Created: 456},
//   }
//   batchInput, err := BatchDeleteItemsInputFromRaw(items)
//   _, err = client.BatchWriteItem(ctx, batchInput)
func BatchDeleteItemsInputFromRaw(items []SchemaItem) (*dynamodb.BatchWriteItemInput, error) {
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
