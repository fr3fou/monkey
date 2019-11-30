package parser

import "github.com/fr3fou/monkey/ast"

// parseIdentifier parses any identifier (variable)
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.tok,
		Value: p.tok.Literal,
	}
}
