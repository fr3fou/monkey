package parser

import (
	"github.com/fr3fou/monkey/ast"
	"github.com/fr3fou/monkey/token"
)

// parseBoolean parses any boolean literal
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.tok, Value: p.tokIs(token.TRUE)}
}
