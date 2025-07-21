package inputs

// ItemInputsTemplate provides marshaling utilities for DynamoDB item operations
const ItemInputsTemplate = `
// ItemInput converts a SchemaItem to DynamoDB AttributeValue map format.
// Uses AWS SDK's attributevalue package for safe and consistent marshaling.
// The resulting map can be used in PutItem, UpdateItem, and other DynamoDB operations.
func ItemInput(item SchemaItem) (map[string]types.AttributeValue, error) {
    attributeValues, err := attributevalue.MarshalMap(item)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal item: %v", err)
    }
    return attributeValues, nil
}

// ItemsInput converts a slice of SchemaItems to DynamoDB AttributeValue maps.
// Efficiently marshals multiple items for batch operations like BatchWriteItem.
// Maintains order and provides detailed error context for debugging failed marshaling.
func ItemsInput(items []SchemaItem) ([]map[string]types.AttributeValue, error) {
    result := make([]map[string]types.AttributeValue, 0, len(items))
    for i, item := range items {
        av, err := ItemInput(item)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal item at index %d: %v", i, err)
        }
        result = append(result, av)
    }
    return result, nil
}
`
