package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"z/ast"
	"z/build"
	"z/lexer"
	"z/object"
	"z/parser"
)

func BuildSourceCode(sourceCode string, sourceFile string) {
	var duration time.Duration
	start := time.Now()
	fmt.Println("begin to build code")
	compiledCode := parseAstToC(sourceCode)
	fileName := filepath.Base(sourceFile)
	outFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	compileC(compiledCode, outFileName)
	duration = time.Since(start)
	fmt.Printf("build execute time is :%s\n", duration)
}

func parseAstToC(sourceCode string) string {
	l := lexer.New(sourceCode)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	compiledCode := generateProgram(program, env)
	return compiledCode
}

func generateProgram(program *ast.Program, env *object.Environment) string {
	var generateCompiledCode string
	var compiledCode string
	for _, statement := range program.Statements {
		code := build.Eval(statement, env)
		compiledCode = compiledCode + code
	}

	generateCompiledCode += "#include <stdio.h>\n"
	generateCompiledCode += "#include <stdlib.h>\n"
	generateCompiledCode += "int main() {\n"
	generateCompiledCode += compiledCode + "\n"
	generateCompiledCode += "}\n"
	return generateCompiledCode
}

func compileC(code string, outFile string) {
	tempFileName := "./temp.c"
	file, err := os.Create(tempFileName)
	if err != nil {
		panic(err)
	}
	_, _ = file.WriteString(code)
	file.Close()
	cmd := exec.Command("gcc", tempFileName, "-o", outFile)
	_, err = cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("build success!")
	}
	os.Remove(tempFileName)
}
