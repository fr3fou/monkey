package parser

import (
	"github.com/fr3fou/monkey/ast"
	"github.com/fr3fou/monkey/lexer"
	"github.com/fr3fou/monkey/token"
)

// Parser is the struct which does all of the parsing
type Parser struct {
	l       *lexer.Lexer
	tok     token.Token
	nextTok token.Token
}

// New returns a pointer to a parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

// nextToken is a helper function that advances through the tokens
func (p *Parser) nextToken() {
	p.tok = p.nextTok
	p.nextTok = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	return nil
}
