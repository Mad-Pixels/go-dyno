package v2

// ConstantsTemplate generates compile-time constants and runtime variables for DynamoDB operations.
// This template creates:
// - String constants for table and index names (compile-time safe)
// - Column name constants for all attributes (prevents typos in queries)
// - Runtime slices and maps for dynamic operations and projections
const ConstantsTemplate = `// Table and Index name constants
// These constants provide compile-time safety when referencing table and index names
// in DynamoDB operations, preventing runtime errors from typos.
const (
   // TableName is the DynamoDB table name used in all AWS SDK operations
   // Example: dynamoClient.Query(ctx, &dynamodb.QueryInput{TableName: aws.String(TableName)})
   TableName = "{{.TableName}}"
   
   {{range .SecondaryIndexes}}
   // Index{{.Name}} is the name constant for the "{{.Name}}" secondary index
   // Example: input.IndexName = aws.String(Index{{.Name}})
   Index{{.Name}} = "{{.Name}}"
   {{- end}}
)

// Column name constants for all table attributes
// These constants ensure type safety when building DynamoDB expressions and prevent
// attribute name typos in queries, filters, and update expressions.
const (
   {{range .AllAttributes}}
   // Column{{ToSafeName .Name | ToUpperCamelCase}} represents the "{{.Name}}" attribute in DynamoDB
   // DynamoDB type: {{.Type}} - Use in expressions: expression.Name(Column{{ToSafeName .Name | ToUpperCamelCase}})
   Column{{ToSafeName .Name | ToUpperCamelCase}} = "{{.Name}}"
   {{- end}}
)

// Runtime variables for dynamic operations
var (
   // AttributeNames contains all attribute names for this table.
   // Useful for operations that need to iterate over all attributes or
   // for building projection expressions that include all fields.
   //
   // Example usage:
   //   projection := expression.NamesList(expression.Name(AttributeNames[0]), ...)
   //   for _, attrName := range AttributeNames {
   //       fmt.Println("Attribute:", attrName)
   //   }
   AttributeNames = []string{
       {{- range .AllAttributes}}
       "{{.Name}}", // {{.Name}} ({{.Type}})
       {{- end}}
   }

   // IndexProjections maps each secondary index name to its projected attributes.
   // This is useful for understanding what attributes are available when querying
   // specific indexes and for building efficient queries.
   //
   // Projection types:
   // - "ALL": All table attributes are projected (includes all AttributeNames)
   // - "KEYS_ONLY": Only key attributes are projected (hash + range keys)
   // - "INCLUDE": Key attributes plus specified non-key attributes
   //
   // Example usage:
   //   projectedAttrs := IndexProjections["{{range .SecondaryIndexes}}{{.Name}}{{break}}{{end}}"]
   //   fmt.Printf("Index projects %d attributes\n", len(projectedAttrs))
   //   canQuery := slices.Contains(projectedAttrs, "some_attribute")
   IndexProjections = map[string][]string{
       {{- range .SecondaryIndexes}}
       // {{.Name}} index projection ({{.ProjectionType}})
       {{- if eq .ProjectionType "ALL"}}
       // Projects ALL attributes - complete table data available in index
       {{- else if eq .ProjectionType "KEYS_ONLY"}}
       // Projects KEYS_ONLY - only key attributes available
       {{- else}}
       // Projects INCLUDE - key attributes plus specified non-key attributes
       {{- end}}
       "{{.Name}}": {
           {{- if eq .ProjectionType "ALL"}}
           {{- range $.AllAttributes}}
           "{{.Name}}", // {{.Name}} ({{.Type}}) - from ALL projection
           {{- end}}
           {{- else}}
           "{{.HashKey}}", // Hash key for {{.Name}} index
           {{- if .RangeKey}}
           "{{.RangeKey}}", // Range key for {{.Name}} index
           {{- end}}
           {{- range .NonKeyAttributes}}
           "{{.}}", // Non-key attribute included in projection
           {{- end}}
           {{- end}}
       },
       {{- end}}
   }
)`
