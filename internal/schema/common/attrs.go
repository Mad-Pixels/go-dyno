package common

// Attribute defines a basic DynamoDB attribute with a name and type.
// Supported types: "S" (string), "N" (number), "B" (boolean).
type Attribute struct {
	// Name is the logical name of the attribute as defined in the schema.
	Name string `json:"name"`

	// Type is the DynamoDB type of the attribute: "S", "N", or "B".
	Type string `json:"type"`
}

// SecondaryIndex describes a Global Secondary Index (GSI) or Local Secondary Index (LSI)
// for a DynamoDB table, including its keys and projection settings.
type SecondaryIndex struct {
	// Name is the identifier for the index used in DynamoDB and code generation.
	Name string `json:"name"`

	// HashKey is the primary partition key for the index.
	// It can be a single attribute or a composite key (joined with #).
	HashKey string `json:"hash_key"`

	// HashKeyParts is the parsed breakdown of HashKey into its parts.
	// Used internally to support composite key generation.
	HashKeyParts []CompositeKeyPart

	// RangeKey is the optional sort key for the index.
	// It can also be composite (e.g., "user#date").
	RangeKey string `json:"range_key"`

	// RangeKeyParts is the parsed breakdown of RangeKey into its parts.
	RangeKeyParts []CompositeKeyPart

	// ProjectionType defines which attributes are included in the index.
	// Valid values: "ALL", "KEYS_ONLY", "INCLUDE".
	ProjectionType string `json:"projection_type"`

	// NonKeyAttributes lists additional attributes included in the projection
	// when ProjectionType is "INCLUDE".
	NonKeyAttributes []string `json:"non_key_attributes"`
}

// CompositeKeyPart represents a part of a composite key.
// It can either be a constant string (e.g., "user") or an attribute (e.g., "user_id").
type CompositeKeyPart struct {
	// IsConstant indicates whether the part is a literal constant or a reference to an attribute.
	IsConstant bool

	// Value is either the constant string or the attribute name.
	Value string
}
