package v2

// CrudOtherKey ...
const CrudOther = `
// CreateKey extracts the primary key from a SchemaItem and returns it as DynamoDB AttributeValue map.
// This is the primary method for creating keys from existing items.
//
// Example usage:
//   item := SchemaItem{Id: "user123", Created: 1640995200, Name: "John"}
//   key, err := CreateKey(item)
//   if err != nil {
//       return err
//   }
//   _, err = dynamoClient.GetItem(ctx, &dynamodb.GetItemInput{
//       TableName: aws.String(TableSchema.TableName),
//       Key: key,
//   })
func CreateKey(item SchemaItem) (map[string]types.AttributeValue, error) {
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

// CreateKeyFromRaw creates a key from raw hash and range key values.
// Use this when you have individual key values rather than a complete SchemaItem.
//
// Example usage:
//   key, err := CreateKeyFromRaw("user123", 1640995200)
//   if err != nil {
//       return err
//   }
//   _, err = dynamoClient.GetItem(ctx, &dynamodb.GetItemInput{
//       TableName: aws.String(TableSchema.TableName),
//       Key: key,
//   })
func CreateKeyFromRaw(hashKeyValue interface{}, rangeKeyValue interface{}) (map[string]types.AttributeValue, error) {
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

// IncrementAttribute creates an UpdateItemInput to increment/decrement a numeric attribute
// Useful for counters, views, likes, etc.
//
// Example usage:
//   // Increment views by 1
//   updateInput, err := IncrementAttribute("user123", 1640995200, "views", 1)
//   
//   // Decrement likes by 1
//   updateInput, err := IncrementAttribute("post456", 1640995200, "likes", -1)
//   
//   if err != nil {
//       return err
//   }
//   _, err = dynamoClient.UpdateItem(ctx, updateInput)
func IncrementAttribute(hashKeyValue interface{}, rangeKeyValue interface{}, attributeName string, incrementValue int) (*dynamodb.UpdateItemInput, error) {
    key, err := CreateKeyFromRaw(hashKeyValue, rangeKeyValue)
    if err != nil {
        return nil, fmt.Errorf("failed to create key for increment: %v", err)
    }
    
    return &dynamodb.UpdateItemInput{
        TableName:        aws.String(TableSchema.TableName),
        Key:              key,
        UpdateExpression: aws.String("ADD #attr :val"),
        ExpressionAttributeNames: map[string]string{
            "#attr": attributeName,
        },
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":val": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", incrementValue)},
        },
    }, nil
}

// AddToSet creates an UpdateItemInput to add values to a string set (SS) or number set (NS)
// Creates the set if it doesn't exist, otherwise adds to existing set
//
// Example usage:
//   // Add tags to string set
//   updateInput, err := AddToSet("user123", 1640995200, "tags", []string{"golang", "backend"})
//   
//   // Add scores to number set  
//   updateInput, err := AddToSet("user123", 1640995200, "scores", []int{95, 87})
//   
//   if err != nil {
//       return err
//   }
//   _, err = dynamoClient.UpdateItem(ctx, updateInput)
func AddToSet(hashKeyValue interface{}, rangeKeyValue interface{}, attributeName string, values interface{}) (*dynamodb.UpdateItemInput, error) {
    key, err := CreateKeyFromRaw(hashKeyValue, rangeKeyValue)
    if err != nil {
        return nil, fmt.Errorf("failed to create key for add to set: %v", err)
    }
    
    var attributeValue types.AttributeValue
    
    switch v := values.(type) {
    case []string:
        if len(v) == 0 {
            return nil, fmt.Errorf("cannot add empty string set")
        }
        attributeValue = &types.AttributeValueMemberSS{Value: v}
    case []int:
        if len(v) == 0 {
            return nil, fmt.Errorf("cannot add empty number set")
        }
        numberStrings := make([]string, len(v))
        for i, num := range v {
            numberStrings[i] = fmt.Sprintf("%d", num)
        }
        attributeValue = &types.AttributeValueMemberNS{Value: numberStrings}
    default:
        return nil, fmt.Errorf("unsupported type for set operation: %T, expected []string or []int", values)
    }
    
    return &dynamodb.UpdateItemInput{
        TableName:        aws.String(TableSchema.TableName),
        Key:              key,
        UpdateExpression: aws.String("ADD #attr :val"),
        ExpressionAttributeNames: map[string]string{
            "#attr": attributeName,
        },
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":val": attributeValue,
        },
    }, nil
}

// RemoveFromSet creates an UpdateItemInput to remove values from a string set (SS) or number set (NS)
//
// Example usage:
//   // Remove tags from string set
//   updateInput, err := RemoveFromSet("user123", 1640995200, "tags", []string{"deprecated"})
//   
//   // Remove scores from number set
//   updateInput, err := RemoveFromSet("user123", 1640995200, "scores", []int{60})
//   
//   if err != nil {
//       return err
//   }
//   _, err = dynamoClient.UpdateItem(ctx, updateInput)
func RemoveFromSet(hashKeyValue interface{}, rangeKeyValue interface{}, attributeName string, values interface{}) (*dynamodb.UpdateItemInput, error) {
    key, err := CreateKeyFromRaw(hashKeyValue, rangeKeyValue)
    if err != nil {
        return nil, fmt.Errorf("failed to create key for remove from set: %v", err)
    }
    
    var attributeValue types.AttributeValue
    
    switch v := values.(type) {
    case []string:
        if len(v) == 0 {
            return nil, fmt.Errorf("cannot remove empty string set")
        }
        attributeValue = &types.AttributeValueMemberSS{Value: v}
    case []int:
        if len(v) == 0 {
            return nil, fmt.Errorf("cannot remove empty number set")
        }
        numberStrings := make([]string, len(v))
        for i, num := range v {
            numberStrings[i] = fmt.Sprintf("%d", num)
        }
        attributeValue = &types.AttributeValueMemberNS{Value: numberStrings}
    default:
        return nil, fmt.Errorf("unsupported type for set operation: %T, expected []string or []int", values)
    }
    
    return &dynamodb.UpdateItemInput{
        TableName:        aws.String(TableSchema.TableName),
        Key:              key,
        UpdateExpression: aws.String("DELETE #attr :val"),
        ExpressionAttributeNames: map[string]string{
            "#attr": attributeName,
        },
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":val": attributeValue,
        },
    }, nil
}
`
