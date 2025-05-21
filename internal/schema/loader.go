package schema

import (
	"github.com/Mad-Pixels/go-dyno/internal/utils/fs"
)

func LoadSchema(path string) (*DynamoSchema, error) {
	var s DynamoSchema

	if err := fs.ReadAndParseJsonFile(path, &s); err != nil {
		return nil, err
	}
	return &s, nil
}
