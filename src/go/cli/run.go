package cli

import (
	"fmt"
	"z/compile"
	"z/evaluator"
	"z/lexer"
	"z/object"
	"z/parser"
	"z/vm"
)

func RunSourceCode(sourceCode string, mode string) {
	l := lexer.New(sourceCode)
	p := parser.New(l)
	program := p.ParseProgram()

	if mode != "vm" {
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
	} else {
		env := object.NewEnvironment()
		evaluator.Eval(program, env)
	}
}
