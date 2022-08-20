package parser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

var keywords = map[string]TokenType{
	"if":    TkIf,
	"elif":  TkElif,
	"else":  TkElse,
	"while": TkWhile,
}

type scanner struct {
	inputScanner     *bufio.Scanner
	lines            [][]rune
	line             int
	tokenStartColumn int
	currentColumn    int
	tokens           []Token
}

func Scan(source io.Reader) ([]Token, [][]rune, error) {
	fileScanner := bufio.NewScanner(source)

	srcScanner := &scanner{
		inputScanner: fileScanner,
		line:         -1,
	}

	err := srcScanner.scan()

	return srcScanner.tokens, srcScanner.lines, err
}

func (s *scanner) scan() error {
	c, err := s.nextCharacter()
	if err != nil {
		return err
	}

	for c != '\000' {
		switch c {
		case '\n':
			s.addToken(TkNewLine)
		case '@':
			s.addToken(TkAt)
		case '(':
			s.addToken(TkOpenParen)
		case ')':
			s.addToken(TkCloseParen)
		case ':':
			s.addToken(TkColon)
		case '.':
			s.addToken(TkDot)
		case ';':
			s.addToken(TkSemicolon)
		case ',':
			s.addToken(TkComma)
		case '+':
			if s.match('=') {
				s.addToken(TkPlusAssign)
			} else {
				s.addToken(TkPlus)
			}
		case '-':
			if s.match('=') {
				s.addToken(TkMinusAssign)
			} else {
				s.addToken(TkMinus)
			}
		case '*':
			if s.match('=') {
				s.addToken(TkMultiplyAssign)
			} else {
				s.addToken(TkMultiply)
			}
		case '/':
			if s.match('/') {
				s.comment()
			} else if s.match('*') {
				err = s.blockComment()
				if err != nil {
					return err
				}
			} else if s.match('=') {
				s.addToken(TkDivideAssign)
			} else {
				s.addToken(TkDivide)
			}
		case '=':
			if s.match('=') {
				s.addToken(TkEqual)
			} else {
				s.addToken(TkAssign)
			}
		case '<':
			if s.match('=') {
				s.addToken(TkLessEqual)
			}
			s.addToken(TkLess)
		case '>':
			if s.match('=') {
				s.addToken(TkGreaterEqual)
			}
			s.addToken(TkGreater)

		case '!':
			if s.match('=') {
				s.addToken(TkNotEqual)
			} else {
				s.addToken(TkBang)
			}

		case '"':
			err = s.string()
			if err != nil {
				return err
			}

		case ' ', '\t':

		default:
			if isDigit(c) {
				s.number()
			} else if isAlpha(c) {
				s.identifier()
			} else {
				return s.newError(fmt.Sprintf("Unexpected character '%c'.", c))
			}
		}

		c, err = s.nextCharacter()
		if err != nil {
			return err
		}
		s.tokenStartColumn = s.currentColumn
	}

	eof := Token{
		Line: s.line,
		Type: TkEOF,
	}
	if s.line >= 0 && s.line < len(s.lines) {
		eof.Column = len(s.lines[s.line])
	}

	s.tokens = append(s.tokens, eof)

	return nil
}

func (s *scanner) identifier() {
	for isAlphaNum(s.peek()) {
		s.nextCharacter()
	}

	name := string(s.lines[s.line][s.tokenStartColumn : s.currentColumn+1])
	if k, ok := keywords[name]; ok {
		s.addToken(k)
	} else {
		s.addToken(TkIdentifier)
	}
}

func (s *scanner) number() {
	for isDigit(s.peek()) {
		s.nextCharacter()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.nextCharacter()
		for isDigit(s.peek()) {
			s.nextCharacter()
		}
	}

	value, _ := strconv.ParseFloat(string(s.lines[s.line][s.tokenStartColumn:s.currentColumn+1]), 64)
	s.addTokenWithValue(TkLiteral, DTInt, value)
}

func (s *scanner) comment() {
	for s.peek() != '\n' {
		s.nextCharacter()
	}
}

func (s *scanner) string() error {
	characters := make([]rune, 0)
	for s.peek() != '"' && s.peek() != '\n' {
		c, _ := s.nextCharacter()
		characters = append(characters, c)
	}
	if !s.match('"') {
		return s.newError("Unterminated string.")
	}
	s.addTokenWithValue(TkLiteral, DTString, string(characters))
	return nil
}

func (s *scanner) blockComment() error {
	nestingLevel := 1
	for nestingLevel > 0 {
		c, err := s.nextCharacter()

		if c == '\000' || err != nil {
			return err
		}
		if c == '/' && s.match('*') {
			nestingLevel++
			continue
		}
		if c == '*' && s.match('/') {
			nestingLevel--
			continue
		}
	}
	return nil
}

func (s *scanner) nextCharacter() (rune, error) {
	s.currentColumn++
	for s.line == -1 || s.currentColumn >= len(s.lines[s.line]) {
		notDone, err := s.nextLine()
		if !notDone {
			return '\000', err
		}
	}

	return s.lines[s.line][s.currentColumn], nil
}

func (s *scanner) peek() rune {
	if s.currentColumn+1 == len(s.lines[s.line]) {
		return '\n'
	}

	return s.lines[s.line][s.currentColumn+1]
}

func (s *scanner) peekNext() rune {
	if s.currentColumn+2 == len(s.lines[s.line]) {
		return '\n'
	}

	return s.lines[s.line][s.currentColumn+2]
}

func (s *scanner) match(char rune) bool {
	if s.peek() != char {
		return false
	}
	s.nextCharacter()
	return true
}

func (s *scanner) nextLine() (bool, error) {
	if !s.inputScanner.Scan() {
		return false, s.inputScanner.Err()
	}
	s.lines = append(s.lines, []rune(s.inputScanner.Text()))
	s.line++
	s.currentColumn = 0
	s.tokenStartColumn = 0

	return true, nil
}

func (s *scanner) addToken(tokenType TokenType) {
	s.tokens = append(s.tokens, Token{
		Line:   s.line,
		Column: s.tokenStartColumn,
		Type:   tokenType,
		Lexeme: string(s.lines[s.line][s.tokenStartColumn : s.currentColumn+1]),
	})
}

func (s *scanner) addTokenWithValue(tokenType TokenType, dataType DataType, value any) {
	s.tokens = append(s.tokens, Token{
		Line:     s.line,
		Column:   s.tokenStartColumn,
		Type:     tokenType,
		Lexeme:   string(s.lines[s.line][s.tokenStartColumn : s.currentColumn+1]),
		DataType: dataType,
		Literal:  value,
	})
}

func isDigit(char rune) bool {
	return char >= '0' && char <= '9'
}

func isAlpha(char rune) bool {
	return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char == '_'
}

func isAlphaNum(char rune) bool {
	return isDigit(char) || isAlpha(char)
}

type ScanError struct {
	Line     int
	LineText []rune
	Column   int
	Message  string
}

func (s ScanError) Error() string {
	return generateErrorText(s.Message, s.LineText, s.Line, s.Column, s.Column+1)
}

func (s *scanner) newError(msg string) error {
	return ScanError{
		Line:     s.line,
		LineText: s.lines[s.line],
		Column:   s.currentColumn,
		Message:  msg,
	}
}
