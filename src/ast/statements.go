package ast

import (
	"github.com/prometheus1400/kel/src/scanner"
)

type VisitStmt interface{
	VisitPrintStmt(stmt *PrintStmt)
	VisitExpressionStmt(stmt *ExpressionStmt)
	VisitReturnStmt(stmt *ReturnStmt)
	VisitIfStmt(stmt *IfStmt)
	VisitBlockStmt(stmt *BlockStmt)
	VisitVarStmt(stmt *VarStmt)
	VisitFnStmt(stmt *FnStmt)
}

type Type struct {
	Token scanner.Token
	IsPointer bool
}
type Param struct {
	Name scanner.Token
	Type Type
}

type IfStmt struct {
	IfCondition Expr
	IfBlock Stmt
	ElifConditions []Expr
	ElifBlocks []Stmt
	ElseBlock Stmt
}
func (e *IfStmt) stmt() {}
func (e *IfStmt) Visit(visitor VisitStmt) {visitor.VisitIfStmt(e)}

type BlockStmt struct {
	Body []Stmt
}
func (e *BlockStmt) stmt() {}
func (e *BlockStmt) Visit(visitor VisitStmt) {visitor.VisitBlockStmt(e)}

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

