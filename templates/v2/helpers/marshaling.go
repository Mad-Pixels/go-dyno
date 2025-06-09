package helpers

// MarshalingHelpersTemplate ...
const MarshalingHelpersTemplate = `
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Float interface {
	~float32 | ~float64
}

// toIntStrings converts any signed or unsigned integer slice to string slice
func toIntStrings[T Signed | Unsigned](nums []T) []string {
	out := make([]string, len(nums))
	for i, n := range nums {
		out[i] = strconv.FormatInt(int64(n), 10)
	}
	return out
}

// toFloatStrings converts any float slice to string slice
func toFloatStrings[F Float](nums []F) []string {
	out := make([]string, len(nums))
	for i, f := range nums {
		// 'g' убирает лишние нули и точку, -1 — максимальная точность
		out[i] = strconv.FormatFloat(float64(f), 'g', -1, 64)
	}
	return out
}

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
        ss, ok := value.([]string)
        if !ok {
            return nil, fmt.Errorf("SS: expected []string, got %T", value)
        }
        return &types.AttributeValueMemberSS{Value: ss}, nil
    case "NS":
        switch v := value.(type) {
        case []int:
            return &types.AttributeValueMemberNS{Value: toIntStrings(v)}, nil
        case []int8:
            return &types.AttributeValueMemberNS{Value: toIntStrings(v)}, nil
        case []int16:
            return &types.AttributeValueMemberNS{Value: toIntStrings(v)}, nil
        case []int32:
            return &types.AttributeValueMemberNS{Value: toIntStrings(v)}, nil
        case []int64:
            return &types.AttributeValueMemberNS{Value: toIntStrings(v)}, nil
        case []uint:
            return &types.AttributeValueMemberNS{Value: toIntStrings(v)}, nil
        case []uint8:
            return &types.AttributeValueMemberNS{Value: toIntStrings(v)}, nil
        case []uint16:
            return &types.AttributeValueMemberNS{Value: toIntStrings(v)}, nil
        case []uint32:
            return &types.AttributeValueMemberNS{Value: toIntStrings(v)}, nil
        case []uint64:
            return &types.AttributeValueMemberNS{Value: toIntStrings(v)}, nil
        case []float32:
            return &types.AttributeValueMemberNS{Value: toFloatStrings(v)}, nil
        case []float64:
            return &types.AttributeValueMemberNS{Value: toFloatStrings(v)}, nil
        default:
            return nil, fmt.Errorf("NS: expected numeric slice, got %T", value)
        }
    default:
        return attributevalue.Marshal(value)
    }
}
`
