package generic

// OperatorsTemplate generates operator definitions for universal query/filter methods
const OperatorsTemplate = `
// OperatorType defines the type of operation for queries and filters
type OperatorType string

const (
    // Equality and comparison operators
    EQ  OperatorType = "="
    NE  OperatorType = "<>"
    GT  OperatorType = ">"
    LT  OperatorType = "<"
    GTE OperatorType = ">="
    LTE OperatorType = "<="

    // Range operator
    BETWEEN OperatorType = "BETWEEN"

    // String operators
    CONTAINS     OperatorType = "contains"
    NOT_CONTAINS OperatorType = "not_contains"
    BEGINS_WITH  OperatorType = "begins_with"

    // Set operators (for scalar values only)
    IN     OperatorType = "IN"
    NOT_IN OperatorType = "NOT_IN"

    // Existence operators
    EXISTS     OperatorType = "attribute_exists"
    NOT_EXISTS OperatorType = "attribute_not_exists"
)

// ConditionType defines whether this is a key condition or filter condition
type ConditionType string

const (
    KeyCondition    ConditionType = "KEY"
    FilterCondition ConditionType = "FILTER"
)

// Condition represents a single query or filter condition
type Condition struct {
    Field     string
    Operator  OperatorType
    Values    []interface{}
    Type      ConditionType
}

// Type-safe handler functions for different expression types
type KeyOperatorHandler func(expression.KeyBuilder, []interface{}) expression.KeyConditionBuilder
type ConditionOperatorHandler func(expression.NameBuilder, []interface{}) expression.ConditionBuilder

// Package-level constants for O(1) handler lookup - initialized once
var keyOperatorHandlers = map[OperatorType]KeyOperatorHandler{
    EQ: func(field expression.KeyBuilder, values []interface{}) expression.KeyConditionBuilder {
        return field.Equal(expression.Value(values[0]))
    },
    GT: func(field expression.KeyBuilder, values []interface{}) expression.KeyConditionBuilder {
        return field.GreaterThan(expression.Value(values[0]))
    },
    LT: func(field expression.KeyBuilder, values []interface{}) expression.KeyConditionBuilder {
        return field.LessThan(expression.Value(values[0]))
    },
    GTE: func(field expression.KeyBuilder, values []interface{}) expression.KeyConditionBuilder {
        return field.GreaterThanEqual(expression.Value(values[0]))
    },
    LTE: func(field expression.KeyBuilder, values []interface{}) expression.KeyConditionBuilder {
        return field.LessThanEqual(expression.Value(values[0]))
    },
    BETWEEN: func(field expression.KeyBuilder, values []interface{}) expression.KeyConditionBuilder {
        return field.Between(expression.Value(values[0]), expression.Value(values[1]))
    },
}

var conditionOperatorHandlers = map[OperatorType]ConditionOperatorHandler{
    // Basic operators (same as for keys, but with NameBuilder)
    EQ: func(field expression.NameBuilder, values []interface{}) expression.ConditionBuilder {
        return field.Equal(expression.Value(values[0]))
    },
    NE: func(field expression.NameBuilder, values []interface{}) expression.ConditionBuilder {
        return field.NotEqual(expression.Value(values[0]))
    },
    GT: func(field expression.NameBuilder, values []interface{}) expression.ConditionBuilder {
        return field.GreaterThan(expression.Value(values[0]))
    },
    LT: func(field expression.NameBuilder, values []interface{}) expression.ConditionBuilder {
        return field.LessThan(expression.Value(values[0]))
    },
    GTE: func(field expression.NameBuilder, values []interface{}) expression.ConditionBuilder {
        return field.GreaterThanEqual(expression.Value(values[0]))
    },
    LTE: func(field expression.NameBuilder, values []interface{}) expression.ConditionBuilder {
        return field.LessThanEqual(expression.Value(values[0]))
    },
    BETWEEN: func(field expression.NameBuilder, values []interface{}) expression.ConditionBuilder {
        return field.Between(expression.Value(values[0]), expression.Value(values[1]))
    },
    
    // Filter-only operators
    CONTAINS: func(field expression.NameBuilder, values []interface{}) expression.ConditionBuilder {
        return field.Contains(fmt.Sprintf("%v", values[0]))
    },
    NOT_CONTAINS: func(field expression.NameBuilder, values []interface{}) expression.ConditionBuilder {
        return expression.Not(field.Contains(fmt.Sprintf("%v", values[0])))
    },
    BEGINS_WITH: func(field expression.NameBuilder, values []interface{}) expression.ConditionBuilder {
        return field.BeginsWith(fmt.Sprintf("%v", values[0]))
    },
    IN: func(field expression.NameBuilder, values []interface{}) expression.ConditionBuilder {
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
    NOT_IN: func(field expression.NameBuilder, values []interface{}) expression.ConditionBuilder {
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
    EXISTS: func(field expression.NameBuilder, values []interface{}) expression.ConditionBuilder {
        return expression.AttributeExists(field)
    },
    NOT_EXISTS: func(field expression.NameBuilder, values []interface{}) expression.ConditionBuilder {
        return expression.AttributeNotExists(field)
    },
}

// ValidateOperator checks if operator is valid for the given field using pre-computed cache
func ValidateOperator(fieldName string, op OperatorType) bool {
    if fi, ok := TableSchema.FieldsMap[fieldName]; ok {
        return fi.SupportsOperator(op)
    }
    return false
}

// ValidateValues checks if the number of values is correct for the operator
func ValidateValues(op OperatorType, values []interface{}) bool {
    switch op {
    case EQ, NE, GT, LT, GTE, LTE, CONTAINS, NOT_CONTAINS, BEGINS_WITH:
        return len(values) == 1
    case BETWEEN:
        return len(values) == 2
    case IN, NOT_IN:
        return len(values) >= 1
    case EXISTS, NOT_EXISTS:
        return len(values) == 0
    default:
        return false
    }
}

// IsKeyConditionOperator checks if operator can be used in key conditions
func IsKeyConditionOperator(op OperatorType) bool {
    switch op {
    case EQ, GT, LT, GTE, LTE, BETWEEN:
        return true
    default:
        return false
    }
}

// BuildConditionExpression converts operator to DynamoDB filter expression
func BuildConditionExpression(field string, op OperatorType, values []interface{}) (expression.ConditionBuilder, error) {
    if !ValidateValues(op, values) {
        return expression.ConditionBuilder{}, fmt.Errorf("invalid number of values for operator %s", op)
    }
    
    handler, ok := conditionOperatorHandlers[op]
    if !ok {
        return expression.ConditionBuilder{}, fmt.Errorf("unsupported operator %s for filter conditions", op)
    }
    
    fieldExpr := expression.Name(field)
    result := handler(fieldExpr, values)
    
    return result, nil
}

// BuildKeyConditionExpression converts operator to DynamoDB key condition
func BuildKeyConditionExpression(field string, op OperatorType, values []interface{}) (expression.KeyConditionBuilder, error) {
    if !ValidateValues(op, values) {
        return expression.KeyConditionBuilder{}, fmt.Errorf("invalid number of values for operator %s", op)
    }
    
    // Key conditions have limited operators
    if !IsKeyConditionOperator(op) {
        return expression.KeyConditionBuilder{}, fmt.Errorf("operator %s not supported for key conditions", op)
    }
    
    handler, ok := keyOperatorHandlers[op]
    if !ok {
        return expression.KeyConditionBuilder{}, fmt.Errorf("unsupported operator %s for key conditions", op)
    }
    
    fieldExpr := expression.Key(field)
    result := handler(fieldExpr, values)
    
    return result, nil
}
`
