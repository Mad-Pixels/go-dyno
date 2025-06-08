package inputs

// KeyInputsTemplate ...
const KeyInputsTemplate = `
// KeyInput creates key from SchemaItem with validation
func KeyInput(item SchemaItem) (map[string]types.AttributeValue, error) {
    key := make(map[string]types.AttributeValue)
   
    var hashKeyValue interface{}
    {{range .AllAttributes}}{{if eq .Name $.HashKey}}
    hashKeyValue = item.{{ToSafeName .Name | ToUpperCamelCase}}
    {{end}}{{end}}
   
    if err := validateHashKey(hashKeyValue); err != nil {
        return nil, fmt.Errorf("invalid hash key: %v", err)
    }
   
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
       
        if err := validateRangeKey(rangeKeyValue); err != nil {
            return nil, fmt.Errorf("invalid range key: %v", err)
        }
       
        rangeKeyAV, err := attributevalue.Marshal(rangeKeyValue)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal range key: %v", err)
        }
        key[TableSchema.RangeKey] = rangeKeyAV
    }
   
    return key, nil
}

// KeyInputFromRaw creates key from raw values with validation
func KeyInputFromRaw(hashKeyValue interface{}, rangeKeyValue interface{}) (map[string]types.AttributeValue, error) {
    if err := validateHashKey(hashKeyValue); err != nil {
        return nil, fmt.Errorf("invalid hash key: %v", err)
    }
    
    if err := validateRangeKey(rangeKeyValue); err != nil {
        return nil, fmt.Errorf("invalid range key: %v", err)
    }
    
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
