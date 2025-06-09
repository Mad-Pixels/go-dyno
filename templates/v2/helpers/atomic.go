package helpers

// AtomicHelpersTemplate ...
const AtomicHelpersTemplate = `
// IncrementAttribute ...
func IncrementAttribute(hashKeyValue interface{}, rangeKeyValue interface{}, attributeName string, incrementValue int) (*dynamodb.UpdateItemInput, error) {
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

// AddToSet ...
func AddToSet(hashKeyValue interface{}, rangeKeyValue interface{}, attributeName string, values interface{}) (*dynamodb.UpdateItemInput, error) {
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

// RemoveFromSet ...
func RemoveFromSet(hashKeyValue interface{}, rangeKeyValue interface{}, attributeName string, values interface{}) (*dynamodb.UpdateItemInput, error) {
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
