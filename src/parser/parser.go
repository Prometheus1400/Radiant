package parser

import (
	"fmt"
	"strconv"

	"github.com/prometheus1400/kel/src/ast"
	"github.com/prometheus1400/kel/src/scanner"
)

type Precedence int

const (
	PREC_NONE       Precedence = iota
	PREC_ASSIGNMENT            // =
	PREC_OR                    // or
	PREC_AND                   // and
	PREC_EQUALITY              // ==
	PREC_COMPARISON            // < > <= >=
	PREC_TERM                  // + -
	PREC_FACTOR                // * /
	PREC_UNARY                 // ! -
	PREC_CALL                  // . ()
	PREC_PRIMARY
)

type ParseRule struct {
	PrefixRule func(p *Parser) (ast.Expr, error)
	InfixRule  func(p *Parser, left ast.Expr) (ast.Expr, error) 
	Precedence Precedence
}

type ParseTable struct {
	table map[scanner.TokenType]ParseRule
}

func NewParseTable() *ParseTable {
	return &ParseTable{
		table: map[scanner.TokenType]ParseRule{
			scanner.LEFT_PAREN:  {grouping, call, PREC_CALL},
			scanner.RIGHT_PAREN: {nil, nil, PREC_NONE},
			scanner.LEFT_BRACE:  {nil, nil, PREC_NONE},
			scanner.RIGHT_BRACE: {nil, nil, PREC_NONE},
			scanner.LEFT_BRACK:  {nil, nil, PREC_NONE},
			scanner.RIGHT_BRACK: {nil, nil, PREC_NONE},
			scanner.SEMI_COLON:  {nil, nil, PREC_NONE},
			scanner.COMMA:       {nil, nil, PREC_NONE},
			scanner.PLUS:        {nil, binary, PREC_TERM},
			scanner.MINUS:       {unary, binary, PREC_TERM},
			scanner.STAR:        {unary, binary, PREC_FACTOR},
			scanner.SLASH:       {nil, binary, PREC_FACTOR},
			scanner.BANG:        {unary, nil, PREC_UNARY},
			scanner.ADDRESS:     {unary, nil, PREC_UNARY},
			scanner.ASSIGN:      {nil, nil, PREC_NONE},
			scanner.DOT:         {nil, nil, PREC_NONE},
			scanner.PLUSPLUS:    {nil, nil, PREC_NONE},
			scanner.MINUSMINUS:  {nil, nil, PREC_NONE},
			scanner.EQUAL:       {nil, binary, PREC_EQUALITY},
			scanner.NOT_EQUAL:   {nil, binary, PREC_EQUALITY},
			scanner.LESS:        {nil, binary, PREC_COMPARISON},
			scanner.GREATER:     {nil, binary, PREC_COMPARISON},
			scanner.LESS_EQ:     {nil, binary, PREC_COMPARISON},
			scanner.GREATER_EQ:  {nil, binary, PREC_COMPARISON},
			scanner.NUMBER:      {number, nil, PREC_PRIMARY},
			scanner.STRING:      {string_, nil, PREC_PRIMARY},
			scanner.CHAR:        {char, nil, PREC_PRIMARY},
			scanner.BOOL:        {nil, nil, PREC_PRIMARY},
			scanner.IDENTIFIER:  {variable, nil, PREC_PRIMARY},
			scanner.LET:         {nil, nil, PREC_NONE},
			scanner.TRUE:        {boolean, nil, PREC_NONE},
			scanner.FALSE:       {boolean, nil, PREC_NONE},
			scanner.FN:          {nil, nil, PREC_NONE},
			scanner.IF:          {nil, nil, PREC_NONE},
			scanner.ELSE:        {nil, nil, PREC_NONE},
			scanner.RETURN:      {nil, nil, PREC_NONE},
			scanner.EOF:         {nil, nil, PREC_NONE},
			// scanner.DOTDOT:      {nil, nil, PREC_NONE},
			// scanner.DOTDOTDOT:   {nil, nil, PREC_NONE},
			// scanner.STRUCT:      {nil, nil, PREC_NONE},
			// scanner.ENUM:        {nil, nil, PREC_NONE},
			// scanner.ELIF:        {nil, nil, PREC_NONE},
			// scanner.PUB:         {nil, nil, PREC_NONE},
			// scanner.FOR:         {nil, nil, PREC_NONE},
			// scanner.WHILE:       {nil, nil, PREC_NONE},
			// scanner.BREAK:       {nil, nil, PREC_NONE},
			// scanner.CONTINUE:    {nil, nil, PREC_NONE},
			// scanner.IMPORT:      {nil, nil, PREC_NONE},
		},
	}
}

func (t *ParseTable) GetRule(type_ scanner.TokenType) ParseRule {
	rule, exists := t.table[type_]
	if !exists {
		panic("undefined rule")
	}
	// assert.Condition(exists == true)
	return rule
}

type Parser struct {
	HadError   bool
	Errors     []error
	tokens     []scanner.Token
	start      int
	current    int
	parseTable *ParseTable
}

func NewParser() *Parser {
	var parser Parser
	parser.Init()
	return &parser
}

func (p *Parser) Init() {
	p.tokens = nil
	p.start = 0
	p.current = 0
	p.HadError = false
	p.Errors = make([]error, 0)
	if p.parseTable == nil {
		p.parseTable = NewParseTable()
	}
}

func (p *Parser) Parse(tokens []scanner.Token) []ast.Stmt {
	p.Init()
	p.tokens = tokens

	stmts := make([]ast.Stmt, 0)
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			p.synchronize()
		} else {
			stmts = append(stmts, stmt)
		}
	}

	return stmts
}

func (p *Parser) ReportErrors() {
	for _, err := range p.Errors {
		fmt.Println(err)
	}
}

func (p *Parser) declaration() (ast.Stmt, error) {
	if p.match(scanner.LET) {
		return p.varDeclaration()
	} else if p.match(scanner.FN) {
		return p.fnDeclaration()
	} else {
		return p.statement()
	}
}

func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(scanner.IF) {
		return p.ifStmt()
	} else if p.match(scanner.LEFT_BRACE) {
		return p.blockStmt()
	} else if p.match(scanner.RETURN) {
		return p.returnStmt()
	} else {
		return p.expressionStmt()
	}
}

func (p *Parser) varDeclaration() (ast.Stmt, error) {
	name, err := p.consume(scanner.IDENTIFIER, "expect variable name")
	if err != nil {
		return nil, err
	}

	var typeToken scanner.Token = scanner.Token{Type: scanner.TYPE, Lexeme: "auto"}
	isPointer := false
	if p.match(scanner.STAR) {
		isPointer = true
	}
	if p.check(scanner.TYPE) || p.check(scanner.IDENTIFIER) {
		typeToken, err = p.consumeType("expected type after variable name in variable declaration")
		if err != nil {
			return nil, err
		}
	}

	var initializer ast.Expr = nil
	if p.match(scanner.ASSIGN) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	if typeToken.Lexeme == "auto" && isPointer {
		p.errorAtCurrent("need to supply type with '*'")
	}
	if typeToken.Lexeme == "auto" && initializer == nil {
		p.errorAtCurrent("cannot infer type without initializer")
	}

	_, err = p.consume(scanner.SEMI_COLON, "expect semicolon after variable declaration")
	if err != nil {
		return nil, err
	}

	return &ast.VarStmt{Name: name, Type: ast.Type{Token: typeToken, IsPointer: isPointer}, Initializer: initializer}, nil
}

func (p *Parser) fnDeclaration() (ast.Stmt, error) {
	name, err := p.consume(scanner.IDENTIFIER, "expect function name")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(scanner.LEFT_PAREN, "expect '(' after function identifier")
	if err != nil {
		return nil, err
	}

	params := make([]ast.Param, 0)
	if !p.check(scanner.RIGHT_PAREN) {
		for {
			paramName, err := p.consume(scanner.IDENTIFIER, "expected identifier as parameter name")
			if err != nil {
				return nil, err
			}
			isPointer := false
			if p.match(scanner.STAR) {
				isPointer = true
			}
			paramType := p.advance()
			// TODO more error checking here

			params = append(params, ast.Param{Name: paramName, Type: ast.Type{Token: paramType, IsPointer: isPointer}})
			if !p.match(scanner.COMMA) || p.isAtEnd() {
				break
			}
		}
		if p.isAtEnd() {
			return nil, fmt.Errorf("expected ')' to close function parameter list")
		}
	}

	_, err = p.consume(scanner.RIGHT_PAREN, "expected ')' to close function parameter list")
	if err != nil {
		return nil, err
	}

	isPointer := false
	if p.match(scanner.STAR) {
		isPointer = true
	}
	var returnType scanner.Token = scanner.Token{Type: scanner.TYPE, Lexeme: "void"}
	if !p.check(scanner.LEFT_BRACE) {
		returnType, err = p.consumeType("expected valid type for function return")
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(scanner.LEFT_BRACE, "expect function body after function declaration")
	if err != nil {
		return nil, err
	}
	body, err := p.blockStmt()
	if err != nil {
		return nil, err
	}
	return &ast.FnStmt{Name: name, Params: params, Body: body, Return: ast.Type{Token: returnType, IsPointer: isPointer}}, nil
}

func (p *Parser) returnStmt() (ast.Stmt, error) {
	if p.match(scanner.SEMI_COLON) {
		return &ast.ReturnStmt{Expression: nil}, nil
	}
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(scanner.SEMI_COLON, "expect ';' after return statement")
	if err != nil {
		return nil, err
	}

	return &ast.ReturnStmt{Expression: expr}, nil
}

func (p *Parser) ifStmt() (ast.Stmt, error) {
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(scanner.LEFT_BRACE, "expect block after if condition")
	if err != nil {
		return nil, err
	}
	ifBlock, err := p.blockStmt()
	if err != nil {
		return nil, err
	}
	var elseBlock ast.Stmt = nil
	if p.match(scanner.ELSE) {
		_, err = p.consume(scanner.LEFT_BRACE, "expect block after else condition")
		if err != nil {
			return nil, err
		}
		elseBlock, err = p.blockStmt()
		if err != nil {
			return nil, err
		}
	}
	return &ast.IfStmt{Condition: condition, IfBlock: ifBlock, ElseBlock: elseBlock}, nil
}

func (p *Parser) blockStmt() (ast.Stmt, error) {
	stmts := make([]ast.Stmt, 0)
	for !p.check(scanner.RIGHT_BRACE) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	_, err := p.consume(scanner.RIGHT_BRACE, "expect '}' to end block statement")
	if err != nil {
		return nil, err
	}
	return &ast.BlockStmt{Body: stmts}, nil
}

func (p *Parser) expressionStmt() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(scanner.SEMI_COLON, "Expected ';'")
	if err != nil {
		return nil, err
	}
	return &ast.ExpressionStmt{Expression: expr}, nil
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.prattParse(PREC_ASSIGNMENT)
}

func grouping(p *Parser) (ast.Expr, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(scanner.RIGHT_PAREN, "expected closing paren")
	return &ast.GroupingExpr{Expression: expr}, nil
}

func call(p *Parser, left ast.Expr) (ast.Expr, error) {
	args := make([]ast.Expr, 0)
	for !p.check(scanner.RIGHT_PAREN) && !p.isAtEnd() {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		p.match(scanner.COMMA)
		// p.consume(scanner.COMMA, "expected comma to seperate function call arguments")
		args = append(args, expr)
	}
	p.consume(scanner.RIGHT_PAREN, "expected closing paren")
	return &ast.CallExpr{Callee: left, Args: args}, nil
}

func number(p *Parser) (ast.Expr, error) {
	value := p.prev().Literal.(float64)
	return &ast.NumberExpr{
		Value: value,
	}, nil
}

func string_(p *Parser) (ast.Expr, error) {
	value := p.prev().Literal.(string)
	return &ast.StringExpr{
		Value: value,
	}, nil
}

func char(p *Parser) (ast.Expr, error) {
	value := p.prev().Literal.(int8)
	return &ast.CharExpr{
		Value: value,
	}, nil
}

func boolean(p *Parser) (ast.Expr, error) {
	value, err := strconv.ParseBool(p.prev().Lexeme)
	if err != nil {
		return nil, err
	}
	return &ast.BoolExpr{
		Value: value,
	}, nil
}

func variable(p *Parser) (ast.Expr, error) {
	token := p.prev()
	return &ast.IdentifierExpr{Value: token}, nil
}

func binary(p *Parser, left ast.Expr) (ast.Expr, error) {
	operator := p.prev()
	operatorPrecedence := p.tokenPrecedence(operator.Type)
	right, err := p.prattParse(operatorPrecedence)
	if err != nil {
		return nil, err
	}
	return &ast.BinaryExpr{
		Left:     left,
		Operator: operator,
		Right:    right,
	}, nil
}

func unary(p *Parser) (ast.Expr, error) {
	operator := p.prev()
	operatorPrecedence := p.tokenPrecedence(operator.Type)
	right, err := p.prattParse(operatorPrecedence)
	if err != nil {
		return nil, err
	}
	return &ast.UnaryExpr{
		Operator: operator,
		Right:    right,
	}, nil
}

func (p *Parser) prattParse(precedence Precedence) (ast.Expr, error) {
	token := p.advance()
	prefixFn := p.parseTable.GetRule(token.Type).PrefixRule
	if prefixFn == nil {
		return nil, p.errorAtCurrent(fmt.Sprintf("no prefix parse expression for lexeme '%s'", string(token.Lexeme)))
	}

	left, err := prefixFn(p)
	if err != nil {
		return nil, err
	}
	for precedence < p.currentTokenPrecedence() {
		token := p.advance()
		infixFn := p.parseTable.GetRule(token.Type).InfixRule
		if infixFn == nil {
			return left, nil
		}
		left, err = infixFn(p, left)
		if err != nil {
			return nil, err
		}
	}
	return left, nil
}

func (p *Parser) advance() scanner.Token {
	previous := p.peek()
	p.current++
	return previous
}
func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.prev().Type == scanner.SEMI_COLON {
			return
		}
		switch p.peek().Type {
		case scanner.LET | scanner.FN | scanner.RETURN | scanner.IF:
			return
		}
		p.advance()
	}
}
func (p *Parser) peek() scanner.Token {
	return p.tokens[p.current]
}
func (p *Parser) prev() scanner.Token {
	return p.tokens[p.current-1]
}
func (p *Parser) check(type_ scanner.TokenType) bool {
	return p.peek().Type == type_
}
func (p *Parser) match(types ...scanner.TokenType) bool {
	for _, type_ := range types {
		if p.check(type_) {
			p.advance()
			return true
		}
	}
	return false
}
func (p *Parser) consume(type_ scanner.TokenType, message string) (scanner.Token, error) {
	curToken := p.peek()
	if curToken.Type == type_ {
		return p.advance(), nil
	}
	message = fmt.Sprintf("%s - instead consumed '%s'", message, scanner.TokenKindString(&curToken))
	return scanner.Token{}, p.errorAtCurrent(message)
}

// seperate helper function because need to handle primite + user defined types
func (p *Parser) consumeType(msg string) (scanner.Token, error) {
	if p.match(scanner.TYPE) {
		// must be primitive type
		return p.prev(), nil
	} else if p.match(scanner.IDENTIFIER) {
		// must be user defined type - need to change from identifier to type
		typeToken := p.prev()
		typeToken.Type = scanner.TYPE
		return typeToken, nil
	} else {
		return scanner.Token{}, p.errorAtCurrent(msg)
	}
}
func (p *Parser) currentTokenRule() ParseRule {
	return p.parseTable.GetRule(p.peek().Type)
}
func (p *Parser) currentTokenPrecedence() Precedence {
	return p.currentTokenRule().Precedence
}
func (p *Parser) tokenPrecedence(type_ scanner.TokenType) Precedence {
	return p.parseTable.GetRule(type_).Precedence
}
func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens) || p.peek().Type == scanner.EOF
}
func (p *Parser) errorAtCurrent(message string) error {
	err := fmt.Errorf("error: near line [%d]. cause: %s", p.prev().Line, message)
	p.Errors = append(p.Errors, err)
	p.HadError = true
	return err
}
