package parser

type TokenType int

const (
	TkNewLine TokenType = iota
	TkAt
	TkOpenParen
	TkCloseParen
	TkColon
	TkDot
	TkSemicolon
	TkComma
	TkBang

	TkPlus
	TkMinus
	TkMultiply
	TkDivide
	TkPlusAssign
	TkMinusAssign
	TkMultiplyAssign
	TkDivideAssign

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

	TkIdentifier
	TkLiteral

	TkEOF
)

type DataType int

const (
	DTInt DataType = iota
	DTBool
	DTString
)

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
	Column int

	DataType DataType
	Literal  any
}
