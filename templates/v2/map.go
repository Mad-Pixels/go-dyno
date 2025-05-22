package v2

import "github.com/Mad-Pixels/go-dyno/internal/schema/common"

type TemplateMapV2 struct {
	PackageName      string
	TableName        string
	HashKey          string
	RangeKey         string
	Attributes       []common.Attribute
	CommonAttributes []common.Attribute
	AllAttributes    []common.Attribute
	SecondaryIndexes []common.SecondaryIndex
}
