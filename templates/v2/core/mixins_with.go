package core

// FilterMixinSugarTemplate provides convenience Filter methods (only for ALL mode)
const FilterMixinSugarTemplate = `
// CONVENIENCE METHODS - Only available in ALL mode

// FilterEQ adds equality filter condition.
func (fm *FilterMixin) FilterEQ(field string, value any) {
    fm.Filter(field, EQ, value)
}

// FilterContains adds contains filter for strings or sets.
func (fm *FilterMixin) FilterContains(field string, value any) {
    fm.Filter(field, CONTAINS, value)
}

// FilterNotContains adds not contains filter for strings or sets.
func (fm *FilterMixin) FilterNotContains(field string, value any) {
    fm.Filter(field, NOT_CONTAINS, value)
}

// FilterBeginsWith adds begins_with filter for strings.
func (fm *FilterMixin) FilterBeginsWith(field string, value any) {
    fm.Filter(field, BEGINS_WITH, value)
}

// FilterBetween adds range filter for comparable values.
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
`
