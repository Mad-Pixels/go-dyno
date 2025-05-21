package generator

type DynamoSchema struct {
	TableName        string           `json:"table_name"`
	HashKey          string           `json:"hash_key"`
	RangeKey         string           `json:"range_key"`
	Attributes       []Attribute      `json:"attributes"`
	CommonAttributes []Attribute      `json:"common_attributes"`
	SecondaryIndexes []SecondaryIndex `json:"secondary_indexes"`
}

type Attribute struct {
	Name string `json:"name"`
	Type string `json:"type"`
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
