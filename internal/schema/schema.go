package schema

import (
	"github.com/Mad-Pixels/go-dyno/internal/utils"
)

type dynamoSchema struct {
	TableName        string            `json:"table_name"`
	HashKey          string            `json:"hash_key"`
	RangeKey         string            `json:"range_key"`
	Attributes       []utils.Attribute `json:"attributes"`
	CommonAttributes []utils.Attribute `json:"common_attributes"`
	SecondaryIndexes []SecondaryIndex  `json:"secondary_indexes"`
}

type DynamoSchema struct {
	schema dynamoSchema
}

func (ds DynamoSchema) TableName() string {
	return utils.ToUpperCamelCase(ds.schema.TableName)
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
