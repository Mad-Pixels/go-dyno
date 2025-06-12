package core

// ConstantsTemplate define constants.
const ConstantsTemplate = `
const (
    // TableName is the DynamoDB table name for all operations.
    TableName = "{{.TableName}}"
   
    {{range .SecondaryIndexes}}
    // Index{{.Name}} is the {{if .IsGSI}}GSI{{else}}LSI{{end}} index name for alternative queries.
    Index{{.Name}} = "{{.Name}}"
    {{- end}}

    {{range .AllAttributes}}
    // Column{{ToSafeName .Name | ToUpperCamelCase}} is the attribute name for "{{.Name}}" field.
    Column{{ToSafeName .Name | ToUpperCamelCase}} = "{{.Name}}"
    {{- end}}
)

var (
    // AttributeNames contains all table attribute names for projection expressions.
    // Example: expression.NamesList(expression.Name(AttributeNames[0]))
    AttributeNames = []string{
        {{- range .AllAttributes}}
        "{{.Name}}",
        {{- end}}
    }

    // KeyAttributeNames contains primary key attributes for key operations.
    // Example: validateKeys(item, KeyAttributeNames)
    KeyAttributeNames = []string{
        "{{.HashKey}}",
        {{- if .RangeKey}}
        "{{.RangeKey}}",
        {{- end}}
    }
)
`
