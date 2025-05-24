// Package generate provides a CLI command for generating static Go code from
// a DynamoDB JSON schema definition.
package generate

const usageTemplate = `
{{.Command}} generates static Go code based on a DynamoDB JSON schema definition.

Required Flags:
  --{{.FlagCfg}}, -c     Path to the JSON config schema.
                         Example: --{{.FlagCfg}} ./schema/user.json

  --{{.FlagDst}}, -d    Output directory where the Go file will be written.
                         Example: --{{.FlagDst}} ./gen_output

Environment Variables:
  {{.EnvPrefix}}_{{.FlagCfg | ToUpper}}      Path to the schema file.
  {{.EnvPrefix}}_{{.FlagDst | ToUpper}}     Output directory.

Example Usage:

  $ go-dyno {{.Command}} \
    --{{.FlagCfg}} ./schema/{{.ExampleJSON}} \
    --{{.FlagDst}} ./gen

This will generate Go code from the given schema into: ./gen/{{.ExampleJSON | ToSafeName}}.go

The generated file contains:
  - TableName constant
  - Index constants
  - All attributes as constants
  - Go struct representing each record
  - Query builder with typed methods
`
