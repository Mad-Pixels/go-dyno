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
    updates := make(map[string]types.AttributeValue, len(allAttributes)-2)
    
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
    
    if conditionNames != nil {
        for key, value := range conditionNames {
            baseNames[key] = value
        }
    }
    
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
            av, err := marshalValueByType(value, fieldInfo.DynamoType)
            if err != nil {
                return nil, fmt.Errorf("failed to marshal field %s: %v", fieldName, err)
            }
            result[fieldName] = av
        } else {
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
        {{- $nsTypes := GetUsedNumericSetTypes .AllAttributes}}
        {{- if gt (len $nsTypes) 0}}
        switch v := value.(type) {
        {{- range $nsTypes}}
        case {{.}}:
            {{- if IsFloatType (Slice . 2)}}
            return &types.AttributeValueMemberNS{Value: toFloatStrings(v)}, nil
            {{- else}}
            return &types.AttributeValueMemberNS{Value: toIntStrings(v)}, nil
            {{- end}}
        {{- end}}
        default:
            return nil, fmt.Errorf("NS: expected numeric slice, got %T", value)
        }
        {{- else}}
        return nil, fmt.Errorf("NS: no numeric set types defined in schema")
        {{- end}}
    default:
        return attributevalue.Marshal(value)
    }
}
`
