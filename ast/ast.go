package ast

import (
	"bytes"
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
