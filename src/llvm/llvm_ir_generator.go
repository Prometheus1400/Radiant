package llvm

import (
	"fmt"
	"os"

	"github.com/prometheus1400/kel/src/ast"
	"github.com/prometheus1400/kel/src/environment"
	"github.com/prometheus1400/kel/src/scanner"
	"tinygo.org/x/go-llvm"
)

type IRGenerator struct {
	ctx               llvm.Context
	module            llvm.Module
	builder           llvm.Builder
	depth             int
	environment       *environment.Environment[llvm.Value]
	identifierAddress bool
	currentFunction   llvm.Value
}

func (g *IRGenerator) Init() {
	g.depth = 0
	g.environment = environment.NewEnvironment[llvm.Value](nil)
	g.identifierAddress = false

	llvm.InitializeAllTargetInfos()
	llvm.InitializeAllTargets()
	llvm.InitializeAllTargetMCs()
	llvm.InitializeAllAsmParsers()
	llvm.InitializeAllAsmPrinters()
}

func NewIRGenerator() *IRGenerator {
	gen := &IRGenerator{}
	// gen.Init()
	return gen
}

func (g *IRGenerator) GenerateIR(stmts []ast.Stmt, outputFile string) {
	g.Init()
	g.ctx = llvm.NewContext()
	g.module = g.ctx.NewModule("example")
	g.builder = g.ctx.NewBuilder()
	defer g.ctx.Dispose()
	defer g.module.Dispose()
	defer g.builder.Dispose()

	g.defineBuiltInTypes()
	g.declareExternalFuncs()

	for _, stmt := range stmts {
		g.execute(stmt)
	}

	// g.module.Dump()
	file, _ := os.OpenFile("build/"+outputFile+".ll", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	file.WriteString(g.module.String())
	// fmt.Fprintf(file, g.module.String())
	file.Close()
}

func (g *IRGenerator) defineBuiltInTypes() {
}

func (g *IRGenerator) declareExternalFuncs() {
	printfType := llvm.FunctionType(g.ctx.Int32Type(), []llvm.Type{llvm.PointerType(g.ctx.Int8Type(), 0)}, true)
	printf := llvm.AddFunction(g.module, "printf", printfType)
	g.environment.Define("printf")
	g.environment.Set("printf", printf)
}

func (g *IRGenerator) VisitBlockStmt(stmt *ast.BlockStmt) {
	oldTable := g.environment
	newTable := environment.NewEnvironment[llvm.Value](oldTable)
	g.environment = newTable
	g.depth++
	for _, stmt_ := range stmt.Body {
		g.execute(stmt_)
	}
	g.depth--
	g.environment = oldTable
}

func (g *IRGenerator) VisitFnStmt(stmt *ast.FnStmt) {
	returnType := g.llvmTypeFromAstType(stmt.Return)
	paramTypes := make([]llvm.Type, 0)
	for _, param := range stmt.Params {
		paramType := g.llvmTypeFromAstType(param.Type)
		paramTypes = append(paramTypes, paramType)
	}
	fnType := llvm.FunctionType(returnType, paramTypes, false)
	fn := llvm.AddFunction(g.module, stmt.Name.Lexeme, fnType)
	entry := llvm.AddBasicBlock(fn, "entry")
	g.builder.SetInsertPointAtEnd(entry)

	g.environment.Define(stmt.Name.Lexeme)
	g.environment.Set(stmt.Name.Lexeme, fn)

	prevEnv := g.environment
	g.environment = environment.NewEnvironment[llvm.Value](prevEnv)
	for i, param := range stmt.Params {
		fnParam := fn.Param(i)
		fnParam.SetName(param.Name.Lexeme)
		g.environment.Define(param.Name.Lexeme)
		g.environment.Set(param.Name.Lexeme, fnParam)
	}
	g.currentFunction = fn
	g.execute(stmt.Body)
	g.environment = prevEnv
}

func (g *IRGenerator) VisitReturnStmt(stmt *ast.ReturnStmt) {
	if stmt.Expression == nil {
		g.builder.CreateRetVoid()
		return
	}
	res := g.evaluate(stmt.Expression)
	g.builder.CreateRet(res)
}

func (g *IRGenerator) VisitVarStmt(stmt *ast.VarStmt) {
	// assuming type checking pass has already been done by this point
	initializer := g.evaluate(stmt.Initializer)
	// initializerVal := initializer
	// if initializer.Type().TypeKind() == llvm.PointerTypeKind {
	// 	initializerVal = g.builder.CreateLoad(initializer.AllocatedType(), initializer, "val")
	// }

	var llvmType llvm.Type
	if stmt.Type.Token.Lexeme == "auto" {
		llvmType = initializer.Type()
	} else {
		llvmType = g.llvmTypeFromAstType(stmt.Type)
	}
	var varPtr llvm.Value
	if g.depth == 0 {
		varPtr = llvm.AddGlobal(g.module, llvmType, stmt.Name.Lexeme)
		varPtr.SetInitializer(initializer)
	} else {
		varPtr = g.builder.CreateAlloca(llvmType, stmt.Name.Lexeme)
		g.builder.CreateStore(initializer, varPtr)
	}
	g.environment.Define(stmt.Name.Lexeme)
	g.environment.Set(stmt.Name.Lexeme, varPtr)
}

func (g *IRGenerator) VisitIfStmt(stmt *ast.IfStmt) {
	conditionVal := g.evaluate(stmt.IfCondition)
	ifBlock := llvm.AddBasicBlock(g.currentFunction, "ifBlock")
	elseBlock := llvm.AddBasicBlock(g.currentFunction, "elseBlock")
	mergeBlock := llvm.AddBasicBlock(g.currentFunction, "mergeBlock")

	g.builder.CreateCondBr(conditionVal, ifBlock, elseBlock)

	g.builder.SetInsertPointAtEnd(ifBlock)
	g.execute(stmt.IfBlock)
	g.builder.CreateBr(mergeBlock)
	g.builder.SetInsertPointAtEnd(elseBlock)
    for i, _ := range stmt.ElifConditions {
        elifBlock := llvm.AddBasicBlock(g.currentFunction, fmt.Sprint("elifBlock-%d", i))
        elifElseBlock := llvm.AddBasicBlock(g.currentFunction, fmt.Sprint("elifElseBlock-%d", i))

        elifCondition := g.evaluate(stmt.ElifConditions[i])
        elifStmts := stmt.ElifBlocks[i]

        g.builder.CreateCondBr(elifCondition, elifBlock, elifElseBlock)
        g.builder.SetInsertPointAtEnd(elifBlock)
        g.execute(elifStmts)
        g.builder.CreateBr(mergeBlock)
        g.builder.SetInsertPointAtEnd(elifElseBlock)
    }
	if stmt.ElseBlock != nil {
		g.execute(stmt.ElseBlock)
	}
	g.builder.CreateBr(mergeBlock)
	g.builder.SetInsertPointAtEnd(mergeBlock)
}

func (g *IRGenerator) VisitExpressionStmt(stmt *ast.ExpressionStmt) {
	g.evaluate(stmt.Expression)
}

func (g *IRGenerator) VisitPrintStmt(stmt *ast.PrintStmt) {
}

func (g *IRGenerator) VisitNumberExpr(expr *ast.NumberExpr) llvm.Value {
	val := llvm.ConstFloat(g.ctx.DoubleType(), expr.Value)
	return val
}

func (g *IRGenerator) VisitStringExpr(expr *ast.StringExpr) llvm.Value {
	str := llvm.ConstString(expr.Value, true)
	strPtr := g.builder.CreateAlloca(str.Type(), "")
	g.builder.CreateStore(str, strPtr)
	return strPtr
}

func (g *IRGenerator) VisitCharExpr(expr *ast.CharExpr) llvm.Value {
	char := llvm.ConstInt(g.ctx.Int8Type(), uint64(expr.Value), true)
	return char
}

func (g *IRGenerator) VisitBoolExpr(expr *ast.BoolExpr) llvm.Value {
	var boolVal uint64
	if expr.Value {
		boolVal = 1
	} else {
		boolVal = 0
	}
	val := llvm.ConstInt(g.ctx.Int1Type(), boolVal, false)
	return val
}

func (g *IRGenerator) VisitIdentifierExpr(expr *ast.IdentifierExpr) llvm.Value {
	name := expr.Value.Lexeme
	varPtr, exists := g.environment.Get(name)
	if !exists {
		panic("trying to reference undefined identifier")
	}

	if !varPtr.IsAFunction().IsNil() || g.identifierAddress || !varPtr.IsAArgument().IsNil() {
		return varPtr
	}

	varVal := g.builder.CreateLoad(varPtr.AllocatedType(), varPtr, "")
	return varVal
}

func (g *IRGenerator) VisitGroupingExpr(expr *ast.GroupingExpr) llvm.Value {
	return llvm.Value{}
}

func (g *IRGenerator) VisitCallExpr(expr *ast.CallExpr) llvm.Value {
	fn := g.evaluate(expr.Callee)
	args := make([]llvm.Value, 0)
	for _, arg := range expr.Args {
		argTmp := g.evaluate(arg)
		args = append(args, argTmp)
	}
	return g.builder.CreateCall(fn.GlobalValueType(), fn, args, "callRes")
}

func (g *IRGenerator) VisitBinaryExpr(expr *ast.BinaryExpr) llvm.Value {
	lhsVal := expr.Left.Visit(g)
	rhsVal := expr.Right.Visit(g)
	switch expr.Operator.Type {
	case scanner.PLUS:
		return g.builder.CreateFAdd(lhsVal, rhsVal, "add")
	case scanner.MINUS:
		return g.builder.CreateFSub(lhsVal, rhsVal, "subtract")
	case scanner.STAR:
		return g.builder.CreateFMul(lhsVal, rhsVal, "multiply")
	case scanner.SLASH:
		return g.builder.CreateExactSDiv(lhsVal, rhsVal, "divide")
	case scanner.LESS:
		return g.builder.CreateFCmp(llvm.FloatOLT, lhsVal, rhsVal, "less than")
	case scanner.LESS_EQ:
		return g.builder.CreateFCmp(llvm.FloatOLE, lhsVal, rhsVal, "less than or equal to")
	case scanner.GREATER:
		return g.builder.CreateFCmp(llvm.FloatOGT, lhsVal, rhsVal, "greater than")
	case scanner.GREATER_EQ:
		return g.builder.CreateFCmp(llvm.FloatOGE, lhsVal, rhsVal, "greater than or equal to")
	case scanner.EQUAL:
		return g.builder.CreateFCmp(llvm.FloatUEQ, lhsVal, rhsVal, "equal")
	case scanner.NOT_EQUAL:
		return g.builder.CreateFCmp(llvm.FloatUNE, lhsVal, rhsVal, "not equal")
	default:
		panic(fmt.Sprintf("can't handle operator '%s' in binary expression", expr.Operator.Lexeme))
	}
}
func (g *IRGenerator) VisitUnaryExpr(expr *ast.UnaryExpr) llvm.Value {
	switch expr.Operator.Type {
	case scanner.MINUS:
		right := g.evaluate(expr.Right)
		return g.builder.CreateFNeg(right, "negate")
	case scanner.ADDRESS:
		g.identifierAddress = true
		right := g.evaluate(expr.Right)
		g.identifierAddress = false
		return right
	case scanner.STAR:
		fmt.Println("here")
		right := g.evaluate(expr.Right)
		return g.builder.CreateLoad(right.Type(), right, "dereference")
	default:
		panic(fmt.Sprintf("unhandled unary operator '%s'", expr.Operator.Lexeme))
	}
}

func (g *IRGenerator) execute(stmt ast.Stmt) {
	stmt.Visit(g)
}

func (g *IRGenerator) evaluate(expr ast.Expr) llvm.Value {
	return expr.Visit(g)
}

func (g *IRGenerator) llvmTypeFromAstType(langType ast.Type) llvm.Type {
	// assume it's always a TYPE token
	var llvmType llvm.Type
	if langType.Token.IsPrimitiveType() {
		switch langType.Token.Lexeme {
		case "number":
			llvmType = g.ctx.DoubleType()
		case "string":
			llvmType = llvm.PointerType(g.ctx.Int8Type(), 0)
		case "bool":
			llvmType = g.ctx.Int1Type()
		case "char":
			fmt.Println("here1")
			llvmType = g.ctx.Int8Type()
		case "void":
			llvmType = g.ctx.VoidType()
		}
	} else {
		// TODO handle lookups of custom types
	}
	if langType.IsPointer {
		fmt.Println("here2")
		llvmType = llvm.PointerType(llvmType, 0)
	}

	return llvmType
}
