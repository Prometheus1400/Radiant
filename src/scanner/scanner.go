package scanner

import (
	"fmt"
	"strconv"
)

type Scanner struct {
	source   []byte
	start    int
	current  int
	line     int
	Tokens   []Token
	keywords map[string]TokenType
	Errors   []error
	HadError bool
}

func (s *Scanner) Init() {
	s.source = nil
	s.start = 0
	s.current = 0
	s.line = 1
	s.Tokens = make([]Token, 0, 16)
	s.keywords = getKeywords()
	s.Errors = nil
	s.HadError = false
}

func NewScanner() *Scanner {
	scanner := &Scanner{}
	scanner.Init()
	return scanner
}

func (s *Scanner) Scan(source []byte) []Token {
	s.Init()
	s.source = source
	s.scanTokens()
	return s.Tokens
}

func (s *Scanner) PrintTokens() {
	for _, token := range s.Tokens {
		token.DebugPrint()
	}
}

func (s *Scanner) ReportErrors() {
	for _, err := range s.Errors {
		fmt.Println(err)
	}
}

func getKeywords() map[string]TokenType {
	return map[string]TokenType{
		"let":    LET,
		"true":   TRUE,
		"false":  FALSE,
		"fn":     FN,
		"if":     IF,
		"else":   ELSE,
		"return": RETURN,
		"number": TYPE,
		"string": TYPE,
		"bool":   TYPE,
		"char":   TYPE,
		// "nil":    NIL,
		// "elif":     ELIF,
		// "pub":      PUB,
		// "for":      FOR,
		// "while":    WHILE,
		// "struct":   STRUCT,
		// "enum":     ENUM,
		// "break":    BREAK,
		// "continue": CONTINUE,
		// "import":   IMPORT,
		// "print":    PRINT,
	}
}

func (s *Scanner) scanTokens() {
	for !s.isAtEnd() {
		s.start = s.current
		s.advance()
		c := s.prev()

		switch c {
		// ignore whitespace
		case '\t', '\r', ' ':
			continue
		case '\n':
			s.line++
			continue
		case '(':
			s.addToken(LEFT_PAREN)
		case ')':
			s.addToken(RIGHT_PAREN)
		case '{':
			s.addToken(LEFT_BRACE)
		case '}':
			s.addToken(RIGHT_BRACE)
		case '[':
			s.addToken(LEFT_BRACK)
		case ']':
			s.addToken(RIGHT_BRACK)
		case ',':
			s.addToken(COMMA)
		case '*':
			s.addToken(STAR)
		case ';':
			s.addToken(SEMI_COLON)
		case '&':
			s.addToken(ADDRESS)
		case '<':
			if s.peek() == '=' {
				s.advance()
				s.addToken(LESS_EQ)
			} else {
				s.addToken(LESS)
			}
		case '>':
			if s.peek() == '=' {
				s.advance()
				s.addToken(GREATER_EQ)
			} else {
				s.addToken(GREATER)
			}
		case '.':
			// if s.peek() == '.' && s.peekNext() == '.' {
			// 	s.advance()
			// 	s.advance()
			// 	s.addToken(DOTDOTDOT)
			// } else if s.peek() == '.' {
			// 	s.advance()
			// 	s.addToken(DOTDOT)
			// } else {
			s.addToken(DOT)
			// }
		case '+':
			if s.peek() == '+' {
				s.advance()
				s.addToken(PLUSPLUS)
			} else {
				s.addToken(PLUS)
			}
		case '-':
			if s.peek() == '-' {
				s.advance()
				s.addToken(MINUSMINUS)
			} else {
				s.addToken(MINUS)
			}
		case '!':
			if s.peek() == '=' {
				s.advance()
				s.addToken(NOT_EQUAL)
			} else {
				s.addToken(BANG)
			}
		case '=':
			if s.peek() == '=' {
				s.advance()
				s.addToken(EQUAL)
			} else {
				s.addToken(ASSIGN)
			}
		case '/':
			// matching for comments
			if s.match('/') {
				// consume till end of line - not including '\n'
				for s.peek() != '\n' {
					s.advance()
				}
			} else {
				s.addToken(SLASH)
			}
		case '"':
			s.string()
		case '\'':
			s.char()
		default:
			if isDigit(c) {
				s.number()
			} else if isAlpha(c) {
				s.identifier()
			} else {
				s.errorAtCurrent(fmt.Sprintf("unrecognized character '%c'", c))
			}
		}
	}
	s.addToken(EOF)
}

func (s *Scanner) addToken(type_ TokenType) {
	s.addTokenWithLiteral(type_, nil)
}

func (s *Scanner) addTokenWithLiteral(type_ TokenType, literal LiteralValue) {
	lexeme := string(s.source[s.start:s.current])
	token := NewToken(type_, literal, lexeme, s.line)
	s.Tokens = append(s.Tokens, *token)
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.errorAtCurrent("multiline strings are not supported")
			s.advance()
			break
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.errorAtCurrent("unterminated string")
	}
	s.advance()
	valueStr := string(s.source[s.start+1 : s.current-1])
	s.addTokenWithLiteral(STRING, valueStr)
}

func (s *Scanner) char() {
	for s.peek() != '\'' && !s.isAtEnd() {
		s.advance()
	}
	if s.isAtEnd() {
		s.errorAtCurrent("unterminated character")
	}
	s.advance()
	if s.current-s.start > 3 {
		s.errorAtCurrent("can only specify 1 character inside of single quotes")
	}

	valueChar := int8(s.source[s.start+1])
	s.addTokenWithLiteral(CHAR, valueChar)
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	valueStr := s.source[s.start:s.current]
	value, _ := strconv.ParseFloat(string(valueStr), 64)
	s.addTokenWithLiteral(NUMBER, value)
}

func (s *Scanner) identifier() {
	for isAlpha(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	type_, ok := s.keywords[string(text)]
	if !ok {
		s.addToken(IDENTIFIER)
		return
	}
	s.addToken(type_)
}

func (s *Scanner) advance() {
	s.current++
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func (s *Scanner) prev() byte {
	return s.source[s.current-1]
}

func (s *Scanner) match(c byte) bool {
	if s.source[s.current] == c {
		s.advance()
		return true
	}
	return false
}
func (s *Scanner) errorAtCurrent(message string) error {
	err := fmt.Errorf("error: near line [%d]. cause: %s", s.line, message)
	s.Errors = append(s.Errors, err)
	return err
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func isDigit(c byte) bool {
	if '0' <= c && c <= '9' {
		return true
	}
	return false
}

func isAlpha(c byte) bool {
	if c == '_' {
		return true
	}
	if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') {
		return true
	}
	return false
}
