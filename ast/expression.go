package ast

import (
	"bytes"

	"github.com/fr3fou/monkey/token"
)

// PrefixExpression is any expression for prefix expressions
// -5
// --5
// !false
type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral returns the PrefixExpression token literal
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// InfixExpression is any expression for infix expressions
// 5 + 5
type InfixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
	Left     Expression
}

func (ie *InfixExpression) expressionNode() {}

// TokenLiteral returns the InfixExpression token literal
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}
