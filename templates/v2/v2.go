// Package v2 provides the data structures and code templates used for generating DynamoDB integration code.
//
// It defines a flexible schema representation (`TemplateMapV2`) that can be rendered into Go code, enabling developers to quickly scaffold
// strongly-typed models, constants, and query builders for AWS DynamoDB tables.
//
// generated code work with AWS-SDK-V2 version.
package v2

import "github.com/Mad-Pixels/go-dyno/internal/schema/common"

// TemplateMapV2 defines the full set of metadata used to generate DynamoDB-related code.
// It acts as the main input structure for the Go code template engine.
type TemplateMapV2 struct {
	// PackageName is the Go package name to use in the generated file.
	PackageName string

	// TableName is the name of the DynamoDB table.
	TableName string

	// HashKey is the primary partition key of the table.
	HashKey string

	// RangeKey is the optional sort key of the table.
	RangeKey string

	// Attributes are the table-specific attributes defined in the schema.
	Attributes []common.Attribute

	// CommonAttributes are shared attributes used across multiple tables.
	CommonAttributes []common.Attribute

	// AllAttributes is a union of Attributes and CommonAttributes, used in template rendering.
	AllAttributes []common.Attribute

	// SecondaryIndexes defines all global and local secondary indexes for the table.
	SecondaryIndexes []common.SecondaryIndex
}
