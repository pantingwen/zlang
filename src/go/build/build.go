package build

import (
	"z/ast"
	"z/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) string {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		return "char " + node.Name.Value + "[]=" + val + ";\n"
	case *ast.Identifier:
		return node.Value
	case *ast.FunctionLiteral:
		return node.Name
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		callString := ""
		switch function {
		case "puts":
			for _, argument := range node.Arguments {
				argumentRes := Eval(argument, env)
				callString += "printf(\"%s\"," + argumentRes + ");\n"
			}
		}
		return callString
	case *ast.StringLiteral:
		return string("\"" + node.Value + "\"")
	}
	return ""
}

func evalProgram(program *ast.Program, env *object.Environment) string {
	var result string
	for _, statement := range program.Statements {
		result += Eval(statement, env)
	}
	return result
}

func nativeBoolToBooleanObject(input bool) string {
	if input {
		return "true"
	}
	return "false"
}
