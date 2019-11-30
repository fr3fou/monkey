package object

import "fmt"

// Integer represents an integer value
type Integer struct {
	Value int64
}

// Inspect is used for debugging
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// Type returns the integer type
func (i *Integer) Type() Type {
	return INTEGER_OBJ
}
