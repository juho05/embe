package parser

import "fmt"

type TokenType int

const (
	TkNewLine TokenType = iota
	TkAt
	TkOpenParen
	TkCloseParen
	TkColon
	TkDot
	TkComma

	TkBang
	TkOr
	TkAnd

	TkPlus
	TkPlusAssign
	TkMinus
	TkMinusAssign
	TkMultiply
	TkMultiplyAssign
	TkDivide
	TkDivideAssign
	TkModulus
	TkModulusAssign

	TkAssign
	TkEqual
	TkNotEqual
	TkLess
	TkGreater
	TkLessEqual
	TkGreaterEqual

	TkIf
	TkElif
	TkElse
	TkWhile
	TkFor
	TkVar
	TkConst

	TkIdentifier
	TkLiteral
	TkType

	TkEOF
)

type DataType string

const (
	DTNumber DataType = "number"
	DTBool   DataType = "boolean"
	DTString DataType = "string"
)

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
	Column int
	Indent int

	DataType DataType
	Literal  any
}

func (t Token) String() string {
	return fmt.Sprintf("([%d:%d:%d]%s)", t.Line+1, t.Column+1, t.Indent, t.Lexeme)
}
