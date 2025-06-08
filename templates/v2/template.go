package v2

import (
	"github.com/Mad-Pixels/go-dyno/templates/v2/core"
	"github.com/Mad-Pixels/go-dyno/templates/v2/generic"
	"github.com/Mad-Pixels/go-dyno/templates/v2/helpers"
	"github.com/Mad-Pixels/go-dyno/templates/v2/inputs"
	"github.com/Mad-Pixels/go-dyno/templates/v2/query"
	"github.com/Mad-Pixels/go-dyno/templates/v2/scan"
)

// CodeTemplate with mixins and optimized operators
const CodeTemplate = `
package {{.PackageName}}

` + core.ImportsTemplate + `

` + core.ConstantsTemplate + `

` + generic.OperatorsTemplate + `

` + core.SchemaTemplate + `

` + core.MixinsTemplate + `

` + query.QueryBuilderTemplate + query.QueryBuilderBuildTemplate + query.QueryBuilderUtilsTemplate + `

` + scan.ScanBuilderTemplate + scan.ScanBuilderBuildTemplate + `

` + inputs.ItemInputsTemplate + inputs.UpdateInputsTemplate + inputs.DeleteInputsTemplate + inputs.KeyInputsTemplate + `

` + helpers.AtomicHelpersTemplate + helpers.StreamHelpersTemplate + helpers.ConverterHelpersTemplate + helpers.MarshalingHelpersTemplate + helpers.ValidationHelpersTemplate + `
`
