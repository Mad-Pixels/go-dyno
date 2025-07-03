package validate

const usageTemplate = `
🔍 {{.Command}} verifies the correctness and integrity of a DynamoDB JSON schema definition.

This command performs comprehensive validation of your schema including:
  • 📝 Schema structure and required fields validation
  • 🗃️ DynamoDB attribute types and constraints verification  
  • 🔑 Primary key (hash_key/range_key) existence and validity
  • 📊 Secondary indexes (GSI/LSI) configuration and attribute references
  • 🔗 Composite key syntax and attribute mapping validation
  • 🐹 Go identifier safety checks for generated code compatibility

The validator ensures your schema is ready for code generation and will work 
correctly with AWS DynamoDB before you generate any Go code. 🚀

EXAMPLES:
   $ {{.EnvPrefix}}_{{.FlagSchemaPath}}=./schema.json godyno {{.Command}}
   $ godyno {{.Command}} --{{.FlagSchemaPath}} ./configs/user-posts.json
   $ godyno {{.Command}} -s ./schemas/orders.json

VALIDATION CHECKS:
   ✅ JSON syntax and structure
   ✅ Required fields presence (table_name, hash_key, attributes)
   ✅ DynamoDB type compatibility (S, N, B, SS, NS, BS, etc.)
   ✅ Index key references to existing attributes
   ✅ Composite key format and attribute resolution
   ✅ Go naming conventions and reserved keyword conflicts
`
