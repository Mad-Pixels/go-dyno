package scan

// ScanBuilderFilterTemplate provides Filter methods for scan conditions
const ScanBuilderFilterTemplate = `
// Filter adds a filter condition and returns ScanBuilder for method chaining.
// Wraps FilterMixin.Filter with fluent interface support.
func (sb *ScanBuilder) Filter(field string, op OperatorType, values ...any) *ScanBuilder {
    sb.FilterMixin.Filter(field, op, values...)
    return sb
}
`

// ScanBuilderFilterSugarTemplate provides convenience Filter methods (only for ALL mode)
const ScanBuilderFilterSugarTemplate = `
// CONVENIENCE METHODS - Only available in ALL mode

// FilterEQ adds equality filter and returns ScanBuilder for method chaining.
func (sb *ScanBuilder) FilterEQ(field string, value any) *ScanBuilder {
    sb.FilterMixin.FilterEQ(field, value)
    return sb
}

// FilterContains adds contains filter and returns ScanBuilder for method chaining.
// Works with String attributes (substring) and Set attributes (membership).
func (sb *ScanBuilder) FilterContains(field string, value any) *ScanBuilder {
    sb.FilterMixin.FilterContains(field, value)
    return sb
}

// FilterNotContains adds not contains filter and returns ScanBuilder for method chaining.
// Opposite of FilterContains for exclusion filtering.
func (sb *ScanBuilder) FilterNotContains(field string, value any) *ScanBuilder {
    sb.FilterMixin.FilterNotContains(field, value)
    return sb
}

// FilterBeginsWith adds begins_with filter and returns ScanBuilder for method chaining.
// Only works with String attributes for prefix matching.
func (sb *ScanBuilder) FilterBeginsWith(field string, value any) *ScanBuilder {
    sb.FilterMixin.FilterBeginsWith(field, value)
    return sb
}

// FilterBetween adds range filter and returns ScanBuilder for method chaining.
// Works with comparable types for inclusive range filtering.
func (sb *ScanBuilder) FilterBetween(field string, start, end any) *ScanBuilder {
    sb.FilterMixin.FilterBetween(field, start, end)
    return sb
}

// FilterGT adds greater than filter and returns ScanBuilder for method chaining.
func (sb *ScanBuilder) FilterGT(field string, value any) *ScanBuilder {
    sb.FilterMixin.FilterGT(field, value)
    return sb
}

// FilterLT adds less than filter and returns ScanBuilder for method chaining.
func (sb *ScanBuilder) FilterLT(field string, value any) *ScanBuilder {
    sb.FilterMixin.FilterLT(field, value)
    return sb
}

// FilterGTE adds greater than or equal filter and returns ScanBuilder for method chaining.
func (sb *ScanBuilder) FilterGTE(field string, value any) *ScanBuilder {
    sb.FilterMixin.FilterGTE(field, value)
    return sb
}

// FilterLTE adds less than or equal filter and returns ScanBuilder for method chaining.
func (sb *ScanBuilder) FilterLTE(field string, value any) *ScanBuilder {
    sb.FilterMixin.FilterLTE(field, value)
    return sb
}

// FilterExists adds attribute exists filter and returns ScanBuilder for method chaining.
// Checks if the specified attribute exists in the item.
func (sb *ScanBuilder) FilterExists(field string) *ScanBuilder {
    sb.FilterMixin.FilterExists(field)
    return sb
}

// FilterNotExists adds attribute not exists filter and returns ScanBuilder for method chaining.
// Checks if the specified attribute does not exist in the item.
func (sb *ScanBuilder) FilterNotExists(field string) *ScanBuilder {
    sb.FilterMixin.FilterNotExists(field)
    return sb
}

// FilterNE adds not equal filter and returns ScanBuilder for method chaining.
func (sb *ScanBuilder) FilterNE(field string, value any) *ScanBuilder {
    sb.FilterMixin.FilterNE(field, value)
    return sb
}

// FilterIn adds IN filter and returns ScanBuilder for method chaining.
// For scalar values only - use FilterContains for DynamoDB Sets.
func (sb *ScanBuilder) FilterIn(field string, values ...any) *ScanBuilder {
    sb.FilterMixin.FilterIn(field, values...)
    return sb
}

// FilterNotIn adds NOT_IN filter and returns ScanBuilder for method chaining.
// For scalar values only - use FilterNotContains for DynamoDB Sets.
func (sb *ScanBuilder) FilterNotIn(field string, values ...any) *ScanBuilder {
    sb.FilterMixin.FilterNotIn(field, values...)
    return sb
}
`
