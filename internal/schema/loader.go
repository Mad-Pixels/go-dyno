package schema

import (
	"fmt"
	"strings"

	"github.com/Mad-Pixels/go-dyno/internal/schema/common"
	"github.com/Mad-Pixels/go-dyno/internal/utils"
)

// LoadSchema loads a DynamoDB table schema from a JSON file and validates it comprehensively.
// It performs the following operations:
// 1. Reads and parses the JSON schema file
// 2. Validates all attributes for correctness
// 3. Validates and processes secondary indexes (both GSI and LSI)
// 4. Ensures DynamoDB constraints are met (LSI limits, unique names, etc.)
// 5. Parses composite keys for advanced query capabilities
// 6. Returns a fully validated DynamoSchema ready for code generation
//
// Example:
//
//	schema, err := schema.LoadSchema("schemas/user_posts.json")
//	if err != nil {
//	    log.Fatal(err) // Schema validation failed
//	}
//	// schema is now ready for Terraform and code generation
func LoadSchema(path string) (*DynamoSchema, error) {
	var schema DynamoSchema

	// Load and parse JSON
	if err := utils.ReadAndParseJSON(path, &schema.schema); err != nil {
		return nil, fmt.Errorf("failed to parse schema file '%s': %v", path, err)
	}

	// Validate basic table structure
	if err := validateTableStructure(&schema); err != nil {
		return nil, fmt.Errorf("invalid table structure: %v", err)
	}

	// Validate all attributes
	if err := validateAttributes(&schema); err != nil {
		return nil, err
	}

	// Validate and process secondary indexes
	if err := validateAndProcessIndexes(&schema); err != nil {
		return nil, err
	}

	return &schema, nil
}

// validateTableStructure validates basic table properties
func validateTableStructure(schema *DynamoSchema) error {
	if schema.schema.TableName == "" {
		return fmt.Errorf("table_name cannot be empty")
	}

	if schema.schema.HashKey == "" {
		return fmt.Errorf("hash_key cannot be empty")
	}

	// range_key is optional, so no validation needed

	return nil
}

// validateAttributes validates all table attributes
func validateAttributes(schema *DynamoSchema) error {
	// Validate primary attributes
	for _, attr := range schema.schema.Attributes {
		if err := attr.Validate(); err != nil {
			return fmt.Errorf("invalid attribute in 'attributes': %v", err)
		}
	}

	// Validate common attributes
	for _, attr := range schema.schema.CommonAttributes {
		if err := attr.Validate(); err != nil {
			return fmt.Errorf("invalid attribute in 'common_attributes': %v", err)
		}
	}

	// Ensure hash_key and range_key are defined in attributes
	allAttributes := schema.AllAttributes()
	if !isAttributeDefined(schema.schema.HashKey, allAttributes) {
		return fmt.Errorf("hash_key '%s' is not defined in attributes", schema.schema.HashKey)
	}

	if schema.schema.RangeKey != "" && !isAttributeDefined(schema.schema.RangeKey, allAttributes) {
		return fmt.Errorf("range_key '%s' is not defined in attributes", schema.schema.RangeKey)
	}

	return nil
}

// validateAndProcessIndexes validates and processes all secondary indexes
func validateAndProcessIndexes(schema *DynamoSchema) error {
	if len(schema.schema.SecondaryIndexes) == 0 {
		return nil // No indexes to validate
	}

	// Check for duplicate index names
	if err := schema.ValidateIndexNames(); err != nil {
		return err
	}

	// Count LSI indexes (DynamoDB limit is 10)
	lsiCount := 0
	allAttributes := schema.AllAttributes()

	for i := range schema.schema.SecondaryIndexes {
		idx := &schema.schema.SecondaryIndexes[i]

		// Set default type to GSI for backward compatibility
		if idx.Type == "" {
			idx.Type = common.GSI
		}

		if idx.IsLSI() {
			idx.HashKey = schema.schema.HashKey // Автоматически устанавливаем hash_key для LSI
		}

		// Validate the index
		if err := idx.Validate(schema.schema.HashKey, schema.schema.RangeKey); err != nil {
			return fmt.Errorf("invalid secondary index '%s': %v", idx.Name, err)
		}

		if idx.IsLSI() {
			idx.HashKey = schema.schema.HashKey // ← ВОТ КЛЮЧЕВАЯ СТРОКА!
		}

		// Validate that referenced attributes exist
		if err := validateIndexAttributes(idx, allAttributes); err != nil {
			return fmt.Errorf("index '%s' references undefined attributes: %v", idx.Name, err)
		}

		// Count and validate LSI limits
		if idx.IsLSI() {
			lsiCount++
			if lsiCount > 10 {
				return fmt.Errorf("too many LSI indexes: %d (DynamoDB limit is 10)", lsiCount)
			}
		}

		// Parse composite keys
		if err := parseIndexCompositeKeys(idx, allAttributes); err != nil {
			return fmt.Errorf("failed to parse composite keys for index '%s': %v", idx.Name, err)
		}
	}

	return nil
}

// validateIndexAttributes ensures all attributes referenced by the index are defined
func validateIndexAttributes(idx *common.SecondaryIndex, allAttributes []common.Attribute) error {
	// Validate hash_key (only for GSI)
	if idx.IsGSI() && idx.HashKey != "" {
		if !isAttributeDefined(idx.HashKey, allAttributes) && !isCompositeKey(idx.HashKey) {
			return fmt.Errorf("hash_key '%s' is not defined", idx.HashKey)
		}
	}

	// Validate range_key
	if idx.RangeKey != "" {
		if !isAttributeDefined(idx.RangeKey, allAttributes) && !isCompositeKey(idx.RangeKey) {
			return fmt.Errorf("range_key '%s' is not defined", idx.RangeKey)
		}
	}

	// Validate non_key_attributes
	for _, attr := range idx.NonKeyAttributes {
		if !isAttributeDefined(attr, allAttributes) {
			return fmt.Errorf("non_key_attribute '%s' is not defined", attr)
		}
	}

	return nil
}

// parseIndexCompositeKeys parses composite keys for the index
func parseIndexCompositeKeys(idx *common.SecondaryIndex, allAttributes []common.Attribute) error {
	// Parse hash key composite parts (only for GSI)
	if idx.IsGSI() && isCompositeKey(idx.HashKey) {
		parts := parseCompositeKeys(idx.HashKey, allAttributes)
		if len(parts) == 0 {
			return fmt.Errorf("failed to parse composite hash_key '%s'", idx.HashKey)
		}
		idx.HashKeyParts = parts
	}

	// Parse range key composite parts
	if isCompositeKey(idx.RangeKey) {
		parts := parseCompositeKeys(idx.RangeKey, allAttributes)
		if len(parts) == 0 {
			return fmt.Errorf("failed to parse composite range_key '%s'", idx.RangeKey)
		}
		idx.RangeKeyParts = parts
	}

	return nil
}

// isAttributeDefined checks if an attribute name is defined in the attribute list
func isAttributeDefined(name string, attributes []common.Attribute) bool {
	for _, attr := range attributes {
		if attr.Name == name {
			return true
		}
	}
	return false
}

// isCompositeKey returns true if the key contains the composite separator '#'
func isCompositeKey(key string) bool {
	return strings.Contains(key, "#")
}

// parseCompositeKeys splits a composite key string into its constituent parts.
// Each part is classified as either a constant string or a reference to an attribute.
//
// Composite key format: "part1#part2#part3"
// - If a part matches an existing attribute name → treated as dynamic (IsConstant = false)
// - If a part doesn't match any attribute → treated as constant (IsConstant = true)
//
// Returns nil if the key is empty or doesn't contain the "#" separator (simple key).
//
// Example with mixed parts:
//
//	attrs := []common.Attribute{
//	    {Name: "user_id", Type: "S"},
//	    {Name: "year", Type: "N"},
//	}
//	parts := parseCompositeKeys("user_id#2023#active", attrs)
//
//	// Result:
//	// [
//	//   {IsConstant: false, Value: "user_id"}, // matches attribute
//	//   {IsConstant: true,  Value: "2023"},    // constant value
//	//   {IsConstant: true,  Value: "active"}   // constant value
//	// ]
func parseCompositeKeys(key string, attrs []common.Attribute) []common.CompositeKeyPart {
	if key == "" || !strings.Contains(key, "#") {
		return nil
	}

	var (
		parts  = strings.Split(key, "#")
		result []common.CompositeKeyPart
	)
	for _, part := range parts {
		isAttr := isAttribute(part, attrs)

		if isAttr {
			result = append(result, common.CompositeKeyPart{IsConstant: false, Value: part})
		} else {
			result = append(result, common.CompositeKeyPart{IsConstant: true, Value: part})
		}
	}
	return result
}

// isAttribute checks if a given name matches any attribute defined in the schema.
// Performs case-sensitive string comparison against attribute names.
//
// Parameters:
//   - name: The string to check against attribute names
//   - attrs: Slice of attributes to search through
//
// Returns true if the name matches any attribute, false otherwise.
//
// Example:
//
//	attrs := []common.Attribute{
//	    {Name: "user_id", Type: "S"},
//	    {Name: "created_at", Type: "N"},
//	}
//
//	isAttribute("user_id", attrs)    // → true
//	isAttribute("unknown", attrs)    // → false
//	isAttribute("USER_ID", attrs)    // → false (case-sensitive)
func isAttribute(name string, attrs []common.Attribute) bool {
	for _, attr := range attrs {
		if attr.Name == name {
			return true
		}
	}
	return false
}
