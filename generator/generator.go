package generator

import (
	"fmt"

	"github.com/Bananenpro/embe/blocks"
	"github.com/Bananenpro/embe/parser"
)

func GenerateBlocks(statements []parser.Stmt, lines [][]rune) (map[string]*blocks.Block, []error) {
	g := &generator{
		blocks: make(map[string]*blocks.Block),
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
	blocks map[string]*blocks.Block
	parent string
	lines  [][]rune

	blockID  string
	dataType parser.DataType
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
	g.parent = block.ID
	for _, s := range stmt.Body {
		err = s.Accept(g)
		if err != nil {
			return err
		}
		g.parent = g.blockID
	}
	return nil
}

func (g *generator) VisitFuncCall(stmt *parser.StmtFuncCall) error {
	fn, ok := funcCalls[stmt.Name.Lexeme]
	if !ok {
		return g.newError("Unknown function.", stmt.Name)
	}
	block, err := fn(g, stmt, g.parent)
	if err != nil {
		return err
	}
	g.blocks[block.ID] = block
	g.blockID = block.ID
	return nil
}

func (g *generator) VisitAssignment(stmt *parser.StmtAssignment) error {
	return nil
}

func (g *generator) VisitIf(stmt *parser.StmtIf) error {
	var block *blocks.Block
	if stmt.ElseBody == nil {
		block = blocks.NewBlock(blocks.If, g.parent)
	} else {
		block = blocks.NewBlock(blocks.IfElse, g.parent)
	}
	g.blocks[g.parent].Next = &block.ID
	g.blocks[block.ID] = block
	g.parent = block.ID

	err := stmt.Condition.Accept(g)
	if err != nil {
		return err
	}
	if g.dataType != parser.DTBool {
		return g.newError("The condition must be a boolean.", stmt.Keyword)
	}
	block.Inputs["CONDITION"] = []any{2, g.blockID}

	for i, s := range stmt.Body {
		s.Accept(g)
		if i == 0 {
			block.Inputs["SUBSTACK"] = []any{2, g.blockID}
		}
		g.parent = g.blockID
	}

	for i, s := range stmt.ElseBody {
		s.Accept(g)
		if i == 0 {
			block.Inputs["SUBSTACK2"] = []any{2, g.blockID}
		}
		g.parent = g.blockID
	}

	block.Next = nil
	g.blockID = block.ID
	return nil
}

func (g *generator) VisitLoop(stmt *parser.StmtLoop) error {
	var block *blocks.Block
	var err error
	parent := g.parent
	if stmt.Condition == nil {
		block = blocks.NewBlock(blocks.RepeatForever, parent)
	} else if stmt.Keyword.Type == parser.TkWhile {
		block = blocks.NewBlock(blocks.RepeatUntil, parent)
		g.blocks[block.ID] = block
		g.parent = block.ID
		block.Inputs["CONDITION"], err = g.value(parent, stmt.Keyword, stmt.Condition, parser.DTBool)
		if err != nil {
			return err
		}
	} else if stmt.Keyword.Type == parser.TkFor {
		block = blocks.NewBlock(blocks.Repeat, parent)
		g.blocks[block.ID] = block
		g.parent = block.ID
		block.Inputs["TIMES"], err = g.value(parent, stmt.Keyword, stmt.Condition, parser.DTNumber)
		if err != nil {
			return err
		}
	} else {
		return g.newError("Unknown loop type.", stmt.Keyword)
	}
	g.blocks[parent].Next = &block.ID
	g.blocks[block.ID] = block
	g.parent = block.ID
	for i, s := range stmt.Body {
		s.Accept(g)
		if i == 0 {
			block.Inputs["SUBSTACK"] = []any{2, g.blockID}
		}
		g.parent = g.blockID
	}
	block.Next = nil

	g.blockID = block.ID
	return nil
}

func (g *generator) VisitIdentifier(expr *parser.ExprIdentifier) error {
	return nil
}

func (g *generator) VisitLiteral(expr *parser.ExprLiteral) error {
	g.blockID = "literal"
	g.dataType = expr.Token.DataType
	return nil
}

func (g *generator) VisitUnary(expr *parser.ExprUnary) error {
	var block *blocks.Block
	var dataType parser.DataType
	switch expr.Operator.Type {
	case parser.TkBang:
		dataType = parser.DTBool
		block = blocks.NewBlock(blocks.OpNot, g.parent)
	}
	g.blocks[block.ID] = block
	g.parent = block.ID
	input, err := g.value(g.parent, expr.Operator, expr.Right, dataType)
	if err != nil {
		return err
	}
	block.Inputs["OPERAND"] = input

	block.Next = nil
	g.blockID = block.ID
	return nil
}

func (g *generator) VisitBinary(expr *parser.ExprBinary) error {
	var block *blocks.Block
	var operandDataType parser.DataType
	retDataType := parser.DTBool
	operandName := "OPERAND"
	switch expr.Operator.Type {
	case parser.TkLess:
		block = blocks.NewBlock(blocks.OpLessThan, g.parent)
		operandDataType = parser.DTNumber
	case parser.TkGreater:
		block = blocks.NewBlock(blocks.OpGreaterThan, g.parent)
		operandDataType = parser.DTNumber
	case parser.TkEqual:
		block = blocks.NewBlock(blocks.OpEquals, g.parent)
		operandDataType = parser.DTNumber
	case parser.TkAnd:
		block = blocks.NewBlock(blocks.OpAnd, g.parent)
		operandDataType = parser.DTBool
	case parser.TkOr:
		block = blocks.NewBlock(blocks.OpOr, g.parent)
		operandDataType = parser.DTBool
	default:
		retDataType = parser.DTNumber
		operandDataType = parser.DTNumber
		operandName = "NUM"
		switch expr.Operator.Type {
		case parser.TkPlus:
			block = blocks.NewBlock(blocks.OpAdd, g.parent)
		case parser.TkMinus:
			block = blocks.NewBlock(blocks.OpSubtract, g.parent)
		case parser.TkMultiply:
			block = blocks.NewBlock(blocks.OpMultiply, g.parent)
		case parser.TkDivide:
			block = blocks.NewBlock(blocks.OpDivide, g.parent)
		default:
			return g.newError("Unknown binary operator.", expr.Operator)
		}
	}
	g.blocks[block.ID] = block

	left, err := g.value(block.ID, expr.Operator, expr.Left, operandDataType)
	if err != nil {
		return err
	}
	block.Inputs[operandName+"1"] = left

	right, err := g.value(block.ID, expr.Operator, expr.Right, operandDataType)
	if err != nil {
		return err
	}
	block.Inputs[operandName+"2"] = right

	block.Next = nil
	g.dataType = retDataType
	g.blockID = block.ID
	return nil
}

func (g *generator) value(parent string, token parser.Token, expr parser.Expr, dataType parser.DataType) ([]any, error) {
	if literal, ok := expr.(*parser.ExprLiteral); ok {
		if literal.Token.DataType != dataType {
			return nil, g.newError(fmt.Sprintf("The value must be of type %s.", dataType), literal.Token)
		}
		return []any{1, []any{4, fmt.Sprintf("%v", literal.Token.Literal)}}, nil
	} else {
		g.parent = parent
		err := expr.Accept(g)
		if err != nil {
			return nil, err
		}
		if g.dataType != dataType {
			return nil, g.newError(fmt.Sprintf("The value must be of type %s.", dataType), token)
		}
		if dataType == parser.DTBool {
			return []any{2, g.blockID}, nil
		}
		return []any{3, g.blockID, []any{4, ""}}, nil
	}
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
