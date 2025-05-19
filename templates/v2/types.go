package templates

var AttributeUtilsTemplate = `
func ConvertMapToAttributeValues(input map[string]interface{}) (map[string]types.AttributeValue, error) {
    result := make(map[string]types.AttributeValue)
    
    for key, value := range input {
        switch v := value.(type) {
        case string:
            result[key] = &types.AttributeValueMemberS{Value: v}
        case float64:
            // Числа в JSON обычно представлены как float64
            if v == float64(int64(v)) {
                // Целое число
                result[key] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", int64(v))}
            } else {
                result[key] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%g", v)}
            }
        case bool:
            if v {
                result[key] = &types.AttributeValueMemberN{Value: "1"} // Для совместимости с вашим BoolToInt
            } else {
                result[key] = &types.AttributeValueMemberN{Value: "0"}
            }
        case nil:
            result[key] = &types.AttributeValueMemberNULL{Value: true}
        case map[string]interface{}:
            // Обработка вложенных объектов
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
            // Обработка массивов
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
