package inputs

// KeyInputsTemplate ...
const KeyInputsTemplate = `
// KeyInput ...
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

// KeyInputFromRaw ...
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
