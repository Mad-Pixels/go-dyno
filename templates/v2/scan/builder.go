package scan

// ScanBuilderTemplate provides the main ScanBuilder struct for DynamoDB table scanning
const ScanBuilderTemplate = `
// ScanBuilder provides a fluent interface for building DynamoDB scan operations.
// Scans read every item in a table or index, applying filters after data is read.
// Use Query for efficient key-based access; use Scan for full table analysis.
// Combines FilterMixin and PaginationMixin for comprehensive scan functionality.
type ScanBuilder struct {
    FilterMixin                               // Filter conditions applied after reading items
    PaginationMixin                           // Limit and pagination support
    IndexName            string               // Optional secondary index to scan
    ProjectionAttributes []string             // Specific attributes to return
    ParallelScanConfig   *ParallelScanConfig  // Parallel scan configuration
}

// ParallelScanConfig configures parallel scan operations for improved throughput.
// Divides the table into segments that can be scanned concurrently.
// Each worker scans one segment, reducing overall scan time for large tables.
type ParallelScanConfig struct {
    TotalSegments int  // Total number of segments to divide the table into
    Segment       int  // Which segment this scan worker should process (0-based)
}

// NewScanBuilder creates a new ScanBuilder instance with initialized mixins.
// All mixins are properly initialized for immediate use.
func NewScanBuilder() *ScanBuilder {
    return &ScanBuilder{
        FilterMixin:     NewFilterMixin(),
        PaginationMixin: NewPaginationMixin(),
    }
}

// Limit sets the maximum number of items and returns ScanBuilder for method chaining.
// Controls the number of items returned in a single scan request.
func (sb *ScanBuilder) Limit(limit int) *ScanBuilder {
    sb.PaginationMixin.Limit(limit)
    return sb
}

// StartFrom sets the exclusive start key and returns ScanBuilder for method chaining.
// Use LastEvaluatedKey from previous response for pagination.
func (sb *ScanBuilder) StartFrom(lastEvaluatedKey map[string]types.AttributeValue) *ScanBuilder {
    sb.PaginationMixin.StartFrom(lastEvaluatedKey)
    return sb
}

// WithIndex sets the index name for scanning a secondary index.
// Allows scanning GSI or LSI instead of the main table.
// Index must exist and be in ACTIVE state.
func (sb *ScanBuilder) WithIndex(indexName string) *ScanBuilder {
    sb.IndexName = indexName
    return sb
}

// WithProjection sets the projection attributes to return specific fields only.
// Reduces network traffic and costs by returning only needed attributes.
// Pass attribute names that should be included in the response.
func (sb *ScanBuilder) WithProjection(attributes []string) *ScanBuilder {
    sb.ProjectionAttributes = attributes
    return sb
}

// WithParallelScan configures parallel scan settings for improved throughput.
// Divides the table into segments for concurrent processing by multiple workers.
// totalSegments: how many segments to divide the table (typically number of workers)
// segment: which segment this worker processes (0-based, must be < totalSegments)
func (sb *ScanBuilder) WithParallelScan(totalSegments, segment int) *ScanBuilder {
    sb.ParallelScanConfig = &ParallelScanConfig{
        TotalSegments: totalSegments,
        Segment:       segment,
    }
    return sb
}
`
