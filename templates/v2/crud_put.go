package v2

// CrudPutTemplate ...
const CrudPutTemplate = `
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
`
