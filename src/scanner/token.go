package scanner

import "fmt"

type LiteralValue interface{}

type TokenType int

const (
	// special tokens
	EOF TokenType = iota

	literal_beg
	NUMBER
	STRING // maybe can't be primitive type but must be struct with methods
	CHAR
	BOOL
	IDENTIFIER
	literal_end

	operator_beg
	LEFT_PAREN  // (
	RIGHT_PAREN // )
	LEFT_BRACE  // {
	RIGHT_BRACE // }
	LEFT_BRACK  // [
	RIGHT_BRACK // ]
	SEMI_COLON  // ;
	COMMA       // ,
	PLUS        // +
	MINUS       // -
	STAR        // *
	SLASH       // /
	BANG        // !
	ADDRESS     // &
	ASSIGN      // =
	DOT         // .
	PLUSPLUS    // ++
	MINUSMINUS  // --
	EQUAL       // ==
	NOT_EQUAL   // !=
	LESS        // <
	GREATER     // >
	LESS_EQ     // <=
	GREATER_EQ  // >=
	operator_end

	keyword_beg
	TRUE
	FALSE
	LET
	FN
	RETURN
	IF
    ELIF
	ELSE
	// STRING_TYPE
	// NUMBER_TYPE
	// BOOL_TYPE
	NIL
	keyword_end

	placeholders_beg
	TYPE
	// VOID
	placeholders_end

	// TODO: implement these more difficult concepts
	// DOTDOT                       // ..
	// DOTDOTDOT                    // ...
	// STRUCT
	// ENUM
	// PUB
	// FOR
	// WHILE
	// BREAK
	// CONTINUE
	// IMPORT
	// PRINT
)

type Token struct {
	Type    TokenType
	Literal LiteralValue // might not need this one
	Lexeme  string       // actual string from source code
	Line    int          // line the token appears in
}

func NewToken(type_ TokenType, literal LiteralValue, lexeme string, line int) *Token {
	return &Token{type_, literal, lexeme, line}
}

func (t *Token) DebugPrint() {
	if t.isOneOf(IDENTIFIER, STRING, NUMBER, BOOL, TYPE) {
		fmt.Printf("%s (%s)\n", TokenKindString(t), t.Lexeme)
	} else {
		fmt.Printf("%s\n", TokenKindString(t))
	}
}

func (t *Token) isOneOf(tokens ...TokenType) bool {
	for _, tokentype := range tokens {
		if t.Type == tokentype {
			return true
		}
	}
	return false
}

var typeKeywords = map[string]bool{
	"number": true,
	"string": true,
	"bool":   true,
	"char":   true,
	// placeholders - created parser; should not be used from code
	"void": true,
	"auto": true,
}

func (t *Token) IsPrimitiveType() bool {
	return typeKeywords[t.Lexeme]
}

func TokenKindString(token *Token) string {
	switch token.Type {
	case NUMBER:
		return "number"
	case STRING:
		return "string"
	case BOOL:
		return "bool"
	case CHAR:
		return "char"
	case LEFT_PAREN:
		return "left_paren"
	case RIGHT_PAREN:
		return "right_paren"
	case LEFT_BRACE:
		return "left_brace"
	case RIGHT_BRACE:
		return "right_brace"
	case LEFT_BRACK:
		return "left_brack"
	case RIGHT_BRACK:
		return "right_brack"
	case SEMI_COLON:
		return "semi_colon"
	case COMMA:
		return "comma"
	case PLUS:
		return "plus"
	case MINUS:
		return "minus"
	case STAR:
		return "star"
	case SLASH:
		return "slash"
	case BANG:
		return "bang"
	case ADDRESS:
		return "address"
	case ASSIGN:
		return "assign"
	case DOT:
		return "dot"
	case PLUSPLUS:
		return "plusplus"
	case MINUSMINUS:
		return "minusminus"
	case EQUAL:
		return "equal"
	case NOT_EQUAL:
		return "not_equal"
	case LESS:
		return "less"
	case GREATER:
		return "greater"
	case LESS_EQ:
		return "less_eq"
	case GREATER_EQ:
		return "greater_eq"
	case IDENTIFIER:
		return "identifier"
	case LET:
		return "let"
	case TRUE:
		return "true"
	case FALSE:
		return "false"
	case FN:
		return "fn"
	case IF:
		return "if"
	case ELIF:
		return "elif"
	case ELSE:
		return "else"
	case RETURN:
		return "return"
	case EOF:
		return "eof"
	case TYPE:
		return "type"
	case NIL:
		return "nil"
		// case DOTDOT:
		// 	return "dotdot"
		// case DOTDOTDOT:
		// 	return "dotdotdot"
		// case STRUCT:
		// 	return "struct"
		// case ENUM:
		// 	return "enum"
		// case ELIF:
		// 	return "elif"
		// case PUB:
		// 	return "pub"
		// case FOR:
		// 	return "for"
		// case WHILE:
		// 	return "while"
		// case BREAK:
		// 	return "break"
		// case CONTINUE:
		// 	return "continue"
		// case IMPORT:
		// 	return "import"
		// case PRINT:
		// 	return "print"
		// default:
	}
	panic(fmt.Sprintf("no case to handle TokenType with enum value '%d', and lexeme '%s'", token.Type, token.Lexeme))
}
