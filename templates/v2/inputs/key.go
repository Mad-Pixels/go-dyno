package inputs

// KeyInputsTemplate ...
const KeyInputsTemplate = `
// KeyInput creates key from SchemaItem with validation
func KeyInput(item SchemaItem) (map[string]types.AttributeValue, error) {
    var hashKeyValue any
    {{range .AllAttributes}}{{if eq .Name $.HashKey}}
    hashKeyValue = item.{{ToSafeName .Name | ToUpperCamelCase}}
    {{end}}{{end}}
    
    var rangeKeyValue any
    {{if .RangeKey}}{{range .AllAttributes}}{{if eq .Name $.RangeKey}}
    rangeKeyValue = item.{{ToSafeName .Name | ToUpperCamelCase}}
    {{end}}{{end}}{{end}}
    
    // Single validation call at the beginning
    if err := validateKeyInputs(hashKeyValue, rangeKeyValue); err != nil {
        return nil, err
    }
    
    // Pure business logic after validation
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

// KeyInputFromRaw creates key from raw values (assumes validation already done)
func KeyInputFromRaw(hashKeyValue any, rangeKeyValue any) (map[string]types.AttributeValue, error) {
    // Pure business logic - validation should be done by caller
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
