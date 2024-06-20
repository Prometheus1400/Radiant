package ast

import (
	"github.com/prometheus1400/kel/src/scanner"
)

type VisitStmt interface{
	VisitFnStmt(stmt *FnStmt)
	VisitPrintStmt(stmt *PrintStmt)
	VisitExpressionStmt(stmt *ExpressionStmt)
	VisitReturnStmt(stmt *ReturnStmt)
	VisitIfStmt(stmt *IfStmt)
	VisitBlockStmt(stmt *BlockStmt)
	VisitVarStmt(stmt *VarStmt)
}

type Type struct {
	Token scanner.Token
	IsPointer bool
}
type Param struct {
	Name scanner.Token
	Type Type
}

type VarStmt struct {
	Name scanner.Token
	Type Type
	Initializer Expr
}
func (e *VarStmt) stmt() {}
func (e *VarStmt) Visit(visitor VisitStmt) {visitor.VisitVarStmt(e)}

type FnStmt struct {
	Name scanner.Token
	Params []Param
	Body Stmt
	Return Type
}
func (e *FnStmt) stmt() {}
func (e *FnStmt) Visit(visitor VisitStmt) {visitor.VisitFnStmt(e)}

type PrintStmt struct {
	Expression Expr
}
func (e *PrintStmt) stmt() {}
func (e *PrintStmt) Visit(visitor VisitStmt) {visitor.VisitPrintStmt(e)}

type ExpressionStmt struct {
	Expression Expr
}
func (e *ExpressionStmt) stmt() {}
func (e *ExpressionStmt) Visit(visitor VisitStmt) {visitor.VisitExpressionStmt(e)}

type ReturnStmt struct {
	Expression Expr
}
func (e *ReturnStmt) stmt() {}
func (e *ReturnStmt) Visit(visitor VisitStmt) {visitor.VisitReturnStmt(e)}

type IfStmt struct {
	Condition Expr
	IfBlock Stmt
	ElseBlock Stmt
}
func (e *IfStmt) stmt() {}
func (e *IfStmt) Visit(visitor VisitStmt) {visitor.VisitIfStmt(e)}

type BlockStmt struct {
	Body []Stmt
}
func (e *BlockStmt) stmt() {}
func (e *BlockStmt) Visit(visitor VisitStmt) {visitor.VisitBlockStmt(e)}

