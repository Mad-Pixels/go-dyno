package helpers

// ValidationHelpersTemplate provides comprehensive validation for DynamoDB operations
const ValidationHelpersTemplate = `
// validateKeyPart checks if key part (hash or range) value is valid for DynamoDB.
// Hash keys are required and cannot be nil/empty, range keys are optional.
// Supports string, numeric types commonly used as DynamoDB keys.
func validateKeyPart(partName string, value any) error {
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
        // numbers are always valid for keys
    case float32, float64:
        // floats are valid but unusual for keys
    default:
        return fmt.Errorf("unsupported %s key type: %T", partName, value)
    }
    
    return nil
}

// validateHashKey checks if hash key value is valid for DynamoDB operations.
// Hash key is required for all DynamoDB operations and cannot be nil or empty.
// Example: validateHashKey("user123") -> nil, validateHashKey("") -> error
func validateHashKey(value any) error {
    return validateKeyPart("hash", value)
}

// validateRangeKey checks if range key value is valid (nil is allowed).
// Range key is optional - tables can have simple (hash only) or composite keys.
// Example: validateRangeKey(nil) -> nil, validateRangeKey("timestamp") -> nil
func validateRangeKey(value any) error {
    return validateKeyPart("range", value)
}

// validateAttributeName checks if attribute name meets DynamoDB requirements.
// DynamoDB limits: non-empty, max 255 characters.
// Used to prevent API errors from invalid attribute names.
func validateAttributeName(name string) error {
    if name == "" {
        return fmt.Errorf("attribute name cannot be empty")
    }
    
    if len(name) > 255 {
        return fmt.Errorf("attribute name too long: %d chars (max 255)", len(name))
    }
    
    return nil
}

// validateUpdatesMap checks if updates map is valid for UpdateItem operations.
// Ensures non-empty map with valid attribute names and non-nil values.
// Prevents wasted API calls and provides clear error messages.
func validateUpdatesMap(updates map[string]any) error {
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

// validateBatchSize checks if batch size is within DynamoDB limits.
// DynamoDB batch operations (BatchGetItem, BatchWriteItem) have a 25 item limit.
// Prevents API errors and guides proper batch partitioning.
// Example: validateBatchSize(30, "write") -> error about exceeding limit
func validateBatchSize(size int, operation string) error {
    if size == 0 {
        return fmt.Errorf("%s batch cannot be empty", operation)
    }
    
    if size > 25 {
        return fmt.Errorf("%s batch size %d exceeds DynamoDB limit of 25", operation, size)
    }
    
    return nil
}

// validateSetValues checks if set values are valid for AddToSet/RemoveFromSet operations.
// DynamoDB sets cannot be empty and string sets cannot contain empty strings.
// Validates both string sets (SS) and numeric sets (NS) with proper type checking.
func validateSetValues(values any) error {
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
    case []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64, []float32, []float64:
        // Use reflection to check length for all numeric types
        rv := reflect.ValueOf(v)
        if rv.Len() == 0 {
            return fmt.Errorf("number set cannot be empty")
        }
    default:
        return fmt.Errorf("unsupported set type: %T, expected []string or numeric slice", values)
    }
    
    return nil
}

// validateConditionExpression checks if condition expression meets DynamoDB limits.
// DynamoDB condition expressions have a 4KB size limit.
// Helps prevent API errors from oversized expressions.
func validateConditionExpression(expr string) error {
    if expr == "" {
        return fmt.Errorf("condition expression cannot be empty")
    }
    
    if len(expr) > 4096 {
        return fmt.Errorf("condition expression too long: %d chars (max 4096)", len(expr))
    }
    
    return nil
}

// validateIncrementValue checks if increment value is valid for atomic operations.
// DynamoDB ADD operation accepts any integer value (positive or negative).
// Function maintained for API consistency and future validation needs.
func validateIncrementValue(value int) error {
    // DynamoDB supports any int value for ADD operation
    // No specific validation needed, but we keep the function for consistency
    return nil
}

// validateKeyInputs validates both hash and range key inputs for DynamoDB operations.
// Comprehensive validation for all key-based operations (GetItem, UpdateItem, etc.).
// Provides clear error context for debugging key-related issues.
// Example: validateKeyInputs("user123", "2023-01-01") -> nil
func validateKeyInputs(hashKeyValue, rangeKeyValue any) error {
    if err := validateHashKey(hashKeyValue); err != nil {
        return fmt.Errorf("invalid hash key: %v", err)
    }
    
    if err := validateRangeKey(rangeKeyValue); err != nil {
        return fmt.Errorf("invalid range key: %v", err)
    }
    
    return nil
}
`
