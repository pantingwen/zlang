package main

import (
	"flag"
	"fmt"
	"z/compile"
	"z/evaluator"
	"z/lexer"
	"z/object"
	"z/parser"
	"z/vm"
	"time"
)

var engine = flag.String("engine", "vm", "use 'vm' or 'eval'")

var input = `
let fibonacci = fn(x) {
	if (x ==0 ) {
		0
	} else {
		if (x == 1) {
			return 1;
		} else {
			fibonacci(x - 1) + fibonacci(x - 2);
		}
	}
};
fibonacci(35);
`

func main() {
	flag.Parse()

	var duration time.Duration
	var result object.Object

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if *engine == "vm" {
		comp := compile.New()
		start := time.Now()
		err := comp.Compile(program)

		if err != nil {
			fmt.Printf("compile error: %s", err)
		}

		machine := vm.New(comp.Bytecode())

		err = machine.Run()

		if err != nil {
			fmt.Printf("vm error: %s", err)
		}
		duration = time.Since(start)
		result = machine.LastPoppedStackElem()
	} else {
		env := object.NewEnvironment()
		start := time.Now()
		result = evaluator.Eval(program, env)
		duration = time.Since(start)
	}

	fmt.Printf("engine=%s, result=%s, duration=%s\n", *engine, result.Inspect(), duration)

}
