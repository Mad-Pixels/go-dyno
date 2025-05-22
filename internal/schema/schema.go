package schema

import (
	"github.com/Mad-Pixels/go-dyno/internal/schema/common"
	"github.com/Mad-Pixels/go-dyno/internal/utils"
)

type dynamoSchema struct {
	TableName        string                  `json:"table_name"`
	HashKey          string                  `json:"hash_key"`
	RangeKey         string                  `json:"range_key"`
	Attributes       []common.Attribute      `json:"attributes"`
	CommonAttributes []common.Attribute      `json:"common_attributes"`
	SecondaryIndexes []common.SecondaryIndex `json:"secondary_indexes"`
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

func (ds DynamoSchema) Directory() string {
	return utils.ToSafeName(ds.schema.TableName)
}

func (ds DynamoSchema) Filename() string {
	return utils.ToSafeName(ds.schema.TableName) + ".go"
}

func (ds DynamoSchema) Attributes() []common.Attribute {
	return ds.schema.Attributes
}

func (ds DynamoSchema) CommonAttributes() []common.Attribute {
	return ds.schema.CommonAttributes
}

func (ds DynamoSchema) SecondaryIndexes() []common.SecondaryIndex {
	return ds.schema.SecondaryIndexes
}

func (ds DynamoSchema) AllAtributes() []common.Attribute {
	return append(ds.Attributes(), ds.CommonAttributes()...)
}
