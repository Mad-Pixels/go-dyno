package scan

// ScanBuilderTemplate with mixins
const ScanBuilderTemplate = `
// ScanBuilder provides a fluent interface for building DynamoDB scans
type ScanBuilder struct {
    FilterMixin
    PaginationMixin
    IndexName            string
    ProjectionAttributes []string
    ParallelScanConfig   *ParallelScanConfig
}

// ParallelScanConfig configures parallel scan operations
type ParallelScanConfig struct {
    TotalSegments int
    Segment       int
}

// NewScanBuilder creates a new ScanBuilder instance
func NewScanBuilder() *ScanBuilder {
    return &ScanBuilder{
        FilterMixin:     NewFilterMixin(),
        PaginationMixin: NewPaginationMixin(),
    }
}

// Filter adds a filter condition and returns ScanBuilder for chaining
func (sb *ScanBuilder) Filter(field string, op OperatorType, values ...interface{}) *ScanBuilder {
    sb.FilterMixin.Filter(field, op, values...)
    return sb
}

// FilterEQ adds equality filter and returns ScanBuilder for chaining
func (sb *ScanBuilder) FilterEQ(field string, value interface{}) *ScanBuilder {
    sb.FilterMixin.FilterEQ(field, value)
    return sb
}

// FilterContains adds contains filter and returns ScanBuilder for chaining
func (sb *ScanBuilder) FilterContains(field string, value interface{}) *ScanBuilder {
    sb.FilterMixin.FilterContains(field, value)
    return sb
}

// FilterNotContains adds not contains filter and returns ScanBuilder for chaining
func (sb *ScanBuilder) FilterNotContains(field string, value interface{}) *ScanBuilder {
    sb.FilterMixin.FilterNotContains(field, value)
    return sb
}

// FilterBeginsWith adds begins_with filter and returns ScanBuilder for chaining
func (sb *ScanBuilder) FilterBeginsWith(field string, value interface{}) *ScanBuilder {
    sb.FilterMixin.FilterBeginsWith(field, value)
    return sb
}

// FilterBetween adds range filter and returns ScanBuilder for chaining
func (sb *ScanBuilder) FilterBetween(field string, start, end interface{}) *ScanBuilder {
    sb.FilterMixin.FilterBetween(field, start, end)
    return sb
}

// FilterGT adds greater than filter and returns ScanBuilder for chaining
func (sb *ScanBuilder) FilterGT(field string, value interface{}) *ScanBuilder {
    sb.FilterMixin.FilterGT(field, value)
    return sb
}

// FilterLT adds less than filter and returns ScanBuilder for chaining
func (sb *ScanBuilder) FilterLT(field string, value interface{}) *ScanBuilder {
    sb.FilterMixin.FilterLT(field, value)
    return sb
}

// FilterGTE adds greater than or equal filter and returns ScanBuilder for chaining
func (sb *ScanBuilder) FilterGTE(field string, value interface{}) *ScanBuilder {
    sb.FilterMixin.FilterGTE(field, value)
    return sb
}

// FilterLTE adds less than or equal filter and returns ScanBuilder for chaining
func (sb *ScanBuilder) FilterLTE(field string, value interface{}) *ScanBuilder {
    sb.FilterMixin.FilterLTE(field, value)
    return sb
}

// FilterExists adds attribute exists filter and returns ScanBuilder for chaining
func (sb *ScanBuilder) FilterExists(field string) *ScanBuilder {
    sb.FilterMixin.FilterExists(field)
    return sb
}

// FilterNotExists adds attribute not exists filter and returns ScanBuilder for chaining
func (sb *ScanBuilder) FilterNotExists(field string) *ScanBuilder {
    sb.FilterMixin.FilterNotExists(field)
    return sb
}

// FilterNE adds not equal filter and returns ScanBuilder for chaining
func (sb *ScanBuilder) FilterNE(field string, value interface{}) *ScanBuilder {
    sb.FilterMixin.FilterNE(field, value)
    return sb
}

// FilterIn adds IN filter and returns ScanBuilder for chaining
func (sb *ScanBuilder) FilterIn(field string, values ...interface{}) *ScanBuilder {
    sb.FilterMixin.FilterIn(field, values...)
    return sb
}

// FilterNotIn adds NOT_IN filter and returns ScanBuilder for chaining
func (sb *ScanBuilder) FilterNotIn(field string, values ...interface{}) *ScanBuilder {
    sb.FilterMixin.FilterNotIn(field, values...)
    return sb
}

// Limit sets the maximum number of items and returns ScanBuilder for chaining
func (sb *ScanBuilder) Limit(limit int) *ScanBuilder {
    sb.PaginationMixin.Limit(limit)
    return sb
}

// StartFrom sets the exclusive start key and returns ScanBuilder for chaining
func (sb *ScanBuilder) StartFrom(lastEvaluatedKey map[string]types.AttributeValue) *ScanBuilder {
    sb.PaginationMixin.StartFrom(lastEvaluatedKey)
    return sb
}

// WithIndex sets the index name for the scan
func (sb *ScanBuilder) WithIndex(indexName string) *ScanBuilder {
    sb.IndexName = indexName
    return sb
}

// WithProjection sets the projection attributes
func (sb *ScanBuilder) WithProjection(attributes []string) *ScanBuilder {
    sb.ProjectionAttributes = attributes
    return sb
}

// WithParallelScan configures parallel scan settings
func (sb *ScanBuilder) WithParallelScan(totalSegments, segment int) *ScanBuilder {
    sb.ParallelScanConfig = &ParallelScanConfig{
        TotalSegments: totalSegments,
        Segment:       segment,
    }
    return sb
}
`
