package parser

import (
	"github.com/fr3fou/monkey/ast"
	"github.com/fr3fou/monkey/token"
)

// parseExpression is a general function for parsing expressions
// it checks if the current token type has a matching function in our
// map for prefix / infix expressions
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.tok.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.tok.Type)
		return nil
	}

	leftExp := prefix()
	for !p.nextTokIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.nextTok.Type]

		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

// parsePrefixExpression parses any expression that has a prefix
// -5
// ++5
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.tok,
		Operator: p.tok.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)
	return expression
}

// parseInfixExpression parses any infix expression
// 5 + 5
// 5 / 5
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.tok,
		Operator: p.tok.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression

}

// parseGroupedExpression parses grouped expressions
// (5 + 5) * 2
func (p *Parser) parseGroupedExpression() ast.Expression {
	// carry on with the next token
	p.nextToken()

	// kick off parsing expresions like you normally would
	exp := p.parseExpression(LOWEST)

	// after the expression has been parsed
	// check if it ends with a `)`
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

// parseIfExpression parses any if expression
// if (foo) { x } else { y }
func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{Token: p.tok}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// carry on with the condition
	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	// parse consequence
	exp.Consequence = p.parseBlockStatement()

	// check if there is an alternative
	if p.nextTokIs(token.ELSE) {
		p.nextToken()

		// TODO: check if there is another if
		// handle this:
		// else if () {}

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		exp.Alternative = p.parseBlockStatement()
	}

	return exp
}
