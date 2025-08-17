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
{{if IsALL .Mode}}
` + core.FilterMixinSugarTemplate + core.KeyConditionMixinSugarTemplate + `
{{end}}

` + query.QueryBuilderTemplate + query.QueryBuilderWithTemplate + query.QueryBuilderFilterTemplate + `
{{if IsALL .Mode}}
` + query.QueryBuilderWithSugarTemplate + query.QueryBuilderFilterSugarTemplate + `
{{end}}
` + query.QueryBuilderBuildTemplate + query.QueryBuilderUtilsTemplate + `

` + scan.ScanBuilderTemplate + scan.ScanBuilderFilterTemplate + `
{{if IsALL .Mode}}
` + scan.ScanBuilderFilterSugarTemplate + `
{{end}}
` + scan.ScanBuilderBuildTemplate + `

` + inputs.ItemInputsTemplate + inputs.UpdateInputsTemplate + inputs.DeleteInputsTemplate + inputs.KeyInputsTemplate + `

` + helpers.AtomicHelpersTemplate + `
{{if .UseStreamEvents}}
` + helpers.StreamHelpersTemplate + `
{{end}}
` + helpers.ConverterHelpersTemplate + helpers.MarshalingHelpersTemplate + helpers.ValidationHelpersTemplate + `
`
