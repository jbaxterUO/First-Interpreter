package parser

import (
	"WeekTwo/ast"
	"WeekTwo/lexer"
	"fmt"
	"testing"
)

func testLetStatement(testInput *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		testInput.Fatalf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		testInput.Fatalf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		testInput.Fatalf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		testInput.Fatalf("s.Name not '%s'. got=%s", name, letStmt.Name)
		return false
	}

	return true
}

func checkParserErrors(testInput *testing.T, parser *Parser) {
	errors := parser.Errors()
	if len(errors) == 0 {
		return
	}

	testInput.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		testInput.Errorf("parser error: %q", msg)
	}
	testInput.FailNow()
}

func testIntegerLiteral(testInput *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		testInput.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		testInput.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		testInput.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func TestLetStatements(testInput *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`

	lex := lexer.New(input)
	p := New(lex)

	program := p.ParseProgram()
	checkParserErrors(testInput, p)

	if program == nil {
		testInput.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		testInput.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(testInput, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatements(testInput *testing.T) {
	input := `
	return 5;
	return 10;
	return 993322;
	`

	lex := lexer.New(input)
	p := New(lex)

	program := p.ParseProgram()
	checkParserErrors(testInput, p)

	if program == nil {
		testInput.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		testInput.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			testInput.Fatalf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			testInput.Fatalf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(testInput *testing.T) {
	input := "foobar;"

	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	checkParserErrors(testInput, parser)

	if len(program.Statements) != 1 {
		testInput.Fatalf("program does not have enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		testInput.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)

	if !ok {
		testInput.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		testInput.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		testInput.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(testInput *testing.T) {
	input := "5;"

	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	checkParserErrors(testInput, parser)

	if len(program.Statements) != 1 {
		testInput.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		testInput.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		testInput.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != 5 {
		testInput.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		testInput.Errorf("literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(testInput *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!10;", "!", 10},
		{"-22;", "-", 22},
	}

	for _, tt := range prefixTests {
		lex := lexer.New(tt.input)
		parser := New(lex)
		program := parser.ParseProgram()
		checkParserErrors(testInput, parser)

		if len(program.Statements) != 1 {
			testInput.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			testInput.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			testInput.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			testInput.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(testInput, exp.Right, tt.integerValue) {
			return
		}
	}
}
