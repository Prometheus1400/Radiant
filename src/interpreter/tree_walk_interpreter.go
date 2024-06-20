package interpreter

import (
	"fmt"

	"github.com/prometheus1400/kel/src/ast"
	"github.com/prometheus1400/kel/src/scanner"
)

// implements Statement and Expression Visitor
type TreeWalkInterpreter struct {
}

func NewTreeWalkInterpreter() *TreeWalkInterpreter {
	return &TreeWalkInterpreter{}
}

func (i *TreeWalkInterpreter) Interpret(stmt *ast.BlockStmt) {
	i.execute(stmt)
}

func (i *TreeWalkInterpreter) VisitBlockStmt(stmt *ast.BlockStmt) interface{} {
	for _, statement := range stmt.Body {
		i.execute(statement)
	}
	return nil
}
func (i *TreeWalkInterpreter) VisitExpressionStmt(stmt *ast.ExpressionStmt) interface{} {
	i.evaluate(stmt.Expression)
	return nil
}
func (i *TreeWalkInterpreter) VisitPrintStmt(stmt *ast.PrintStmt) interface{} {
	res := i.evaluate(stmt.Expression)
	fmt.Println(res)
	return nil
}

func (i *TreeWalkInterpreter) VisitNumberExpr(expr *ast.NumberExpr) interface{} {
	return expr.Value
}
func (i *TreeWalkInterpreter) VisitStringExpr(expr *ast.StringExpr) interface{} {
	return nil
}
func (i *TreeWalkInterpreter) VisitIdentifierExpr(expr *ast.IdentifierExpr) interface{} {
	return nil
}
func (i *TreeWalkInterpreter) VisitBinaryExpr(expr *ast.BinaryExpr) interface{} {
	left := i.evaluate(expr.Left).(float64)
	right := i.evaluate(expr.Right).(float64)
	switch expr.Operator.Type {
	case scanner.PLUS:
		return left + right
	}
	return nil
}
func (i *TreeWalkInterpreter) VisitUnaryExpr(expr *ast.UnaryExpr) interface{} {
	return nil
}

func (i *TreeWalkInterpreter) execute(stmt ast.Stmt) {
	stmt.Visit(i)
}

func (i *TreeWalkInterpreter) evaluate(expr ast.Expr) interface{} {
	return expr.Visit(i)
}
