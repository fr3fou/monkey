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
	token.LPAREN:   CALL,
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
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.infixParseFns = make(map[token.Type]infixParseFn)

	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	return p
}

// nextToken is a helper function that advances through the tokens
func (p *Parser) nextToken() {
	p.tok = p.nextTok
	p.nextTok = p.l.NextToken()
}

// ParseProgram starts parsing the program using our lexer
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for !p.tokIs(token.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
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
