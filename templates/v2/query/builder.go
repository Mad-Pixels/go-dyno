package query

// QueryBuilderTemplate with mixins
const QueryBuilderTemplate = `
// QueryBuilder provides a fluent interface for building DynamoDB queries
type QueryBuilder struct {
    FilterMixin
    PaginationMixin
    KeyConditionMixin
    IndexName string
}

// NewQueryBuilder creates a new QueryBuilder instance
func NewQueryBuilder() *QueryBuilder {
    return &QueryBuilder{
        FilterMixin:       NewFilterMixin(),
        PaginationMixin:   NewPaginationMixin(),
        KeyConditionMixin: NewKeyConditionMixin(),
    }
}

// Filter adds a filter condition and returns QueryBuilder for chaining
func (qb *QueryBuilder) Filter(field string, op OperatorType, values ...any) *QueryBuilder {
    qb.FilterMixin.Filter(field, op, values...)
    return qb
}

// FilterEQ adds equality filter and returns QueryBuilder for chaining
func (qb *QueryBuilder) FilterEQ(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterEQ(field, value)
    return qb
}

// FilterContains adds contains filter and returns QueryBuilder for chaining
func (qb *QueryBuilder) FilterContains(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterContains(field, value)
    return qb
}

// FilterNotContains adds not contains filter and returns QueryBuilder for chaining
func (qb *QueryBuilder) FilterNotContains(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterNotContains(field, value)
    return qb
}

// FilterBeginsWith adds begins_with filter and returns QueryBuilder for chaining
func (qb *QueryBuilder) FilterBeginsWith(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterBeginsWith(field, value)
    return qb
}

// FilterBetween adds range filter and returns QueryBuilder for chaining
func (qb *QueryBuilder) FilterBetween(field string, start, end any) *QueryBuilder {
    qb.FilterMixin.FilterBetween(field, start, end)
    return qb
}

// FilterGT adds greater than filter and returns QueryBuilder for chaining
func (qb *QueryBuilder) FilterGT(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterGT(field, value)
    return qb
}

// FilterLT adds less than filter and returns QueryBuilder for chaining
func (qb *QueryBuilder) FilterLT(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterLT(field, value)
    return qb
}

// FilterGTE adds greater than or equal filter and returns QueryBuilder for chaining
func (qb *QueryBuilder) FilterGTE(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterGTE(field, value)
    return qb
}

// FilterLTE adds less than or equal filter and returns QueryBuilder for chaining
func (qb *QueryBuilder) FilterLTE(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterLTE(field, value)
    return qb
}

// FilterExists adds attribute exists filter and returns QueryBuilder for chaining
func (qb *QueryBuilder) FilterExists(field string) *QueryBuilder {
    qb.FilterMixin.FilterExists(field)
    return qb
}

// FilterNotExists adds attribute not exists filter and returns QueryBuilder for chaining
func (qb *QueryBuilder) FilterNotExists(field string) *QueryBuilder {
    qb.FilterMixin.FilterNotExists(field)
    return qb
}

// FilterNE adds not equal filter and returns QueryBuilder for chaining
func (qb *QueryBuilder) FilterNE(field string, value any) *QueryBuilder {
    qb.FilterMixin.FilterNE(field, value)
    return qb
}

// FilterIn adds IN filter and returns QueryBuilder for chaining
func (qb *QueryBuilder) FilterIn(field string, values ...any) *QueryBuilder {
    qb.FilterMixin.FilterIn(field, values...)
    return qb
}

// FilterNotIn adds NOT_IN filter and returns QueryBuilder for chaining
func (qb *QueryBuilder) FilterNotIn(field string, values ...any) *QueryBuilder {
    qb.FilterMixin.FilterNotIn(field, values...)
    return qb
}

// With adds key condition and returns QueryBuilder for chaining
func (qb *QueryBuilder) With(field string, op OperatorType, values ...any) *QueryBuilder {
    qb.KeyConditionMixin.With(field, op, values...)
    if op == EQ && len(values) == 1 {
        qb.Attributes[field] = values[0]
        qb.UsedKeys[field] = true
    }
    return qb
}

// WithEQ adds equality key condition and returns QueryBuilder for chaining
func (qb *QueryBuilder) WithEQ(field string, value any) *QueryBuilder {
    fmt.Printf("DEBUG WithEQ: field='%s', value='%v'\n", field, value)
    qb.KeyConditionMixin.WithEQ(field, value)
    qb.Attributes[field] = value
    qb.UsedKeys[field] = true
    fmt.Printf("DEBUG WithEQ AFTER: UsedKeys[%s]=%v, Attributes[%s]=%v\n", 
        field, qb.UsedKeys[field], field, qb.Attributes[field])
    return qb
}

// WithBetween adds range key condition and returns QueryBuilder for chaining
func (qb *QueryBuilder) WithBetween(field string, start, end any) *QueryBuilder {
    qb.KeyConditionMixin.WithBetween(field, start, end)
    qb.Attributes[field+"_start"] = start
    qb.Attributes[field+"_end"] = end
    qb.UsedKeys[field] = true
    return qb
}

// WithGT adds greater than key condition and returns QueryBuilder for chaining
func (qb *QueryBuilder) WithGT(field string, value any) *QueryBuilder {
    qb.KeyConditionMixin.WithGT(field, value)
    qb.Attributes[field] = value 
    qb.UsedKeys[field] = true
    return qb
}

// WithGTE adds greater than or equal key condition and returns QueryBuilder for chaining
func (qb *QueryBuilder) WithGTE(field string, value any) *QueryBuilder {
    qb.KeyConditionMixin.WithGTE(field, value)
    qb.Attributes[field] = value
    qb.UsedKeys[field] = true
    return qb
}

// WithLT adds less than key condition and returns QueryBuilder for chaining
func (qb *QueryBuilder) WithLT(field string, value any) *QueryBuilder {
    qb.KeyConditionMixin.WithLT(field, value)
    qb.Attributes[field] = value
    qb.UsedKeys[field] = true
    return qb
}

// WithLTE adds less than or equal key condition and returns QueryBuilder for chaining
func (qb *QueryBuilder) WithLTE(field string, value any) *QueryBuilder {
    qb.KeyConditionMixin.WithLTE(field, value)
    qb.Attributes[field] = value
    qb.UsedKeys[field] = true
    return qb
}

// WithPreferredSortKey sets the preferred sort key and returns QueryBuilder for chaining
func (qb *QueryBuilder) WithPreferredSortKey(key string) *QueryBuilder {
    qb.KeyConditionMixin.WithPreferredSortKey(key)
    return qb
}

// OrderByDesc sets descending sort order and returns QueryBuilder for chaining
func (qb *QueryBuilder) OrderByDesc() *QueryBuilder {
    qb.KeyConditionMixin.OrderByDesc()
    return qb
}

// OrderByAsc sets ascending sort order and returns QueryBuilder for chaining
func (qb *QueryBuilder) OrderByAsc() *QueryBuilder {
    qb.KeyConditionMixin.OrderByAsc()
    return qb
}

// Limit sets the maximum number of items and returns QueryBuilder for chaining
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
    qb.PaginationMixin.Limit(limit)
    return qb
}

// StartFrom sets the exclusive start key and returns QueryBuilder for chaining
func (qb *QueryBuilder) StartFrom(lastEvaluatedKey map[string]types.AttributeValue) *QueryBuilder {
    qb.PaginationMixin.StartFrom(lastEvaluatedKey)
    return qb
}

{{range .SecondaryIndexes}}
{{if gt (len .HashKeyParts) 0}}
{{- $hasNonConstant := false -}}
{{- range .HashKeyParts -}}{{- if not .IsConstant -}}{{- $hasNonConstant = true -}}{{- end -}}{{- end -}}
{{- if $hasNonConstant}}
// With{{ToUpperCamelCase .Name}}HashKey sets composite hash key for {{.Name}} index
func (qb *QueryBuilder) With{{ToUpperCamelCase .Name}}HashKey({{range $i, $part := .HashKeyParts}}{{if not $part.IsConstant}}{{if $i}}, {{end}}{{$part.Value | ToLowerCamelCase}} {{ToGolangAttrType $part.Value $.AllAttributes}}{{end}}{{end}}) *QueryBuilder {
    {{range .HashKeyParts}}{{if not .IsConstant}}
    qb.Attributes["{{.Value}}"] = {{.Value | ToLowerCamelCase}}
    qb.UsedKeys["{{.Value}}"] = true
    {{end}}{{end}}
    
    compositeValue := qb.buildCompositeKeyValue([]CompositeKeyPart{
        {{range .HashKeyParts}}
        {{if .IsConstant}}
        {IsConstant: true, Value: "{{.Value}}"},
        {{else}}
        {IsConstant: false, Value: "{{.Value}}"},
        {{end}}
        {{end}}
    })
    
    qb.Attributes["{{.HashKey}}"] = compositeValue
    qb.UsedKeys["{{.HashKey}}"] = true
    qb.KeyConditions["{{.HashKey}}"] = expression.Key("{{.HashKey}}").Equal(expression.Value(compositeValue))
    return qb
}
{{end}}
{{else if .HashKey}}
// With{{ToUpperCamelCase .Name}}HashKey sets hash key for {{.Name}} index
func (qb *QueryBuilder) With{{ToUpperCamelCase .Name}}HashKey({{.HashKey | ToLowerCamelCase}} {{ToGolangAttrType .HashKey $.AllAttributes}}) *QueryBuilder {
    qb.Attributes["{{.HashKey}}"] = {{.HashKey | ToLowerCamelCase}}
    qb.UsedKeys["{{.HashKey}}"] = true
    qb.KeyConditions["{{.HashKey}}"] = expression.Key("{{.HashKey}}").Equal(expression.Value({{.HashKey | ToLowerCamelCase}}))
    return qb
}
{{end}}
{{end}}

{{range .SecondaryIndexes}}
{{if gt (len .RangeKeyParts) 0}}
{{- $hasNonConstant := false -}}
{{- range .RangeKeyParts -}}{{- if not .IsConstant -}}{{- $hasNonConstant = true -}}{{- end -}}{{- end -}}
{{- if $hasNonConstant}}
// With{{ToUpperCamelCase .Name}}RangeKey sets composite range key for {{.Name}} index
func (qb *QueryBuilder) With{{ToUpperCamelCase .Name}}RangeKey({{range $i, $part := .RangeKeyParts}}{{if not $part.IsConstant}}{{if $i}}, {{end}}{{$part.Value | ToLowerCamelCase}} {{ToGolangAttrType $part.Value $.AllAttributes}}{{end}}{{end}}) *QueryBuilder {
    {{range .RangeKeyParts}}{{if not .IsConstant}}
    qb.Attributes["{{.Value}}"] = {{.Value | ToLowerCamelCase}}
    qb.UsedKeys["{{.Value}}"] = true
    {{end}}{{end}}
    return qb
}
{{end}}
{{else if .RangeKey}}
// With{{ToUpperCamelCase .Name}}RangeKey sets range key for {{.Name}} index
func (qb *QueryBuilder) With{{ToUpperCamelCase .Name}}RangeKey({{.RangeKey | ToLowerCamelCase}} {{ToGolangAttrType .RangeKey $.AllAttributes}}) *QueryBuilder {
    qb.Attributes["{{.RangeKey}}"] = {{.RangeKey | ToLowerCamelCase}}
    qb.UsedKeys["{{.RangeKey}}"] = true
    qb.KeyConditions["{{.RangeKey}}"] = expression.Key("{{.RangeKey}}").Equal(expression.Value({{.RangeKey | ToLowerCamelCase}}))
    return qb
}
{{end}}
{{end}}
`
