package schema

import (
	"strings"

	"github.com/Mad-Pixels/go-dyno/internal/schema/common"
	"github.com/Mad-Pixels/go-dyno/internal/utils"
)

// LoadSchema loads a DynamoDB table schema from a JSON file and parses it into a DynamoSchema structure.
// It also expands composite hash/range keys in secondary indexes into their parts.
// Returns a pointer to DynamoSchema or an error if reading or parsing fails.
//
// Example:
//
//	schema, err := schema.LoadSchema("schemas/user_activity.json")
func LoadSchema(path string) (*DynamoSchema, error) {
	var schema DynamoSchema

	if err := utils.ReadAndParseJSON(path, &schema.schema); err != nil {
		return nil, err
	}
	for i, idx := range schema.SecondaryIndexes() {
		schema.schema.SecondaryIndexes[i].HashKeyParts = parseCompositeKeys(idx.HashKey, schema.AllAttributes())
		schema.schema.SecondaryIndexes[i].RangeKeyParts = parseCompositeKeys(idx.RangeKey, schema.AllAttributes())
	}

	return &schema, nil
}

// parseCompositeKeys splits a composite key string like "user#2023" into parts.
// If a part matches an existing attribute name, it's treated as dynamic (IsConstant = false),
// otherwise it's treated as a constant string value.
//
// Example:
//
//	attrs := []common.Attribute{{Name: "user"}, {Name: "year"}}
//	parts := parseCompositeKeys("user#year#fixed", attrs)
//
//	// Result:
//	// [
//	//   { IsConstant: false, Value: "user" },
//	//   { IsConstant: false, Value: "year" },
//	//   { IsConstant: true,  Value: "fixed" }
//	// ]
func parseCompositeKeys(key string, attrs []common.Attribute) []common.CompositeKeyPart {
	if key == "" {
		return nil
	}

	var (
		parts  = strings.Split(key, "#")
		result []common.CompositeKeyPart
	)

	for _, part := range parts {
		if isAttribute(part, attrs) {
			result = append(result, common.CompositeKeyPart{IsConstant: false, Value: part})
			continue
		}
		result = append(result, common.CompositeKeyPart{IsConstant: true, Value: part})
	}
	return result
}

// isAttribute checks whether a string matches one of the attribute names defined in the schema.
//
// Example:
//
//	isAttribute("user_id", []common.Attribute{{Name: "user_id"}}) // → true
func isAttribute(name string, attrs []common.Attribute) bool {
	for _, attr := range attrs {
		if attr.Name == name {
			return true
		}
	}
	return false
}
