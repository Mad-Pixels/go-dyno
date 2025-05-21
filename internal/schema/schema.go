package schema

import (
	"path/filepath"
	"strings"

	"github.com/Mad-Pixels/go-dyno/internal/utils"
)

type DynamoSchema struct {
	TableName        string            `json:"table_name"`
	HashKey          string            `json:"hash_key"`
	RangeKey         string            `json:"range_key"`
	Attributes       []utils.Attribute `json:"attributes"`
	CommonAttributes []utils.Attribute `json:"common_attributes"`
	SecondaryIndexes []SecondaryIndex  `json:"secondary_indexes"`
}

func (ds DynamoSchema) PackageName() string {
	return strings.ToLower(utils.ToSafeName(ds.TableName))
}

func (ds DynamoSchema) packageDir(root string) string {
	return filepath.Join(root, ds.PackageName())
}

type CompositeKeyPart struct {
	IsConstant bool
	Value      string
}

type SecondaryIndex struct {
	Name             string `json:"name"`
	HashKey          string `json:"hash_key"`
	HashKeyParts     []CompositeKeyPart
	RangeKey         string `json:"range_key"`
	RangeKeyParts    []CompositeKeyPart
	ProjectionType   string   `json:"projection_type"`
	NonKeyAttributes []string `json:"non_key_attributes"`
}
