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

// UpdateItem creates an UpdateItemInput for DynamoDB update operation
// Uses SET action to update specific attributes while preserving others
//
// Example usage:
//   updates := map[string]interface{}{
//       "name": "Updated Name",
//       "age": 30,
//       "is_active": true,
//   }
//   updateInput, err := UpdateItem("user123", 1640995200, updates)
//   if err != nil {
//       return err
//   }
//   _, err = dynamoClient.UpdateItem(ctx, updateInput)
func UpdateItem(hashKeyValue interface{}, rangeKeyValue interface{}, updates map[string]interface{}) (*dynamodb.UpdateItemInput, error) {
    key, err := CreateKey(hashKeyValue, rangeKeyValue)
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

// UpdateItemFromItem creates an UpdateItemInput using an existing SchemaItem
// Updates all non-key attributes from the provided item
//
// Example usage:
//   item := SchemaItem{
//       Id: "user123", Created: 1640995200,
//       Name: "Updated Name", Age: 30,
//   }
//   updateInput, err := UpdateItemFromItem(item)
//   if err != nil {
//       return err
//   }
//   _, err = dynamoClient.UpdateItem(ctx, updateInput)
func UpdateItemFromItem(item SchemaItem) (*dynamodb.UpdateItemInput, error) {
    key, err := CreateKeyFromItem(item)
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
    
    i := 0
    for attrName, attrValue := range updates {
        nameKey := fmt.Sprintf("#attr%d", i)
        valueKey := fmt.Sprintf(":val%d", i)
        
        updateExpressionParts = append(updateExpressionParts, fmt.Sprintf("%s = %s", nameKey, valueKey))
        expressionAttributeNames[nameKey] = attrName
        updates[valueKey] = attrValue
        delete(updates, attrName)
        i++
    }
    
    updateExpression := "SET " + strings.Join(updateExpressionParts, ", ")
    
    return &dynamodb.UpdateItemInput{
        TableName:                 aws.String(TableSchema.TableName),
        Key:                       key,
        UpdateExpression:          aws.String(updateExpression),
        ExpressionAttributeNames:  expressionAttributeNames,
        ExpressionAttributeValues: updates,
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
    updateInput, err := UpdateItem(hashKeyValue, rangeKeyValue, updates)
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
    key, err := CreateKey(hashKeyValue, rangeKeyValue)
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
    key, err := CreateKey(hashKeyValue, rangeKeyValue)
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
    key, err := CreateKey(hashKeyValue, rangeKeyValue)
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
    key, err := CreateKey(hashKeyValue, rangeKeyValue)
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
