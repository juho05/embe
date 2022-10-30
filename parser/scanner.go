package parser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var keywords = map[string]TokenType{
	"if":    TkIf,
	"elif":  TkElif,
	"else":  TkElse,
	"while": TkWhile,
	"for":   TkFor,
	"var":   TkVar,
	"const": TkConst,
	"func":  TkFunc,
	"event": TkEvent,
}

var types = map[string]DataType{
	"boolean": DTBool,
	"number":  DTNumber,
	"string":  DTString,
	"image":   DTImage,
}

type scanner struct {
	inputScanner      *bufio.Scanner
	lines             [][]rune
	line              int
	tokenStartColumn  int
	currentColumn     int
	tokens            []Token
	lineContainsToken bool
	errors            []error
	path              string
}

func Scan(source io.Reader, path string) ([]Token, [][]rune, []error) {
	fileScanner := bufio.NewScanner(source)

	srcScanner := &scanner{
		inputScanner: fileScanner,
		line:         -1,
		errors:       make([]error, 0),
		path:         path,
	}

	srcScanner.scan()
	return srcScanner.tokens, srcScanner.lines, srcScanner.errors
}

func (s *scanner) scan() {
	c, err := s.nextCharacter()
	if err != nil {
		s.errors = append(s.errors, err)
		return
	}

	for c != '\000' {
		switch c {
		case '@':
			s.addToken(TkAt)
		case '(':
			s.addToken(TkOpenParen)
		case ')':
			s.addToken(TkCloseParen)
		case '[':
			s.addToken(TkOpenBracket)
		case ']':
			s.addToken(TkCloseBracket)
		case ':':
			s.addToken(TkColon)
		case '.':
			s.addToken(TkDot)
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
				s.blockComment()
			} else if s.match('=') {
				s.addToken(TkDivideAssign)
			} else {
				s.addToken(TkDivide)
			}
		case '%':
			if s.match('=') {
				s.addToken(TkModulusAssign)
			} else {
				s.addToken(TkModulus)
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
			} else {
				s.addToken(TkLess)
			}
		case '>':
			if s.match('=') {
				s.addToken(TkGreaterEqual)
			} else {
				s.addToken(TkGreater)
			}

		case '!':
			if s.match('=') {
				s.addToken(TkNotEqual)
			} else {
				s.addToken(TkBang)
			}

		case '|':
			if !s.match('|') {
				s.errors = append(s.errors, s.newError(fmt.Sprintf("Unexpected character '%c'.", c)))
				s.synchronize(true)
				break
			}
			s.addToken(TkOr)
		case '&':
			if !s.match('&') {
				s.errors = append(s.errors, s.newError(fmt.Sprintf("Unexpected character '%c'.", c)))
				s.synchronize(true)
				break
			}
			s.addToken(TkAnd)

		case '"':
			s.string()

		case ' ', '\t':

		case '#':
			s.preprocessor()

		default:
			if isDigit(c, 10) {
				s.number()
			} else if isAlpha(c) {
				s.identifier()
			} else {
				s.errors = append(s.errors, s.newError(fmt.Sprintf("Unexpected character '%c'.", c)))
				s.synchronize(true)
			}
		}

		c, err = s.nextCharacter()
		if err != nil {
			s.errors = append(s.errors, err)
			return
		}
		s.tokenStartColumn = s.currentColumn
	}

	if len(s.tokens) > 0 {
		if s.tokens[len(s.tokens)-1].Type != TkNewLine {
			s.addToken(TkNewLine)
		}
	}

	eof := Token{
		Type: TkEOF,
		Pos: Position{
			Path: s.path,
		},
	}
	if s.line >= 0 && s.line < len(s.lines) {
		eof.Pos.Column = len(s.lines[s.line])
	}

	s.tokens = append(s.tokens, eof)
}

func (s *scanner) identifier() {
	for {
		for isAlphaNum(s.peek()) {
			s.nextCharacter()
		}
		if s.peek() != '.' || !isAlphaNum(s.peekNext()) {
			break
		}
		s.nextCharacter()
	}

	name := string(s.lines[s.line][s.tokenStartColumn : s.currentColumn+1])
	if t, ok := types[name]; ok {
		if name != "boolean" && s.peek() == '[' && s.peekNext() == ']' {
			s.nextCharacter()
			s.nextCharacter()
			t += "[]"
		}
		s.addTokenWithValue(TkType, t, nil)
	} else if k, ok := keywords[name]; ok {
		s.addToken(k)
	} else if name == "true" || name == "false" {
		v, err := strconv.ParseBool(name)
		if err != nil {
			panic(err)
		}
		s.addTokenWithValue(TkLiteral, DTBool, v)
	} else {
		s.addToken(TkIdentifier)
	}
}

func (s *scanner) number() {
	base := 10

	if string(s.lines[s.line][s.currentColumn:s.currentColumn+1]) == "0" {
		switch s.peek() {
		case 'x':
			base = 16
			s.nextCharacter()
		case 'o':
			base = 8
			s.nextCharacter()
		case 'b':
			base = 2
			s.nextCharacter()
		}
	}

	for isDigit(s.peek(), base) {
		s.nextCharacter()
	}

	if base == 10 && s.peek() == '.' && isDigit(s.peekNext(), base) {
		s.nextCharacter()
		for isDigit(s.peek(), base) {
			s.nextCharacter()
		}
		value, _ := strconv.ParseFloat(string(s.lines[s.line][s.tokenStartColumn:s.currentColumn+1]), 64)
		s.addTokenWithValue(TkLiteral, DTNumber, value)
		return
	}

	lexeme := string(s.lines[s.line][s.tokenStartColumn : s.currentColumn+1])
	switch base {
	case 16:
		lexeme = strings.TrimPrefix(lexeme, "0x")
	case 8:
		lexeme = strings.TrimPrefix(lexeme, "0o")
	case 2:
		lexeme = strings.TrimPrefix(lexeme, "0b")
	}
	if lexeme == "" {
		s.errors = append(s.errors, s.newError("There must be at least one digit after a number prefix."))
		s.synchronize(true)
		return
	}
	value, _ := strconv.ParseInt(lexeme, base, 64)
	s.addTokenWithValue(TkLiteral, DTNumber, float64(value))
}

func (s *scanner) comment() {
	for s.peek() != '\n' {
		s.nextCharacter()
	}
}

func (s *scanner) string() {
	characters := make([]rune, 0)
	for s.peek() != '"' && s.peek() != '\n' {
		c, _ := s.nextCharacter()
		characters = append(characters, c)
	}
	if !s.match('"') {
		s.errors = append(s.errors, s.newError("Unterminated string."))
		s.synchronize(false)
		return
	}
	s.addTokenWithValue(TkLiteral, DTString, string(characters))
}

func (s *scanner) preprocessor() {
	for isAlphaNum(s.peek()) {
		s.nextCharacter()
	}
	s.addToken(TkPreprocessor)
}

func (s *scanner) blockComment() {
	nestingLevel := 1
	for nestingLevel > 0 {
		c, err := s.nextCharacter()

		if c == '\000' || err != nil {
			return
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
}

func (s *scanner) nextCharacter() (rune, error) {
	if s.line != -1 && s.peek() == '\n' && s.lineContainsToken {
		s.addToken(TkNewLine)
		s.lineContainsToken = false
	}
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
	if s.currentColumn+1 >= len(s.lines[s.line]) {
		return '\n'
	}

	return s.lines[s.line][s.currentColumn+1]
}

func (s *scanner) peekNext() rune {
	if s.currentColumn+2 >= len(s.lines[s.line]) {
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
	s.lines = append(s.lines, []rune(strings.ReplaceAll(s.inputScanner.Text(), "\t", "    ")))
	s.line++
	s.currentColumn = 0
	s.tokenStartColumn = 0

	return true, nil
}

func (s *scanner) synchronize(spaceAsSeparator bool) {
	for s.peek() != '\n' && s.peek() != '\000' {
		if spaceAsSeparator && s.peek() == ' ' {
			return
		}
		s.nextCharacter()
	}
}

func (s *scanner) addToken(tokenType TokenType) {
	lexeme := string(s.lines[s.line][s.tokenStartColumn : s.currentColumn+1])
	if tokenType == TkNewLine {
		prev := s.tokens[len(s.tokens)-1]
		s.tokenStartColumn = prev.Pos.Column + len(prev.Lexeme)
		lexeme = "\\n"
	}
	s.tokens = append(s.tokens, Token{
		Pos: Position{
			Line:   s.line,
			Column: s.tokenStartColumn,
			Path:   s.path,
		},
		EndPos: Position{
			Line:   s.line,
			Column: s.tokenStartColumn + len(lexeme) - 1,
			Path:   s.path,
		},
		LineAfterInclude: s.line,
		Type:             tokenType,
		Lexeme:           lexeme,
		Indent:           getIndentation(s.lines[s.line]),
	})
	if tokenType != TkNewLine && tokenType != TkEOF {
		s.lineContainsToken = true
	}
}

func (s *scanner) addTokenWithValue(tokenType TokenType, dataType DataType, value any) {
	lexeme := string(s.lines[s.line][s.tokenStartColumn : s.currentColumn+1])
	s.tokens = append(s.tokens, Token{
		Pos: Position{
			Line:   s.line,
			Column: s.tokenStartColumn,
			Path:   s.path,
		},
		EndPos: Position{
			Line:   s.line,
			Column: s.tokenStartColumn + len(lexeme) - 1,
			Path:   s.path,
		},
		LineAfterInclude: s.line,
		Type:             tokenType,
		Lexeme:           lexeme,
		DataType:         dataType,
		Literal:          value,
		Indent:           getIndentation(s.lines[s.line]),
	})
	if tokenType != TkNewLine && tokenType != TkEOF {
		s.lineContainsToken = true
	}
}

func getIndentation(line []rune) int {
	level := 0
	for ; level < len(line) && line[level] == ' '; level++ {
	}
	return level
}

func isDigit(r rune, base int) bool {
	switch base {
	case 16:
		return r >= '0' && r <= '9' || r >= 'a' && r <= 'f' || r >= 'A' && r <= 'F'
	case 10:
		return r >= '0' && r <= '9'
	case 8:
		return r >= '0' && r <= '7'
	case 2:
		return r == '0' || r == '1'
	default:
		return false
	}
}

func isAlpha(char rune) bool {
	return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char == '_'
}

func isAlphaNum(char rune) bool {
	return isDigit(char, 10) || isAlpha(char)
}

type ScanError struct {
	Pos     Position
	Message string
}

func (s ScanError) Error() string {
	return "ERROR: " + s.Message
}

func (s *scanner) newError(msg string) error {
	return ScanError{
		Pos: Position{
			Line:   s.line,
			Column: s.currentColumn,
		},
		Message: msg,
	}
}
