package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"z/parser"
	"z/lexer"
	"z/object"
	"z/ast"
	"z/build"
)

func BuildSourceCode(sourceCode string, sourceFile string) {
	fmt.Println("begin to build code")
	compiledCode := parseAstToC(sourceCode)
	fmt.Println(compiledCode)
	fileName := filepath.Base(sourceFile)
	outFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	compileC(compiledCode, outFileName)
}

func parseAstToC(sourceCode string) string {
	l := lexer.New(sourceCode)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	//compiledCode := generateCode(program, env)
	compiledCode := generateProgram(program, env);
	return compiledCode
}

func generateCode(node ast.Node, env *object.Environment) string {
	switch node := node.(type) {
	case *ast.Program:
		return generateProgram(node, env);
	case *ast.LetStatement:
		return generateLetCode(node, env)
	case *ast.StringLiteral:
		return node.Value
	case *ast.CallExpression:
		fmt.Println("call")
		functionName := generateCode(node.Function, env)
		var callString string;
		callString += functionName + "("
		for _, param := range node.Arguments {
			paramName := generateCode(param, env)
			callString += paramName
		}
		callString += ");\n"
		return callString
	case *ast.FunctionLiteral:
		return node.Name
	case *ast.Identifier:
		return node.Value
	case *ast.InfixExpression:
		return "xx"
	default:
		return "default" + node.TokenLiteral()
	}
}

func generateLetCode(node *ast.LetStatement, env *object.Environment) string{
	code := generateCode(node.Value, env)
	retCode := "char " + node.Name.Value + "[]=\"" + code + "\";"
	return retCode
}

func generateProgram(program *ast.Program, env *object.Environment) string {
	var generateCompiledCode string;
	var compiledCode string;
	for _, statement := range program.Statements {
		fmt.Println("statement:" + statement.String())
		code := build.Eval(statement, env)
		compiledCode = compiledCode + code
	}

	generateCompiledCode += "#include <stdio.h>\n"
	generateCompiledCode += "#include <stdlib.h>\n"
	generateCompiledCode += "int main() {\n"
	generateCompiledCode += compiledCode + "\n"
	generateCompiledCode += "}\n"
	return generateCompiledCode;
}

func compileC(code string, outFile string) {
	tempFileName := "./temp.c";
	file, err := os.Create(tempFileName)
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(code)
	file.Close()
	cmd := exec.Command("gcc", tempFileName, "-o", outFile)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(stdout))
	os.Remove(tempFileName)
}