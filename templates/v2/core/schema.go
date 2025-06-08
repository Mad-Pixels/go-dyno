package core

// SchemaTemplate ...
const SchemaTemplate = `
// FieldInfo contains metadata about a schema field
type FieldInfo struct {
    DynamoType string
    IsKey      bool
    IsHashKey  bool
    IsRangeKey bool
}

// DynamoSchema ...
type DynamoSchema struct {
    TableName        string
    HashKey          string
    RangeKey         string
    Attributes       []Attribute
    CommonAttributes []Attribute
    SecondaryIndexes []SecondaryIndex
    // NEW: Быстрый поиск полей O(1)
    FieldsMap        map[string]FieldInfo
}

// Остальные типы без изменений...
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

type SchemaItem struct {
{{- range .AllAttributes}}
    {{ToSafeName .Name | ToUpperCamelCase}} {{if eq .Type "SS"}}[]string{{else if eq .Type "NS"}}[]int{{else if eq .Type "BS"}}[][]byte{{else}}{{ToGolangBaseType .}}{{end}} ` + "`dynamodbav:\"{{.Name}}\"`" + `
{{- end}}
}

// TableSchema ...
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
            {{- if .HashKeyParts}}
            HashKeyParts: []CompositeKeyPart{
                {{- range .HashKeyParts}}
                {IsConstant: {{.IsConstant}}, Value: "{{.Value}}"},
                {{- end}}
            },
            {{- end}}
            RangeKey:       "{{.RangeKey}}",
            {{- if .RangeKeyParts}}
            RangeKeyParts: []CompositeKeyPart{
                {{- range .RangeKeyParts}}
                {IsConstant: {{.IsConstant}}, Value: "{{.Value}}"},
                {{- end}}
            },
            {{- end}}
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
    
    // NEW: Предварительно собранная map для O(1) поиска
    FieldsMap: map[string]FieldInfo{
        {{- range .AllAttributes}}
        "{{.Name}}": {
            DynamoType: "{{.Type}}",
            IsKey:      {{if or (eq .Name $.HashKey) (eq .Name $.RangeKey)}}true{{else}}false{{end}},
            IsHashKey:  {{if eq .Name $.HashKey}}true{{else}}false{{end}},
            IsRangeKey: {{if eq .Name $.RangeKey}}true{{else}}false{{end}},
        },
        {{- end}}
    },
}
`
