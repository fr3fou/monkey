package parser

import (
	"github.com/fr3fou/monkey/ast"
	"github.com/fr3fou/monkey/token"
)

// parseStatement is a helper function that parses the current token
// with the appropriate parsing function
func (p *Parser) parseStatement() ast.Statement {
	switch p.tok.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseReturnStatement parses any return statment (return '5')
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{
		Token: p.tok,
	}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.nextTokIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseLetStatement parses any let statment (let foo = 5)
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{
		Token: p.tok,
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.tok,
		Value: p.tok.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.nextTokIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseExpressionStatement parses any expression statements (foobar;)
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token: p.tok,
	}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.nextTokIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseBlockStatement parses any block statement
// {
//   let a = 5;
//   let b = 4;
// }
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	// create the block
	block := &ast.BlockStatement{
		Token:      p.tok,
		Statements: []ast.Statement{},
	}

	p.nextToken()

	for !p.tokIs(token.RBRACE) && !p.tokIs(token.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		p.nextToken()
	}

	return block
}
