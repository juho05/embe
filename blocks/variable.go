package blocks

import "github.com/Bananenpro/embe/parser"

type Variable struct {
	ID       string
	Name     parser.Token
	DataType parser.DataType
}
