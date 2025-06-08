package v2

import (
	"github.com/Mad-Pixels/go-dyno/templates/v2/core"
	"github.com/Mad-Pixels/go-dyno/templates/v2/helpers"
	"github.com/Mad-Pixels/go-dyno/templates/v2/inputs"
	"github.com/Mad-Pixels/go-dyno/templates/v2/query"
	"github.com/Mad-Pixels/go-dyno/templates/v2/scan"
)

// CodeTemplate ...
const CodeTemplate = `
package {{.PackageName}}

` + core.ImportsTemplate + `

` + core.ConstantsTemplate + `

` + core.SchemaTemplate + `

` + query.QueryBuilderTemplate + query.QueryBuilderBuildTemplate + query.QueryBuilderUtilsTemplate + `

` + scan.ScanBuilderTemplate + scan.ScanBuilderBuildTemplate + `

` + inputs.ItemInputsTemplate + inputs.UpdateInputsTemplate + inputs.DeleteInputsTemplate + inputs.KeyInputsTemplate + `

` + helpers.AtomicHelpersTemplate + helpers.StreamHelpersTemplate + helpers.ConverterHelpersTemplate + `
`
