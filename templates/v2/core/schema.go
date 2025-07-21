package core

// SchemaTemplate with pre-computed allowed operators
const SchemaTemplate = `
// FieldInfo contains metadata about a schema field with operator validation.
type FieldInfo struct {
    DynamoType       string
    IsKey            bool
    IsHashKey        bool
    IsRangeKey       bool
    AllowedOperators map[OperatorType]bool
}

// SupportsOperator checks if this field supports the given operator.
// Returns false for invalid operator/type combinations.
func (fi FieldInfo) SupportsOperator(op OperatorType) bool {
    return fi.AllowedOperators[op]
}

// buildAllowedOperators returns the set of allowed operators for a DynamoDB type.
// Implements DynamoDB operator compatibility rules for each data type.
func buildAllowedOperators(dynamoType string) map[OperatorType]bool {
    allowed := make(map[OperatorType]bool)
    
    switch dynamoType {
    case "S": // String - supports all comparison and string operations
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
        
    case "N": // Number - supports comparison operations, no string functions
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
        
    case "BOOL": // Boolean - only equality and existence checks
        allowed[EQ] = true
        allowed[NE] = true
        allowed[EXISTS] = true
        allowed[NOT_EXISTS] = true
        
    case "SS": // String Set - membership operations only, not IN/NOT_IN
        allowed[CONTAINS] = true
        allowed[NOT_CONTAINS] = true
        allowed[EXISTS] = true
        allowed[NOT_EXISTS] = true
        
    case "NS": // Number Set - membership operations only, not IN/NOT_IN
        allowed[CONTAINS] = true
        allowed[NOT_CONTAINS] = true
        allowed[EXISTS] = true
        allowed[NOT_EXISTS] = true
        
    case "BS": // Binary Set - membership operations only
        allowed[CONTAINS] = true
        allowed[NOT_CONTAINS] = true
        allowed[EXISTS] = true
        allowed[NOT_EXISTS] = true
        
    case "L": // List - only existence checks
        allowed[EXISTS] = true
        allowed[NOT_EXISTS] = true
        
    case "M": // Map - only existence checks
        allowed[EXISTS] = true
        allowed[NOT_EXISTS] = true
        
    case "NULL": // Null - only existence checks
        allowed[EXISTS] = true
        allowed[NOT_EXISTS] = true
        
    default:
        // Unknown types - basic operations only
        allowed[EQ] = true
        allowed[NE] = true
        allowed[EXISTS] = true
        allowed[NOT_EXISTS] = true
    }
    return allowed
}

// DynamoSchema represents the complete table schema with indexes and metadata.
type DynamoSchema struct {
    TableName        string
    HashKey          string
    RangeKey         string
    Attributes       []Attribute
    CommonAttributes []Attribute
    SecondaryIndexes []SecondaryIndex
    FieldsMap        map[string]FieldInfo
}

// Attribute represents a DynamoDB table attribute with its type.
type Attribute struct {
    Name string  // Attribute name
    Type string  // DynamoDB type (S, N, BOOL, SS, NS, etc.)
}

// CompositeKeyPart represents a part of a composite key structure.
// Used for complex key patterns in GSI/LSI definitions.
type CompositeKeyPart struct {
    IsConstant bool    // true if this part is a constant value
    Value      string  // the constant value or attribute name
}

// SecondaryIndex represents a GSI or LSI with optional composite keys.
// Supports both simple and composite key structures for advanced access patterns.
type SecondaryIndex struct {
    Name             string
    HashKey          string
    RangeKey         string
    ProjectionType   string
    HashKeyParts     []CompositeKeyPart  // for composite hash keys
    RangeKeyParts    []CompositeKeyPart  // for composite range keys
    NonKeyAttributes []string            // projected attributes for INCLUDE
}

// SchemaItem represents a single DynamoDB item with all table attributes.
// All fields are properly tagged for AWS SDK marshaling/unmarshaling.
type SchemaItem struct {
{{- range .AllAttributes}}
    {{ToSafeName .Name | ToUpperCamelCase}} {{ToGolangBaseType .}} ` + "`{{ToDynamoDBStructTag .}}`" + `
{{- end}}
}

// TableSchema contains the complete schema definition with pre-computed metadata.
// Used throughout the generated code for validation and operator checking.
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
