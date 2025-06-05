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
        {{$goType := ToGolangBaseType .}}
        {{if eq $goType "string"}}
        item.{{ToSafeName .Name | ToUpperCamelCase}} = val.String()
        {{else if eq $goType "time.Time"}}
        if t, err := time.Parse(time.RFC3339, val.String()); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = t
        }
        {{else if eq $goType "uuid.UUID"}}
        if u, err := uuid.Parse(val.String()); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = u
        }
        {{else}}
        item.{{ToSafeName .Name | ToUpperCamelCase}} = val.String()
        {{end}}
        {{else if eq .Type "N"}}
        {{$goType := ToGolangBaseType .}}
        {{if eq $goType "int"}}
        if n, err := strconv.Atoi(val.Number()); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = n
        }
        {{else if eq $goType "int8"}}
        if n, err := strconv.ParseInt(val.Number(), 10, 8); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = int8(n)
        }
        {{else if eq $goType "int16"}}
        if n, err := strconv.ParseInt(val.Number(), 10, 16); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = int16(n)
        }
        {{else if eq $goType "int32"}}
        if n, err := strconv.ParseInt(val.Number(), 10, 32); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = int32(n)
        }
        {{else if eq $goType "int64"}}
        if n, err := strconv.ParseInt(val.Number(), 10, 64); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = n
        }
        {{else if eq $goType "uint"}}
        if n, err := strconv.ParseUint(val.Number(), 10, 0); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = uint(n)
        }
        {{else if eq $goType "uint8"}}
        if n, err := strconv.ParseUint(val.Number(), 10, 8); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = uint8(n)
        }
        {{else if eq $goType "uint16"}}
        if n, err := strconv.ParseUint(val.Number(), 10, 16); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = uint16(n)
        }
        {{else if eq $goType "uint32"}}
        if n, err := strconv.ParseUint(val.Number(), 10, 32); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = uint32(n)
        }
        {{else if eq $goType "uint64"}}
        if n, err := strconv.ParseUint(val.Number(), 10, 64); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = n
        }
        {{else if eq $goType "float32"}}
        if f, err := strconv.ParseFloat(val.Number(), 32); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = float32(f)
        }
        {{else if eq $goType "float64"}}
        if f, err := strconv.ParseFloat(val.Number(), 64); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = f
        }
        {{else if eq $goType "*big.Int"}}
        if bigInt, ok := new(big.Int).SetString(val.Number(), 10); ok {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = bigInt
        }
        {{else if eq $goType "*decimal.Decimal"}}
        if dec, err := decimal.NewFromString(val.Number()); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = &dec
        }
        {{else}}
        // Default fallback to float64 for unknown numeric types
        if f, err := strconv.ParseFloat(val.Number(), 64); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = f
        }
        {{end}}
        {{else if eq .Type "B"}}
        item.{{ToSafeName .Name | ToUpperCamelCase}} = val.Boolean()
        {{else if eq .Type "BS"}}
        {{$goType := ToGolangBaseType .}}
        {{if eq $goType "[]byte"}}
        if data, err := base64.StdEncoding.DecodeString(val.String()); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = data
        }
        {{else}}
        item.{{ToSafeName .Name | ToUpperCamelCase}} = val.String()
        {{end}}
        {{else if or (eq .Type "SS") (eq .Type "NS") (eq .Type "L") (eq .Type "M")}}
        // Complex types (String Set, Number Set, List, Map) - store as JSON string
        if jsonData, err := json.Marshal(val); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = string(jsonData)
        }
        {{else if eq .Type "NULL"}}
        // NULL type - leave field as zero value
        {{else}}
        // Unknown type - store as string representation
        item.{{ToSafeName .Name | ToUpperCamelCase}} = val.String()
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
