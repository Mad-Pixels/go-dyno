package generic

// OperatorsTemplate generates operator definitions for universal query/filter methods
const OperatorsTemplate = `
// OperatorType defines the type of operation for queries and filters.
// Provides type-safe operator constants for DynamoDB expressions.
type OperatorType string

const (
    // Equality and comparison operators - work with all comparable types
    EQ  OperatorType = "="      // Equal to
    NE  OperatorType = "<>"     // Not equal to
    GT  OperatorType = ">"      // Greater than
    LT  OperatorType = "<"      // Less than
    GTE OperatorType = ">="     // Greater than or equal
    LTE OperatorType = "<="     // Less than or equal

    // Range operator for between comparisons
    BETWEEN OperatorType = "BETWEEN"  // Between two values (inclusive)

    // String operators - work with String types and Sets
    CONTAINS     OperatorType = "contains"      // Contains substring or set member
    NOT_CONTAINS OperatorType = "not_contains"  // Does not contain substring or member
    BEGINS_WITH  OperatorType = "begins_with"   // String starts with prefix

    // Set operators for scalar values only (not DynamoDB Sets SS/NS)
    IN     OperatorType = "IN"      // Value is in list of values
    NOT_IN OperatorType = "NOT_IN"  // Value is not in list of values

    // Existence operators - work with all types
    EXISTS     OperatorType = "attribute_exists"      // Attribute exists
    NOT_EXISTS OperatorType = "attribute_not_exists"  // Attribute does not exist
)

// ConditionType defines whether this is a key condition or filter condition.
// Key conditions are used in Query operations, filters in both Query and Scan.
type ConditionType string

const (
    KeyCondition    ConditionType = "KEY"     // For partition/sort key conditions
    FilterCondition ConditionType = "FILTER"  // For non-key attribute filtering
)

// Condition represents a single query or filter condition with validation metadata.
type Condition struct {
    Field     string         // Attribute name
    Operator  OperatorType   // Operation type
    Values    []any          // Operation values
    Type      ConditionType  // Key or filter condition
}

// Type-safe handler functions for different expression types.
// Provides compile-time safety for DynamoDB expression building.
type KeyOperatorHandler func(expression.KeyBuilder, []any) expression.KeyConditionBuilder
type ConditionOperatorHandler func(expression.NameBuilder, []any) expression.ConditionBuilder

// keyOperatorHandlers provides O(1) lookup for key condition operations.
// Only includes operators valid for key conditions (partition/sort keys).
var keyOperatorHandlers = map[OperatorType]KeyOperatorHandler{
    EQ: func(field expression.KeyBuilder, values []any) expression.KeyConditionBuilder {
        return field.Equal(expression.Value(values[0]))
    },
    GT: func(field expression.KeyBuilder, values []any) expression.KeyConditionBuilder {
        return field.GreaterThan(expression.Value(values[0]))
    },
    LT: func(field expression.KeyBuilder, values []any) expression.KeyConditionBuilder {
        return field.LessThan(expression.Value(values[0]))
    },
    GTE: func(field expression.KeyBuilder, values []any) expression.KeyConditionBuilder {
        return field.GreaterThanEqual(expression.Value(values[0]))
    },
    LTE: func(field expression.KeyBuilder, values []any) expression.KeyConditionBuilder {
        return field.LessThanEqual(expression.Value(values[0]))
    },
    BETWEEN: func(field expression.KeyBuilder, values []any) expression.KeyConditionBuilder {
        return field.Between(expression.Value(values[0]), expression.Value(values[1]))
    },
}

// allowedKeyConditionOperators defines operators valid for key conditions.
// Single source of truth for key condition validation.
var allowedKeyConditionOperators = map[OperatorType]bool{
    EQ:      true,  // Required for partition key
    GT:      true,  // Valid for sort key
    LT:      true,  // Valid for sort key
    GTE:     true,  // Valid for sort key
    LTE:     true,  // Valid for sort key
    BETWEEN: true,  // Valid for sort key
}

// conditionOperatorHandlers provides O(1) lookup for filter operations.
// Includes all operators supported in filter expressions.
var conditionOperatorHandlers = map[OperatorType]ConditionOperatorHandler{
    // Basic comparison operators
    EQ: func(field expression.NameBuilder, values []any) expression.ConditionBuilder {
        return field.Equal(expression.Value(values[0]))
    },
    NE: func(field expression.NameBuilder, values []any) expression.ConditionBuilder {
        return field.NotEqual(expression.Value(values[0]))
    },
    GT: func(field expression.NameBuilder, values []any) expression.ConditionBuilder {
        return field.GreaterThan(expression.Value(values[0]))
    },
    LT: func(field expression.NameBuilder, values []any) expression.ConditionBuilder {
        return field.LessThan(expression.Value(values[0]))
    },
    GTE: func(field expression.NameBuilder, values []any) expression.ConditionBuilder {
        return field.GreaterThanEqual(expression.Value(values[0]))
    },
    LTE: func(field expression.NameBuilder, values []any) expression.ConditionBuilder {
        return field.LessThanEqual(expression.Value(values[0]))
    },
    BETWEEN: func(field expression.NameBuilder, values []any) expression.ConditionBuilder {
        return field.Between(expression.Value(values[0]), expression.Value(values[1]))
    },
    
    // String and set operations
    CONTAINS: func(field expression.NameBuilder, values []any) expression.ConditionBuilder {
        return field.Contains(fmt.Sprintf("%v", values[0]))
    },
    NOT_CONTAINS: func(field expression.NameBuilder, values []any) expression.ConditionBuilder {
        return expression.Not(field.Contains(fmt.Sprintf("%v", values[0])))
    },
    BEGINS_WITH: func(field expression.NameBuilder, values []any) expression.ConditionBuilder {
        return field.BeginsWith(fmt.Sprintf("%v", values[0]))
    },
    
    // Scalar value list operations (not for DynamoDB Sets)
    IN: func(field expression.NameBuilder, values []any) expression.ConditionBuilder {
        if len(values) == 0 {
            return expression.AttributeNotExists(field)
        }
        if len(values) == 1 {
            return field.Equal(expression.Value(values[0]))
        }
        operands := make([]expression.OperandBuilder, len(values))
        for i, v := range values {
            operands[i] = expression.Value(v)
        }
        return field.In(operands[0], operands[1:]...)
    },
    NOT_IN: func(field expression.NameBuilder, values []any) expression.ConditionBuilder {
        if len(values) == 0 {
            return expression.AttributeExists(field)
        }
        if len(values) == 1 {
            return field.NotEqual(expression.Value(values[0]))
        }
        operands := make([]expression.OperandBuilder, len(values))
        for i, v := range values {
            operands[i] = expression.Value(v)
        }
        return expression.Not(field.In(operands[0], operands[1:]...))
    },
    
    // Existence checks
    EXISTS: func(field expression.NameBuilder, values []any) expression.ConditionBuilder {
        return expression.AttributeExists(field)
    },
    NOT_EXISTS: func(field expression.NameBuilder, values []any) expression.ConditionBuilder {
        return expression.AttributeNotExists(field)
    },
}

// ValidateValues checks if the number of values is correct for the operator.
// Prevents runtime errors by validating value count at build time.
func ValidateValues(op OperatorType, values []any) bool {
    switch op {
    case EQ, NE, GT, LT, GTE, LTE, CONTAINS, NOT_CONTAINS, BEGINS_WITH:
        return len(values) == 1  // Single value operators
    case BETWEEN:
        return len(values) == 2  // Start and end values
    case IN, NOT_IN:
        return len(values) >= 1  // At least one value required
    case EXISTS, NOT_EXISTS:
        return len(values) == 0  // No values needed
    default:
        return false
    }
}

// IsKeyConditionOperator checks if operator can be used in key conditions.
// Key conditions have stricter rules than filter conditions.
func IsKeyConditionOperator(op OperatorType) bool {
    return allowedKeyConditionOperators[op]
}

// ValidateOperator checks if operator is valid for the given field using schema.
// Provides type-safe operator validation based on DynamoDB field types.
func ValidateOperator(fieldName string, op OperatorType) bool {
    if fi, ok := TableSchema.FieldsMap[fieldName]; ok {
        return fi.SupportsOperator(op)
    }
    return false
}

// BuildConditionExpression converts operator to DynamoDB filter expression.
// Creates type-safe filter conditions with full validation.
// Example: BuildConditionExpression("name", EQ, []any{"John"})
func BuildConditionExpression(field string, op OperatorType, values []any) (expression.ConditionBuilder, error) {
    // Check if field exists in schema
    fieldInfo, exists := TableSchema.FieldsMap[field]
    if !exists {
        return expression.ConditionBuilder{}, fmt.Errorf("field %s not found in schema", field)
    }
    
    // Check if operator is supported for this field type
    if !fieldInfo.SupportsOperator(op) {
        return expression.ConditionBuilder{}, fmt.Errorf("operator %s not supported for field %s (type %s)", op, field, fieldInfo.DynamoType)
    }
    
    if !ValidateValues(op, values) {
        return expression.ConditionBuilder{}, fmt.Errorf("invalid number of values for operator %s", op)
    }
    
    handler := conditionOperatorHandlers[op]
    fieldExpr := expression.Name(field)
    result := handler(fieldExpr, values)
    
    return result, nil
}

// BuildKeyConditionExpression converts operator to DynamoDB key condition.
// Creates type-safe key conditions for Query operations only.
// Example: BuildKeyConditionExpression("user_id", EQ, []any{"123"})
func BuildKeyConditionExpression(field string, op OperatorType, values []any) (expression.KeyConditionBuilder, error) {
    // Check if field exists in schema
    fieldInfo, exists := TableSchema.FieldsMap[field]
    if !exists {
        return expression.KeyConditionBuilder{}, fmt.Errorf("field %s not found in schema", field)
    }
    
    // Check if field is actually a key
    if !fieldInfo.IsKey {
        return expression.KeyConditionBuilder{}, fmt.Errorf("field %s is not a key field", field)
    }
    
    // Check if operator is supported for this field type
    if !fieldInfo.SupportsOperator(op) {
        return expression.KeyConditionBuilder{}, fmt.Errorf("operator %s not supported for field %s (type %s)", op, field, fieldInfo.DynamoType)
    }
    
    if !ValidateValues(op, values) {
        return expression.KeyConditionBuilder{}, fmt.Errorf("invalid number of values for operator %s", op)
    }
    
    handler := keyOperatorHandlers[op]
    fieldExpr := expression.Key(field)
    result := handler(fieldExpr, values)
    
    return result, nil
}
`
