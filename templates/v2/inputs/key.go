package inputs

// KeyInputsTemplate provides key extraction utilities for DynamoDB operations
const KeyInputsTemplate = `
// KeyInput creates a DynamoDB key map from a SchemaItem with full validation.
// Extracts the primary key (hash + range) from the item and validates values.
// Use when you have a complete item and need to create a key for operations.
// Handles both simple (hash only) and composite (hash + range) keys automatically.
func KeyInput(item SchemaItem) (map[string]types.AttributeValue, error) {
    var hashKeyValue any
    {{range .AllAttributes}}{{if eq .Name $.HashKey}}
    hashKeyValue = item.{{ToSafeName .Name | ToUpperCamelCase}}
    {{end}}{{end}}
    
    var rangeKeyValue any
    {{if .RangeKey}}{{range .AllAttributes}}{{if eq .Name $.RangeKey}}
    rangeKeyValue = item.{{ToSafeName .Name | ToUpperCamelCase}}
    {{end}}{{end}}{{end}}
    
    if err := validateKeyInputs(hashKeyValue, rangeKeyValue); err != nil {
        return nil, err
    }
    key := make(map[string]types.AttributeValue)
   
    hashKeyAV, err := attributevalue.Marshal(hashKeyValue)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal hash key: %v", err)
    }
    key[TableSchema.HashKey] = hashKeyAV
   
    {{if .RangeKey}}
    if TableSchema.RangeKey != "" && rangeKeyValue != nil {
        rangeKeyAV, err := attributevalue.Marshal(rangeKeyValue)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal range key: %v", err)
        }
        key[TableSchema.RangeKey] = rangeKeyAV
    }
    {{end}}
    return key, nil
}

// KeyInputFromRaw creates a DynamoDB key map from raw key values without validation.
// More efficient than KeyInput when you already have validated key values.
// Assumes validation has been done by the caller - use with caution.
// Handles both simple (hash only) and composite (hash + range) keys automatically.
func KeyInputFromRaw(hashKeyValue any, rangeKeyValue any) (map[string]types.AttributeValue, error) {
    key := make(map[string]types.AttributeValue)
   
    hashKeyAV, err := attributevalue.Marshal(hashKeyValue)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal hash key: %v", err)
    }
    key[TableSchema.HashKey] = hashKeyAV
   
    {{if .RangeKey}}
    if TableSchema.RangeKey != "" && rangeKeyValue != nil {
        rangeKeyAV, err := attributevalue.Marshal(rangeKeyValue)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal range key: %v", err)
        }
        key[TableSchema.RangeKey] = rangeKeyAV
    }
    {{end}}
    return key, nil
}
`
