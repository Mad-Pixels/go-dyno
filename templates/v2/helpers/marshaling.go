package helpers

// MarshalingHelpersTemplate ...
const MarshalingHelpersTemplate = `
// marshalItemToMap converts SchemaItem to AttributeValue map (internal helper)
func marshalItemToMap(item SchemaItem) (map[string]types.AttributeValue, error) {
    return attributevalue.MarshalMap(item)
}

// extractNonKeyAttributes filters out key attributes from the map
func extractNonKeyAttributes(allAttributes map[string]types.AttributeValue) map[string]types.AttributeValue {
    updates := make(map[string]types.AttributeValue, len(allAttributes)-2) // pre-allocate minus keys
    
    for attrName, attrValue := range allAttributes {
        if attrName != TableSchema.HashKey && attrName != TableSchema.RangeKey {
            updates[attrName] = attrValue
        }
    }
    
    return updates
}

// buildUpdateExpression creates SET expression from attribute map
func buildUpdateExpression(updates map[string]types.AttributeValue) (string, map[string]string, map[string]types.AttributeValue) {
    if len(updates) == 0 {
        return "", nil, nil
    }
    
    updateParts := make([]string, 0, len(updates))
    attrNames := make(map[string]string, len(updates))
    attrValues := make(map[string]types.AttributeValue, len(updates))
    
    i := 0
    for attrName, attrValue := range updates {
        nameKey := fmt.Sprintf("#attr%d", i)
        valueKey := fmt.Sprintf(":val%d", i)
        
        updateParts = append(updateParts, fmt.Sprintf("%s = %s", nameKey, valueKey))
        attrNames[nameKey] = attrName
        attrValues[valueKey] = attrValue
        i++
    }
    
    return "SET " + strings.Join(updateParts, ", "), attrNames, attrValues
}

// marshalRawUpdates converts raw map to AttributeValue map
func marshalRawUpdates(updates map[string]interface{}) (map[string]types.AttributeValue, error) {
    result := make(map[string]types.AttributeValue, len(updates))
    
    for attrName, value := range updates {
        av, err := attributevalue.Marshal(value)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal update value for %s: %v", attrName, err)
        }
        result[attrName] = av
    }
    
    return result, nil
}

// mergeExpressionAttributes merges condition attributes into existing maps
func mergeExpressionAttributes(
    baseNames map[string]string, 
    baseValues map[string]types.AttributeValue,
    conditionNames map[string]string, 
    conditionValues map[string]types.AttributeValue,
) (map[string]string, map[string]types.AttributeValue) {
    
    // Merge names
    if conditionNames != nil {
        for key, value := range conditionNames {
            baseNames[key] = value
        }
    }
    
    // Merge values  
    if conditionValues != nil {
        for key, value := range conditionValues {
            baseValues[key] = value
        }
    }
    
    return baseNames, baseValues
}
`
