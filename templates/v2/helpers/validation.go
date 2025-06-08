package helpers

// ValidationHelpersTemplate ...
const ValidationHelpersTemplate = `
// validateKeyPart checks if key part (hash or range) value is valid
func validateKeyPart(partName string, value interface{}) error {
    if value == nil {
        // Only hash key cannot be nil, range key can be nil
        if partName == "hash" {
            return fmt.Errorf("hash key cannot be nil")
        }
        return nil // range key can be nil
    }
    
    switch v := value.(type) {
    case string:
        if v == "" && partName == "hash" {
            return fmt.Errorf("hash key string cannot be empty")
        }
    case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
        // numbers are always valid
    case float32, float64:
        // floats are valid but unusual for keys
    default:
        return fmt.Errorf("unsupported %s key type: %T", partName, value)
    }
    
    return nil
}

// validateHashKey checks if hash key value is valid
func validateHashKey(value interface{}) error {
    return validateKeyPart("hash", value)
}

// validateRangeKey checks if range key value is valid (nil is allowed)
func validateRangeKey(value interface{}) error {
    return validateKeyPart("range", value)
}

// validateAttributeName checks if attribute name is valid
func validateAttributeName(name string) error {
    if name == "" {
        return fmt.Errorf("attribute name cannot be empty")
    }
    
    if len(name) > 255 {
        return fmt.Errorf("attribute name too long: %d chars (max 255)", len(name))
    }
    
    return nil
}

// validateUpdatesMap checks if updates map is valid
func validateUpdatesMap(updates map[string]interface{}) error {
    if len(updates) == 0 {
        return fmt.Errorf("updates map cannot be empty")
    }
    
    for attrName, value := range updates {
        if err := validateAttributeName(attrName); err != nil {
            return fmt.Errorf("invalid attribute name '%s': %v", attrName, err)
        }
        
        if value == nil {
            return fmt.Errorf("update value for '%s' cannot be nil", attrName)
        }
    }
    
    return nil
}

// validateBatchSize checks if batch size is within DynamoDB limits
func validateBatchSize(size int, operation string) error {
    if size == 0 {
        return fmt.Errorf("%s batch cannot be empty", operation)
    }
    
    if size > 25 {
        return fmt.Errorf("%s batch size %d exceeds DynamoDB limit of 25", operation, size)
    }
    
    return nil
}

// validateSetValues checks if set values are valid for AddToSet/RemoveFromSet
func validateSetValues(values interface{}) error {
    if values == nil {
        return fmt.Errorf("set values cannot be nil")
    }
    
    switch v := values.(type) {
    case []string:
        if len(v) == 0 {
            return fmt.Errorf("string set cannot be empty")
        }
        for i, str := range v {
            if str == "" {
                return fmt.Errorf("string set item %d cannot be empty", i)
            }
        }
    case []int:
        if len(v) == 0 {
            return fmt.Errorf("number set cannot be empty")
        }
    default:
        return fmt.Errorf("unsupported set type: %T, expected []string or []int", values)
    }
    
    return nil
}

// validateConditionExpression checks if condition expression is valid
func validateConditionExpression(expr string) error {
    if expr == "" {
        return fmt.Errorf("condition expression cannot be empty")
    }
    
    if len(expr) > 4096 {
        return fmt.Errorf("condition expression too long: %d chars (max 4096)", len(expr))
    }
    
    return nil
}

// validateIncrementValue checks if increment value is valid
func validateIncrementValue(value int) error {
    // DynamoDB supports any int value for ADD operation
    // No specific validation needed, but we keep the function for consistency
    return nil
}

// validateKeyInputs validates both hash and range key inputs for operations
func validateKeyInputs(hashKeyValue, rangeKeyValue interface{}) error {
    if err := validateHashKey(hashKeyValue); err != nil {
        return fmt.Errorf("invalid hash key: %v", err)
    }
    
    if err := validateRangeKey(rangeKeyValue); err != nil {
        return fmt.Errorf("invalid range key: %v", err)
    }
    
    return nil
}
`
