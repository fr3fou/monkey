package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/fr3fou/monkey/parser"

	"github.com/fr3fou/monkey/lexer"
)

const prompt = ">> "

// Start takes in an io.Reader and an io.Writer
// and starts prompting inside and writing to it
func Start(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)

	fmt.Fprint(w, prompt)
	for scanner.Scan() {
		line := scanner.Text()

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(w, p.Errors())
			continue
		}

		io.WriteString(w, program.String())
		io.WriteString(w, "\n")
		fmt.Fprint(w, prompt)
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
