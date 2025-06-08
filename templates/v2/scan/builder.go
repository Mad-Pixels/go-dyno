package scan

// ScanBuilderTemplate ...
const ScanBuilderTemplate = `
// ScanBuilder ...
type ScanBuilder struct {
    IndexName           string
    FilterConditions    []expression.ConditionBuilder
    UsedKeys            map[string]bool
    Attributes          map[string]interface{}
    LimitValue          *int
    ExclusiveStartKey   map[string]types.AttributeValue
    ProjectionAttributes []string
    ParallelScanConfig  *ParallelScanConfig
}

// ParallelScanConfig ...
type ParallelScanConfig struct {
    TotalSegments int
    Segment int
}

// NewScanBuilder ...
func NewScanBuilder() *ScanBuilder {
    return &ScanBuilder{
        FilterConditions: make([]expression.ConditionBuilder, 0),
        UsedKeys:         make(map[string]bool),
        Attributes:       make(map[string]interface{}),
    }
}

// WithIndex ...
func (sb *ScanBuilder) WithIndex(indexName string) *ScanBuilder {
    sb.IndexName = indexName
    return sb
}

{{range .AllAttributes}}
// Filter{{ToSafeName .Name | ToUpperCamelCase}} ...
func (sb *ScanBuilder) Filter{{ToSafeName .Name | ToUpperCamelCase}}({{ToSafeName .Name | ToLowerCamelCase}} {{ToGolangBaseType .}}) *ScanBuilder {
    condition := expression.Name("{{.Name}}").Equal(expression.Value({{ToSafeName .Name | ToLowerCamelCase}}))
    sb.FilterConditions = append(sb.FilterConditions, condition)
    sb.Attributes["{{.Name}}"] = {{ToSafeName .Name | ToLowerCamelCase}}
    sb.UsedKeys["{{.Name}}"] = true
    return sb
}
{{end}}

{{range .AllAttributes}}
{{if IsNumericAttr .}}
// Filter{{ToSafeName .Name | ToUpperCamelCase}}Between ...
func (sb *ScanBuilder) Filter{{ToSafeName .Name | ToUpperCamelCase}}Between(start, end {{ToGolangBaseType .}}) *ScanBuilder {
    condition := expression.Name("{{.Name}}").Between(expression.Value(start), expression.Value(end))
    sb.FilterConditions = append(sb.FilterConditions, condition)
    return sb
}

// Filter{{ToSafeName .Name | ToUpperCamelCase}}GreaterThan ...
func (sb *ScanBuilder) Filter{{ToSafeName .Name | ToUpperCamelCase}}GreaterThan(value {{ToGolangBaseType .}}) *ScanBuilder {
    condition := expression.Name("{{.Name}}").GreaterThan(expression.Value(value))
    sb.FilterConditions = append(sb.FilterConditions, condition)
    return sb
}

// Filter{{ToSafeName .Name | ToUpperCamelCase}}LessThan ... 
func (sb *ScanBuilder) Filter{{ToSafeName .Name | ToUpperCamelCase}}LessThan(value {{ToGolangBaseType .}}) *ScanBuilder {
    condition := expression.Name("{{.Name}}").LessThan(expression.Value(value))
    sb.FilterConditions = append(sb.FilterConditions, condition)
    return sb
}

// Filter{{ToSafeName .Name | ToUpperCamelCase}}GreaterThanOrEqual ...
func (sb *ScanBuilder) Filter{{ToSafeName .Name | ToUpperCamelCase}}GreaterThanOrEqual(value {{ToGolangBaseType .}}) *ScanBuilder {
    condition := expression.Name("{{.Name}}").GreaterThanEqual(expression.Value(value))
    sb.FilterConditions = append(sb.FilterConditions, condition)
    return sb
}

// Filter{{ToSafeName .Name | ToUpperCamelCase}}LessThanOrEqual ...
func (sb *ScanBuilder) Filter{{ToSafeName .Name | ToUpperCamelCase}}LessThanOrEqual(value {{ToGolangBaseType .}}) *ScanBuilder {
    condition := expression.Name("{{.Name}}").LessThanEqual(expression.Value(value))
    sb.FilterConditions = append(sb.FilterConditions, condition)
    return sb
}
{{end}}
{{end}}

{{range .AllAttributes}}
{{if eq (ToGolangBaseType .) "string"}}
// Filter{{ToSafeName .Name | ToUpperCamelCase}}Contains ...
func (sb *ScanBuilder) Filter{{ToSafeName .Name | ToUpperCamelCase}}Contains(value {{ToGolangBaseType .}}) *ScanBuilder {
    condition := expression.Name("{{.Name}}").Contains(value)
    sb.FilterConditions = append(sb.FilterConditions, condition)
    return sb
}

// Filter{{ToSafeName .Name | ToUpperCamelCase}}BeginsWith ...
func (sb *ScanBuilder) Filter{{ToSafeName .Name | ToUpperCamelCase}}BeginsWith(value {{ToGolangBaseType .}}) *ScanBuilder {
    condition := expression.Name("{{.Name}}").BeginsWith(value)
    sb.FilterConditions = append(sb.FilterConditions, condition)
    return sb
}
{{end}}
{{end}}

// Limit ...
func (sb *ScanBuilder) Limit(limit int) *ScanBuilder {
    sb.LimitValue = &limit
    return sb
}

// StartFrom ...
func (sb *ScanBuilder) StartFrom(lastEvaluatedKey map[string]types.AttributeValue) *ScanBuilder {
    sb.ExclusiveStartKey = lastEvaluatedKey
    return sb
}

// WithProjection ...
func (sb *ScanBuilder) WithProjection(attributes []string) *ScanBuilder {
    sb.ProjectionAttributes = attributes
    return sb
}

// WithParallelScan ...
func (sb *ScanBuilder) WithParallelScan(totalSegments, segment int) *ScanBuilder {
    sb.ParallelScanConfig = &ParallelScanConfig{
        TotalSegments: totalSegments,
        Segment:       segment,
    }
    return sb
}
`
