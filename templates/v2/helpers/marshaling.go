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

// marshalUpdatesWithSchema marshals updates map considering field types from schema
func marshalUpdatesWithSchema(updates map[string]interface{}) (map[string]types.AttributeValue, error) {
    result := make(map[string]types.AttributeValue, len(updates))
    
    for fieldName, value := range updates {
        if fieldInfo, exists := TableSchema.FieldsMap[fieldName]; exists {
            // Use schema-aware marshaling for known fields
            av, err := marshalValueByType(value, fieldInfo.DynamoType)
            if err != nil {
                return nil, fmt.Errorf("failed to marshal field %s: %v", fieldName, err)
            }
            result[fieldName] = av
        } else {
            // Use default marshaling for unknown fields
            av, err := attributevalue.Marshal(value)
            if err != nil {
                return nil, fmt.Errorf("failed to marshal field %s: %v", fieldName, err)
            }
            result[fieldName] = av
        }
    }
    
    return result, nil
}

// marshalValueByType marshals value according to specific DynamoDB type
func marshalValueByType(value interface{}, dynamoType string) (types.AttributeValue, error) {
    switch dynamoType {
    case "SS":
        if strSlice, ok := value.([]string); ok {
            return &types.AttributeValueMemberSS{Value: strSlice}, nil
        }
        return nil, fmt.Errorf("expected []string for SS type, got %T", value)
    case "NS":
        if intSlice, ok := value.([]int); ok {
            numbers := make([]string, len(intSlice))
            for i, num := range intSlice {
                numbers[i] = fmt.Sprintf("%d", num)
            }
            return &types.AttributeValueMemberNS{Value: numbers}, nil
        }
        return nil, fmt.Errorf("expected []int for NS type, got %T", value)
    default:
        return attributevalue.Marshal(value)
    }
}
`
