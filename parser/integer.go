package parser

import (
	"fmt"
	"strconv"

	"github.com/fr3fou/monkey/ast"
)

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
