package query

// QueryBuilderFilterTemplate provides Filter methods for query conditions
const QueryBuilderFilterTemplate = `
// Filter adds a filter condition and returns QueryBuilder for method chaining.
// Wraps FilterMixin.Filter with fluent interface support.
func (qb *QueryBuilder) Filter(field string, op OperatorType, values ...any) *QueryBuilder {
    qb.FilterMixin.Filter(field, op, values...)
    return qb
}
`

// QueryBuilderFilterSugarTemplate provides convenience Filter methods (only for ALL mode)
const QueryBuilderFilterSugarTemplate = `
// CONVENIENCE METHODS - Only available in ALL mode

// FilterEQ adds equality filter and returns QueryBuilder for method chaining.
func (qb *QueryBuilder) FilterEQ(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterEQ(field, value)
    return qb
}

// FilterContains adds contains filter and returns QueryBuilder for method chaining.
// Works with String attributes (substring) and Set attributes (membership).
func (qb *QueryBuilder) FilterContains(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterContains(field, value)
    return qb
}

// FilterNotContains adds not contains filter and returns QueryBuilder for method chaining.
// Opposite of FilterContains for exclusion filtering.
func (qb *QueryBuilder) FilterNotContains(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterNotContains(field, value)
    return qb
}

// FilterBeginsWith adds begins_with filter and returns QueryBuilder for method chaining.
// Only works with String attributes for prefix matching.
func (qb *QueryBuilder) FilterBeginsWith(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterBeginsWith(field, value)
    return qb
}

// FilterBetween adds range filter and returns QueryBuilder for method chaining.
// Works with comparable types for inclusive range filtering.
func (qb *QueryBuilder) FilterBetween(field string, start, end any) *QueryBuilder {
    qb.FilterMixin.FilterBetween(field, start, end)
    return qb
}

// FilterGT adds greater than filter and returns QueryBuilder for method chaining.
func (qb *QueryBuilder) FilterGT(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterGT(field, value)
    return qb
}

// FilterLT adds less than filter and returns QueryBuilder for method chaining.
func (qb *QueryBuilder) FilterLT(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterLT(field, value)
    return qb
}

// FilterGTE adds greater than or equal filter and returns QueryBuilder for method chaining.
func (qb *QueryBuilder) FilterGTE(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterGTE(field, value)
    return qb
}

// FilterLTE adds less than or equal filter and returns QueryBuilder for method chaining.
func (qb *QueryBuilder) FilterLTE(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterLTE(field, value)
    return qb
}

// FilterExists adds attribute exists filter and returns QueryBuilder for method chaining.
// Checks if the specified attribute exists in the item.
func (qb *QueryBuilder) FilterExists(field string) *QueryBuilder {
    qb.FilterMixin.FilterExists(field)
    return qb
}

// FilterNotExists adds attribute not exists filter and returns QueryBuilder for method chaining.
// Checks if the specified attribute does not exist in the item.
func (qb *QueryBuilder) FilterNotExists(field string) *QueryBuilder {
    qb.FilterMixin.FilterNotExists(field)
    return qb
}

// FilterNE adds not equal filter and returns QueryBuilder for method chaining.
func (qb *QueryBuilder) FilterNE(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterNE(field, value)
    return qb
}

// FilterIn adds IN filter and returns QueryBuilder for method chaining.
// For scalar values only - use FilterContains for DynamoDB Sets.
func (qb *QueryBuilder) FilterIn(field string, values ...any) *QueryBuilder {
    qb.FilterMixin.FilterIn(field, values...)
    return qb
}

// FilterNotIn adds NOT_IN filter and returns QueryBuilder for method chaining.
// For scalar values only - use FilterNotContains for DynamoDB Sets.
func (qb *QueryBuilder) FilterNotIn(field string, values ...any) *QueryBuilder {
    qb.FilterMixin.FilterNotIn(field, values...)
    return qb
}
`
