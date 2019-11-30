package ast

import "github.com/fr3fou/monkey/token"

// Boolean is any literal that contains only a bool
// true;
// false;
// let foobar = true;
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}

// TokenLiteral returns the boolean token literal
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return b.Token.Literal
}
