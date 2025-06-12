// Package schema provides a high-level representation of a DynamoDB table schema,
// including attributes, primary keys, and secondary indexes (GSI and LSI).
//
// It offers utilities for loading schema definitions from JSON files, validating
// schema correctness, resolving optimal indexes for queries, and generating
// Go-safe identifiers such as package names and filenames.
//
// This package is used as an internal model in the code generation pipeline.
package schema

import (
	"github.com/Mad-Pixels/go-dyno/internal/generator/attribute"
	"github.com/Mad-Pixels/go-dyno/internal/generator/index"
	"github.com/Mad-Pixels/go-dyno/internal/logger"
	"github.com/Mad-Pixels/go-dyno/internal/utils/conv"
	"github.com/Mad-Pixels/go-dyno/internal/utils/fs"
)

// Schema wraps the raw schema definition.
type Schema struct {
	raw schema
}

// NewSchema loads and parses a schema definition from the given file path.
func NewSchema(path string) (*Schema, error) {
	var spec Schema

	if err := fs.ReadAndParseJSON(path, &spec.raw); err != nil {
		return nil, err
	}
	return &spec, nil
}

// TableName returns the logical name of the DynamoDB table defined in the schema.
func (s Schema) TableName() string {
	return s.raw.TableName
}

// HashKey returns the primary partition key of the table.
func (s Schema) HashKey() string {
	return s.raw.HashKey
}

// RangeKey returns the primary sort key of the table, if defined.
func (s Schema) RangeKey() string {
	return s.raw.RangeKey
}

// PackageName returns a Go-safe lowercase package name derived from the table name.
func (s Schema) PackageName() string {
	return conv.ToLowerInlineCase(s.raw.TableName)
}

// Filename returns the default Go filename for generated code based on the table name.
func (s Schema) Filename() string {
	return fs.AddFileExt(s.PackageName(), ".go")
}

// Attributes returns the key attributes defined in the schema.
func (s Schema) Attributes() []attribute.Attribute {
	return s.raw.Attributes
}

// CommonAttributes returns attributes that are not used in keys or indexes.
func (s Schema) CommonAttributes() []attribute.Attribute {
	return s.raw.CommonAttributes
}

// AllAttributes returns both key and common attributes combined.
func (s Schema) AllAttributes() []attribute.Attribute {
	return append(s.Attributes(), s.CommonAttributes()...)
}

// SecondaryIndexes returns all secondary indexes (GSI and LSI) defined in the schema.
func (s Schema) SecondaryIndexes() []index.Index {
	return s.raw.SecondaryIndexes
}

// GlobalSecondaryIndexes returns only the GSIs (Global Secondary Indexes).
func (s Schema) GlobalSecondaryIndexes() []index.Index {
	return s.filterIndexesByType(func(idx index.Index) bool { return idx.IsGSI() })
}

// LocalSecondaryIndexes returns only the LSIs (Local Secondary Indexes).
func (s Schema) LocalSecondaryIndexes() []index.Index {
	return s.filterIndexesByType(func(idx index.Index) bool { return idx.IsLSI() })
}

// HasGSI returns true if the schema defines any Global Secondary Indexes.
func (s Schema) HasGSI() bool {
	return len(s.GlobalSecondaryIndexes()) > 0
}

// HasLSI returns true if the schema defines any Local Secondary Indexes.
func (s Schema) HasLSI() bool {
	return len(s.LocalSecondaryIndexes()) > 0
}

// HasSecondaryIndexes returns true if any secondary indexes (GSI or LSI) are defined.
func (s Schema) HasSecondaryIndexes() bool {
	return len(s.raw.SecondaryIndexes) > 0
}

// GetIndexByName returns a pointer to the index with the given name, or nil if not found.
func (s Schema) GetIndexByName(name string) *index.Index {
	for i := range s.raw.SecondaryIndexes {
		if s.raw.SecondaryIndexes[i].Name == name {
			return &s.raw.SecondaryIndexes[i]
		}
	}
	return nil
}

// ValidateIndexNames checks for duplicate index names.
func (s Schema) ValidateIndexNames() error {
	seen := make(map[string]bool)

	for _, idx := range s.SecondaryIndexes() {
		if seen[idx.Name] {
			return logger.NewFailure("duplicate index name", nil).
				With("name", idx.Name)
		}
		seen[idx.Name] = true
	}
	return nil
}

// GetOptimalIndexForQuery returns the most efficient index for a query with the given keys.
//
// Index selection priority:
//  1. LSI if hash key matches main table and range key is present
//  2. GSI if hash key matches a GSI
//  3. nil (fallback to primary table)
func (s Schema) GetOptimalIndexForQuery(hashKey string) *index.Index {
	var hasHashKey = hashKey != ""

	for _, idx := range s.LocalSecondaryIndexes() {
		if idx.SupportsQuery(s.HashKey()) {
			if hasHashKey && hashKey == s.HashKey() {
				return &idx
			}
		}
	}
	for _, idx := range s.GlobalSecondaryIndexes() {
		if idx.SupportsQuery(s.HashKey()) {
			if hasHashKey && hashKey == idx.HashKey {
				return &idx
			}
		}
	}
	return nil
}

type schema struct {
	// TableName defines the logical name of the DynamoDB table.
	// It is used for generating Go package, filename, and identifiers in the template.
	// Go package and filename can be overridden by generator config.
	TableName string `json:"table_name"`

	// HashKey is the primary partition key of the DynamoDB table.
	// This field is required and must match one of the attribute names.
	HashKey string `json:"hash_key"`

	// RangeKey is the optional sort key of the DynamoDB table.
	// If defined, it must also match one of the attribute names.
	RangeKey string `json:"range_key"`

	// Attributes define the key attributes that can be used in primary keys
	// and secondary indexes (hash_key, range_key for GSI/LSI).
	// These fields must be defined for DynamoDB key operations.
	Attributes []attribute.Attribute `json:"attributes"`

	// CommonAttributes define non-key attributes that are only stored as data
	// and cannot be used as keys in indexes. These are regular table fields.
	CommonAttributes []attribute.Attribute `json:"common_attributes"`

	// SecondaryIndexes defines Global or Local Secondary Indexes (GSI/LSI)
	// used for advanced querying in DynamoDB. Each index has its own keys and projection.
	SecondaryIndexes []index.Index `json:"secondary_indexes"`
}

func (s Schema) filterIndexesByType(predicate func(index.Index) bool) []index.Index {
	var result []index.Index
	for _, idx := range s.raw.SecondaryIndexes {
		if predicate(idx) {
			result = append(result, idx)
		}
	}
	return result
}
