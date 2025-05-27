package v2

const QueryBuilderTemplate = `type QueryBuilder struct {
    IndexName           string
    KeyConditions       map[string]expression.KeyConditionBuilder
    FilterConditions    []expression.ConditionBuilder
    UsedKeys            map[string]bool
    Attributes          map[string]interface{}
    SortDescending      bool
    LimitValue          *int
    ExclusiveStartKey   map[string]types.AttributeValue
    PreferredSortKey    string
}

func NewQueryBuilder() *QueryBuilder {
    return &QueryBuilder{
        KeyConditions:   make(map[string]expression.KeyConditionBuilder),
        UsedKeys:        make(map[string]bool),
        Attributes:      make(map[string]interface{}),
    }
}

{{range .AllAttributes}}
func (qb *QueryBuilder) With{{ToSafeName .Name | ToUpperCamelCase}}({{ToSafeName .Name | ToLowerCamelCase}} {{ToGolangBaseType .Type}}) *QueryBuilder {
    qb.Attributes["{{.Name}}"] = {{ToSafeName .Name | ToLowerCamelCase}}
    qb.UsedKeys["{{.Name}}"] = true
    return qb
}
{{end}}

{{range .SecondaryIndexes}}
{{if gt (len .HashKeyParts) 0}}
{{- $hasNonConstant := false -}}
{{- range .HashKeyParts -}}{{- if not .IsConstant -}}{{- $hasNonConstant = true -}}{{- end -}}{{- end -}}
{{- if $hasNonConstant}}
func (qb *QueryBuilder) With{{ToUpperCamelCase .Name}}HashKey({{range $i, $part := .HashKeyParts}}{{if not $part.IsConstant}}{{if $i}}, {{end}}{{$part.Value | ToLowerCamelCase}} {{ToGolangAttrType $part.Value $.AllAttributes}}{{end}}{{end}}) *QueryBuilder {
    // Сохраняем отдельные значения для построения композитного ключа
    {{range .HashKeyParts}}{{if not .IsConstant}}
    qb.Attributes["{{.Value}}"] = {{.Value | ToLowerCamelCase}}
    qb.UsedKeys["{{.Value}}"] = true
    {{end}}{{end}}
    
    // Строим композитный ключ
    compositeValue := qb.buildCompositeKeyValue([]CompositeKeyPart{
        {{range .HashKeyParts}}
        {{if .IsConstant}}
        {IsConstant: true, Value: "{{.Value}}"},
        {{else}}
        {IsConstant: false, Value: "{{.Value}}"},
        {{end}}
        {{end}}
    })
    
    // Устанавливаем композитное значение и условие
    qb.Attributes["{{.HashKey}}"] = compositeValue
    qb.UsedKeys["{{.HashKey}}"] = true
    qb.KeyConditions["{{.HashKey}}"] = expression.Key("{{.HashKey}}").Equal(expression.Value(compositeValue))
    
    return qb
}
{{end}}
{{else if .HashKey}}
func (qb *QueryBuilder) With{{ToUpperCamelCase .Name}}HashKey({{.HashKey | ToLowerCamelCase}} {{ToGolangAttrType .HashKey $.AllAttributes}}) *QueryBuilder {
    qb.Attributes["{{.HashKey}}"] = {{.HashKey | ToLowerCamelCase}}
    qb.UsedKeys["{{.HashKey}}"] = true
    qb.KeyConditions["{{.HashKey}}"] = expression.Key("{{.HashKey}}").Equal(expression.Value({{.HashKey | ToLowerCamelCase}}))
    return qb
}
{{end}}
{{end}}

func (qb *QueryBuilder) WithPreferredSortKey(key string) *QueryBuilder {
    qb.PreferredSortKey = key
    return qb
}

{{range .SecondaryIndexes}}
{{if gt (len .RangeKeyParts) 0}}
{{- $hasNonConstant := false -}}
{{- range .RangeKeyParts -}}{{- if not .IsConstant -}}{{- $hasNonConstant = true -}}{{- end -}}{{- end -}}
{{- if $hasNonConstant}}
func (qb *QueryBuilder) With{{ToUpperCamelCase .Name}}RangeKey({{range $i, $part := .RangeKeyParts}}{{if not $part.IsConstant}}{{if $i}}, {{end}}{{$part.Value | ToLowerCamelCase}} {{ToGolangAttrType $part.Value $.AllAttributes}}{{end}}{{end}}) *QueryBuilder {
    {{range .RangeKeyParts}}{{if not .IsConstant}}
    qb.Attributes["{{.Value}}"] = {{.Value | ToLowerCamelCase}}
    qb.UsedKeys["{{.Value}}"] = true
    {{end}}{{end}}
    return qb
}
{{end}}
{{else if .RangeKey}}
func (qb *QueryBuilder) With{{ToUpperCamelCase .Name}}RangeKey({{.RangeKey | ToLowerCamelCase}} {{ToGolangAttrType .RangeKey $.AllAttributes}}) *QueryBuilder {
    qb.Attributes["{{.RangeKey}}"] = {{.RangeKey | ToLowerCamelCase}}
    qb.UsedKeys["{{.RangeKey}}"] = true
    qb.KeyConditions["{{.RangeKey}}"] = expression.Key("{{.RangeKey}}").Equal(expression.Value({{.RangeKey | ToLowerCamelCase}}))
    return qb
}
{{end}}
{{end}}

{{range .AllAttributes}}
{{if eq (ToGolangBaseType .Type) "int"}}
func (qb *QueryBuilder) With{{ToSafeName .Name | ToUpperCamelCase}}Between(start, end {{ToGolangBaseType .Type}}) *QueryBuilder {
    qb.KeyConditions["{{.Name}}"] = expression.Key("{{.Name}}").Between(expression.Value(start), expression.Value(end))
    qb.UsedKeys["{{.Name}}"] = true
    return qb
}

func (qb *QueryBuilder) With{{ToSafeName .Name | ToUpperCamelCase}}GreaterThan(value {{ToGolangBaseType .Type}}) *QueryBuilder {
    qb.KeyConditions["{{.Name}}"] = expression.Key("{{.Name}}").GreaterThan(expression.Value(value))
    qb.UsedKeys["{{.Name}}"] = true
    return qb
}

func (qb *QueryBuilder) With{{ToSafeName .Name | ToUpperCamelCase}}LessThan(value {{ToGolangBaseType .Type}}) *QueryBuilder {
    qb.KeyConditions["{{.Name}}"] = expression.Key("{{.Name}}").LessThan(expression.Value(value))
    qb.UsedKeys["{{.Name}}"] = true
    return qb
}
{{end}}
{{end}}

func (qb *QueryBuilder) OrderByDesc() *QueryBuilder {
    qb.SortDescending = true
    return qb
}

func (qb *QueryBuilder) OrderByAsc() *QueryBuilder {
    qb.SortDescending = false
    return qb
}

func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
    qb.LimitValue = &limit
    return qb
}

func (qb *QueryBuilder) StartFrom(lastEvaluatedKey map[string]types.AttributeValue) *QueryBuilder {
    qb.ExclusiveStartKey = lastEvaluatedKey
    return qb
}` + QueryBuilderBuildTemplate + QueryBuilderUtilsTemplate
