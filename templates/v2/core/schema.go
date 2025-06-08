package core

// SchemaTemplate with pre-computed allowed operators
const SchemaTemplate = `
// FieldInfo contains metadata about a schema field
type FieldInfo struct {
    DynamoType       string
    IsKey            bool
    IsHashKey        bool
    IsRangeKey       bool
    AllowedOperators map[OperatorType]bool
}

// SupportsOperator checks if this field supports the given operator
func (fi FieldInfo) SupportsOperator(op OperatorType) bool {
    return fi.AllowedOperators[op]
}

// buildAllowedOperators returns the set of allowed operators for a DynamoDB type
func buildAllowedOperators(dynamoType string) map[OperatorType]bool {
    allowed := make(map[OperatorType]bool)
    
    switch dynamoType {
    case "S": // String
        allowed[EQ] = true
        allowed[NE] = true
        allowed[GT] = true
        allowed[LT] = true
        allowed[GTE] = true
        allowed[LTE] = true
        allowed[BETWEEN] = true
        allowed[CONTAINS] = true
        allowed[NOT_CONTAINS] = true
        allowed[BEGINS_WITH] = true
        allowed[IN] = true
        allowed[NOT_IN] = true
        allowed[EXISTS] = true
        allowed[NOT_EXISTS] = true
        
    case "N": // Number  
        allowed[EQ] = true
        allowed[NE] = true
        allowed[GT] = true
        allowed[LT] = true
        allowed[GTE] = true
        allowed[LTE] = true
        allowed[BETWEEN] = true
        allowed[IN] = true
        allowed[NOT_IN] = true
        allowed[EXISTS] = true
        allowed[NOT_EXISTS] = true
        
    case "BOOL": // Boolean
        allowed[EQ] = true
        allowed[NE] = true
        allowed[EXISTS] = true
        allowed[NOT_EXISTS] = true
        
    case "SS": // String Set - only CONTAINS/NOT_CONTAINS, not IN/NOT_IN
        allowed[CONTAINS] = true
        allowed[NOT_CONTAINS] = true
        allowed[EXISTS] = true
        allowed[NOT_EXISTS] = true
        
    case "NS": // Number Set - only CONTAINS/NOT_CONTAINS, not IN/NOT_IN
        allowed[CONTAINS] = true
        allowed[NOT_CONTAINS] = true
        allowed[EXISTS] = true
        allowed[NOT_EXISTS] = true
        
    case "BS": // Binary Set (rare)
        allowed[CONTAINS] = true
        allowed[NOT_CONTAINS] = true
        allowed[EXISTS] = true
        allowed[NOT_EXISTS] = true
        
    case "L": // List
        allowed[EXISTS] = true
        allowed[NOT_EXISTS] = true
        
    case "M": // Map
        allowed[EXISTS] = true
        allowed[NOT_EXISTS] = true
        
    case "NULL": // Null
        allowed[EXISTS] = true
        allowed[NOT_EXISTS] = true
        
    default:
        // For unknown types allow only basic operations
        allowed[EQ] = true
        allowed[NE] = true
        allowed[EXISTS] = true
        allowed[NOT_EXISTS] = true
    }
    
    return allowed
}

// DynamoSchema ...
type DynamoSchema struct {
    TableName        string
    HashKey          string
    RangeKey         string
    Attributes       []Attribute
    CommonAttributes []Attribute
    SecondaryIndexes []SecondaryIndex
    // Быстрый поиск полей O(1) с предвычисленными операторами
    FieldsMap        map[string]FieldInfo
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

type SchemaItem struct {
{{- range .AllAttributes}}
    {{ToSafeName .Name | ToUpperCamelCase}} {{ToGolangBaseType .}} ` + "`dynamodbav:\"{{.Name}}\"`" + `
{{- end}}
}

// TableSchema with pre-computed FieldsMap including allowed operators
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
    
    FieldsMap: map[string]FieldInfo{
        {{- range .AllAttributes}}
        "{{.Name}}": {
            DynamoType:       "{{.Type}}",
            IsKey:            {{if or (eq .Name $.HashKey) (eq .Name $.RangeKey)}}true{{else}}false{{end}},
            IsHashKey:        {{if eq .Name $.HashKey}}true{{else}}false{{end}},
            IsRangeKey:       {{if eq .Name $.RangeKey}}true{{else}}false{{end}},
            AllowedOperators: buildAllowedOperators("{{.Type}}"),
        },
        {{- end}}
    },
}
`
