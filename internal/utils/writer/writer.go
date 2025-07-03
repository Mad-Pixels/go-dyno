package writer

// Writer defines the interface for writing generated code.
type Writer interface {
	Write(data []byte) error
	Type() string
}
