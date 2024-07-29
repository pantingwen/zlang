// parser/parser_test.go

package parser

import (
	"fmt"
	"testing"
	"z/ast"
	"z/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = x + 5;
let foobar = 34;
let name = "sevenpan";
y = y + 1;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 5 {
		t.Fatalf("program.Statements does not contain 5 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
		{"name"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 5 return;
	return 10;
	return 993 322;
	`

	l := lexer.New(input)

	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements, got=%d", len(program.Statements))
	}

	for _, stms := range program.Statements {
		returnStms, ok := stms.(*ast.ReturnStatement)

		if !ok {
			t.Errorf("stms not *ast.ReturnStatement, got=%T", stms)
			continue
		}

		if returnStms.TokenLiteral() != "return" {
			t.Errorf("returnStms.Tokenliteral not 'return', got %q", returnStms.TokenLiteral())
		}
	}
}

func TestIndentifierExpression(t *testing.T) {
	input := "foobar"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements, got=%d", len(program.Statements))
	}

	stms, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	ident, ok := stms.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf("exp not *st.Identifier. got=%T", stms.Expression)
	}

	if ident.Value != "foobar" {
		t.Fatalf("ident.Value not %s, got=%s", "foobar", ident.Value)
	}

	if ident.Token.Literal != "foobar" {
		t.Errorf("ident.TokenLiteral not %s, got=%s", "foobar", ident.TokenLiteral())
	}

}

func TestIntegerExpression(t *testing.T) {
	input := "5"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements, got=%d", len(program.Statements))
	}

	stms, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	ident, ok := stms.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("exp not *st.IntegerLiteral. got=%T", stms.Expression)
	}

	if ident.Value != 5 {
		t.Fatalf("ident.Value not %d, got=%d", 5, ident.Value)
	}

	if ident.Token.Literal != "5" {
		t.Errorf("ident.TokenLiteral not %s, got=%s", "5", ident.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)

		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements, got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", stmt)
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression, got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5+5", 5, "+", 5},
		{"5-5", 5, "-", 5},
		{"5*5", 5, "*", 5},
		{"5/5", 5, "/", 5},
		{"5>5", 5, ">", 5},
		{"5<5", 5, "<", 5},
		{"5==5", 5, "==", 5},
		{"5!=5", 5, "!=", 5},
		{"5<=5", 5, "<=", 5},
		{"5>=5", 5, ">=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)

		if !ok {
			t.Fatalf("exp is not ast.InfixExpression, got %T", stmt.Expression)
		}

		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s', got=%s", tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}

}

func TestBoolExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("expected program.Statements len is %d, but got=%d", 1, len(program.Statements))
		}

		stms, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
		}

		ident, ok := stms.Expression.(*ast.Boolean)

		if !ok {
			t.Fatalf("exp not *st.IntegerLiteral. got=%T", stms.Expression)
		}

		if ident.Value != tt.expected {
			t.Fatalf("ident.Value not %t, got=%t", tt.expected, ident.Value)
		}
	}
}
func TestWhileExpression(t *testing.T) {
	input := "while (x < y) {x}"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	fmt.Println(len(program.Statements))

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statement, got=%d", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.WhileExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not While expression")
	}
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Body.Statements) != 1 {
		t.Fatalf("exp.Body.Statements len is not 1, got=%d", len(exp.Body.Statements))
	}
	body, ok := exp.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.body.Statement[0] is not ast.ExpressionStatement")
	}

	if !testIdentifier(t, body.Expression, "x") {
		return
	}
}
func TestIfExpression(t *testing.T) {
	input := `if (x < y) {x}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements, got=%d", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression, got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("Consequence is not 1 statements, got=%d", len(exp.Consequence.Statements))
	}
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil, got=%+v", exp.Alternative)
	}

}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) {x} else {y}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements, got=%d", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression, got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("Consequence is not 1 statements, got=%d", len(exp.Consequence.Statements))
	}
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Consequence Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative == nil {
		t.Errorf("exp.Alternative.Statements was nil")
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("Alternative is not 1 statements, got=%d", len(exp.Consequence.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Alternative Statements[0] is not ast.ExpressionStatement, got=%T", exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y ; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d Statements, got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stms.Expression is not ast.Functional, got=%T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong, want 2, got=%d\n", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function body has not 1 statements, got=%d", len(function.Body.Statements))
	}

	bodystmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement, got=%T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodystmt.Expression, "x", "+", "y")
}

func TestFunctionParamterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length paramters wrong ,want %d, got=%d", len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5)`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d Statements, got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	epx, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stms.Expression is not ast.Functional, got=%T", stmt.Expression)
	}

	if !testIdentifier(t, epx.Function, "add") {
		return
	}

	if len(epx.Arguments) != 3 {
		t.Fatalf("function literal parameters wrong, want 2, got=%d\n", len(epx.Arguments))
	}

	testLiteralExpression(t, epx.Arguments[0], 1)
	testInfixExpression(t, epx.Arguments[1], 2, "*", 3)
	testInfixExpression(t, epx.Arguments[2], 4, "+", 5)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world"`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q, got=%q", "hello world", literal.Value)
	}

	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 3, 3 + 3]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement")
	}
	array, ok := stmt.Expression.(*ast.ArrayLiteral)

	if !ok {
		t.Fatalf("expt not ast.ArrayLiteral, got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}
	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 3)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpression(t *testing.T) {
	input := "myArr[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp not *ast.ExpressionStatement")
	}

	indexExp, ok := stmt.Expression.(*ast.IndexExpression)

	if !ok {
		t.Fatalf("expt not *ast.IndexExpress. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArr") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("ext is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash pairs was wrong length. got=%d", len(hash.Pairs))
	}
	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)

		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}

		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestFunctionLiteralWithName(t *testing.T) {
	input := `let myFunction = fn(){};`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("program.statements[0] is not set ast.LetStatement.got=%T", program.Statements[0])
	}

	function, ok := stmt.Value.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Value is not ast.Functionaleral. got=%T", stmt.Value)
	}

	if function.Name != "myFunction" {
		t.Fatalf("function literal name wrong, want 'myFunction', got=%q\n", function.Name)
	}
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression, got=%T", exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not %s, got=%s", operator, opExp.Operator)
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}
	return false
}
func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral, got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d, got=%d", value, integ.Value)
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d, got=%s", value, integ.TokenLiteral())
		return false
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)

	if !ok {
		t.Errorf("exp not *ast.Indentifier, got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s,got=%s", value, ident.TokenLiteral())
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parse error: %q", msg)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}
