package parser

import (
	"github.com/fr3fou/monkey/ast"
	"github.com/fr3fou/monkey/token"
)

// parseFunctionLiteral parses any function literals
// fun(x,y) { x }
func (p *Parser) parseFunctionLiteral() ast.Expression {
	exp := &ast.FunctionLiteral{
		Token: p.tok,
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	exp.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	exp.Body = p.parseBlockStatement()

	return exp
}

// parseFunctionParameters parses params to functions
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	// handle functions with no arguments
	if p.nextTokIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	// skip the first `(`
	p.nextToken()

	ident := &ast.Identifier{Token: p.tok, Value: p.tok.Literal}
	identifiers = append(identifiers, ident)

	// make sure there are commas in between each param
	// (x,y,z,f)
	for p.nextTokIs(token.COMMA) {
		// skip comma
		p.nextToken()
		p.nextToken()

		ident := &ast.Identifier{Token: p.tok, Value: p.tok.Literal}
		identifiers = append(identifiers, ident)
	}

	// if there isn't a closing `)`, return nil
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

// parseCallExpression parses any call to a func
// foo(3,"hi")
func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.tok, Function: fn}
	exp.Arguments = p.parseCallArguments()
	return exp
}

// parseCallArguments parses the args to a function call
// foo(3,1,"foo")
func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	// handle functions with no arguments
	if p.nextTokIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	// skip the first `(`
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	// make sure there are commas in between each param
	// (x,y,z,f)
	for p.nextTokIs(token.COMMA) {
		// skip comma
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	// if there isn't a closing `)`, return nil
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}
