package schema

import (
	"github.com/Mad-Pixels/go-dyno/internal/utils"
)

func LoadSchema(path string) (*DynamoSchema, error) {
	var schema DynamoSchema

	if err := utils.ReadAndParseJsonFile(path, &schema); err != nil {
		return nil, err
	}
	return &schema, nil
}
