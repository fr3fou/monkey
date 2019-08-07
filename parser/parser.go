package parser

import (
	"fmt"

	"github.com/fr3fou/monkey/ast"
	"github.com/fr3fou/monkey/lexer"
	"github.com/fr3fou/monkey/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Precendence order
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

// Parser is the struct which does all of the parsing
type Parser struct {
	l              *lexer.Lexer
	tok            token.Token
	nextTok        token.Token
	errors         []string
	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

// New returns a pointer to a parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)

	return p
}

// Errors is a function that returns all of the errors
// that the parser met during parsing
func (p *Parser) Errors() []string {
	return p.errors
}

// nextToken is a helper function that advances through the tokens
func (p *Parser) nextToken() {
	p.tok = p.nextTok
	p.nextTok = p.l.NextToken()
}

// ParseProgram starts parsing the program using our lexer
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

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

	// TODO: parse expressions
	// We're skipping the expressions until we
	// encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement parses any return statment (return '5')
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{
		Token: p.tok,
	}

	p.nextToken()

	// TODO: parse expressions
	// We're skipping the expressions until we
	// encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
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

// parseExpression is a general function for parsing expression
// it checks if the current token type has a matching function in our
// map for prefix / infix expressions
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.tok.Type]

	if prefix == nil {
		return nil
	}

	leftExp := prefix()

	return leftExp
}

// parseIdentifier parses any identifier (variable)
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.tok,
		Value: p.tok.Literal,
	}
}

// curTokenIs is a helper function that checks if the current
// token type matches the one provided
func (p *Parser) curTokenIs(t token.Type) bool {
	return p.tok.Type == t
}

// nextTokIs is a helper function that checks if the next
// token type matches the one provided
func (p *Parser) nextTokIs(t token.Type) bool {
	return p.nextTok.Type == t
}

// expectPeek is a helper function that checks if the next
// token type matches the one provided and if yes, it calls p.nextToken()
func (p *Parser) expectPeek(t token.Type) bool {
	if p.nextTokIs(t) {
		p.nextToken()
		return true
	}

	p.peekError(t)
	return false
}

// peekError is a helper function that
// is used when meeting an unexpected token when peeking
func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.nextTok.Type)
	p.errors = append(p.errors, msg)
}

// registerPrefix is a helper function that adds the provided function in the map
// for prefix functions
func (p *Parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfix is a helper function that adds the provided function in the map
// for infix functions
func (p *Parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
