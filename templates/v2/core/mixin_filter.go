package core

// KeyConditionMixinSugarTemplate provides convenience With methods (only for ALL mode)
const KeyConditionMixinSugarTemplate = `
// CONVENIENCE METHODS - Only available in ALL mode

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
`
