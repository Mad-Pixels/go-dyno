package core

const SchemaTemplate = `
// DynamoSchema ...
type DynamoSchema struct {
   	TableName        string
   	HashKey          string
   	RangeKey         string
   	Attributes       []Attribute
   	CommonAttributes []Attribute
   	SecondaryIndexes []SecondaryIndex
}

// Attribute ...
type Attribute struct {
   	Name string
   	Type string
}

// CompositeKeyPart ...
type CompositeKeyPart struct {
   	IsConstant bool
   	Value      string
}

// SecondaryIndex ...
type SecondaryIndex struct {
   	Name             string
   	HashKey          string
   	HashKeyParts     []CompositeKeyPart
   	RangeKey         string
   	RangeKeyParts    []CompositeKeyPart
   	ProjectionType   string
   	NonKeyAttributes []string
}

// SchemaItem ...
type SchemaItem struct {
{{- range .AllAttributes}}
   	{{ToSafeName .Name | ToUpperCamelCase}} {{ToGolangBaseType .}} ` + "`dynamodbav:\"{{.Name}}\"`" + `
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
}
`
