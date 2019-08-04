package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/fr3fou/monkey/token"

	"github.com/fr3fou/monkey/lexer"
)

const prompt = ">> "

// Start takes in an io.Reader and an io.Writer
// and starts prompting inside and writing to it
func Start(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)

	fmt.Fprint(w, prompt)
	for scanner.Scan() {
		l := lexer.NewLexer(scanner.Text())

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(w, "%+v\n", tok)
		}
		fmt.Fprint(w, prompt)
	}
}
