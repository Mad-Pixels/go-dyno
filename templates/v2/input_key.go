package v2

// KeyInputTemplate generates utility functions for creating DynamoDB primary keys.
// This template creates helper functions for extracting keys from items or raw values.
const KeyInputTemplate = `
// KeyInput extracts primary key attributes from a SchemaItem for DynamoDB operations.
//
// Parameters:
//   - item: SchemaItem containing hash and range key values
//
// Returns:
//   - map[string]types.AttributeValue: Primary key ready for GetItem/DeleteItem
//   - error: If key marshaling fails
//
// Example:
//   item := SchemaItem{Id: "user123", Created: 1640995200}
//   key, err := KeyInput(item)
//   _, err = client.GetItem(ctx, &dynamodb.GetItemInput{
//       TableName: aws.String(TableName),
//       Key: key,
//   })
func KeyInput(item SchemaItem) (map[string]types.AttributeValue, error) {
   key := make(map[string]types.AttributeValue)
   
   var hashKeyValue interface{}
   {{range .AllAttributes}}{{if eq .Name $.HashKey}}
   hashKeyValue = item.{{ToSafeName .Name | ToUpperCamelCase}}
   {{end}}{{end}}
   
   hashKeyAV, err := attributevalue.Marshal(hashKeyValue)
   if err != nil {
       return nil, fmt.Errorf("failed to marshal hash key: %v", err)
   }
   key[TableSchema.HashKey] = hashKeyAV
   
   if TableSchema.RangeKey != "" {
       var rangeKeyValue interface{}
       {{range .AllAttributes}}{{if eq .Name $.RangeKey}}
       rangeKeyValue = item.{{ToSafeName .Name | ToUpperCamelCase}}
       {{end}}{{end}}
       
       rangeKeyAV, err := attributevalue.Marshal(rangeKeyValue)
       if err != nil {
           return nil, fmt.Errorf("failed to marshal range key: %v", err)
       }
       key[TableSchema.RangeKey] = rangeKeyAV
   }
   
   return key, nil
}

// KeyInputFromRaw creates primary key from individual hash/range values.
//
// Parameters:
//   - hashKeyValue: Hash key value (any type)
//   - rangeKeyValue: Range key value (any type, nil if no range key)
//
// Returns:
//   - map[string]types.AttributeValue: Primary key ready for DynamoDB operations
//   - error: If marshaling fails
//
// Example:
//   key, err := KeyInputFromRaw("user123", 1640995200)
//   _, err = client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
//       TableName: aws.String(TableName),
//       Key: key,
//   })
func KeyInputFromRaw(hashKeyValue interface{}, rangeKeyValue interface{}) (map[string]types.AttributeValue, error) {
   key := make(map[string]types.AttributeValue)
   
   hashKeyAV, err := attributevalue.Marshal(hashKeyValue)
   if err != nil {
       return nil, fmt.Errorf("failed to marshal hash key: %v", err)
   }
   key[TableSchema.HashKey] = hashKeyAV
   
   if TableSchema.RangeKey != "" && rangeKeyValue != nil {
       rangeKeyAV, err := attributevalue.Marshal(rangeKeyValue)
       if err != nil {
           return nil, fmt.Errorf("failed to marshal range key: %v", err)
       }
       key[TableSchema.RangeKey] = rangeKeyAV
   }
   
   return key, nil
}
`
