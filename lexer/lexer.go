package lexer

import "github.com/fr3fou/monkey/token"

// Lexer is a struct that
// takes in an input and returns
// all of the valid tokens
type Lexer struct {
	input   string // should be io.Reader instead
	pos     int    // "pointer" to the the current char
	nextPos int    // "poitner" to the position after the current char (allows us to "peek")
	ch      byte   // current char (should be rune instead - support UTF)
}

// NewLexer returns a pointer to
// a Lexer struct with the given input
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input: input,
	}

	// Call it immediately after creation
	// so we can make sure we have initialized our lexer
	l.readChar()

	return l
}

// readChar is a helper function that goes advances
// through our input by incrementing pos and nextPos
// (and updating our current character)
func (l *Lexer) readChar() {
	// Peek into the next position and check len to prevent out of bounds
	if l.nextPos >= len(l.input) {
		// 0 is NUL in ASCII
		l.ch = 0
	} else {
		// Otherwise, set the char to the next char
		l.ch = l.input[l.nextPos]
	}

	// Update the "pointers"
	l.pos = l.nextPos
	l.nextPos++
}

// NextToken advances through our input
// and returns the actual Token struct for the
// given character
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}

	l.readChar()
	return tok
}

// newToken is a helper function that simplifies creating tokens
func newToken(tokenType token.Type, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}
