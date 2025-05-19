package templates

var SchemaStructsTemplate = `
type DynamoSchema struct {
    TableName        string
    HashKey          string
    RangeKey         string
    Attributes       []Attribute
    CommonAttributes []Attribute
    SecondaryIndexes []SecondaryIndex
}

type Attribute struct {
    Name string
    Type string
}

type CompositeKeyPart struct {
    IsConstant bool
    Value      string
}

type SecondaryIndex struct {
    Name             string
    HashKey          string
    HashKeyParts     []CompositeKeyPart
    RangeKey         string
    RangeKeyParts    []CompositeKeyPart
    ProjectionType   string
    NonKeyAttributes []string
}

// SchemaItem represents an item in "{{.TableName}}"
type SchemaItem struct {
    {{range .AllAttributes}}
    {{SafeName .Name | ToCamelCase}} {{TypeGo .Type}} ` + "`dynamodbav:\"{{.Name}}\"`" + `
    {{end}}
}

var TableSchema = DynamoSchema{
    TableName: "{{.TableName}}",
    HashKey:   "{{.HashKey}}",
    RangeKey:  "{{.RangeKey}}",
    Attributes: []Attribute{
        {{- range .Attributes}}
        {Name: "{{.Name}}", Type: "{{.Type}}"},
        {{- end}}
    },
    CommonAttributes: []Attribute{
        {{- range .CommonAttributes}}
        {Name: "{{.Name}}", Type: "{{.Type}}"},
        {{- end}}
    },
    SecondaryIndexes: []SecondaryIndex{
        {{- range .SecondaryIndexes}}
        {
            Name:           "{{.Name}}",
            HashKey:        "{{.HashKey}}",
            RangeKey:       "{{.RangeKey}}",
            ProjectionType: "{{.ProjectionType}}",
            {{- if .NonKeyAttributes}}
            NonKeyAttributes: []string{
                {{- range .NonKeyAttributes}}
                "{{.}}",
                {{- end}}
            },
            {{- end}}
        },
        {{- end}}
    },
}
`
