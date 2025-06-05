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

func BoolToInt(b bool) int {
    if b {
        return 1
    }
    return 0
}

func IntToBool(i int) bool {
    return i != 0
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
        {{else if eq .Type "B"}}
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
            if v {
                result[key] = &types.AttributeValueMemberN{Value: "1"}
            } else {
                result[key] = &types.AttributeValueMemberN{Value: "0"}
            }
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
