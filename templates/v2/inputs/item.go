package inputs

// ItemInputsTemplate ...
const ItemInputsTemplate = `
// ItemInput ...
func ItemInput(item SchemaItem) (map[string]types.AttributeValue, error) {
    attributeValues, err := attributevalue.MarshalMap(item)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal item: %v", err)
    }
    return attributeValues, nil
}

// ItemsInput ...
func ItemsInput(items []SchemaItem) ([]map[string]types.AttributeValue, error) {
    result := make([]map[string]types.AttributeValue, 0, len(items))
    for _, item := range items {
        av, err := ItemInput(item)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal item: %v", err)
        }
        result = append(result, av)
    }
    return result, nil
}
`
