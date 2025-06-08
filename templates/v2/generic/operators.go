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

// ExpressionScope defines whether this is a key condition or filter condition
type ExpressionScope string

const (
    ScopeKey        ExpressionScope = "KEY"
    ScopeCondition  ExpressionScope = "CONDITION"
)

// Condition represents a single query or filter condition
type Condition struct {
    Field     string
    Operator  OperatorType
    Values    []interface{}
    Scope     ExpressionScope
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

// Pre-allocated map to avoid allocations on every check
var keyOperators = map[OperatorType]bool{
    EQ:      true,
    GT:      true,
    LT:      true,
    GTE:     true,
    LTE:     true,
    BETWEEN: true,
}

// isKeyOperator checks if operator can be used in key conditions
func isKeyOperator(op OperatorType) bool {
    return keyOperators[op]
}

// validateExpressionInput performs common validation for both key and condition expressions
func validateExpressionInput(field string, op OperatorType, values []interface{}, scope ExpressionScope) (string, error) {
    // Common validation for both key and condition expressions
    if !ValidateValues(op, values) {
        return "", fmt.Errorf("invalid number of values for operator %s", op)
    }
    
    // Look up field type automatically from schema
    fieldInfo, exists := TableSchema.FieldsMap[field]
    if !exists {
        return "", fmt.Errorf("field %s not found in schema", field)
    }
    
    // Validate operator compatibility with DynamoDB field type
    if !ValidateOperator(fieldInfo.DynamoType, op) {
        return "", fmt.Errorf("operator %s not supported for DynamoDB type %s", op, fieldInfo.DynamoType)
    }
    
    // Key expressions have limited operators
    if scope == ScopeKey {
        if !fieldInfo.IsKey {
            return "", fmt.Errorf("field %s is not a key field", field)
        }
        if !isKeyOperator(op) {
            return "", fmt.Errorf("operator %s not supported for key conditions", op)
        }
    }
    
    return fieldInfo.DynamoType, nil
}

// buildKeyConditionInternal builds key condition expressions with proper validation
func buildKeyConditionInternal(field string, op OperatorType, values []interface{}) (expression.KeyConditionBuilder, error) {
    _, err := validateExpressionInput(field, op, values, ScopeKey)
    if err != nil {
        return expression.KeyConditionBuilder{}, err
    }
    
    return buildKeyExpressionByOperator(field, op, values)
}

// buildConditionInternal builds filter condition expressions with proper validation  
func buildConditionInternal(field string, op OperatorType, values []interface{}) (expression.ConditionBuilder, error) {
    _, err := validateExpressionInput(field, op, values, ScopeCondition)
    if err != nil {
        return expression.ConditionBuilder{}, err
    }
    
    return buildFilterExpressionByOperator(field, op, values)
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
// Automatically looks up field type from TableSchema
func BuildConditionExpression(field string, op OperatorType, values []interface{}) (expression.ConditionBuilder, error) {
    return buildConditionInternal(field, op, values)
}

// BuildKeyConditionExpression converts operator to DynamoDB key condition
// Automatically looks up field type from TableSchema
func BuildKeyConditionExpression(field string, op OperatorType, values []interface{}) (expression.KeyConditionBuilder, error) {
    return buildKeyConditionInternal(field, op, values)
}

// BuildConditionExpressionWithType converts operator to DynamoDB filter expression with explicit type
// Use this for fields not in TableSchema or to override type validation
func BuildConditionExpressionWithType(field string, op OperatorType, values []interface{}, dynamoType string) (expression.ConditionBuilder, error) {
    // Manual validation when type is provided explicitly
    if !ValidateValues(op, values) {
        return expression.ConditionBuilder{}, fmt.Errorf("invalid number of values for operator %s", op)
    }
    
    if !ValidateOperator(dynamoType, op) {
        return expression.ConditionBuilder{}, fmt.Errorf("operator %s not supported for DynamoDB type %s", op, dynamoType)
    }
    
    return buildFilterExpressionByOperator(field, op, values)
}

// BuildKeyConditionExpressionWithType converts operator to DynamoDB key condition with explicit type
// Use this for fields not in TableSchema or to override type validation
func BuildKeyConditionExpressionWithType(field string, op OperatorType, values []interface{}, dynamoType string) (expression.KeyConditionBuilder, error) {
    // Manual validation when type is provided explicitly
    if !ValidateValues(op, values) {
        return expression.KeyConditionBuilder{}, fmt.Errorf("invalid number of values for operator %s", op)
    }
    
    if !ValidateOperator(dynamoType, op) {
        return expression.KeyConditionBuilder{}, fmt.Errorf("operator %s not supported for DynamoDB type %s", op, dynamoType)
    }
    
    if !isKeyOperator(op) {
        return expression.KeyConditionBuilder{}, fmt.Errorf("operator %s not supported for key conditions", op)
    }
    
    return buildKeyExpressionByOperator(field, op, values)
}
`
