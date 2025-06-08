package generic

// OperatorsTemplate generates operator definitions for universal query/filter methods
const OperatorsTemplate = `
// OperatorType defines the type of operation for queries and filters
type OperatorType string

const (
	// Equality and comparison operators
	EQ  OperatorType = "="
	GT  OperatorType = ">"
	LT  OperatorType = "<"
	GTE OperatorType = ">="
	LTE OperatorType = "<="

	// Range operator
	BETWEEN OperatorType = "BETWEEN"

	// String operators
	CONTAINS    OperatorType = "contains"
	BEGINS_WITH OperatorType = "begins_with"

	// Set operators
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

// ValidateOperator checks if operator is valid for the given DynamoDB type
func ValidateOperator(dynamoType string, op OperatorType) bool {
	switch dynamoType {
	case "S": // String
		return op == EQ || op == GT || op == LT || op == GTE || op == LTE || 
			   op == BETWEEN || op == CONTAINS || op == BEGINS_WITH || op == IN
	case "N": // Number  
		return op == EQ || op == GT || op == LT || op == GTE || op == LTE || 
			   op == BETWEEN || op == IN
	case "BOOL": // Boolean
		return op == EQ
	case "SS", "NS": // Sets
		return op == CONTAINS || op == EXISTS || op == NOT_EXISTS
	default:
		return op == EQ || op == EXISTS || op == NOT_EXISTS
	}
}

// ValidateValues checks if the number of values is correct for the operator
func ValidateValues(op OperatorType, values []interface{}) bool {
	switch op {
	case EQ, GT, LT, GTE, LTE, CONTAINS, BEGINS_WITH:
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

// BuildConditionExpression converts operator to DynamoDB expression
func BuildConditionExpression(field string, op OperatorType, values []interface{}) (expression.ConditionBuilder, error) {
	if !ValidateValues(op, values) {
		return expression.ConditionBuilder{}, fmt.Errorf("invalid number of values for operator %s", op)
	}

	fieldExpr := expression.Name(field)
	
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
	case CONTAINS:
		return fieldExpr.Contains(fmt.Sprintf("%v", values[0])), nil
	case BEGINS_WITH:
		return fieldExpr.BeginsWith(fmt.Sprintf("%v", values[0])), nil
	case EXISTS:
		return expression.AttributeExists(fieldExpr), nil
	case NOT_EXISTS:
		return expression.AttributeNotExists(fieldExpr), nil
	default:
		return expression.ConditionBuilder{}, fmt.Errorf("unsupported operator: %s", op)
	}
}

// BuildKeyConditionExpression converts operator to DynamoDB key condition
func BuildKeyConditionExpression(field string, op OperatorType, values []interface{}) (expression.KeyConditionBuilder, error) {
	if !ValidateValues(op, values) {
		return expression.KeyConditionBuilder{}, fmt.Errorf("invalid number of values for operator %s", op)
	}

	// Key conditions have limited operators
	if op != EQ && op != GT && op != LT && op != GTE && op != LTE && op != BETWEEN {
		return expression.KeyConditionBuilder{}, fmt.Errorf("operator %s not supported for key conditions", op)
	}

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
		return expression.KeyConditionBuilder{}, fmt.Errorf("unsupported key operator: %s", op)
	}
}
`
