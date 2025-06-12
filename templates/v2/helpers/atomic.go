package helpers

// AtomicHelpersTemplate provides atomic update operations for DynamoDB
const AtomicHelpersTemplate = `
// IncrementAttribute atomically increments a numeric attribute by a specified value.
// Uses DynamoDB's ADD operation to ensure thread-safe increments without race conditions.
// Creates the attribute with the increment value if it doesn't exist.
// Example: IncrementAttribute("user123", nil, "view_count", 1)
func IncrementAttribute(hashKeyValue any, rangeKeyValue any, attributeName string, incrementValue int) (*dynamodb.UpdateItemInput, error) {
    if err := validateKeyInputs(hashKeyValue, rangeKeyValue); err != nil {
        return nil, err
    }
    if err := validateAttributeName(attributeName); err != nil {
        return nil, err
    }
    if err := validateIncrementValue(incrementValue); err != nil {
        return nil, err
    }

    key, err := KeyInputFromRaw(hashKeyValue, rangeKeyValue)
    if err != nil {
        return nil, fmt.Errorf("failed to create key for increment: %v", err)
    }
   
    return &dynamodb.UpdateItemInput{
        TableName:        aws.String(TableSchema.TableName),
        Key:              key,
        UpdateExpression: aws.String("ADD #attr :val"),
        ExpressionAttributeNames: map[string]string{
            "#attr": attributeName,
        },
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":val": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", incrementValue)},
        },
    }, nil
}

// AddToSet atomically adds values to a DynamoDB Set (SS or NS).
// Uses DynamoDB's ADD operation for sets - duplicate values are automatically ignored.
// Creates the set with provided values if the attribute doesn't exist.
// Supports string sets ([]string) and numeric sets ([]int, []float64, etc.).
// Example: AddToSet("user123", nil, "tags", []string{"premium", "verified"})
func AddToSet(hashKeyValue any, rangeKeyValue any, attributeName string, values any) (*dynamodb.UpdateItemInput, error) {
    if err := validateKeyInputs(hashKeyValue, rangeKeyValue); err != nil {
        return nil, err
    }
    if err := validateAttributeName(attributeName); err != nil {
        return nil, err
    }
    if err := validateSetValues(values); err != nil {
        return nil, err
    }

    key, err := KeyInputFromRaw(hashKeyValue, rangeKeyValue)
    if err != nil {
        return nil, fmt.Errorf("failed to create key for add to set: %v", err)
    }
   
    var attributeValue types.AttributeValue
   
    {{- $nsTypes := GetUsedNumericSetTypes .AllAttributes}}
    switch v := values.(type) {
    case []string:
        attributeValue = &types.AttributeValueMemberSS{Value: v}
    {{- if gt (len $nsTypes) 0}}
    {{- range $nsTypes}}
    case {{.}}:
        {{- if IsFloatType (Slice . 2)}}
        attributeValue = &types.AttributeValueMemberNS{Value: toFloatStrings(v)}
        {{- else}}
        attributeValue = &types.AttributeValueMemberNS{Value: toIntStrings(v)}
        {{- end}}
    {{- end}}
    {{- end}}
    default:
        return nil, fmt.Errorf("unsupported type for set operation: %T, expected []string{{if gt (len $nsTypes) 0}} or numeric slice{{end}}", values)
    }
   
    return &dynamodb.UpdateItemInput{
        TableName:        aws.String(TableSchema.TableName),
        Key:              key,
        UpdateExpression: aws.String("ADD #attr :val"),
        ExpressionAttributeNames: map[string]string{
            "#attr": attributeName,
        },
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":val": attributeValue,
        },
    }, nil
}

// RemoveFromSet atomically removes values from a DynamoDB Set (SS or NS).
// Uses DynamoDB's DELETE operation for sets - non-existent values are ignored.
// If all values are removed, the attribute is deleted from the item.
// Supports string sets ([]string) and numeric sets ([]int, []float64, etc.).
// Example: RemoveFromSet("user123", nil, "tags", []string{"temporary"})
func RemoveFromSet(hashKeyValue any, rangeKeyValue any, attributeName string, values any) (*dynamodb.UpdateItemInput, error) {
    if err := validateKeyInputs(hashKeyValue, rangeKeyValue); err != nil {
        return nil, err
    }
    if err := validateAttributeName(attributeName); err != nil {
        return nil, err
    }
    if err := validateSetValues(values); err != nil {
        return nil, err
    }
    
    key, err := KeyInputFromRaw(hashKeyValue, rangeKeyValue)
    if err != nil {
        return nil, fmt.Errorf("failed to create key for remove from set: %v", err)
    }
   
    var attributeValue types.AttributeValue
   
    {{- $nsTypes := GetUsedNumericSetTypes .AllAttributes}}
    switch v := values.(type) {
    case []string:
        attributeValue = &types.AttributeValueMemberSS{Value: v}
    {{- if gt (len $nsTypes) 0}}
    {{- range $nsTypes}}
    case {{.}}:
        {{- if IsFloatType (Slice . 2)}}
        attributeValue = &types.AttributeValueMemberNS{Value: toFloatStrings(v)}
        {{- else}}
        attributeValue = &types.AttributeValueMemberNS{Value: toIntStrings(v)}
        {{- end}}
    {{- end}}
    {{- end}}
    default:
        return nil, fmt.Errorf("unsupported type for set operation: %T, expected []string{{if gt (len $nsTypes) 0}} or numeric slice{{end}}", values)
    }
   
    return &dynamodb.UpdateItemInput{
        TableName:        aws.String(TableSchema.TableName),
        Key:              key,
        UpdateExpression: aws.String("DELETE #attr :val"),
        ExpressionAttributeNames: map[string]string{
            "#attr": attributeName,
        },
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":val": attributeValue,
        },
    }, nil
}
`
