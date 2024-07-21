package build

import (
	"fmt"
	"z/ast"
	"z/object"
	//"strings"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

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
	// case *ast.PrefixExpression:
	// 	right := Eval(node.Right, env)
	// 	if isError(right) {
	// 		return right
	// 	}
	// 	fmt.Println("pre...")
	// 	return evalPrefixExpression(node.Operator, right)
	case *ast.Identifier:
		fmt.Println("idt ....")
		return node.Value;
	// case *ast.InfixExpression:
	// 	left := Eval(node.Left, env)
	// 	if isError(left) {
	// 		return left
	// 	}
	// 	right := Eval(node.Right, env)
	// 	if isError(right) {
	// 		return right
	// 	}
	// 	return evalInfixExpression(node.Operator, left, right)
	// case *ast.IfExpression:
	// 	return evalIfExpression(node, env)
	// case *ast.BlockStatement:
	// 	return evalBlockStatement(node, env)
	//case *ast.IntegerLiteral:
	//	return string(node.Value)
	// case *ast.ReturnStatement:
	// 	val := Eval(node.ReturnValue, env)
	// 	if isError(val) {
	// 		return val
	// 	}
	// 	return &object.ReturnValue{Value: val}
	case *ast.FunctionLiteral:
		return node.Name
	case *ast.CallExpression:
		fmt.Println("call ....")
		function := Eval(node.Function, env)
		callString := ""
		switch function {
		case "puts":
			for _, argument := range node.Arguments {
				argumentRes := Eval(argument, env)
				callString += "printf(" + argumentRes + ");\n"
			}
		}
		// var callString string = ""
		// callString += function
		// fmt.Println("function name" + callString)
		// callString += "("
		// for _, argument := range node.Arguments {
		// 	callString += Eval(argument, env)
		// 	callString += ","
		// }
		// callString = strings.TrimRight(callString, ",")
		// callString += ");\n"
		return callString

	case *ast.StringLiteral:
		fmt.Println("string....")
		return string("\"" + node.Value + "\"")
	// case *ast.ArrayLiteral:
	// 	elements := evalExpressions(node.Elements, env)
	// 	if len(elements) == 1 && isError(elements[0]) {
	// 		return elements[0]
	// 	}
	// 	return &object.Array{Elements: elements}
	// case *ast.IndexExpression:
	// 	left := Eval(node.Left, env)
	// 	if isError(left) {
	// 		return left
	// 	}
	// 	index := Eval(node.Index, env)
	// 	if isError(index) {
	// 		return index
	// 	}
	// 	return evalIndexExpression(left, index)
	// case *ast.HashLiteral:
	// 	return evalHashLiteral(node, env)
}
	return ""
}


func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s ", left.Type())
	}
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)
	key, ok := index.(object.Hashable)

	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]

	if !ok {
		return NULL
	}
	return pair.Value
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObj := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObj.Elements) - 1)
	if idx < 0 || idx > max {
		return NULL
	}
	return arrayObj.Elements[idx]
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnviroment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}


func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case FALSE:
		return false
	case TRUE:
		return true
	default:
		return true
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	return &object.String{Value: leftVal + rightVal}
}


func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperationExpression(right)
	default:
		return newError("unkown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperationExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
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
