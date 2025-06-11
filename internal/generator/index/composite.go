package index

// CompositeKey represent a part of a composite key.
type CompositeKey struct {
	// IsConstant indicates whether the part is a literal constant or a reference to an attribute.
	IsConstant bool

	// Value is either the constant string or the attribute name.
	Value string
}
