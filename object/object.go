package object

// Object represents values from our language in go
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Type is used to determine the object variant
type Type string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
)
