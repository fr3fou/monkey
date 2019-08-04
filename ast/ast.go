package ast

import "github.com/fr3fou/monkey/token"

// Node is the main interface / component of our AST,
// everything has to implement it
type Node interface {
	TokenLiteral() string
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

// Identifier is a name of a variable / function
// It implements the Expression interface, so that it can make things
// easier for us in the future:
// let a = 5;
// let b = a;
// in this case we can see that a IS an expression (it returns a value)
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral returns the identifier token literal (variable name)
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// TokenLiteral returns the root token literal
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}
