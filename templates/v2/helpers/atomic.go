package helpers

// AtomicHelpersTemplate ...
const AtomicHelpersTemplate = `
// IncrementAttribute ...
func IncrementAttribute(hashKeyValue interface{}, rangeKeyValue interface{}, attributeName string, incrementValue int) (*dynamodb.UpdateItemInput, error) {
    // All validations at the beginning
    if err := validateKeyInputs(hashKeyValue, rangeKeyValue); err != nil {
        return nil, err
    }
    if err := validateAttributeName(attributeName); err != nil {
        return nil, err
    }
    if err := validateIncrementValue(incrementValue); err != nil {
        return nil, err
    }

    // Pure business logic after validation
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
    // All validations at the beginning
    if err := validateKeyInputs(hashKeyValue, rangeKeyValue); err != nil {
        return nil, err
    }
    if err := validateAttributeName(attributeName); err != nil {
        return nil, err
    }
    if err := validateSetValues(values); err != nil {
        return nil, err
    }

    // Pure business logic after validation
    key, err := KeyInputFromRaw(hashKeyValue, rangeKeyValue)
    if err != nil {
        return nil, fmt.Errorf("failed to create key for add to set: %v", err)
    }
   
    var attributeValue types.AttributeValue
   
    switch v := values.(type) {
    case []string:
        attributeValue = &types.AttributeValueMemberSS{Value: v}
    case []int:
        attributeValue = &types.AttributeValueMemberNS{Value: toIntStrings(v)}
    case []int8:
        attributeValue = &types.AttributeValueMemberNS{Value: toIntStrings(v)}
    case []int16:
        attributeValue = &types.AttributeValueMemberNS{Value: toIntStrings(v)}
    case []int32:
        attributeValue = &types.AttributeValueMemberNS{Value: toIntStrings(v)}
    case []int64:
        attributeValue = &types.AttributeValueMemberNS{Value: toIntStrings(v)}
    case []uint:
        attributeValue = &types.AttributeValueMemberNS{Value: toIntStrings(v)}
    case []uint8:
        attributeValue = &types.AttributeValueMemberNS{Value: toIntStrings(v)}
    case []uint16:
        attributeValue = &types.AttributeValueMemberNS{Value: toIntStrings(v)}
    case []uint32:
        attributeValue = &types.AttributeValueMemberNS{Value: toIntStrings(v)}
    case []uint64:
        attributeValue = &types.AttributeValueMemberNS{Value: toIntStrings(v)}
    case []float32:
        attributeValue = &types.AttributeValueMemberNS{Value: toFloatStrings(v)}
    case []float64:
        attributeValue = &types.AttributeValueMemberNS{Value: toFloatStrings(v)}
    default:
        return nil, fmt.Errorf("unsupported type for set operation: %T, expected []string or numeric slice", values)
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
    // All validations at the beginning
    if err := validateKeyInputs(hashKeyValue, rangeKeyValue); err != nil {
        return nil, err
    }
    if err := validateAttributeName(attributeName); err != nil {
        return nil, err
    }
    if err := validateSetValues(values); err != nil {
        return nil, err
    }
    
    // Pure business logic after validation
    key, err := KeyInputFromRaw(hashKeyValue, rangeKeyValue)
    if err != nil {
        return nil, fmt.Errorf("failed to create key for remove from set: %v", err)
    }
   
    var attributeValue types.AttributeValue
   
    switch v := values.(type) {
    case []string:
        attributeValue = &types.AttributeValueMemberSS{Value: v}
    case []int:
        attributeValue = &types.AttributeValueMemberNS{Value: toIntStrings(v)}
    case []int8:
        attributeValue = &types.AttributeValueMemberNS{Value: toIntStrings(v)}
    case []int16:
        attributeValue = &types.AttributeValueMemberNS{Value: toIntStrings(v)}
    case []int32:
        attributeValue = &types.AttributeValueMemberNS{Value: toIntStrings(v)}
    case []int64:
        attributeValue = &types.AttributeValueMemberNS{Value: toIntStrings(v)}
    case []uint:
        attributeValue = &types.AttributeValueMemberNS{Value: toIntStrings(v)}
    case []uint8:
        attributeValue = &types.AttributeValueMemberNS{Value: toIntStrings(v)}
    case []uint16:
        attributeValue = &types.AttributeValueMemberNS{Value: toIntStrings(v)}
    case []uint32:
        attributeValue = &types.AttributeValueMemberNS{Value: toIntStrings(v)}
    case []uint64:
        attributeValue = &types.AttributeValueMemberNS{Value: toIntStrings(v)}
    case []float32:
        attributeValue = &types.AttributeValueMemberNS{Value: toFloatStrings(v)}
    case []float64:
        attributeValue = &types.AttributeValueMemberNS{Value: toFloatStrings(v)}
    default:
        return nil, fmt.Errorf("unsupported type for set operation: %T, expected []string or numeric slice", values)
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
