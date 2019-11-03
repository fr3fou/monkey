package parser

import (
	"fmt"
	"strconv"

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

// Precendence order
var precedences = map[token.Type]int{
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

// Parser is the struct which does all of the parsing
type Parser struct {
	l              *lexer.Lexer
	tok            token.Token
	nextTok        token.Token
	errors         []string
	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
	// TODO: implement postfix too
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
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	p.infixParseFns = make(map[token.Type]infixParseFn)

	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	return p
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

	for !p.tokIs(token.EOF) {
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

	if !p.peekNext(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.tok,
		Value: p.tok.Literal,
	}

	if !p.peekNext(token.ASSIGN) {
		return nil
	}

	// TODO: parse expressions
	// We're skipping the expressions until we
	// encounter a semicolon
	for !p.tokIs(token.SEMICOLON) {
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
	for !p.tokIs(token.SEMICOLON) {
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

	// TODO: parse expressions
	// We're skipping the expressions until we
	// encounter a semicolon
	if p.nextTokIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
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

// parseIdentifier parses any identifier (variable)
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.tok,
		Value: p.tok.Literal,
	}
}

// parseIntegerLiteral parses any integer literal and applies
// strconv.ParseInt on it
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{
		Token: p.tok,
	}

	value, err := strconv.ParseInt(p.tok.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.tok.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

// tokIs is a helper function that checks if the current
// token type matches the one provided
func (p *Parser) tokIs(t token.Type) bool {
	return p.tok.Type == t
}

// nextTokIs is a helper function that checks if the next
// token type matches the one provided
func (p *Parser) nextTokIs(t token.Type) bool {
	return p.nextTok.Type == t
}

// peekNext is a helper function that checks if the next
// token type matches the one provided and if yes, it calls p.nextToken()
func (p *Parser) peekNext(t token.Type) bool {
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

// noPrefixParseFnError is a helper function
// that formats a better error when missing a prefix fn
func (p *Parser) noPrefixParseFnError(t token.Type) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// curPrecedence is a helper function that returns the precedence
// of the next token type
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.nextTok.Type]; ok {
		return p
	}

	return LOWEST
}

// curPrecedence is a helper function that returns the precedence
// of the current token type
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.tok.Type]; ok {
		return p
	}

	return LOWEST
}

// Errors is a function that returns all of the errors
// that the parser met during parsing
func (p *Parser) Errors() []string {
	return p.errors
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
