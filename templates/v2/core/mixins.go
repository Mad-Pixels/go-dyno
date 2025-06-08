package core

// MixinsTemplate defines common components for QueryBuilder and ScanBuilder
const MixinsTemplate = `
// FilterMixin contains common filtering logic for Query and Scan
type FilterMixin struct {
    FilterConditions  []expression.ConditionBuilder
    UsedKeys          map[string]bool
    Attributes        map[string]interface{}
}

// NewFilterMixin creates a new FilterMixin
func NewFilterMixin() FilterMixin {
    return FilterMixin{
        FilterConditions: make([]expression.ConditionBuilder, 0),
        UsedKeys:         make(map[string]bool),
        Attributes:       make(map[string]interface{}),
    }
}

// Filter adds a filter condition using the universal operator system with proper type validation
func (fm *FilterMixin) Filter(field string, op OperatorType, values ...interface{}) {
    if !ValidateValues(op, values) {
        return
    }

    fieldInfo, exists := TableSchema.FieldsMap[field]
    if !exists {
        return
    }

    // Use type-aware function for better validation
    filterCond, err := BuildConditionExpressionWithType(field, op, values, fieldInfo.DynamoType)
    if err != nil {
        return
    }

    fm.FilterConditions = append(fm.FilterConditions, filterCond)
    fm.UsedKeys[field] = true

    if op == EQ && len(values) == 1 {
        fm.Attributes[field] = values[0]
    }
}

// FilterEQ is a convenience method for equality filters
func (fm *FilterMixin) FilterEQ(field string, value interface{}) {
    fm.Filter(field, EQ, value)
}

// FilterContains is a convenience method for contains filters
func (fm *FilterMixin) FilterContains(field string, value interface{}) {
    fm.Filter(field, CONTAINS, value)
}

// FilterNotContains is a convenience method for not contains filters
func (fm *FilterMixin) FilterNotContains(field string, value interface{}) {
    fm.Filter(field, NOT_CONTAINS, value)
}

// FilterBeginsWith is a convenience method for begins_with filters
func (fm *FilterMixin) FilterBeginsWith(field string, value interface{}) {
    fm.Filter(field, BEGINS_WITH, value)
}

// FilterBetween is a convenience method for range filters
func (fm *FilterMixin) FilterBetween(field string, start, end interface{}) {
    fm.Filter(field, BETWEEN, start, end)
}

// FilterGT is a convenience method for greater than filters
func (fm *FilterMixin) FilterGT(field string, value interface{}) {
    fm.Filter(field, GT, value)
}

// FilterLT is a convenience method for less than filters
func (fm *FilterMixin) FilterLT(field string, value interface{}) {
    fm.Filter(field, LT, value)
}

// FilterGTE is a convenience method for greater than or equal filters
func (fm *FilterMixin) FilterGTE(field string, value interface{}) {
    fm.Filter(field, GTE, value)
}

// FilterLTE is a convenience method for less than or equal filters
func (fm *FilterMixin) FilterLTE(field string, value interface{}) {
    fm.Filter(field, LTE, value)
}

// FilterExists is a convenience method for attribute exists filters
func (fm *FilterMixin) FilterExists(field string) {
    fm.Filter(field, EXISTS)
}

// FilterNotExists is a convenience method for attribute not exists filters
func (fm *FilterMixin) FilterNotExists(field string) {
    fm.Filter(field, NOT_EXISTS)
}

// FilterNE is a convenience method for not equal filters
func (fm *FilterMixin) FilterNE(field string, value interface{}) {
    fm.Filter(field, NE, value)
}

// FilterIn is a convenience method for IN filters (for scalar values only)
// For checking membership in DynamoDB Sets (SS/NS), use FilterContains instead
func (fm *FilterMixin) FilterIn(field string, values ...interface{}) {
    if len(values) == 0 {
        return
    }
    fm.Filter(field, IN, values...)
}

// FilterNotIn is a convenience method for NOT_IN filters (for scalar values only)
// For checking non-membership in DynamoDB Sets (SS/NS), use FilterNotContains instead
func (fm *FilterMixin) FilterNotIn(field string, values ...interface{}) {
    if len(values) == 0 {
        return
    }
    fm.Filter(field, NOT_IN, values...)
}

// PaginationMixin contains common pagination logic
type PaginationMixin struct {
    LimitValue        *int
    ExclusiveStartKey map[string]types.AttributeValue
}

// NewPaginationMixin creates a new PaginationMixin
func NewPaginationMixin() PaginationMixin {
    return PaginationMixin{}
}

// Limit sets the maximum number of items to return
func (pm *PaginationMixin) Limit(limit int) {
    pm.LimitValue = &limit
}

// StartFrom sets the exclusive start key for pagination
func (pm *PaginationMixin) StartFrom(lastEvaluatedKey map[string]types.AttributeValue) {
    pm.ExclusiveStartKey = lastEvaluatedKey
}

// KeyConditionMixin contains logic for key conditions (Query only)
type KeyConditionMixin struct {
    KeyConditions    map[string]expression.KeyConditionBuilder
    SortDescending   bool
    PreferredSortKey string
}

// NewKeyConditionMixin creates a new KeyConditionMixin
func NewKeyConditionMixin() KeyConditionMixin {
    return KeyConditionMixin{
        KeyConditions: make(map[string]expression.KeyConditionBuilder),
    }
}

// With adds a key condition using the universal operator system with proper type validation
func (kcm *KeyConditionMixin) With(field string, op OperatorType, values ...interface{}) {
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

    // Use type-aware function for better validation
    keyCond, err := BuildKeyConditionExpressionWithType(field, op, values, fieldInfo.DynamoType)
    if err != nil {
        return
    }

    kcm.KeyConditions[field] = keyCond
}

// WithEQ is a convenience method for equality key conditions
func (kcm *KeyConditionMixin) WithEQ(field string, value interface{}) {
    kcm.With(field, EQ, value)
}

// WithBetween is a convenience method for range key conditions
func (kcm *KeyConditionMixin) WithBetween(field string, start, end interface{}) {
    kcm.With(field, BETWEEN, start, end)
}

// WithGT is a convenience method for greater than key conditions
func (kcm *KeyConditionMixin) WithGT(field string, value interface{}) {
    kcm.With(field, GT, value)
}

// WithGTE is a convenience method for greater than or equal key conditions
func (kcm *KeyConditionMixin) WithGTE(field string, value interface{}) {
    kcm.With(field, GTE, value)
}

// WithLT is a convenience method for less than key conditions
func (kcm *KeyConditionMixin) WithLT(field string, value interface{}) {
    kcm.With(field, LT, value)
}

// WithLTE is a convenience method for less than or equal key conditions
func (kcm *KeyConditionMixin) WithLTE(field string, value interface{}) {
    kcm.With(field, LTE, value)
}

// WithPreferredSortKey sets the preferred sort key for index selection
func (kcm *KeyConditionMixin) WithPreferredSortKey(key string) {
    kcm.PreferredSortKey = key
}

// OrderByDesc sets descending sort order
func (kcm *KeyConditionMixin) OrderByDesc() {
    kcm.SortDescending = true
}

// OrderByAsc sets ascending sort order
func (kcm *KeyConditionMixin) OrderByAsc() {
    kcm.SortDescending = false
}
`
