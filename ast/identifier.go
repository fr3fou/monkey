package ast

import "github.com/fr3fou/monkey/token"

// Identifier is a name of a variable / function
// It implements the Expression interface, so that it can make things
// easier for us in the future:
// let a = 5;
// let b = a;
// in this case we can see that a IS an expression (it returns a value)
//
// Value is technically the same as Token.Literal or TokenLiteral()
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral returns the identifier token literal (variable name)
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}
