package core

// ConstantsTemplate ...
const ConstantsTemplate = `
const (
    // TableName ...
    TableName = "{{.TableName}}"
   
    {{range .SecondaryIndexes}}
    // Index{{.Name}} ...
    Index{{.Name}} = "{{.Name}}"
    {{- end}}

    {{range .AllAttributes}}
    // Column{{ToSafeName .Name | ToUpperCamelCase}} ...
    Column{{ToSafeName .Name | ToUpperCamelCase}} = "{{.Name}}"
    {{- end}}
)

var (
    // AttributeNames contains all table attribute names
    AttributeNames = []string{
        {{- range .AllAttributes}}
        "{{.Name}}",
        {{- end}}
    }

    // KeyAttributeNames contains only key attribute names  
    KeyAttributeNames = []string{
        "{{.HashKey}}",
        {{- if .RangeKey}}
        "{{.RangeKey}}",
        {{- end}}
    }
)
`
