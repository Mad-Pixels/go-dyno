package common

type Attribute struct {
	Name string `json:"name"`
	Type string `json:"type"`
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

type CompositeKeyPart struct {
	IsConstant bool
	Value      string
}
