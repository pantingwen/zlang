// repl/repl.go

package repl

import (
	"bufio"
	"fmt"
	"io"
	"z/compile"
	"z/lexer"
	"z/object"
	"z/parser"
	"z/vm"
	"z/util"
)

func CStart(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	fmt.Printf("run as compile vm mode\n")

	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalSize)
	symbolTable := compile.NewSymbolTable()

	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}
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

		comp := compile.NewWithState(symbolTable, constants)
		err := comp.Compile(program)

		if err != nil {
			fmt.Fprintf(out, "woops! compilation failed:\n  %s\n", err)
			continue
		}

		machine := vm.NewWithGlobalsStore(comp.Bytecode(), globals)
		err = machine.Run()

		if err != nil {
			fmt.Fprintf(out, "woops! execute bytecode faild\n %s\n", err)
			continue
		}

		stackTop := machine.LastPoppedStackElem()

		io.WriteString(out, stackTop.Inspect())
		io.WriteString(out, "\n")
	}
}
