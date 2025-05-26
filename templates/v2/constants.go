package v2

// ConstantsTemplate ...
const ConstantsTemplate = `const (
    TableName = "{{.TableName}}"
    {{range .SecondaryIndexes}}
    Index{{.Name}} = "{{.Name}}"
    {{- end}}
)

const (
    {{range .AllAttributes}}
    Column{{ToSafeName .Name | ToUpperCamelCase}} = "{{.Name}}"
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
)`
