package token

// All possible token variants
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals

	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 1343456

	// Operators

	ASSIGN = "="
	PLUS   = "+"

	// Delimiters

	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	// Keywords

	FUNCTION = "FUNCTION"
	LET      = "LET"
)

// Type is used to determine the token variant
type Type string

// Token is the data structure that holds
// all the necessary fields that our
// lexer is going to output
type Token struct {
	Type    Type
	Literal string // This can be optimized by changing it to byte or int
}
