package ast

import (
	"bytes"

	"github.com/fr3fou/monkey/token"
)

// Node is the main interface / component of our AST,
// everything has to implement it
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement is the interface for all statements
type Statement interface {
	Node
	statementNode()
}

// Expression is the interface for all expressions
type Expression interface {
	Node
	expressionNode()
}

// Program is going to be the root node of our AST
type Program struct {
	Statements []Statement
}

// TokenLiteral returns the root token literal
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// LetStatement is any statment for declaring variables (let x = "foo")
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

// TokenLiteral returns the `let` token literal
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

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

// ReturnStatement is any statment that returns from a function (return 5)
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral returns the `return` token literal
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")
	return out.String()
}

// ExpressionStatement is any statement that consists solely of an expression
// let a = 5;
// a + 5
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// TokenLiteral returns the expression token literal
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}
