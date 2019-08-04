package lexer

import "github.com/fr3fou/monkey/token"

// Lexer is a struct that
// takes in an input and returns
// all of the valid tokens
type Lexer struct {
	input string // TODO: should be io.Reader instead

	// pos is a "pointer" to the the current char
	pos int

	// nextPos is a "pointer" to the position after the current char (allows us to "peek")
	nextPos int

	// ch is the current character in our lexer
	ch byte // TODO: should be rune instead - support UTF
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

// NextToken advances through our input
// and returns the actual Token struct for the
// given character
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	// eat any whitespace chars
	l.eatWhitespace()

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
	default:
		// handle identifiers (variable names) and keywords
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)

			// return early on to prevent calling l.readChar() again
			return tok
		}

		// handle numbers (integers)
		if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()

			// return early on to prevent calling l.readChar() again
			return tok
		}

		// anything else is illegal and unknown to our parser
		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
}

// readChar is a helper function that advances
// through our input by incrementing pos and nextPos by 1
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

// readIdentifier starts reading up from the current position
// until any character that is NOT a letter (for example space)
func (l *Lexer) readIdentifier() string {
	pos := l.pos

	// keep advancing through our input
	// and check if ch is still a letter
	for isLetter(l.ch) {
		l.readChar()
	}

	// return the identifier
	return l.input[pos:l.pos]
}

// readIdentifier starts reading up from the current position
// until any character that is NOT a digit
// TODO: handle floats, non-decimal notation, etc
func (l *Lexer) readNumber() string {
	pos := l.pos

	// keep advancing through our input
	// and check if ch is still a digit
	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.pos]
}

// isDigit checks if the given character matches [0-9]
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// isLetter checks if the given character matches [a-zA-Z_]
func isLetter(ch byte) bool {
	// TODO: maybe support for UTF?
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// eatWhitespace is a helper function that keeps advancing
// the current position until it meets a non whitespace character
func (l *Lexer) eatWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// newToken is a helper function that simplifies creating tokens
func newToken(tokenType token.Type, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}
