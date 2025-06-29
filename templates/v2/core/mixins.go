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

// FilterEQ adds equality filter condition.
// Example: .FilterEQ("status", "active")
func (fm *FilterMixin) FilterEQ(field string, value any) {
    fm.Filter(field, EQ, value)
}

// FilterContains adds contains filter for strings or sets.
// Example: .FilterContains("tags", "important")
func (fm *FilterMixin) FilterContains(field string, value any) {
    fm.Filter(field, CONTAINS, value)
}

// FilterNotContains adds not contains filter for strings or sets.
func (fm *FilterMixin) FilterNotContains(field string, value any) {
    fm.Filter(field, NOT_CONTAINS, value)
}

// FilterBeginsWith adds begins_with filter for strings.
// Example: .FilterBeginsWith("email", "admin@")
func (fm *FilterMixin) FilterBeginsWith(field string, value any) {
    fm.Filter(field, BEGINS_WITH, value)
}

// FilterBetween adds range filter for comparable values.
// Example: .FilterBetween("price", 10, 100)
func (fm *FilterMixin) FilterBetween(field string, start, end any) {
    fm.Filter(field, BETWEEN, start, end)
}

// FilterGT adds greater than filter.
func (fm *FilterMixin) FilterGT(field string, value any) {
    fm.Filter(field, GT, value)
}

// FilterLT adds less than filter.
func (fm *FilterMixin) FilterLT(field string, value any) {
    fm.Filter(field, LT, value)
}

// FilterGTE adds greater than or equal filter.
func (fm *FilterMixin) FilterGTE(field string, value any) {
    fm.Filter(field, GTE, value)
}

// FilterLTE adds less than or equal filter.
func (fm *FilterMixin) FilterLTE(field string, value any) {
    fm.Filter(field, LTE, value)
}

// FilterExists checks if attribute exists.
// Example: .FilterExists("optional_field")
func (fm *FilterMixin) FilterExists(field string) {
    fm.Filter(field, EXISTS)
}

// FilterNotExists checks if attribute does not exist.
func (fm *FilterMixin) FilterNotExists(field string) {
    fm.Filter(field, NOT_EXISTS)
}

// FilterNE adds not equal filter.
func (fm *FilterMixin) FilterNE(field string, value any) {
    fm.Filter(field, NE, value)
}

// FilterIn adds IN filter for scalar values.
// For DynamoDB Sets (SS/NS), use FilterContains instead.
// Example: .FilterIn("category", "books", "electronics")
func (fm *FilterMixin) FilterIn(field string, values ...any) {
    if len(values) == 0 {
        return
    }
    fm.Filter(field, IN, values...)
}

// FilterNotIn adds NOT_IN filter for scalar values.
// For DynamoDB Sets (SS/NS), use FilterNotContains instead.
func (fm *FilterMixin) FilterNotIn(field string, values ...any) {
    if len(values) == 0 {
        return
    }
    fm.Filter(field, NOT_IN, values...)
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
// Example: .Limit(25)
func (pm *PaginationMixin) Limit(limit int) {
    pm.LimitValue = &limit
}

// StartFrom sets the exclusive start key for pagination.
// Use LastEvaluatedKey from previous response for next page.
// Example: .StartFrom(previousResponse.LastEvaluatedKey)
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

// WithEQ adds equality key condition.
// Required for partition key, optional for sort key.
// Example: .WithEQ("user_id", "123")
func (kcm *KeyConditionMixin) WithEQ(field string, value any) {
    kcm.With(field, EQ, value)
}

// WithBetween adds range key condition for sort keys.
// Example: .WithBetween("created_at", start_time, end_time)
func (kcm *KeyConditionMixin) WithBetween(field string, start, end any) {
    kcm.With(field, BETWEEN, start, end)
}

// WithGT adds greater than key condition for sort keys.
func (kcm *KeyConditionMixin) WithGT(field string, value any) {
    kcm.With(field, GT, value)
}

// WithGTE adds greater than or equal key condition for sort keys.
func (kcm *KeyConditionMixin) WithGTE(field string, value any) {
    kcm.With(field, GTE, value)
}

// WithLT adds less than key condition for sort keys.
func (kcm *KeyConditionMixin) WithLT(field string, value any) {
    kcm.With(field, LT, value)
}

// WithLTE adds less than or equal key condition for sort keys.
func (kcm *KeyConditionMixin) WithLTE(field string, value any) {
    kcm.With(field, LTE, value)
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
