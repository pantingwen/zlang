// repl/repl.go

package repl

import (
	"bufio"
	"fmt"
	"io"
	"z/evaluator"
	"z/lexer"
	"z/object"
	"z/parser"
	"z/util"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			util.PrintErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)

		// io.WriteString(out, program.String())
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}