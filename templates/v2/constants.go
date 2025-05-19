package templates

var ConstantsTemplate = `
const (
    TableName = "{{.TableName}}"
    {{range .SecondaryIndexes}}
    Index{{.Name}} = "{{.Name}}"
    {{- end}}
)

const (
    {{range .AllAttributes}}
    Column{{SafeName .Name | ToCamelCase}} = "{{.Name}}"
    {{- end}}
)

var (
    AttributeNames = []string{
        {{- range .AllAttributes}}
        "{{.Name}}",
        {{- end}}
    }

    IndexProjections = map[string][]string{
        {{- range .SecondaryIndexes}}
        "{{.Name}}": {
            {{- if eq .ProjectionType "ALL"}}
            {{- range $.AllAttributes}}
            "{{.Name}}",
            {{- end}}
            {{- else}}
            "{{.HashKey}}", {{if .RangeKey}}"{{.RangeKey}}",{{end}}
            {{- range .NonKeyAttributes}}
            "{{.}}",
            {{- end}}
            {{- end}}
        },
        {{- end}}
    }
)
`
