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

// ExpressionType defines whether this is a key condition or filter condition
type ExpressionType string

const (
    KeyExpression    ExpressionType = "KEY"
    FilterExpression ExpressionType = "FILTER"
)

// Condition represents a single query or filter condition
type Condition struct {
    Field     string
    Operator  OperatorType
    Values    []interface{}
    Type      ExpressionType
}

// ValidateOperator checks if operator is valid for the given DynamoDB type
func ValidateOperator(dynamoType string, op OperatorType) bool {
    switch dynamoType {
    case "S": // String
        return op == EQ || op == NE || op == GT || op == LT || op == GTE || op == LTE || 
               op == BETWEEN || op == CONTAINS || op == NOT_CONTAINS || op == BEGINS_WITH || op == IN || op == NOT_IN ||
               op == EXISTS || op == NOT_EXISTS
               
    case "N": // Number  
        return op == EQ || op == NE || op == GT || op == LT || op == GTE || op == LTE || 
               op == BETWEEN || op == IN || op == NOT_IN ||
               op == EXISTS || op == NOT_EXISTS
               
    case "BOOL": // Boolean
        return op == EQ || op == NE || op == EXISTS || op == NOT_EXISTS
        
    case "SS": // String Set - only CONTAINS/NOT_CONTAINS, not IN/NOT_IN
        return op == CONTAINS || op == NOT_CONTAINS || op == EXISTS || op == NOT_EXISTS
        
    case "NS": // Number Set - only CONTAINS/NOT_CONTAINS, not IN/NOT_IN
        return op == CONTAINS || op == NOT_CONTAINS || op == EXISTS || op == NOT_EXISTS
        
    case "BS": // Binary Set (rare)
        return op == CONTAINS || op == NOT_CONTAINS || op == EXISTS || op == NOT_EXISTS
        
    case "L": // List
        return op == EXISTS || op == NOT_EXISTS
        
    case "M": // Map
        return op == EXISTS || op == NOT_EXISTS
        
    case "NULL": // Null
        return op == EXISTS || op == NOT_EXISTS
        
    default:
        // For unknown types allow only basic operations
        return op == EQ || op == NE || op == EXISTS || op == NOT_EXISTS
    }
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

// getKeyOperators returns operators that can be used in key conditions
func getKeyOperators() map[OperatorType]bool {
    return map[OperatorType]bool{
        EQ:      true,
        GT:      true,
        LT:      true,
        GTE:     true,
        LTE:     true,
        BETWEEN: true,
    }
}

// isKeyOperator checks if operator can be used in key conditions
func isKeyOperator(op OperatorType) bool {
    keyOps := getKeyOperators()
    return keyOps[op]
}

// buildExpressionInternal unified internal function for building expressions
func buildExpressionInternal(field string, op OperatorType, values []interface{}, exprType ExpressionType) (interface{}, error) {
    // Common validation for both key and filter expressions
    if !ValidateValues(op, values) {
        return nil, fmt.Errorf("invalid number of values for operator %s", op)
    }
    
    // Key expressions have limited operators
    if exprType == KeyExpression && !isKeyOperator(op) {
        return nil, fmt.Errorf("operator %s not supported for key conditions", op)
    }
    
    if exprType == KeyExpression {
        return buildKeyExpressionByOperator(field, op, values)
    } else {
        return buildFilterExpressionByOperator(field, op, values)
    }
}

// buildKeyExpressionByOperator builds key condition expressions
func buildKeyExpressionByOperator(field string, op OperatorType, values []interface{}) (expression.KeyConditionBuilder, error) {
    fieldExpr := expression.Key(field)
    
    switch op {
    case EQ:
        return fieldExpr.Equal(expression.Value(values[0])), nil
    case GT:
        return fieldExpr.GreaterThan(expression.Value(values[0])), nil
    case LT:
        return fieldExpr.LessThan(expression.Value(values[0])), nil
    case GTE:
        return fieldExpr.GreaterThanEqual(expression.Value(values[0])), nil
    case LTE:
        return fieldExpr.LessThanEqual(expression.Value(values[0])), nil
    case BETWEEN:
        return fieldExpr.Between(expression.Value(values[0]), expression.Value(values[1])), nil
    default:
        return expression.KeyConditionBuilder{}, fmt.Errorf("unsupported key operator %s", op)
    }
}

// buildFilterExpressionByOperator builds filter condition expressions
func buildFilterExpressionByOperator(field string, op OperatorType, values []interface{}) (expression.ConditionBuilder, error) {
    fieldExpr := expression.Name(field)
    
    switch op {
    case EQ:
        return fieldExpr.Equal(expression.Value(values[0])), nil
    case NE:
        return fieldExpr.NotEqual(expression.Value(values[0])), nil
    case GT:
        return fieldExpr.GreaterThan(expression.Value(values[0])), nil
    case LT:
        return fieldExpr.LessThan(expression.Value(values[0])), nil
    case GTE:
        return fieldExpr.GreaterThanEqual(expression.Value(values[0])), nil
    case LTE:
        return fieldExpr.LessThanEqual(expression.Value(values[0])), nil
    case BETWEEN:
        return fieldExpr.Between(expression.Value(values[0]), expression.Value(values[1])), nil
    case CONTAINS:
        return fieldExpr.Contains(fmt.Sprintf("%v", values[0])), nil
    case NOT_CONTAINS:
        return expression.Not(fieldExpr.Contains(fmt.Sprintf("%v", values[0]))), nil
    case BEGINS_WITH:
        return fieldExpr.BeginsWith(fmt.Sprintf("%v", values[0])), nil
    case IN:
        if len(values) == 0 {
            return expression.AttributeNotExists(fieldExpr), nil
        }
        if len(values) == 1 {
            return fieldExpr.Equal(expression.Value(values[0])), nil
        }
        operands := make([]expression.OperandBuilder, len(values))
        for i, v := range values {
            operands[i] = expression.Value(v)
        }
        return fieldExpr.In(operands[0], operands[1:]...), nil
    case NOT_IN:
        if len(values) == 0 {
            return expression.AttributeExists(fieldExpr), nil
        }
        if len(values) == 1 {
            return fieldExpr.NotEqual(expression.Value(values[0])), nil
        }
        operands := make([]expression.OperandBuilder, len(values))
        for i, v := range values {
            operands[i] = expression.Value(v)
        }
        return expression.Not(fieldExpr.In(operands[0], operands[1:]...)), nil
    case EXISTS:
        return expression.AttributeExists(fieldExpr), nil
    case NOT_EXISTS:
        return expression.AttributeNotExists(fieldExpr), nil
    default:
        return expression.ConditionBuilder{}, fmt.Errorf("unsupported filter operator %s", op)
    }
}

// BuildConditionExpression converts operator to DynamoDB filter expression
func BuildConditionExpression(field string, op OperatorType, values []interface{}) (expression.ConditionBuilder, error) {
    result, err := buildExpressionInternal(field, op, values, FilterExpression)
    if err != nil {
        return expression.ConditionBuilder{}, err
    }
    
    conditionBuilder, ok := result.(expression.ConditionBuilder)
    if !ok {
        return expression.ConditionBuilder{}, fmt.Errorf("internal error: expected ConditionBuilder")
    }
    
    return conditionBuilder, nil
}

// BuildKeyConditionExpression converts operator to DynamoDB key condition
func BuildKeyConditionExpression(field string, op OperatorType, values []interface{}) (expression.KeyConditionBuilder, error) {
    result, err := buildExpressionInternal(field, op, values, KeyExpression)
    if err != nil {
        return expression.KeyConditionBuilder{}, err
    }
    
    keyConditionBuilder, ok := result.(expression.KeyConditionBuilder)
    if !ok {
        return expression.KeyConditionBuilder{}, fmt.Errorf("internal error: expected KeyConditionBuilder")
    }
    
    return keyConditionBuilder, nil
}
`
