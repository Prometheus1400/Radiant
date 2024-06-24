package ast

import (
"github.com/prometheus1400/kel/src/scanner"
"tinygo.org/x/go-llvm"
)
type VisitExpr interface{
	VisitStringExpr(expr *StringExpr) llvm.Value
	VisitGroupingExpr(expr *GroupingExpr) llvm.Value
	VisitNumberExpr(expr *NumberExpr) llvm.Value
	VisitCharExpr(expr *CharExpr) llvm.Value
	VisitBoolExpr(expr *BoolExpr) llvm.Value
	VisitIdentifierExpr(expr *IdentifierExpr) llvm.Value
	VisitBinaryExpr(expr *BinaryExpr) llvm.Value
	VisitUnaryExpr(expr *UnaryExpr) llvm.Value
	VisitCallExpr(expr *CallExpr) llvm.Value
}
type StringExpr struct {
	Value string
}
func (e *StringExpr) expr() {}
func (e *StringExpr) Visit(visitor VisitExpr) llvm.Value {return visitor.VisitStringExpr(e)}

type GroupingExpr struct {
	Expression Expr
}
func (e *GroupingExpr) expr() {}
func (e *GroupingExpr) Visit(visitor VisitExpr) llvm.Value {return visitor.VisitGroupingExpr(e)}

type BinaryExpr struct {
	Left Expr
	Operator scanner.Token
	Right Expr
}
func (e *BinaryExpr) expr() {}
func (e *BinaryExpr) Visit(visitor VisitExpr) llvm.Value {return visitor.VisitBinaryExpr(e)}

type UnaryExpr struct {
	Operator scanner.Token
	Right Expr
}
func (e *UnaryExpr) expr() {}
func (e *UnaryExpr) Visit(visitor VisitExpr) llvm.Value {return visitor.VisitUnaryExpr(e)}

type CallExpr struct {
	Callee Expr
	Args []Expr
}
func (e *CallExpr) expr() {}
func (e *CallExpr) Visit(visitor VisitExpr) llvm.Value {return visitor.VisitCallExpr(e)}

type NumberExpr struct {
	Value float64
}
func (e *NumberExpr) expr() {}
func (e *NumberExpr) Visit(visitor VisitExpr) llvm.Value {return visitor.VisitNumberExpr(e)}

type CharExpr struct {
	Value int8
}
func (e *CharExpr) expr() {}
func (e *CharExpr) Visit(visitor VisitExpr) llvm.Value {return visitor.VisitCharExpr(e)}

type BoolExpr struct {
	Value bool
}
func (e *BoolExpr) expr() {}
func (e *BoolExpr) Visit(visitor VisitExpr) llvm.Value {return visitor.VisitBoolExpr(e)}

type IdentifierExpr struct {
	Value scanner.Token
}
func (e *IdentifierExpr) expr() {}
func (e *IdentifierExpr) Visit(visitor VisitExpr) llvm.Value {return visitor.VisitIdentifierExpr(e)}

