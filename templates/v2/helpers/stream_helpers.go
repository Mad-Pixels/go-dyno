package helpers

// StreamHelpersTemplate ...
const StreamHelpersTemplate = `
// ExtractFromDynamoDBStreamEvent ...
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
        {{$goType := ToGolangBaseType .}}
        {{if eq $goType "int"}}
        if n, err := strconv.Atoi(val.Number()); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = n
        }
        {{else if eq $goType "int64"}}
        if n, err := strconv.ParseInt(val.Number(), 10, 64); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = n
        }
        {{else if eq $goType "int32"}}
        if n, err := strconv.ParseInt(val.Number(), 10, 32); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = int32(n)
        }
        {{else if eq $goType "int16"}}
        if n, err := strconv.ParseInt(val.Number(), 10, 16); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = int16(n)
        }
        {{else if eq $goType "int8"}}
        if n, err := strconv.ParseInt(val.Number(), 10, 8); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = int8(n)
        }
        {{else if eq $goType "uint"}}
        if n, err := strconv.ParseUint(val.Number(), 10, 0); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = uint(n)
        }
        {{else if eq $goType "uint64"}}
        if n, err := strconv.ParseUint(val.Number(), 10, 64); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = n
        }
        {{else if eq $goType "uint32"}}
        if n, err := strconv.ParseUint(val.Number(), 10, 32); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = uint32(n)
        }
        {{else if eq $goType "uint16"}}
        if n, err := strconv.ParseUint(val.Number(), 10, 16); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = uint16(n)
        }
        {{else if eq $goType "uint8"}}
        if n, err := strconv.ParseUint(val.Number(), 10, 8); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = uint8(n)
        }
        {{else if eq $goType "float32"}}
        if n, err := strconv.ParseFloat(val.Number(), 32); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = float32(n)
        }
        {{else if eq $goType "float64"}}
        if n, err := strconv.ParseFloat(val.Number(), 64); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = n
        }
        {{else}}
        if n, err := strconv.ParseFloat(val.Number(), 64); err == nil {
            item.{{ToSafeName .Name | ToUpperCamelCase}} = {{$goType}}(n)
        }
        {{end}}
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
        _ = val
        {{end}}
    }
    {{end}}
  
    return item, nil
}

// IsFieldModified ...
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

// GetBoolFieldChanged ...
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

// CreateTriggerHandler ...
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
`
