package helpers

const ConverterHelpersTemplate = `
// ConvertMapToAttributeValues ...
func ConvertMapToAttributeValues(input map[string]interface{}) (map[string]types.AttributeValue, error) {
    result := make(map[string]types.AttributeValue)
    
    for key, value := range input {
        switch v := value.(type) {
        case string:
            result[key] = &types.AttributeValueMemberS{Value: v}
        case []string:
            result[key] = &types.AttributeValueMemberSS{Value: v}
        case []int:
            numbers := make([]string, len(v))
            for i, num := range v {
                numbers[i] = fmt.Sprintf("%d", num)
            }
            result[key] = &types.AttributeValueMemberNS{Value: numbers}
        case int:
            result[key] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", v)}
        case float64:
            if v == float64(int64(v)) {
                result[key] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", int64(v))}
            } else {
                result[key] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%g", v)}
            }
        case bool:
            result[key] = &types.AttributeValueMemberBOOL{Value: v}
        case nil:
            result[key] = &types.AttributeValueMemberNULL{Value: true}
        case map[string]interface{}:
            b, err := json.Marshal(v)
            if err != nil {
                return nil, err
            }
            result[key] = &types.AttributeValueMemberM{
                Value: map[string]types.AttributeValue{
                    "json": &types.AttributeValueMemberS{Value: string(b)},
                },
            }
        case []interface{}:
            b, err := json.Marshal(v)
            if err != nil {
                return nil, err
            }
            result[key] = &types.AttributeValueMemberL{
                Value: []types.AttributeValue{
                    &types.AttributeValueMemberS{Value: string(b)},
                },
            }
        default:
            return nil, fmt.Errorf("unsupported type for key %s: %T", key, value)
        }
    }
    
    return result, nil
}
`
