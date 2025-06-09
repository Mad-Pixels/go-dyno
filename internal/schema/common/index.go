package common

import "fmt"

// IndexType defines the type of secondary index
type IndexType string

const (
	// GSI represents a Global Secondary Index
	GSI IndexType = "GSI"
	// LSI represents a Local Secondary Index
	LSI IndexType = "LSI"
)

// String returns the string representation of IndexType
func (it IndexType) String() string {
	return string(it)
}

// SecondaryIndex represents either a Global Secondary Index (GSI) or Local Secondary Index (LSI)
// for a DynamoDB table. The structure is designed to handle both types efficiently
// while maintaining type safety and validation.
type SecondaryIndex struct {
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
	HashKeyParts  []CompositeKeyPart `json:"-"`
	RangeKeyParts []CompositeKeyPart `json:"-"`
}

// IsGSI returns true if this is a Global Secondary Index
func (si SecondaryIndex) IsGSI() bool {
	return si.Type == GSI || si.Type == "" // Default to GSI for backward compatibility
}

// IsLSI returns true if this is a Local Secondary Index
func (si SecondaryIndex) IsLSI() bool {
	return si.Type == LSI
}

// GetEffectiveHashKey returns the hash key that will be used for this index
// For GSI: returns the specified HashKey
// For LSI: returns empty string (uses table's hash key)
func (si SecondaryIndex) GetEffectiveHashKey(tableHashKey string) string {
	if si.IsLSI() {
		return tableHashKey
	}
	return si.HashKey
}

// HasCompositeHashKey returns true if the hash key is composite (contains #)
func (si SecondaryIndex) HasCompositeHashKey() bool {
	return len(si.HashKeyParts) > 0
}

// HasCompositeRangeKey returns true if the range key is composite (contains #)
func (si SecondaryIndex) HasCompositeRangeKey() bool {
	return len(si.RangeKeyParts) > 0
}

// Validate performs comprehensive validation of the secondary index configuration
func (si SecondaryIndex) Validate(tableHashKey, tableRangeKey string) error {
	// Basic validation
	if si.Name == "" {
		return fmt.Errorf("index name cannot be empty")
	}

	// Validate index type
	if si.Type != "" && si.Type != GSI && si.Type != LSI {
		return fmt.Errorf("index type must be 'GSI' or 'LSI', got '%s'", si.Type)
	}

	// Validate projection type
	validProjectionTypes := map[string]bool{
		"ALL":       true,
		"KEYS_ONLY": true,
		"INCLUDE":   true,
	}
	if !validProjectionTypes[si.ProjectionType] {
		return fmt.Errorf("invalid projection type '%s', must be one of: ALL, KEYS_ONLY, INCLUDE", si.ProjectionType)
	}

	// Validate non-key attributes for INCLUDE projection
	if si.ProjectionType == "INCLUDE" && len(si.NonKeyAttributes) == 0 {
		return fmt.Errorf("non_key_attributes must be specified when projection_type is 'INCLUDE'")
	}
	if si.ProjectionType != "INCLUDE" && len(si.NonKeyAttributes) > 0 {
		return fmt.Errorf("non_key_attributes can only be specified when projection_type is 'INCLUDE'")
	}

	// LSI-specific validation
	if si.IsLSI() {
		if err := si.validateLSI(tableHashKey, tableRangeKey); err != nil {
			return fmt.Errorf("LSI validation failed: %v", err)
		}
	}

	// GSI-specific validation
	if si.IsGSI() {
		if err := si.validateGSI(); err != nil {
			return fmt.Errorf("GSI validation failed: %v", err)
		}
	}

	return nil
}

// validateLSI performs LSI-specific validation
func (si SecondaryIndex) validateLSI(tableHashKey, tableRangeKey string) error {
	// Убираем эту проверку - hash_key устанавливается программно
	// if si.HashKey != "" {
	//     return fmt.Errorf("LSI '%s' cannot specify hash_key (automatically uses table's hash key '%s')", si.Name, tableHashKey)
	// }

	// LSI must have a range key
	if si.RangeKey == "" {
		return fmt.Errorf("LSI '%s' must specify range_key", si.Name)
	}

	// LSI range key must be different from table's range key
	if si.RangeKey == tableRangeKey {
		return fmt.Errorf("LSI '%s' range_key '%s' cannot be the same as table's range_key", si.Name, si.RangeKey)
	}

	// LSI cannot specify throughput (uses table's throughput)
	if si.ReadCapacity != nil || si.WriteCapacity != nil {
		return fmt.Errorf("LSI '%s' cannot specify read/write capacity (uses table's provisioned throughput)", si.Name)
	}

	return nil
}

// validateGSI performs GSI-specific validation
func (si SecondaryIndex) validateGSI() error {
	// GSI must have a hash key
	if si.HashKey == "" {
		return fmt.Errorf("GSI '%s' must specify hash_key", si.Name)
	}

	// Note: GSI range_key is optional, so no validation needed
	// Note: GSI throughput is optional (will use on-demand if not specified)

	return nil
}

// SupportsQuery returns true if this index can be used for the given key pattern
func (si SecondaryIndex) SupportsQuery(hasHashKey, hasRangeKey bool, tableHashKey string) bool {
	effectiveHashKey := si.GetEffectiveHashKey(tableHashKey)

	// Must have effective hash key to support queries
	if effectiveHashKey == "" {
		return false
	}

	// If query needs hash key, we must have one
	if hasHashKey && effectiveHashKey == "" {
		return false
	}

	// If query needs range key, index must have one
	if hasRangeKey && si.RangeKey == "" {
		return false
	}

	return true
}

// GetTerraformType returns the Terraform resource type for this index
func (si SecondaryIndex) GetTerraformType() string {
	if si.IsLSI() {
		return "local_secondary_index"
	}
	return "global_secondary_index"
}
