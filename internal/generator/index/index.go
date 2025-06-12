// Package index defines the structure and validation logic for DynamoDB secondary indexes,
// including Global Secondary Indexes (GSI) and Local Secondary Indexes (LSI).
//
// It provides utilities for:
//   - Representing index metadata (keys, projection, capacity)
//   - Validating index configuration
//   - Handling composite key parts (e.g., user#type)
//
// This package is used internally during schema parsing and code generation.
package index

import (
	"strings"

	"github.com/Mad-Pixels/go-dyno/internal/logger"
	"github.com/Mad-Pixels/go-dyno/internal/utils"
)

// IndexType defines the type of secondary index
type IndexType string

const (
	// GSI represents a Global Secondary Index
	GSI IndexType = "GSI"
	// LSI represents a Local Secondary Index
	LSI IndexType = "LSI"
)

var (
	// validIndexesTypes list all supported DynamoDB indexes
	validIndexesTypes = map[string]bool{
		"GSI": true,
		"LSI": true,
	}

	// validProjectionTypes
	validProjectionTypes = map[string]bool{
		"ALL":       true,
		"KEYS_ONLY": true,
		"INCLUDE":   true,
	}
)

// String returns the string representation of IndexType
func (it IndexType) String() string {
	return string(it)
}

// Index represents either a Global Secondary Index (GSI) or Local Secondary Index (LSI)
type Index struct {
	// Name is the unique identifier for the index
	Name string `json:"name"`

	// Type specifies whether this is a "GSI" or "LSI"
	// If empty, defaults to "GSI" for backward compatibility
	Type IndexType `json:"type,omitempty"`

	// HashKey is the partition key for GSI only
	// For LSI, this field is ignored as LSI uses table's hash key
	// Can be a simple attribute name or composite key (e.g., "user#type")
	HashKey string `json:"hash_key,omitempty"`

	// RangeKey is the sort key for the index
	// Required for LSI, optional for GSI
	// Can be a simple attribute name or composite key
	RangeKey string `json:"range_key,omitempty"`

	// ProjectionType defines which attributes are included in the index
	// Valid values: "ALL", "KEYS_ONLY", "INCLUDE"
	ProjectionType string `json:"projection_type"`

	// NonKeyAttributes lists additional attributes when ProjectionType is "INCLUDE"
	NonKeyAttributes []string `json:"non_key_attributes,omitempty"`

	// Throughput settings - only valid for GSI
	// LSI uses the table's provisioned throughput
	ReadCapacity  *int `json:"read_capacity,omitempty"`
	WriteCapacity *int `json:"write_capacity,omitempty"`

	// Parsed composite key parts (populated during schema loading)
	HashKeyParts  []CompositeKey `json:"-"`
	RangeKeyParts []CompositeKey `json:"-"`
}

// SupportsQuery returns true if this index can be used for the given key pattern.
func (i Index) SupportsQuery(hasHashKey, hasRangeKey bool, tableHashKey string) bool {
	return i.GetEffectiveHashKey(tableHashKey) != ""
}

// GetTerraformType returns the Terraform resource type for this index.
func (i Index) GetTerraformType() string {
	if i.IsLSI() {
		return "local_secondary_index"
	}
	return "global_secondary_index"
}

// IsGSI returns true if this is a Global Secondary Index.
func (i Index) IsGSI() bool {
	return i.Type == GSI || i.Type == "" // Default to GSI for backward compatibility
}

// IsLSI returns true if this is a Local Secondary Index
func (i Index) IsLSI() bool {
	return i.Type == LSI
}

// GetEffectiveHashKey returns the hash key that will be used for this index
// For GSI: returns the specified HashKey
// For LSI: returns empty string (uses table's hash key)
func (i Index) GetEffectiveHashKey(tableHashKey string) string {
	if i.IsLSI() {
		return tableHashKey
	}
	return i.HashKey
}

// HasCompositeHashKey returns true if the hash key is composite (contains #)
func (si Index) HasCompositeHashKey() bool {
	return len(si.HashKeyParts) > 0
}

// Validate performs comprehensive validation of the secondary index configuration.
func (i Index) Validate(tableHashKey, tableRangeKey string) error {
	if !validIndexesTypes[strings.ToUpper(string(i.Type))] {
		return logger.NewFailure("invalid index", nil).
			With("name", i.Name).
			With("available", utils.AvailableKeys(validIndexesTypes))
	}
	if !validProjectionTypes[strings.ToUpper(i.ProjectionType)] {
		return logger.NewFailure("invalid projection type", nil).
			With("type", i.ProjectionType).
			With("available", utils.AvailableKeys(validProjectionTypes))
	}

	if strings.ToUpper(i.ProjectionType) == "INCLUDE" && len(i.NonKeyAttributes) == 0 {
		return logger.NewFailure("non_key_attributes must be specified when projection_type is 'INCLUDE'", nil)
	}
	if strings.ToUpper(i.ProjectionType) != "INCLUDE" && len(i.NonKeyAttributes) > 0 {
		return logger.NewFailure("non_key_attributes can only be specified when projection_type is 'INCLUDE'", nil)
	}

	if i.IsLSI() {
		if err := i.validateLSI(tableRangeKey); err != nil {
			return err
		}
	}
	if i.IsGSI() {
		if err := i.validateGSI(); err != nil {
			return err
		}
	}
	return nil
}

func (i Index) validateLSI(tableRangeKey string) error {
	if i.RangeKey == "" {
		return logger.NewFailure("LSI must specify range_key", nil).
			With("name", i.Name)
	}
	if i.RangeKey == tableRangeKey {
		return logger.NewFailure("range_key cannot be the same as table's range_key", nil).
			With("name", i.Name).
			With("key", i.RangeKey)
	}
	if i.ReadCapacity != nil || i.WriteCapacity != nil {
		return logger.NewFailure("LSI cannot specify read/write capacity (uses table's provisioned throughput)", nil).
			With("name", i.Name)
	}
	return nil
}

func (i Index) validateGSI() error {
	if i.HashKey == "" {
		return logger.NewFailure("GSI must specify hash_key", nil).
			With("name", i.Name)
	}
	return nil
}
