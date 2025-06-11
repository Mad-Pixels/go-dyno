package schema

import (
	"github.com/Mad-Pixels/go-dyno/internal/generator/attribute"
	"github.com/Mad-Pixels/go-dyno/internal/generator/index"
	"github.com/Mad-Pixels/go-dyno/internal/logger"
	"github.com/Mad-Pixels/go-dyno/internal/utils"
)

type Schema struct {
	raw schema
}

func NewSchema(path string) (*Schema, error) {
	var spec Schema

	if err := utils.ReadAndParseJSON(path, &spec.raw); err != nil {
		return nil, err
	}
	return &spec, nil
}

func (s Schema) TableName() string {
	return s.raw.TableName
}

func (s Schema) HashKey() string {
	return s.raw.HashKey
}

func (s Schema) RangeKey() string {
	return s.raw.RangeKey
}

func (s Schema) PackageName() string {
	return utils.ToLowerInlineCase(s.raw.TableName)
}

func (s Schema) Filename() string {
	return utils.AddFileExt(s.PackageName(), ".go")
}

func (s Schema) Attributes() []attribute.Attribute {
	return s.raw.Attributes
}

func (s Schema) CommonAttributes() []attribute.Attribute {
	return s.raw.CommonAttributes
}

func (s Schema) AllAttributes() []attribute.Attribute {
	return append(s.Attributes(), s.CommonAttributes()...)
}

func (s Schema) SecondaryIndexes() []index.Index {
	return s.raw.SecondaryIndexes
}

func (s Schema) GlobalSecondaryIndexes() []index.Index {
	var gsi []index.Index

	for _, idx := range s.raw.SecondaryIndexes {
		if idx.IsGSI() {
			gsi = append(gsi, idx)
		}
	}
	return gsi
}

func (s Schema) LocalSecondaryIndexes() []index.Index {
	var lsi []index.Index

	for _, idx := range s.raw.SecondaryIndexes {
		if idx.IsLSI() {
			lsi = append(lsi, idx)
		}
	}
	return lsi
}

func (s Schema) HasGSI() bool {
	return len(s.GlobalSecondaryIndexes()) > 0
}

func (s Schema) HasLSI() bool {
	return len(s.LocalSecondaryIndexes()) > 0
}

func (s Schema) HasSecondaryIndexes() bool {
	return len(s.raw.SecondaryIndexes) > 0
}

func (s Schema) GetIndexByName(name string) *index.Index {
	for i := range s.raw.SecondaryIndexes {
		if s.raw.SecondaryIndexes[i].Name == name {
			return &s.raw.SecondaryIndexes[i]
		}
	}
	return nil
}

func (s Schema) ValidateIndexNames() error {
	seen := make(map[string]bool)

	for _, idx := range s.SecondaryIndexes() {
		if seen[idx.Name] {
			return logger.NewFailure("duplicate index name", nil).
				With("name", idx.Name)
		}
		seen[idx.Name] = true
	}
	return nil
}

// GetOptimalIndexForQuery returns the best index for a query with the given keys
// Priority: LSI (cheaper) -> GSI -> nil (use main table)
func (s Schema) GetOptimalIndexForQuery(hashKey, rangeKey string) *index.Index {
	var (
		hasHashKey  = hashKey != ""
		hasRangeKey = rangeKey != ""
	)

	for _, idx := range s.LocalSecondaryIndexes() {
		if idx.SupportsQuery(hasHashKey, hasRangeKey, s.HashKey()) {
			if hasHashKey && hashKey == s.HashKey() {
				return &idx
			}
		}
	}
	for _, idx := range s.GlobalSecondaryIndexes() {
		if idx.SupportsQuery(hasHashKey, hasRangeKey, s.HashKey()) {
			if hasHashKey && hashKey == idx.HashKey {
				return &idx
			}
		}
	}
	return nil
}

type schema struct {
	// TableName defines the logical name of the DynamoDB table.
	// It is used for generating Go package, filename, and identifiers in the template.
	// Go package and filename can be overridden by generator config.
	TableName string `json:"table_name"`

	// HashKey is the primary partition key of the DynamoDB table.
	// This field is required and must match one of the attribute names.
	HashKey string `json:"hash_key"`

	// RangeKey is the optional sort key of the DynamoDB table.
	// If defined, it must also match one of the attribute names.
	RangeKey string `json:"range_key"`

	// Attributes define the key attributes that can be used in primary keys
	// and secondary indexes (hash_key, range_key for GSI/LSI).
	// These fields must be defined for DynamoDB key operations.
	Attributes []attribute.Attribute `json:"attributes"`

	// CommonAttributes define non-key attributes that are only stored as data
	// and cannot be used as keys in indexes. These are regular table fields.
	CommonAttributes []attribute.Attribute `json:"common_attributes"`

	// SecondaryIndexes defines Global or Local Secondary Indexes (GSI/LSI)
	// used for advanced querying in DynamoDB. Each index has its own keys and projection.
	SecondaryIndexes []index.Index `json:"secondary_indexes"`
}
