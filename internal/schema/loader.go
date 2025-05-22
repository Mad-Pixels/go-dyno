package schema

import (
	"strings"

	"github.com/Mad-Pixels/go-dyno/internal/schema/common"
	"github.com/Mad-Pixels/go-dyno/internal/utils"
)

func LoadSchema(path string) (*DynamoSchema, error) {
	var schema DynamoSchema

	if err := utils.ReadAndParseJsonFile(path, &schema.schema); err != nil {
		return nil, err
	}
	for i, idx := range schema.SecondaryIndexes() {
		schema.schema.SecondaryIndexes[i].HashKeyParts = parseCompositeKeys(idx.HashKey, schema.AllAtributes())
		schema.schema.SecondaryIndexes[i].RangeKeyParts = parseCompositeKeys(idx.RangeKey, schema.AllAtributes())
	}

	return &schema, nil
}

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

func isAttribute(name string, attrs []common.Attribute) bool {
	for _, attr := range attrs {
		if attr.Name == name {
			return true
		}
	}
	return false
}
