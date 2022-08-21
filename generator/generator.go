package generator

import (
	"github.com/Bananenpro/embe/blocks"
	"github.com/Bananenpro/embe/parser"
)

func GenerateBlocks(statements []parser.Stmt, lines [][]rune) (map[string]blocks.Block, []error) {
	g := &generator{
		blocks: make(map[string]blocks.Block),
		lines:  lines,
	}
	errs := make([]error, 0)
	for _, stmt := range statements {
		err := stmt.Accept(g)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return g.blocks, errs
}

type generator struct {
	blocks map[string]blocks.Block
	lines  [][]rune
}

func (g *generator) VisitEvent(stmt *parser.StmtEvent) error {
	fn, ok := events[stmt.Name.Lexeme]
	if !ok {
		return g.newError("Unknown event.", stmt.Name)
	}
	block, err := fn(g, stmt)
	if err != nil {
		return err
	}
	g.blocks[block.ID] = block
	return nil
}

func (g *generator) VisitFuncCall(stmt *parser.StmtFuncCall) error {
	return nil
}

func (g *generator) VisitAssignment(stmt *parser.StmtAssignment) error {
	return nil
}

func (g *generator) VisitIf(stmt *parser.StmtIf) error {
	return nil
}

func (g *generator) VisitLoop(stmt *parser.StmtLoop) error {
	return nil
}

func (g *generator) VisitIdentifier(expr *parser.ExprIdentifier) error {
	return nil
}

func (g *generator) VisitLiteral(expr *parser.ExprLiteral) error {
	return nil
}

func (g *generator) VisitUnary(expr *parser.ExprUnary) error {
	return nil
}

func (g *generator) VisitBinary(expr *parser.ExprBinary) error {
	return nil
}

type GenerateError struct {
	Token   parser.Token
	Message string
	Line    []rune
}

func (p GenerateError) Error() string {
	length := len([]rune(p.Token.Lexeme))
	if p.Token.Type == parser.TkNewLine {
		length = 1
	}
	return generateErrorText(p.Message, p.Line, p.Token.Line, p.Token.Column, p.Token.Column+length)
}

func (g *generator) newError(message string, token parser.Token) error {
	return GenerateError{
		Token:   token,
		Message: message,
		Line:    g.lines[token.Line],
	}
}
