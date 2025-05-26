package v2

// QueryBuilderTemplate ...
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
    attrName := "{{.Name}}"
    qb.Attributes[attrName] = {{ToSafeName .Name | ToLowerCamelCase}}
    qb.UsedKeys[attrName] = true
    return qb
}
{{end}}

{{range .SecondaryIndexes}}
{{if gt (len .HashKeyParts) 0}}
{{ $methodParams := "" }}
{{ range $index, $part := .HashKeyParts }}
    {{ if not $part.IsConstant }}
        {{ $paramName := (ToSafeName $part.Value | ToLowerCamelCase) }}
        {{ $paramType := (ToGolangAttrType $part.Value $.AllAttributes) }}
        {{ if eq $methodParams "" }}
            {{ $methodParams = printf "%s %s" $paramName $paramType }}
        {{ else }}
            {{ $methodParams = printf "%s, %s %s" $methodParams $paramName $paramType }}
        {{ end }}
    {{ end }}
{{ end }}
func (qb *QueryBuilder) With{{ToUpperCamelCase .Name}}HashKey({{ $methodParams }}) *QueryBuilder {
    {{ range $index, $part := .HashKeyParts }}
    {{ if not $part.IsConstant }}
    {
        attrName := "{{ $part.Value }}"
        qb.Attributes[attrName] = {{ $part.Value | ToLowerCamelCase }}
        qb.UsedKeys[attrName] = true
        qb.KeyConditions[attrName] = expression.Key(attrName).Equal(expression.Value({{ $part.Value | ToLowerCamelCase }}))
    }
    {{ end }}
    {{ end }}
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
{{ $methodParams := "" }}
{{ range $index, $part := .RangeKeyParts }}
    {{ if not $part.IsConstant }}
        {{ $paramName := (ToSafeName $part.Value | ToLowerCamelCase) }}
        {{ $paramType := (ToGolangAttrType $part.Value $.AllAttributes) }}
        {{ if eq $methodParams "" }}
            {{ $methodParams = printf "%s %s" $paramName $paramType }}
        {{ else }}
            {{ $methodParams = printf "%s, %s %s" $methodParams $paramName $paramType }}
        {{ end }}
    {{ end }}
{{ end }}
func (qb *QueryBuilder) With{{ToUpperCamelCase .Name}}RangeKey({{ $methodParams }}) *QueryBuilder {
    {{ range .RangeKeyParts }}
    {{ if not .IsConstant }}
    {
        attrName := "{{ .Value }}"
        qb.Attributes[attrName] = {{ .Value | ToLowerCamelCase }}
        qb.UsedKeys[attrName] = true
        qb.KeyConditions[attrName] = expression.Key(attrName).Equal(expression.Value({{ .Value | ToLowerCamelCase }}))
    }
    {{ end }}
    {{ end }}
    return qb
}
{{end}}
{{end}}

{{range .AllAttributes}}
{{if eq (ToGolangBaseType .Type) "int"}}
func (qb *QueryBuilder) With{{ToSafeName .Name | ToUpperCamelCase}}Between(start, end {{ToGolangBaseType .Type}}) *QueryBuilder {
    attrName := "{{.Name}}"
    qb.KeyConditions[attrName] = expression.Key(attrName).Between(expression.Value(start), expression.Value(end))
    qb.UsedKeys[attrName] = true
    return qb
}

func (qb *QueryBuilder) With{{ToSafeName .Name | ToUpperCamelCase}}GreaterThan(value {{ToGolangBaseType .Type}}) *QueryBuilder {
    attrName := "{{.Name}}"
    qb.KeyConditions[attrName] = expression.Key(attrName).GreaterThan(expression.Value(value))
    qb.UsedKeys[attrName] = true
    return qb
}

func (qb *QueryBuilder) With{{ToSafeName .Name | ToUpperCamelCase}}LessThan(value {{ToGolangBaseType .Type}}) *QueryBuilder {
    attrName := "{{.Name}}"
    qb.KeyConditions[attrName] = expression.Key(attrName).LessThan(expression.Value(value))
    qb.UsedKeys[attrName] = true
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
