package parser

import "fmt"

type TokenType int

const (
	TkNewLine TokenType = iota
	TkAt
	TkOpenParen
	TkCloseParen
	TkOpenBracket
	TkCloseBracket
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
	TkFunc

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
	DTImage  DataType = "image"

	DTNumberList DataType = "number[]"
	DTStringList DataType = "string[]"
)

type Position struct {
	Line   int
	Column int
}

type Token struct {
	Type   TokenType
	Lexeme string
	Pos    Position
	Indent int

	DataType DataType
	Literal  any
}

func (t Token) String() string {
	return fmt.Sprintf("([%d:%d:%d]%s)", t.Pos.Line+1, t.Pos.Column+1, t.Indent, t.Lexeme)
}
