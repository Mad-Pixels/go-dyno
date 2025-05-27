package schema

import (
	"strings"

	"github.com/Mad-Pixels/go-dyno/internal/schema/common"
	"github.com/Mad-Pixels/go-dyno/internal/utils"
)

// LoadSchema loads a DynamoDB table schema from a JSON file and parses it into a DynamoSchema structure.
// It automatically parses composite keys in secondary indexes, splitting them into individual parts
// for proper query building and validation.
//
// The function performs the following operations:
// 1. Reads and parses the JSON schema file
// 2. Identifies composite keys (containing "#" separator) in secondary indexes
// 3. Splits composite keys into constituent parts, marking each as constant or attribute reference
// 4. Returns a fully initialized DynamoSchema ready for code generation
//
// Example:
//
//	schema, err := schema.LoadSchema("schemas/user_activity.json")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// schema.SecondaryIndexes will have HashKeyParts and RangeKeyParts populated
func LoadSchema(path string) (*DynamoSchema, error) {
	var schema DynamoSchema

	if err := utils.ReadAndParseJSON(path, &schema.schema); err != nil {
		return nil, err
	}
	for i := range schema.schema.SecondaryIndexes {
		idx := &schema.schema.SecondaryIndexes[i]

		idx.HashKeyParts = parseCompositeKeys(idx.HashKey, schema.AllAttributes())
		idx.RangeKeyParts = parseCompositeKeys(idx.RangeKey, schema.AllAttributes())
	}
	return &schema, nil
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
