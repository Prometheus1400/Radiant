package ast

import (
"github.com/prometheus1400/kel/src/scanner"
"tinygo.org/x/go-llvm"
)
type VisitExpr interface{
	VisitBoolExpr(expr *BoolExpr) llvm.Value
	VisitIdentifierExpr(expr *IdentifierExpr) llvm.Value
	VisitBinaryExpr(expr *BinaryExpr) llvm.Value
	VisitGroupingExpr(expr *GroupingExpr) llvm.Value
	VisitCallExpr(expr *CallExpr) llvm.Value
	VisitNumberExpr(expr *NumberExpr) llvm.Value
	VisitStringExpr(expr *StringExpr) llvm.Value
	VisitCharExpr(expr *CharExpr) llvm.Value
	VisitUnaryExpr(expr *UnaryExpr) llvm.Value
}
type CallExpr struct {
	Callee Expr
	Args []Expr
}
func (e *CallExpr) expr() {}
func (e *CallExpr) Visit(visitor VisitExpr) llvm.Value {return visitor.VisitCallExpr(e)}

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

type BinaryExpr struct {
	Left Expr
	Operator scanner.Token
	Right Expr
}
func (e *BinaryExpr) expr() {}
func (e *BinaryExpr) Visit(visitor VisitExpr) llvm.Value {return visitor.VisitBinaryExpr(e)}

type GroupingExpr struct {
	Expression Expr
}
func (e *GroupingExpr) expr() {}
func (e *GroupingExpr) Visit(visitor VisitExpr) llvm.Value {return visitor.VisitGroupingExpr(e)}

type NumberExpr struct {
	Value float64
}
func (e *NumberExpr) expr() {}
func (e *NumberExpr) Visit(visitor VisitExpr) llvm.Value {return visitor.VisitNumberExpr(e)}

type StringExpr struct {
	Value string
}
func (e *StringExpr) expr() {}
func (e *StringExpr) Visit(visitor VisitExpr) llvm.Value {return visitor.VisitStringExpr(e)}

type CharExpr struct {
	Value int8
}
func (e *CharExpr) expr() {}
func (e *CharExpr) Visit(visitor VisitExpr) llvm.Value {return visitor.VisitCharExpr(e)}

type UnaryExpr struct {
	Operator scanner.Token
	Right Expr
}
func (e *UnaryExpr) expr() {}
func (e *UnaryExpr) Visit(visitor VisitExpr) llvm.Value {return visitor.VisitUnaryExpr(e)}

