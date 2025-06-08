package query

// QueryBuilderTemplate ...
const QueryBuilderTemplate = `
// QueryBuilder provides a fluent interface for building DynamoDB queries
type QueryBuilder struct {
    IndexName         string
    KeyConditions     map[string]expression.KeyConditionBuilder
    FilterConditions  []expression.ConditionBuilder
    UsedKeys          map[string]bool
    Attributes        map[string]interface{}
    SortDescending    bool
    LimitValue        *int
    ExclusiveStartKey map[string]types.AttributeValue
    PreferredSortKey  string
}

// NewQueryBuilder creates a new QueryBuilder instance
func NewQueryBuilder() *QueryBuilder {
    return &QueryBuilder{
        KeyConditions: make(map[string]expression.KeyConditionBuilder),
        UsedKeys:      make(map[string]bool),
        Attributes:    make(map[string]interface{}),
    }
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

// WithPreferredSortKey sets the preferred sort key for index selection
func (qb *QueryBuilder) WithPreferredSortKey(key string) *QueryBuilder {
    qb.PreferredSortKey = key
    return qb
}

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

// OrderByDesc sets descending sort order
func (qb *QueryBuilder) OrderByDesc() *QueryBuilder {
    qb.SortDescending = true
    return qb
}

// OrderByAsc sets ascending sort order  
func (qb *QueryBuilder) OrderByAsc() *QueryBuilder {
    qb.SortDescending = false
    return qb
}

// Limit sets the maximum number of items to return
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
    qb.LimitValue = &limit
    return qb
}

// StartFrom sets the exclusive start key for pagination
func (qb *QueryBuilder) StartFrom(lastEvaluatedKey map[string]types.AttributeValue) *QueryBuilder {
    qb.ExclusiveStartKey = lastEvaluatedKey
    return qb
}
`
