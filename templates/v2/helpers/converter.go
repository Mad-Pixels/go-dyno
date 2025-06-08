package helpers

// ConverterHelpersTemplate ...
const ConverterHelpersTemplate = `
// MarshalMap converts any Go value (map, struct, etc.) to DynamoDB AttributeValue map
// Uses AWS SDK's built-in marshaler for consistent behavior
func MarshalMap(input interface{}) (map[string]types.AttributeValue, error) {
    // Use AWS SDK's built-in marshaler - handles maps, structs, etc.
    result, err := attributevalue.MarshalMap(input)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal to AttributeValue map: %v", err)
    }
    
    return result, nil
}

// Marshal converts a single Go value to DynamoDB AttributeValue
// Uses AWS SDK's built-in marshaler for consistent behavior
func Marshal(input interface{}) (types.AttributeValue, error) {
    // Use AWS SDK's built-in marshaler - handles all Go types correctly
    result, err := attributevalue.Marshal(input)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal to AttributeValue: %v", err)
    }
    
    return result, nil
}
`
