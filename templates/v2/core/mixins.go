package core

// MixinsTemplate defines common components for QueryBuilder and ScanBuilder
const MixinsTemplate = `
// FilterMixin provides common filtering logic for Query and Scan operations.
// Supports all DynamoDB filter operators with type validation.
type FilterMixin struct {
    FilterConditions  []expression.ConditionBuilder
    UsedKeys          map[string]bool
    Attributes        map[string]any
}

// NewFilterMixin creates a new FilterMixin instance with initialized maps.
func NewFilterMixin() FilterMixin {
    return FilterMixin{
        FilterConditions: make([]expression.ConditionBuilder, 0),
        UsedKeys:         make(map[string]bool),
        Attributes:       make(map[string]any),
    }
}

// Filter adds a filter condition using the universal operator system.
// Validates operator compatibility and value types before adding.
func (fm *FilterMixin) Filter(field string, op OperatorType, values ...any) {
    if !ValidateValues(op, values) {
        return
    }
    if !ValidateOperator(field, op) {
        return
    }

    filterCond, err := BuildConditionExpression(field, op, values)
    if err != nil {
        return
    }

    fm.FilterConditions = append(fm.FilterConditions, filterCond)
    fm.UsedKeys[field] = true

    if op == EQ && len(values) == 1 {
        fm.Attributes[field] = values[0]
    }
}

// PaginationMixin provides pagination support for Query and Scan operations.
type PaginationMixin struct {
    LimitValue        *int
    ExclusiveStartKey map[string]types.AttributeValue
}

// NewPaginationMixin creates a new PaginationMixin instance.
func NewPaginationMixin() PaginationMixin {
    return PaginationMixin{}
}

// Limit sets the maximum number of items to return in one request.
func (pm *PaginationMixin) Limit(limit int) {
    pm.LimitValue = &limit
}

// StartFrom sets the exclusive start key for pagination.
// Use LastEvaluatedKey from previous response for next page.
func (pm *PaginationMixin) StartFrom(lastEvaluatedKey map[string]types.AttributeValue) {
    pm.ExclusiveStartKey = lastEvaluatedKey
}

// KeyConditionMixin provides key condition logic for Query operations only.
// Supports partition key and sort key conditions with automatic index selection.
type KeyConditionMixin struct {
    KeyConditions    map[string]expression.KeyConditionBuilder
    SortDescending   bool
    PreferredSortKey string
}

// NewKeyConditionMixin creates a new KeyConditionMixin instance.
func NewKeyConditionMixin() KeyConditionMixin {
    return KeyConditionMixin{
        KeyConditions: make(map[string]expression.KeyConditionBuilder),
    }
}

// With adds a key condition using the universal operator system.
// Only valid for partition and sort key attributes.
func (kcm *KeyConditionMixin) With(field string, op OperatorType, values ...any) {
    if !ValidateValues(op, values) {
        return
    }
    fieldInfo, exists := TableSchema.FieldsMap[field]
    if !exists {
        return
    }
    if !fieldInfo.IsKey {
        return
    }
    if !ValidateOperator(field, op) {
        return
    }

    keyCond, err := BuildKeyConditionExpression(field, op, values)
    if err != nil {
        return
    }
    kcm.KeyConditions[field] = keyCond
}

// WithPreferredSortKey sets preferred sort key for index selection.
// Useful when multiple indexes match the query pattern.
func (kcm *KeyConditionMixin) WithPreferredSortKey(key string) {
    kcm.PreferredSortKey = key
}

// OrderByDesc sets descending sort order for results.
// Only affects sort key ordering, not filter results.
func (kcm *KeyConditionMixin) OrderByDesc() {
    kcm.SortDescending = true
}

// OrderByAsc sets ascending sort order for results (default).
func (kcm *KeyConditionMixin) OrderByAsc() {
    kcm.SortDescending = false
}
`
