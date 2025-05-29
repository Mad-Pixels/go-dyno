package schema

import (
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
//	input: "user-activity" → output: "user_activity"
func (ds DynamoSchema) Directory() string {
	return utils.ToSafeName(ds.schema.TableName)
}

// Filename returns the generated Go filename for the table schema,
// using the safe name of the table.
//
// Example:
//
//	input: "user-activity" → output: "user_activity.go"
func (ds DynamoSchema) Filename() string {
	return utils.ToSafeName(ds.schema.TableName) + ".go"
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

// SecondaryIndexes returns all defined GSI/LSI indexes attached to the table.
//
// Example:
//
//	schema.SecondaryIndexes() → []SecondaryIndex{...}
func (ds DynamoSchema) SecondaryIndexes() []common.SecondaryIndex {
	return ds.schema.SecondaryIndexes
}

// AllAttributes returns the combined list of core and common attributes.
//
// Example:
//
//	schema.AllAttributes() → append(schema.Attributes(), schema.CommonAttributes()...)
func (ds DynamoSchema) AllAttributes() []common.Attribute {
	return append(ds.Attributes(), ds.CommonAttributes()...)
}
