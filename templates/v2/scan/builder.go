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
// Example: scan := NewScanBuilder().FilterEQ("status", "active").Limit(100)
func NewScanBuilder() *ScanBuilder {
    return &ScanBuilder{
        FilterMixin:     NewFilterMixin(),
        PaginationMixin: NewPaginationMixin(),
    }
}

// Filter adds a filter condition and returns ScanBuilder for method chaining.
// Filters are applied after items are read from DynamoDB.
// Example: scan.Filter("score", GT, 80)
func (sb *ScanBuilder) Filter(field string, op OperatorType, values ...any) *ScanBuilder {
    sb.FilterMixin.Filter(field, op, values...)
    return sb
}

// FilterEQ adds equality filter and returns ScanBuilder for method chaining.
// Example: scan.FilterEQ("status", "active")
func (sb *ScanBuilder) FilterEQ(field string, value any) *ScanBuilder {
    sb.FilterMixin.FilterEQ(field, value)
    return sb
}

// FilterContains adds contains filter and returns ScanBuilder for method chaining.
// Works with String attributes (substring) and Set attributes (membership).
// Example: scan.FilterContains("tags", "premium")
func (sb *ScanBuilder) FilterContains(field string, value any) *ScanBuilder {
    sb.FilterMixin.FilterContains(field, value)
    return sb
}

// FilterNotContains adds not contains filter and returns ScanBuilder for method chaining.
// Opposite of FilterContains for exclusion filtering.
func (sb *ScanBuilder) FilterNotContains(field string, value any) *ScanBuilder {
    sb.FilterMixin.FilterNotContains(field, value)
    return sb
}

// FilterBeginsWith adds begins_with filter and returns ScanBuilder for method chaining.
// Only works with String attributes for prefix matching.
// Example: scan.FilterBeginsWith("email", "admin@")
func (sb *ScanBuilder) FilterBeginsWith(field string, value any) *ScanBuilder {
    sb.FilterMixin.FilterBeginsWith(field, value)
    return sb
}

// FilterBetween adds range filter and returns ScanBuilder for method chaining.
// Works with comparable types for inclusive range filtering.
// Example: scan.FilterBetween("score", 80, 100)
func (sb *ScanBuilder) FilterBetween(field string, start, end any) *ScanBuilder {
    sb.FilterMixin.FilterBetween(field, start, end)
    return sb
}

// FilterGT adds greater than filter and returns ScanBuilder for method chaining.
// Example: scan.FilterGT("last_login", cutoffDate)
func (sb *ScanBuilder) FilterGT(field string, value any) *ScanBuilder {
    sb.FilterMixin.FilterGT(field, value)
    return sb
}

// FilterLT adds less than filter and returns ScanBuilder for method chaining.
// Example: scan.FilterLT("attempts", maxAttempts)
func (sb *ScanBuilder) FilterLT(field string, value any) *ScanBuilder {
    sb.FilterMixin.FilterLT(field, value)
    return sb
}

// FilterGTE adds greater than or equal filter and returns ScanBuilder for method chaining.
// Example: scan.FilterGTE("age", minimumAge)
func (sb *ScanBuilder) FilterGTE(field string, value any) *ScanBuilder {
    sb.FilterMixin.FilterGTE(field, value)
    return sb
}

// FilterLTE adds less than or equal filter and returns ScanBuilder for method chaining.
// Example: scan.FilterLTE("file_size", maxFileSize)
func (sb *ScanBuilder) FilterLTE(field string, value any) *ScanBuilder {
    sb.FilterMixin.FilterLTE(field, value)
    return sb
}

// FilterExists adds attribute exists filter and returns ScanBuilder for method chaining.
// Checks if the specified attribute exists in the item.
// Example: scan.FilterExists("optional_field")
func (sb *ScanBuilder) FilterExists(field string) *ScanBuilder {
    sb.FilterMixin.FilterExists(field)
    return sb
}

// FilterNotExists adds attribute not exists filter and returns ScanBuilder for method chaining.
// Checks if the specified attribute does not exist in the item.
func (sb *ScanBuilder) FilterNotExists(field string) *ScanBuilder {
    sb.FilterMixin.FilterNotExists(field)
    return sb
}

// FilterNE adds not equal filter and returns ScanBuilder for method chaining.
// Example: scan.FilterNE("status", "deleted")
func (sb *ScanBuilder) FilterNE(field string, value any) *ScanBuilder {
    sb.FilterMixin.FilterNE(field, value)
    return sb
}

// FilterIn adds IN filter and returns ScanBuilder for method chaining.
// For scalar values only - use FilterContains for DynamoDB Sets.
// Example: scan.FilterIn("category", "books", "electronics", "clothing")
func (sb *ScanBuilder) FilterIn(field string, values ...any) *ScanBuilder {
    sb.FilterMixin.FilterIn(field, values...)
    return sb
}

// FilterNotIn adds NOT_IN filter and returns ScanBuilder for method chaining.
// For scalar values only - use FilterNotContains for DynamoDB Sets.
func (sb *ScanBuilder) FilterNotIn(field string, values ...any) *ScanBuilder {
    sb.FilterMixin.FilterNotIn(field, values...)
    return sb
}

// Limit sets the maximum number of items and returns ScanBuilder for method chaining.
// Controls the number of items returned in a single scan request.
// Note: DynamoDB may return fewer items due to size limits even with this setting.
// Example: scan.Limit(100)
func (sb *ScanBuilder) Limit(limit int) *ScanBuilder {
    sb.PaginationMixin.Limit(limit)
    return sb
}

// StartFrom sets the exclusive start key and returns ScanBuilder for method chaining.
// Use LastEvaluatedKey from previous response for pagination.
// Example: scan.StartFrom(previousResponse.LastEvaluatedKey)
func (sb *ScanBuilder) StartFrom(lastEvaluatedKey map[string]types.AttributeValue) *ScanBuilder {
    sb.PaginationMixin.StartFrom(lastEvaluatedKey)
    return sb
}

// WithIndex sets the index name for scanning a secondary index.
// Allows scanning GSI or LSI instead of the main table.
// Index must exist and be in ACTIVE state.
// Example: scan.WithIndex("status-index")
func (sb *ScanBuilder) WithIndex(indexName string) *ScanBuilder {
    sb.IndexName = indexName
    return sb
}

// WithProjection sets the projection attributes to return specific fields only.
// Reduces network traffic and costs by returning only needed attributes.
// Pass attribute names that should be included in the response.
// Example: scan.WithProjection([]string{"id", "name", "status"})
func (sb *ScanBuilder) WithProjection(attributes []string) *ScanBuilder {
    sb.ProjectionAttributes = attributes
    return sb
}

// WithParallelScan configures parallel scan settings for improved throughput.
// Divides the table into segments for concurrent processing by multiple workers.
// totalSegments: how many segments to divide the table (typically number of workers)
// segment: which segment this worker processes (0-based, must be < totalSegments)
// Example: scan.WithParallelScan(4, 0) // Process segment 0 of 4 total segments
func (sb *ScanBuilder) WithParallelScan(totalSegments, segment int) *ScanBuilder {
    sb.ParallelScanConfig = &ParallelScanConfig{
        TotalSegments: totalSegments,
        Segment:       segment,
    }
    return sb
}
`
