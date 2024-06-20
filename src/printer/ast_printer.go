package printer

import "github.com/prometheus1400/kel/src/ast"


type AstPrinter struct {

}

func (a *AstPrinter) Print(stmts []ast.Stmt) {

}

// func (a *AstPrinter) VisitBlockStmt(stmt *ast.BlockStmt) {
// }

// func (a *AstPrinter) VisitFnStmt(stmt *ast.FnStmt) {
// }

// func (a *AstPrinter) VisitVarStmt(stmt *ast.VarStmt) {
// }

// func (a *AstPrinter) VisitExpressionStmt(stmt *ast.ExpressionStmt) {
// }

// func (a *AstPrinter) VisitPrintStmt(stmt *ast.PrintStmt) {
// }

// func (a *AstPrinter) VisitNumberExpr(expr *ast.NumberExpr) llvm.Value {
// }

// func (a *AstPrinter) VisitStringExpr(expr *ast.StringExpr) llvm.Value {
// }

// func (a *AstPrinter) VisitIdentifierExpr(expr *ast.IdentifierExpr) llvm.Value {
// }

// func (a *AstPrinter) VisitGroupingExpr(expr *ast.GroupingExpr) llvm.Value {
// }

// func (a *AstPrinter) VisitBinaryExpr(expr *ast.BinaryExpr) llvm.Value {
// }

// func (a *AstPrinter) VisitUnaryExpr(expr *ast.UnaryExpr) llvm.Value {
// }

// func (a *AstPrinter) evaluate(stmt ast.Stmt) {
// 	stmt.Visit(a)
// }

// func (a *AstPrinter) execute(expr ast.Expr) {
// 	expr.Visit(a)
// }
