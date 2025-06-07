package v2

// ItemInputTemplate generates utility functions for DynamoDB item operations.
// This template creates type-safe helper functions for converting Go structs
// to DynamoDB AttributeValue maps and batch processing.
const ItemInputTemplate = `
// ItemInput converts a SchemaItem to DynamoDB AttributeValue map for PutItem operations.
//
// Parameters:
//   - item: SchemaItem struct with table data
//
// Returns:
//   - map[string]types.AttributeValue: Ready for DynamoDB PutItem
//   - error: If marshaling fails
//
// Example:
//   item := SchemaItem{Id: "user123", Name: "John"}
//   av, err := ItemInput(item)
//   _, err = client.PutItem(ctx, &dynamodb.PutItemInput{
//       TableName: aws.String(TableName),
//       Item:      av,
//   })
func ItemInput(item SchemaItem) (map[string]types.AttributeValue, error) {
  attributeValues, err := attributevalue.MarshalMap(item)
  if err != nil {
      return nil, fmt.Errorf("failed to marshal item: %v", err)
  }
  return attributeValues, nil
}

// ItemsInput converts multiple SchemaItems for batch operations (max 25 items).
//
// Parameters:
//   - items: Slice of SchemaItem structs
//
// Returns:
//   - []map[string]types.AttributeValue: Ready for BatchWriteItem
//   - error: If any item marshaling fails
//
// Example:
//   items := []SchemaItem{ {Id: "user123", Name: "John"}, {Id: "user321", Name: "Kate"} }
//   batchItems, err := ItemsInput(items)
//   writeRequests := make([]types.WriteRequest, len(batchItems))
//   for i, item := range batchItems {
//       writeRequests[i] = types.WriteRequest{
//           PutRequest: &types.PutRequest{Item: item},
//       }
//   }
//   _, err = client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
//       RequestItems: map[string][]types.WriteRequest{TableName: writeRequests},
//   })
func ItemsInput(items []SchemaItem) ([]map[string]types.AttributeValue, error) {
  result := make([]map[string]types.AttributeValue, 0, len(items))
  for _, item := range items {
      av, err := ItemInput(item)
      if err != nil {
          return nil, fmt.Errorf("failed to marshal item: %v", err)
      }
      result = append(result, av)
  }
  return result, nil
}
`
