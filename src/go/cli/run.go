package cli

import (
	"fmt"
	"strings"
	"z/compile"
	"z/evaluator"
	"z/lexer"
	"z/object"
	"z/parser"
	"z/vm"
)

func RunSourceCode(sourceCode string, mode string, fileName string) {
	l := lexer.New(sourceCode)
	p := parser.New(l)
	filePaths := strings.Split(fileName, "/")
	runSourceDir := strings.Join(filePaths[0:len(filePaths)-1], "/")
	p.SetRunSourceDir(runSourceDir)
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
		result := evaluator.Eval(program, env)
		_, ok := result.(*object.Error)
		if ok {
			fmt.Println(result)
		}
	}
}
