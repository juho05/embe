package generator

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/Bananenpro/embe/blocks"
	"github.com/Bananenpro/embe/parser"
)

func GenerateBlocks(statements []parser.Stmt, lines [][]rune) (map[string]*blocks.Block, map[string]*blocks.Variable, []error) {
	g := &generator{
		blocks:    make(map[string]*blocks.Block),
		variables: make(map[string]*blocks.Variable),
		lines:     lines,
	}
	errs := make([]error, 0)
	for _, stmt := range statements {
		err := stmt.Accept(g)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return g.blocks, g.variables, errs
}

type generator struct {
	blocks map[string]*blocks.Block
	parent string
	lines  [][]rune

	blockID  string
	dataType parser.DataType

	variableInitializer *blocks.Block
	variables           map[string]*blocks.Variable

	noNext       bool
	variableName string
}

func (g *generator) VisitVarDecl(stmt *parser.StmtVarDecl) error {
	if v, ok := g.variables[stmt.Name.Lexeme]; ok {
		return g.newError(fmt.Sprintf("'%s' is already declared in line %d.", stmt.Name.Lexeme, v.Name.Line+1), stmt.Name)
	}

	if g.variableInitializer == nil {
		fn := Events["start"]
		block, _ := fn(g, &parser.StmtEvent{})
		g.variableInitializer = block
		g.blocks[block.ID] = block
	}

	variable := &blocks.Variable{
		ID:       uuid.NewString(),
		Name:     stmt.Name,
		DataType: stmt.DataType,
	}

	if variable.DataType != "" && stmt.Value == nil {
		stmt.AssignToken = parser.Token{
			Type: parser.TkAssign,
			Line: stmt.Name.Line,
		}
		switch variable.DataType {
		case parser.DTNumber:
			stmt.Value = &parser.ExprLiteral{
				Token: parser.Token{
					Type:     parser.TkLiteral,
					Lexeme:   "0",
					Literal:  0,
					Line:     stmt.Name.Line,
					DataType: parser.DTNumber,
				},
			}
		case parser.DTString:
			stmt.Value = &parser.ExprLiteral{
				Token: parser.Token{
					Type:     parser.TkLiteral,
					Lexeme:   "",
					Literal:  "",
					Line:     stmt.Name.Line,
					DataType: parser.DTString,
				},
			}
		default:
			return g.newError("Unknown type.", stmt.Name)
		}
	}

	g.variables[stmt.Name.Lexeme] = variable
	if stmt.Value != nil {
		g.parent = g.variableInitializer.ID
		assign := &parser.StmtAssignment{
			Variable: stmt.Name,
			Operator: stmt.AssignToken,
			Value:    stmt.Value,
		}
		err := assign.Accept(g)
		if err != nil {
			delete(g.variables, stmt.Name.Lexeme)
			return err
		}
		variable.DataType = g.dataType
		g.variableInitializer = g.blocks[g.blockID]
	}

	if variable.DataType == "" {
		delete(g.variables, stmt.Name.Lexeme)
		return g.newError("Cannot infer the data type of the variable. Please explicitly provide type information.", stmt.Name)
	}

	variable.Declared = true

	return nil
}

func (g *generator) VisitEvent(stmt *parser.StmtEvent) error {
	fn, ok := Events[stmt.Name.Lexeme]
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
	fn, ok := FuncCalls[stmt.Name.Lexeme]
	if !ok {
		return g.newError("Unknown function.", stmt.Name)
	}
	block, err := fn(g, stmt)
	if err != nil {
		return err
	}
	g.blockID = block.ID
	return nil
}

func (g *generator) VisitAssignment(stmt *parser.StmtAssignment) error {
	var block *blocks.Block
	if assignment, ok := Assignments[stmt.Variable.Lexeme]; ok {
		blockType := assignment.AssignType
		if stmt.Operator.Type == parser.TkPlusAssign {
			blockType = assignment.IncreaseType
		}

		block = g.NewBlock(blockType, false)
		value, err := g.value(block.ID, stmt.Operator, stmt.Value, assignment.DataType)
		if err != nil {
			return err
		}
		block.Inputs[assignment.InputName] = value
	} else {
		variable, ok := g.variables[stmt.Variable.Lexeme]
		if !ok {
			return g.newError("Unknown variable.", stmt.Variable)
		}
		block = g.NewBlock(blocks.VariableChangeBy, false)

		value, err := g.value(block.ID, stmt.Operator, stmt.Value, variable.DataType)
		if err != nil {
			return err
		}

		if stmt.Operator.Type == parser.TkAssign {
			block.Type = blocks.VariableSetTo
			index := 2
			if value[0].(int) == 1 {
				index = 1
			}
			value[index].([]any)[0] = 10
		}

		block.Inputs["VALUE"] = value
		block.Fields["VARIABLE"] = []any{variable.Name.Lexeme, variable.ID}
	}

	g.blockID = block.ID
	return nil
}

func (g *generator) VisitIf(stmt *parser.StmtIf) error {
	var block *blocks.Block
	if stmt.ElseBody == nil {
		block = g.NewBlock(blocks.ControlIf, false)
	} else {
		block = g.NewBlock(blocks.ControlIfElse, false)
	}
	g.parent = block.ID

	g.noNext = true
	err := stmt.Condition.Accept(g)
	if err != nil {
		return err
	}
	if g.dataType != parser.DTBool {
		return g.newError("The condition must be a boolean.", stmt.Keyword)
	}
	block.Inputs["CONDITION"] = []any{2, g.blockID}

	g.noNext = true
	for i, s := range stmt.Body {
		err = s.Accept(g)
		if err != nil {
			return err
		}
		if i == 0 {
			block.Inputs["SUBSTACK"] = []any{2, g.blockID}
		}
		g.parent = g.blockID
	}

	g.noNext = true
	for i, s := range stmt.ElseBody {
		err = s.Accept(g)
		if err != nil {
			return err
		}
		if i == 0 {
			block.Inputs["SUBSTACK2"] = []any{2, g.blockID}
		}
		g.parent = g.blockID
	}
	g.noNext = false

	g.blockID = block.ID
	return nil
}

func (g *generator) VisitLoop(stmt *parser.StmtLoop) error {
	var block *blocks.Block
	var err error
	parent := g.parent
	if stmt.Condition == nil {
		block = g.NewBlock(blocks.ControlRepeatForever, false)
	} else if stmt.Keyword.Type == parser.TkWhile {
		block = g.NewBlock(blocks.ControlRepeatUntil, false)
		g.parent = block.ID
		block.Inputs["CONDITION"], err = g.value(parent, stmt.Keyword, stmt.Condition, parser.DTBool)
		if err != nil {
			return err
		}
	} else if stmt.Keyword.Type == parser.TkFor {
		block = g.NewBlock(blocks.ControlRepeat, false)
		g.parent = block.ID
		block.Inputs["TIMES"], err = g.value(parent, stmt.Keyword, stmt.Condition, parser.DTNumber)
		if err != nil {
			return err
		}
	} else {
		return g.newError("Unknown loop type.", stmt.Keyword)
	}
	g.parent = block.ID
	g.noNext = true
	for i, s := range stmt.Body {
		err = s.Accept(g)
		if err != nil {
			return err
		}
		if i == 0 {
			block.Inputs["SUBSTACK"] = []any{2, g.blockID}
		}
		g.parent = g.blockID
	}
	g.noNext = false
	g.blockID = block.ID
	return nil
}

func (g *generator) VisitIdentifier(expr *parser.ExprIdentifier) error {
	if v, ok := Variables[expr.Name.Lexeme]; ok {
		block := g.NewBlock(v.blockType, false)
		if v.fields != nil {
			block.Fields = v.fields
		}
		if v.fn != nil {
			v.fn(g, block)
		}
		g.dataType = v.dataType
		g.blockID = block.ID
		return nil
	}

	if variable, ok := g.variables[expr.Name.Lexeme]; ok {
		if !variable.Declared {
			return g.newError("Cannot use variable in its own initializer.", expr.Name)
		}
		g.variableName = variable.Name.Lexeme
		g.dataType = variable.DataType
		return nil
	}

	return g.newError("Unknown identifier.", expr.Name)
}

func (g *generator) VisitExprFuncCall(expr *parser.ExprFuncCall) error {
	fn, ok := ExprFuncCalls[expr.Name.Lexeme]
	if !ok {
		return g.newError("Unknown function.", expr.Name)
	}
	block, dataType, err := fn(g, expr)
	if err != nil {
		return err
	}
	g.dataType = dataType
	g.blockID = block.ID
	return nil
}

func (g *generator) VisitLiteral(expr *parser.ExprLiteral) error {
	return g.newError("Literals are not allowed in this context.", expr.Token)
}

func (g *generator) VisitUnary(expr *parser.ExprUnary) error {
	var block *blocks.Block
	var dataType parser.DataType
	switch expr.Operator.Type {
	case parser.TkBang:
		dataType = parser.DTBool
		block = g.NewBlock(blocks.OpNot, false)
	}
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
		block = g.NewBlock(blocks.OpLessThan, false)
		operandDataType = parser.DTNumber
	case parser.TkGreater:
		block = g.NewBlock(blocks.OpGreaterThan, false)
		operandDataType = parser.DTNumber
	case parser.TkEqual:
		block = g.NewBlock(blocks.OpEquals, false)
		operandDataType = parser.DTNumber
	case parser.TkAnd:
		block = g.NewBlock(blocks.OpAnd, false)
		operandDataType = parser.DTBool
	case parser.TkOr:
		block = g.NewBlock(blocks.OpOr, false)
		operandDataType = parser.DTBool
	default:
		retDataType = parser.DTNumber
		operandDataType = parser.DTNumber
		operandName = "NUM"
		switch expr.Operator.Type {
		case parser.TkPlus:
			block = g.NewBlock(blocks.OpAdd, false)
		case parser.TkMinus:
			block = g.NewBlock(blocks.OpSubtract, false)
		case parser.TkMultiply:
			block = g.NewBlock(blocks.OpMultiply, false)
		case parser.TkDivide:
			block = g.NewBlock(blocks.OpDivide, false)
		case parser.TkModulus:
			block = g.NewBlock(blocks.OpMod, false)
		default:
			return g.newError("Unknown binary operator.", expr.Operator)
		}
	}

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

	g.dataType = retDataType
	g.blockID = block.ID
	return nil
}

func (g *generator) value(parent string, token parser.Token, expr parser.Expr, dataType parser.DataType) ([]any, error) {
	if literal, ok := expr.(*parser.ExprLiteral); ok {
		if dataType != "" && literal.Token.DataType != dataType {
			return nil, g.newError(fmt.Sprintf("The value must be of type %s.", dataType), literal.Token)
		}
		if literal.Token.DataType == parser.DTBool {
			return nil, g.newError("Boolean literals are not allowed in this context.", literal.Token)
		}
		g.dataType = literal.Token.DataType
		return []any{1, []any{4, fmt.Sprintf("%v", literal.Token.Literal)}}, nil
	} else {
		g.parent = parent
		g.noNext = true
		defer func() { g.variableName = "" }()
		err := expr.Accept(g)
		if err != nil {
			return nil, err
		}
		g.noNext = false
		if dataType != "" && g.dataType != dataType {
			return nil, g.newError(fmt.Sprintf("The value must be of type %s.", dataType), token)
		}
		if g.dataType == parser.DTBool {
			return []any{2, g.blockID}, nil
		}
		if g.variableName != "" {
			variable := g.variables[g.variableName]
			return []any{3, []any{12, variable.Name.Lexeme, variable.ID}, []any{4, ""}}, nil
		}
		return []any{3, g.blockID, []any{4, ""}}, nil
	}
}

func (g *generator) literal(token parser.Token, expr parser.Expr, dataType parser.DataType) (any, error) {
	if literal, ok := expr.(*parser.ExprLiteral); ok {
		if literal.Token.DataType != dataType {
			return nil, g.newError(fmt.Sprintf("The value must be of type %s.", dataType), literal.Token)
		}
		return literal.Token.Literal, nil
	}
	return nil, g.newError("Only literals are allowed in this context.", token)
}

func (g *generator) NewBlock(blockType blocks.BlockType, shadow bool) *blocks.Block {
	var block *blocks.Block
	if shadow {
		block = blocks.NewShadowBlock(blockType, g.parent)
	} else {
		block = blocks.NewBlock(blockType, g.parent)
	}
	g.blocks[block.ID] = block
	if !g.noNext {
		g.blocks[g.parent].Next = &block.ID
	}
	g.noNext = false
	return block
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
