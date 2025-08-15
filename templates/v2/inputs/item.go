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
`
