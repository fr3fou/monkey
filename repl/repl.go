package repl

import (
	"fmt"
	"io"

	"github.com/fr3fou/monkey/parser"

	"github.com/fr3fou/monkey/lexer"
)

const prompt = ">> "

// Start takes in an io.Reader and an io.Writer
// and starts prompting inside and writing to it
func Start(r io.Reader, w io.Writer) {

	fmt.Fprint(w, prompt)
	line := "1 + (2+3) + 4"

	l := lexer.New(line)
	p := parser.New(l)

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(w, p.Errors())
	}

	io.WriteString(w, program.String())
	io.WriteString(w, "\n")
	fmt.Fprint(w, prompt)
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
