package cli

import "fmt"
import "z/lexer"
import "z/parser"
import "z/compile"
import "z/vm"

func RunSourceCode(sourceCode string) {
	l := lexer.New(sourceCode)
	p := parser.New(l)
	program := p.ParseProgram()
	comp := compile.New()
	err := comp.Compile(program)

	if err != nil {
		fmt.Printf("compile error: %s", err)
	}

	machine := vm.New(comp.Bytecode())

	err = machine.Run()

	if err != nil {
		fmt.Printf("vm error: %s", err)
	}
}