package v2

// SchemaStructsTemplate generates Go structures for DynamoDB table schema and item representation.
// This template creates type-safe Go code including:
// - Schema metadata structures (DynamoSchema, Attribute, SecondaryIndex, etc.)
// - Item structure with proper field types and DynamoDB tags
// - Pre-initialized TableSchema variable with all table configuration
const SchemaStructsTemplate = `// DynamoSchema represents the complete metadata structure of a DynamoDB table.
// It contains all information needed to interact with the table including keys,
// attributes, and secondary indexes configuration.
type DynamoSchema struct {
   // TableName is the actual DynamoDB table name used in AWS operations
   TableName        string
   
   // HashKey is the partition key attribute name for the main table
   HashKey          string
   
   // RangeKey is the sort key attribute name for the main table (optional)
   RangeKey         string
   
   // Attributes contains table-specific attribute definitions
   Attributes       []Attribute
   
   // CommonAttributes contains reusable attribute definitions shared across tables
   CommonAttributes []Attribute
   
   // SecondaryIndexes contains all GSI and LSI definitions for the table
   SecondaryIndexes []SecondaryIndex
}

// Attribute represents a single DynamoDB attribute with its name and type.
// Type values follow DynamoDB conventions: "S" (String), "N" (Number), "B" (Boolean).
type Attribute struct {
   // Name is the attribute name as it appears in DynamoDB
   Name string
   
   // Type is the DynamoDB attribute type: "S", "N", "B", etc.
   Type string
}

// CompositeKeyPart represents a single component of a composite key.
// Composite keys are formed by joining multiple parts with "#" separator.
// Example: "user#status" where "user" might be dynamic and "status" constant.
type CompositeKeyPart struct {
   // IsConstant indicates if this part is a fixed string (true) or 
   // a reference to an attribute value (false)
   IsConstant bool
   
   // Value contains either the constant string or the attribute name
   Value      string
}

// SecondaryIndex represents a Global Secondary Index (GSI) or Local Secondary Index (LSI).
// It defines alternate access patterns for querying the DynamoDB table.
type SecondaryIndex struct {
   // Name is the index name used in DynamoDB Query operations
   Name             string
   
   // HashKey is the partition key for this index (can be composite)
   HashKey          string
   
   // HashKeyParts contains parsed components if HashKey is composite
   // Empty if HashKey is a simple attribute reference
   HashKeyParts     []CompositeKeyPart
   
   // RangeKey is the sort key for this index (optional, can be composite)
   RangeKey         string
   
   // RangeKeyParts contains parsed components if RangeKey is composite
   // Empty if RangeKey is a simple attribute reference or not present
   RangeKeyParts    []CompositeKeyPart
   
   // ProjectionType defines which attributes are projected into the index
   // Valid values: "ALL", "KEYS_ONLY", "INCLUDE"
   ProjectionType   string
   
   // NonKeyAttributes lists additional attributes included when ProjectionType is "INCLUDE"
   // Empty for "ALL" and "KEYS_ONLY" projection types
   NonKeyAttributes []string
}

// SchemaItem represents a single item/record in the "{{.TableName}}" DynamoDB table.
// Each field corresponds to a table attribute with proper Go types and DynamoDB tags.
// The struct tags enable automatic marshaling/unmarshaling with AWS SDK v2.
type SchemaItem struct {
{{- range .AllAttributes}}
   // {{ToSafeName .Name | ToUpperCamelCase}} corresponds to the "{{.Name}}" attribute in DynamoDB
   // DynamoDB type: {{.Type}} -> Go type: {{ToGolangBaseType .Type}}
   {{ToSafeName .Name | ToUpperCamelCase}} {{ToGolangBaseType .Type}} ` + "`dynamodbav:\"{{.Name}}\"`" + `
{{- end}}
}

// TableSchema is a pre-initialized DynamoSchema instance containing all metadata
// for the "{{.TableName}}" table. Use this for runtime table operations, queries,
// and to access table configuration without hardcoding values.
//
// Example usage:
//   fmt.Println("Table name:", TableSchema.TableName)
//   fmt.Println("Hash key:", TableSchema.HashKey)
//   for _, idx := range TableSchema.SecondaryIndexes {
//       fmt.Println("Index:", idx.Name)
//   }
var TableSchema = DynamoSchema{
   TableName: "{{.TableName}}",
   HashKey:   "{{.HashKey}}",
   RangeKey:  "{{.RangeKey}}",
   
   // Table-specific attributes defined in the schema
   Attributes: []Attribute{
   	{{- range .Attributes}}
   	{Name: "{{.Name}}", Type: "{{.Type}}"}, // {{.Name}} ({{.Type}})
   	{{- end}}
   },
   
   // Common attributes shared across multiple tables (e.g., timestamps)
   CommonAttributes: []Attribute{
   	{{- range .CommonAttributes}}
   	{Name: "{{.Name}}", Type: "{{.Type}}"}, // {{.Name}} ({{.Type}})
   	{{- end}}
   },
   
   // Secondary indexes for alternate query patterns
   SecondaryIndexes: []SecondaryIndex{
   	{{- range .SecondaryIndexes}}
   	{
   		Name:           "{{.Name}}", // Index for querying by {{.HashKey}}{{if .RangeKey}} and {{.RangeKey}}{{end}}
   		HashKey:        "{{.HashKey}}",
   		RangeKey:       "{{.RangeKey}}",
   		ProjectionType: "{{.ProjectionType}}",
   		{{- if .NonKeyAttributes}}
   		// Additional attributes projected into this index
   		NonKeyAttributes: []string{
   			{{- range .NonKeyAttributes}}
   			"{{.}}", // Projected attribute: {{.}}
   			{{- end}}
   		},
   		{{- end}}
   	},
   	{{- end}}
   },
}`
