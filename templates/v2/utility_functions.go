package v2

// UtilityFunctionsTemplate ...
const UtilityFunctionsTemplate = `
func BatchPutItems(items []SchemaItem) ([]map[string]types.AttributeValue, error) {
    result := make([]map[string]types.AttributeValue, 0, len(items))
    for _, item := range items {
        av, err := PutItem(item)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal item: %v", err)
        }
        result = append(result, av)
    }
    return result, nil
}

// PutItem creates an AttributeValues map for PutItem in DynamoDB
func PutItem(item SchemaItem) (map[string]types.AttributeValue, error) {
    attributeValues, err := attributevalue.MarshalMap(item)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal item: %v", err)
    }
    return attributeValues, nil
}

// ExtractFromDynamoDBStreamEvent DynamoDB Stream to SchemaItem
func ExtractFromDynamoDBStreamEvent(dbEvent events.DynamoDBEventRecord) (*SchemaItem, error) {
    if dbEvent.Change.NewImage == nil {
        return nil, fmt.Errorf("new image is nil in the event")
    }
    
    item := &SchemaItem{}
    
    {{range .AllAttributes}}
    if val, ok := dbEvent.Change.NewImage["{{.Name}}"]; ok {
        {{if eq .Type "S"}}
        item.{{ToSafeName .Name | ToUpperCamelCase}} = val.String()
        {{else if eq .Type "N"}}
        if n, err := strconv.Atoi(val.Number()); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = n
        }
        {{else if eq .Type "BOOL"}}
        item.{{ToSafeName .Name | ToUpperCamelCase}} = val.Boolean()
        {{else if eq .Type "SS"}}
        if ss := val.StringSet(); ss != nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = ss
        }
        {{else if eq .Type "NS"}}
        if ns := val.NumberSet(); ns != nil {
            numbers := make([]int, 0, len(ns))
            for _, numStr := range ns {
                if num, err := strconv.Atoi(numStr); err == nil {
                    numbers = append(numbers, num)
                }
            }
            item.{{ToSafeName .Name | ToUpperCamelCase}} = numbers
        }
        {{else}}
        // Unsupported type: {{.Type}} for attribute {{.Name}}
        _ = val // Mark as used to avoid compilation error
        {{end}}
    }
    {{end}}
    
    return item, nil
}

// IsFieldModified DynamoDB Stream to check if field was modified
func IsFieldModified(dbEvent events.DynamoDBEventRecord, fieldName string) bool {
    if dbEvent.EventName != "MODIFY" {
        return false
    }
    
    if dbEvent.Change.OldImage == nil || dbEvent.Change.NewImage == nil {
        return false
    }
    
    oldVal, oldExists := dbEvent.Change.OldImage[fieldName]
    newVal, newExists := dbEvent.Change.NewImage[fieldName]
    
    if !oldExists || !newExists {
        return false
    }
    
    return oldVal.String() != newVal.String()
}

func GetBoolFieldChanged(dbEvent events.DynamoDBEventRecord, fieldName string) bool {
    if dbEvent.EventName != "MODIFY" {
        return false
    }
    
    if dbEvent.Change.OldImage == nil || dbEvent.Change.NewImage == nil {
        return false
    }
    
    oldValue := false
    if oldVal, ok := dbEvent.Change.OldImage[fieldName]; ok {
        oldValue = oldVal.Boolean()
    }
    
    newValue := false
    if newVal, ok := dbEvent.Change.NewImage[fieldName]; ok {
        newValue = newVal.Boolean()
    }
    
    return !oldValue && newValue
}

func CreateKey(hashKeyValue interface{}, rangeKeyValue interface{}) (map[string]types.AttributeValue, error) {
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

func CreateKeyFromItem(item SchemaItem) (map[string]types.AttributeValue, error) {
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

// DeleteItem creates a DeleteItemInput for DynamoDB delete operation
// Returns a configured DeleteItemInput ready for client.DeleteItem() call
//
// Example usage:
//   deleteInput, err := DeleteItem("user123", 1640995200)
//   if err != nil {
//       return err
//   }
//   _, err = dynamoClient.DeleteItem(ctx, deleteInput)
func DeleteItem(hashKeyValue interface{}, rangeKeyValue interface{}) (*dynamodb.DeleteItemInput, error) {
    key, err := CreateKey(hashKeyValue, rangeKeyValue)
    if err != nil {
        return nil, fmt.Errorf("failed to create key for delete: %v", err)
    }
    
    return &dynamodb.DeleteItemInput{
        TableName: aws.String(TableSchema.TableName),
        Key:       key,
    }, nil
}

// DeleteItemFromItem creates a DeleteItemInput using an existing SchemaItem
// Extracts the key from the item and creates delete input
//
// Example usage:
//   item := SchemaItem{Id: "user123", Created: 1640995200}
//   deleteInput, err := DeleteItemFromItem(item)
//   if err != nil {
//       return err
//   }
//   _, err = dynamoClient.DeleteItem(ctx, deleteInput)
func DeleteItemFromItem(item SchemaItem) (*dynamodb.DeleteItemInput, error) {
    key, err := CreateKeyFromItem(item)
    if err != nil {
        return nil, fmt.Errorf("failed to create key from item for delete: %v", err)
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
    key, err := CreateKey(hashKeyValue, rangeKeyValue)
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
        key, err := CreateKeyFromItem(item)
        if err != nil {
            return nil, fmt.Errorf("failed to create key from item: %v", err)
        }
        keys = append(keys, key)
    }
    
    return BatchDeleteItems(keys)
}

func CreateTriggerHandler(
    onInsert func(context.Context, *SchemaItem) error,
    onModify func(context.Context, *SchemaItem, *SchemaItem) error,
    onDelete func(context.Context, map[string]events.DynamoDBAttributeValue) error,
) func(ctx context.Context, event events.DynamoDBEvent) error {
    return func(ctx context.Context, event events.DynamoDBEvent) error {
        for _, record := range event.Records {
            switch record.EventName {
            case "INSERT":
                if onInsert != nil {
                    item, err := ExtractFromDynamoDBStreamEvent(record)
                    if err != nil {
                        return err
                    }
                    if err := onInsert(ctx, item); err != nil {
                        return err
                    }
                }
                
            case "MODIFY":
                if onModify != nil {
                    oldItem, err := ExtractFromDynamoDBStreamEvent(events.DynamoDBEventRecord{
                        Change: events.DynamoDBStreamRecord{
                            NewImage: record.Change.OldImage,
                        },
                    })
                    if err != nil {
                        return err
                    }
                    
                    newItem, err := ExtractFromDynamoDBStreamEvent(record)
                    if err != nil {
                        return err
                    }
                    
                    if err := onModify(ctx, oldItem, newItem); err != nil {
                        return err
                    }
                }
                
            case "REMOVE":
                if onDelete != nil {
                    if err := onDelete(ctx, record.Change.OldImage); err != nil {
                        return err
                    }
                }
            }
        }
        return nil
    }
}

func ConvertMapToAttributeValues(input map[string]interface{}) (map[string]types.AttributeValue, error) {
    result := make(map[string]types.AttributeValue)
    
    for key, value := range input {
        switch v := value.(type) {
        case string:
            result[key] = &types.AttributeValueMemberS{Value: v}
        case []string:
            result[key] = &types.AttributeValueMemberSS{Value: v}
        case []int:
            numbers := make([]string, len(v))
            for i, num := range v {
                numbers[i] = fmt.Sprintf("%d", num)
            }
            result[key] = &types.AttributeValueMemberNS{Value: numbers}
        case int:
            result[key] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", v)}
        case float64:
            if v == float64(int64(v)) {
                result[key] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", int64(v))}
            } else {
                result[key] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%g", v)}
            }
        case bool:
            result[key] = &types.AttributeValueMemberBOOL{Value: v}
        case nil:
            result[key] = &types.AttributeValueMemberNULL{Value: true}
        case map[string]interface{}:
            b, err := json.Marshal(v)
            if err != nil {
                return nil, err
            }
            result[key] = &types.AttributeValueMemberM{
                Value: map[string]types.AttributeValue{
                    "json": &types.AttributeValueMemberS{Value: string(b)},
                },
            }
        case []interface{}:
            b, err := json.Marshal(v)
            if err != nil {
                return nil, err
            }
            result[key] = &types.AttributeValueMemberL{
                Value: []types.AttributeValue{
                    &types.AttributeValueMemberS{Value: string(b)},
                },
            }
        default:
            return nil, fmt.Errorf("unsupported type for key %s: %T", key, value)
        }
    }
    
    return result, nil
}`
