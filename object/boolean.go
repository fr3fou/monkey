package object

import "fmt"

// Boolean represents an bool value
type Boolean struct {
	Value bool
}

// Inspect is used for debugging
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// Type returns the boolean type
func (b *Boolean) Type() Type {
	return BOOLEAN_OBJ
}
