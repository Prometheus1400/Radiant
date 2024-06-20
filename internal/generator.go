package main

import (
	"fmt"
	"os"
	"strings"
)

type Expressions map[string]string
type Statements map[string]string

func writeExpressionVisitorInterface(expressions Expressions, stringBuilder *strings.Builder) {
	stringBuilder.WriteString("package ast\n\nimport (\n\"github.com/prometheus1400/kel/src/scanner\"\n\"tinygo.org/x/go-llvm\"\n)\ntype VisitExpr interface{\n")
	for name := range expressions {
		fmtName := name + "Expr"
		stringBuilder.WriteString(fmt.Sprintf("\tVisit%s(expr *%s) llvm.Value\n", fmtName, fmtName))
	}
	stringBuilder.WriteString("}\n")
}

func writeExpressions(expressions Expressions, stringBuilder *strings.Builder) {
	for name, args := range expressions {
		fmtName := name + "Expr"
		fmtArgs := "\t" + strings.ReplaceAll(args, ", ", "\n\t")
		e := fmt.Sprintf("type %s struct {\n%s\n}\nfunc (e *%s) expr() {}\nfunc (e *%s) Visit(visitor VisitExpr) llvm.Value {return visitor.Visit%s(e)}", fmtName, fmtArgs, fmtName, fmtName, fmtName)
		stringBuilder.WriteString(e)
		stringBuilder.WriteString("\n\n")
	}
}

func writeStatementVisitorInterface(stmts Statements, stringBuilder *strings.Builder) {
	// stringBuilder.WriteString("package ast\n\nimport \"github.com/prometheus1400/kel/src/scanner\"\n\ntype VisitStmt interface{\n")
	stringBuilder.WriteString("package ast\n\n")
	stringBuilder.WriteString("import (\n")
	// stringBuilder.WriteString("\t\"tinygo.org/x/go-llvm\"\n")
	stringBuilder.WriteString("\t\"github.com/prometheus1400/kel/src/scanner\"\n")
	stringBuilder.WriteString(")\n\n")
	stringBuilder.WriteString("type VisitStmt interface{\n")
	for name := range stmts {
		fmtName := name + "Stmt"
		stringBuilder.WriteString(fmt.Sprintf("\tVisit%s(stmt *%s)\n", fmtName, fmtName))
	}
	stringBuilder.WriteString("}\n\n")
	stringBuilder.WriteString("type Type struct {\n")
	stringBuilder.WriteString("\tToken scanner.Token\n")
	stringBuilder.WriteString("\tIsPointer bool\n")
	stringBuilder.WriteString("}\n")
	stringBuilder.WriteString("type Param struct {\n")
	stringBuilder.WriteString("\tName scanner.Token\n")
	stringBuilder.WriteString("\tType Type\n")
	stringBuilder.WriteString("}\n")
	stringBuilder.WriteString("\n")
}

func writeStatements(stmts Statements, stringBuilder *strings.Builder) {
	for name, args := range stmts {
		fmtName := name + "Stmt"
		fmtArgs := "\t" + strings.ReplaceAll(args, ", ", "\n\t")
		e := fmt.Sprintf("type %s struct {\n%s\n}\nfunc (e *%s) stmt() {}\nfunc (e *%s) Visit(visitor VisitStmt) {visitor.Visit%s(e)}", fmtName, fmtArgs, fmtName, fmtName, fmtName)
		stringBuilder.WriteString(e)
		stringBuilder.WriteString("\n\n")
	}
}

func main() {
	exprString := &strings.Builder{}
	expressions := Expressions{
		"Number":     "Value float64",
		"String":     "Value string",
		"Char":       "Value int8",
		"Bool":       "Value bool",
		"Identifier": "Value scanner.Token",
		"Binary":     "Left Expr, Operator scanner.Token, Right Expr",
		"Unary":      "Operator scanner.Token, Right Expr",
		"Grouping":   "Expression Expr",
		"Call":       "Callee Expr, Args []Expr",
	}
	writeExpressionVisitorInterface(expressions, exprString)
	writeExpressions(expressions, exprString)
	file, err := os.OpenFile("./src/ast/expressions.go", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	file.Write([]byte(exprString.String()))

	stmtString := &strings.Builder{}
	stmts := Statements{
		"Block":      "Body []Stmt",
		"Var":        "Name scanner.Token, Type Type, Initializer Expr",
		"Fn":         "Name scanner.Token, Params []Param, Body Stmt, Return Type",
		"Print":      "Expression Expr",
		"Expression": "Expression Expr",
		"Return":     "Expression Expr",
		"If":         "Condition Expr, IfBlock Stmt, ElseBlock Stmt",
	}
	writeStatementVisitorInterface(stmts, stmtString)
	writeStatements(stmts, stmtString)
	file2, err := os.OpenFile("./src/ast/statements.go", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	file2.Write([]byte(stmtString.String()))
}
