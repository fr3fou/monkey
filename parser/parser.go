package parser

import (
	"fmt"

	"github.com/fr3fou/monkey/ast"
	"github.com/fr3fou/monkey/lexer"
	"github.com/fr3fou/monkey/token"
)

// Parser is the struct which does all of the parsing
type Parser struct {
	l       *lexer.Lexer
	tok     token.Token
	nextTok token.Token
	errors  []string
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
		return nil
	}
}

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
