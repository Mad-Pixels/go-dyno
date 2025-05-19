package templates

var QueryBuilderIndexMethodsTemplate = `
{{range .SecondaryIndexes}}
{{if gt (len .HashKeyParts) 0}}
{{ $methodParams := "" }}
{{ range $index, $part := .HashKeyParts }}
    {{ if not $part.IsConstant }}
        {{ $paramName := (SafeName $part.Value | ToLowerCamelCase) }}
        {{ $paramType := (TypeGoAttr $part.Value $.AllAttributes) }}
        {{ if eq $methodParams "" }}
            {{ $methodParams = printf "%s %s" $paramName $paramType }}
        {{ else }}
            {{ $methodParams = printf "%s, %s %s" $methodParams $paramName $paramType }}
        {{ end }}
    {{ end }}
{{ end }}
func (qb *QueryBuilder) With{{ToCamelCase .Name}}HashKey({{ $methodParams }}) *QueryBuilder {
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
        {{ $paramName := (SafeName $part.Value | ToLowerCamelCase) }}
        {{ $paramType := (TypeGoAttr $part.Value $.AllAttributes) }}
        {{ if eq $methodParams "" }}
            {{ $methodParams = printf "%s %s" $paramName $paramType }}
        {{ else }}
            {{ $methodParams = printf "%s, %s %s" $methodParams $paramName $paramType }}
        {{ end }}
    {{ end }}
{{ end }}
func (qb *QueryBuilder) With{{ToCamelCase .Name}}RangeKey({{ $methodParams }}) *QueryBuilder {
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
`
