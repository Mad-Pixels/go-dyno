package schema

import (
	"fmt"

	"github.com/Mad-Pixels/go-dyno/internal/schema/common"
	"github.com/Mad-Pixels/go-dyno/internal/utils"
)

type dynamoSchema struct {
	// TableName defines the logical name of the DynamoDB table.
	// It is used for generating Go package, filename, and identifiers in the template.
	TableName string `json:"table_name"`

	// HashKey is the primary partition key of the DynamoDB table.
	// This field is required and must match one of the attribute names.
	HashKey string `json:"hash_key"`

	// RangeKey is the optional sort key of the DynamoDB table.
	// If defined, it must also match one of the attribute names.
	RangeKey string `json:"range_key"`

	// Attributes are the user-defined fields specific to this table.
	// Each attribute has a name and a type ("S", "N", "B").
	Attributes []common.Attribute `json:"attributes"`

	// CommonAttributes are additional fields shared across multiple tables,
	// such as audit fields (created_at, updated_at).
	CommonAttributes []common.Attribute `json:"common_attributes"`

	// SecondaryIndexes defines Global or Local Secondary Indexes (GSI/LSI)
	// used for advanced querying in DynamoDB. Each index has its own keys and projection.
	SecondaryIndexes []common.SecondaryIndex `json:"secondary_indexes"`
}

// DynamoSchema wraps the internal schema definition and provides utility accessors
// for table name, keys, attributes, and index metadata.
type DynamoSchema struct {
	schema dynamoSchema
}

// TableName returns the table name in UpperCamelCase format,
// used typically as a Go struct or type name.
//
// Example:
//
//	input: "user_activity" → output: "UserActivity"
func (ds DynamoSchema) TableName() string {
	return ds.schema.TableName
}

// HashKey returns the name of the table's partition (hash) key.
//
// Example:
//
//	schema.HashKey() → "user_id"
func (ds DynamoSchema) HashKey() string {
	return ds.schema.HashKey
}

// RangeKey returns the name of the table's sort (range) key, if defined.
//
// Example:
//
//	schema.RangeKey() → "created_at"
func (ds DynamoSchema) RangeKey() string {
	return ds.schema.RangeKey
}

// PackageName returns the normalized lowercased name for Go package generation,
// derived from the table name.
//
// Example:
//
//	input: "user_activity" → output: "useractivity"
func (ds DynamoSchema) PackageName() string {
	return utils.ToLowerInlineCase(ds.schema.TableName)
}

// Directory returns a safe directory name based on the table name.
// Ensures only valid characters are used.
//
// Example:
//
//	input: "user-activity" → output: "useractivity"
//	input: "user_activity" → output: "useractivity"
func (ds DynamoSchema) Directory() string {
	return utils.ToLowerInlineCase(ds.schema.TableName)
}

// Filename returns the generated Go filename for the table schema,
// using the safe name of the table.
//
// Example:
//
//	input: "user-activity" → output: "useractivity.go"
//	input: "user_activity" → output: "useractivity.go"
func (ds DynamoSchema) Filename() string {
	return utils.ToLowerInlineCase(ds.schema.TableName) + ".go"
}

// Attributes returns the core attributes defined for the DynamoDB table schema.
//
// Example:
//
//	schema.Attributes() → []Attribute{{Name: "id", Type: "S"}, ...}
func (ds DynamoSchema) Attributes() []common.Attribute {
	return ds.schema.Attributes
}

// CommonAttributes returns the additional shared attributes used across tables,
// useful for consistent fields like timestamps or audit fields.
//
// Example:
//
//	schema.CommonAttributes() → []Attribute{{Name: "created_at", Type: "S"}}
func (ds DynamoSchema) CommonAttributes() []common.Attribute {
	return ds.schema.CommonAttributes
}

// SecondaryIndexes returns all defined secondary indexes (both GSI and LSI)
func (ds DynamoSchema) SecondaryIndexes() []common.SecondaryIndex {
	return ds.schema.SecondaryIndexes
}

// GlobalSecondaryIndexes returns only Global Secondary Indexes
func (ds DynamoSchema) GlobalSecondaryIndexes() []common.SecondaryIndex {
	var gsiIndexes []common.SecondaryIndex
	for _, idx := range ds.schema.SecondaryIndexes {
		if idx.IsGSI() {
			gsiIndexes = append(gsiIndexes, idx)
		}
	}
	return gsiIndexes
}

// LocalSecondaryIndexes returns only Local Secondary Indexes
func (ds DynamoSchema) LocalSecondaryIndexes() []common.SecondaryIndex {
	var lsiIndexes []common.SecondaryIndex
	for _, idx := range ds.schema.SecondaryIndexes {
		if idx.IsLSI() {
			lsiIndexes = append(lsiIndexes, idx)
		}
	}
	return lsiIndexes
}

// HasGSI returns true if the schema has any Global Secondary Indexes
func (ds DynamoSchema) HasGSI() bool {
	return len(ds.GlobalSecondaryIndexes()) > 0
}

// HasLSI returns true if the schema has any Local Secondary Indexes
func (ds DynamoSchema) HasLSI() bool {
	return len(ds.LocalSecondaryIndexes()) > 0
}

// HasSecondaryIndexes returns true if the schema has any secondary indexes
func (ds DynamoSchema) HasSecondaryIndexes() bool {
	return len(ds.schema.SecondaryIndexes) > 0
}

// GetIndexByName returns the secondary index with the given name, or nil if not found
func (ds DynamoSchema) GetIndexByName(name string) *common.SecondaryIndex {
	for _, idx := range ds.schema.SecondaryIndexes {
		if idx.Name == name {
			return &idx
		}
	}
	return nil
}

// GetOptimalIndexForQuery returns the best index for a query with the given keys
// Priority: LSI (cheaper) -> GSI -> nil (use main table)
func (ds DynamoSchema) GetOptimalIndexForQuery(hashKey, rangeKey string) *common.SecondaryIndex {
	hasHashKey := hashKey != ""
	hasRangeKey := rangeKey != ""

	// First, try LSI (they're cheaper and use table's throughput)
	for _, idx := range ds.LocalSecondaryIndexes() {
		if idx.SupportsQuery(hasHashKey, hasRangeKey, ds.HashKey()) {
			// Additional check for LSI: hash key must match table's hash key
			if hasHashKey && hashKey == ds.HashKey() {
				return &idx
			}
		}
	}

	// Then, try GSI
	for _, idx := range ds.GlobalSecondaryIndexes() {
		if idx.SupportsQuery(hasHashKey, hasRangeKey, ds.HashKey()) {
			// Additional check for GSI: hash key must match index's hash key
			if hasHashKey && hashKey == idx.HashKey {
				return &idx
			}
		}
	}

	// No suitable index found
	return nil
}

// ValidateIndexNames ensures all index names are unique
func (ds DynamoSchema) ValidateIndexNames() error {
	seen := make(map[string]bool)
	for _, idx := range ds.schema.SecondaryIndexes {
		if seen[idx.Name] {
			return fmt.Errorf("duplicate index name: '%s'", idx.Name)
		}
		seen[idx.Name] = true
	}
	return nil
}

// AllAttributes returns the combined list of core and common attributes.
//
// Example:
//
//	schema.AllAttributes() → append(schema.Attributes(), schema.CommonAttributes()...)
func (ds DynamoSchema) AllAttributes() []common.Attribute {
	return append(ds.Attributes(), ds.CommonAttributes()...)
}
