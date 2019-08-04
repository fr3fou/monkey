package token

// Token is the data structure that holds
// all the necessary fields that our
// lexer is going to output
type Token struct {
	Type    Type
	Literal string // TODO: This can be optimized by changing it to byte or int
}

// Type is used to determine the token variant
type Type string

// keywords is a map that contains all the language defined keywords
var keywords = map[string]Type{
	"fun":    FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// All possible token variants
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals

	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 1343456

	// Operators

	ASSIGN   = "="
	EQ       = "=="
	NEQ      = "!="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"

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
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

// LookupIdentifier checks if the given identifier is
// a keyword or a user defined variable / identifier
// and returns the Type
func LookupIdentifier(ident string) Type {
	if token, ok := keywords[ident]; ok {
		return token
	}

	// if the provided identifier is not a part of the keyword
	// map, then it's a user defined one - we should return IDENT
	return IDENT
}
