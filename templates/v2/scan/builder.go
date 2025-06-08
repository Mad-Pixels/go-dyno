package scan

// ScanBuilderTemplate ...
const ScanBuilderTemplate = `
// ScanBuilder provides a fluent interface for building DynamoDB scans
type ScanBuilder struct {
    IndexName            string
    FilterConditions     []expression.ConditionBuilder
    UsedKeys             map[string]bool
    Attributes           map[string]interface{}
    LimitValue           *int
    ExclusiveStartKey    map[string]types.AttributeValue
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
        FilterConditions: make([]expression.ConditionBuilder, 0),
        UsedKeys:         make(map[string]bool),
        Attributes:       make(map[string]interface{}),
    }
}

// WithIndex sets the index name for the scan
func (sb *ScanBuilder) WithIndex(indexName string) *ScanBuilder {
    sb.IndexName = indexName
    return sb
}

// Limit sets the maximum number of items to return
func (sb *ScanBuilder) Limit(limit int) *ScanBuilder {
    sb.LimitValue = &limit
    return sb
}

// StartFrom sets the exclusive start key for pagination
func (sb *ScanBuilder) StartFrom(lastEvaluatedKey map[string]types.AttributeValue) *ScanBuilder {
    sb.ExclusiveStartKey = lastEvaluatedKey
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
