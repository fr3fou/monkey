package object

// Null represents a null value (lack of a value)
type Null struct {
}

// Inspect is used for debugging
func (n *Null) Inspect() string {
	return "null"
}

// Type returns the null type
func (n *Null) Type() Type {
	return NULL_OBJ
}
