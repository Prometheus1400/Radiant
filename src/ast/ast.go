package ast

import "tinygo.org/x/go-llvm"

type Stmt interface {
	stmt()
	Visit(visitor VisitStmt)
}

type Expr interface {
	expr()
	Visit(visitor VisitExpr) llvm.Value
}
