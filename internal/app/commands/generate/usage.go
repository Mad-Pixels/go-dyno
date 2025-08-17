package generate

const usageTemplate = `
{{.Command}} generates production-ready Go code from a DynamoDB JSON schema definition.

This command transforms your schema into a complete Go package with:
  • 🏗️ Strongly-typed structs with proper DynamoDB tags
  • 🔑 Type-safe constants for table/column names and indexes
  • 🚀 High-performance query builders with fluent API
  • 📊 Scan operations with filtering and pagination support
  • 🔄 Atomic update helpers and batch operations
  • 🎯 AWS SDK v2 compatible marshallers and unmarshallers

Generated code is optimized for performance, includes comprehensive error handling,
and follows Go best practices for maintainable production applications. 🎉

EXAMPLES:
   # Generate to stdout
   $ godyno {{.Command}} --{{.FlagSchemaPath}} ./schema.json
   
   # Generate to specific directory
   $ godyno {{.Command}} -s ./schema.json --output-dir ./generated
   
   # Override package and filename
   $ godyno {{.Command}} -s ./schema.json -o ./models --package users --filename user.go
   
   # Using environment variables
   $ {{.EnvPrefix}}_SCHEMA=./schema.json {{.EnvPrefix}}_OUTPUT_DIR=./gen godyno {{.Command}}

   # With DynamoDB stream events methods
   $ godyno {{.Command}} -s ./schema.json --output-dir ./generated --with-stream-events

GENERATED FEATURES:
   ✨ Type-safe structs with dynamodbav tags
   ✨ Table/column/index constants (no magic strings!)
   ✨ Fluent query builders with condition expressions
   ✨ Batch operations and atomic updates
   ✨ DynamoDB Streams event handlers
   ✨ Comprehensive error handling and validation
   ✨ AWS SDK v2 compatible (latest and greatest!)
   ✨ Production-ready with zero dependencies
`
