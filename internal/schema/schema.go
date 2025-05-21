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

func (ds DynamoSchema) HashKey() string {
	return ds.schema.HashKey
}

func (ds DynamoSchema) RangeKey() string {
	return ds.schema.RangeKey
}

func (ds DynamoSchema) PackageName() string {
	return utils.ToLowerInlineCase(ds.schema.TableName)
}

func (ds DynamoSchema) Dictionary() string {
	return utils.ToSafeName(ds.schema.TableName)
}

func (ds DynamoSchema) Attributes() []utils.Attribute {
	return ds.schema.Attributes
}

func (ds DynamoSchema) CommonAttributes() []utils.Attribute {
	return ds.schema.CommonAttributes
}

func (ds DynamoSchema) SecondaryIndexes() []SecondaryIndex {
	return ds.schema.SecondaryIndexes
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
