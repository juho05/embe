package analyzer

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"

	"github.com/Bananenpro/embe/parser"
)

type Variable struct {
	ID       string
	Name     parser.Token
	DataType parser.DataType
	declared bool
	used     bool
}

type List struct {
	ID            string
	Name          parser.Token
	DataType      parser.DataType
	InitialValues []string
	used          bool
}

type Constant struct {
	Name  parser.Token
	Value parser.Token
	Type  parser.DataType
	used  bool
}

type Function struct {
	Name        parser.Token
	Params      []parser.FuncParam
	ProcCode    string
	ArgumentIDs []string
	StartLine   int
	EndLine     int
	used        bool
}

type analyzer struct {
	lines [][]rune

	variables map[string]*Variable
	lists     map[string]*List
	constants map[string]*Constant
	functions map[string]*Function

	variableIsList bool

	currentFunction *Function

	warnings []error

	unreachable bool

	variableInitializers []parser.Stmt
}

type Definitions struct {
	Variables map[string]*Variable
	Lists     map[string]*List
	Constants map[string]*Constant
	Functions map[string]*Function
}

type AnalyzerResult struct {
	Definitions Definitions
	Warnings    []error
	Errors      []error
}

func Analyze(statements []parser.Stmt, lines [][]rune) ([]parser.Stmt, AnalyzerResult) {
	a := &analyzer{
		lines:                lines,
		variables:            make(map[string]*Variable),
		lists:                make(map[string]*List),
		constants:            make(map[string]*Constant),
		functions:            make(map[string]*Function),
		warnings:             make([]error, 0),
		variableInitializers: make([]parser.Stmt, 0),
	}
	errs := make([]error, 0)
	for _, stmt := range statements {
		err := stmt.Accept(a)
		if err != nil {
			errs = append(errs, err)
		}
	}

	for _, v := range a.variables {
		if !v.used {
			a.newWarningTk("This variable is never used.", v.Name)
		}
	}

	for _, l := range a.lists {
		if !l.used {
			a.newWarningTk("This variable is never used.", l.Name)
		}
	}

	for _, c := range a.constants {
		if !c.used {
			a.newWarningTk("This constant is never used.", c.Name)
		}
	}

	for _, f := range a.functions {
		if !f.used {
			a.newWarningTk("This function is never called.", f.Name)
		}
	}

	if len(a.variableInitializers) > 0 {
		newStatements := make([]parser.Stmt, 0, len(statements)+1)
		newStatements = append(newStatements, &parser.StmtEvent{
			Name: parser.Token{
				Lexeme: "launch",
			},
			Body: a.variableInitializers,
		})
		newStatements = append(newStatements, statements...)
		statements = newStatements
	}

	return statements, AnalyzerResult{
		Definitions: Definitions{
			Variables: a.variables,
			Lists:     a.lists,
			Constants: a.constants,
			Functions: a.functions,
		},
		Warnings: a.warnings,
		Errors:   errs,
	}
}

func (a *analyzer) VisitVarDecl(stmt *parser.StmtVarDecl) error {
	if err := a.assertNotDeclared(stmt.Name); err != nil {
		return err
	}

	if _, ok := stmt.Value.(*parser.ExprListInitializer); ok || strings.HasSuffix(string(stmt.DataType), "[]") {
		list := &List{
			ID:            uuid.NewString(),
			Name:          stmt.Name,
			DataType:      stmt.DataType,
			InitialValues: make([]string, 0),
		}

		if stmt.Value == nil {
			stmt.Value = &parser.ExprListInitializer{
				OpenBracket: stmt.Name,
				Values:      make([]parser.Token, 0),
			}
		}

		var init *parser.ExprListInitializer
		if init, ok = stmt.Value.(*parser.ExprListInitializer); !ok {
			return a.newErrorExpr("Expected a list initializer.", stmt.Value)
		}

		valueType := parser.DataType(strings.TrimSuffix(string(stmt.DataType), "[]"))
		for _, v := range init.Values {
			token := v
			if v.Type == parser.TkIdentifier {
				if c, ok := a.constants[v.Lexeme]; ok {
					v = c.Value
				} else {
					return a.newErrorTk("Unknown constant.", v)
				}
			}
			if valueType == "" {
				valueType = v.DataType
			}
			if v.DataType != valueType {
				return a.newErrorTk(fmt.Sprintf("Wrong data type. Expected %s.", valueType), token)
			}
			list.InitialValues = append(list.InitialValues, fmt.Sprintf("%v", v.Literal))
		}
		if valueType != "" {
			list.DataType = valueType + "[]"
		}

		if list.DataType == "" {
			return a.newErrorTk("Cannot infer the data type of the variable. Please explicitly provide type information.", stmt.Name)
		}

		a.lists[list.Name.Lexeme] = list
	} else {
		variable := &Variable{
			ID:       uuid.NewString(),
			Name:     stmt.Name,
			DataType: stmt.DataType,
		}

		if variable.DataType != "" && stmt.Value == nil {
			stmt.AssignToken = parser.Token{
				Type: parser.TkAssign,
			}
			switch variable.DataType {
			case parser.DTNumber:
				stmt.Value = &parser.ExprLiteral{
					Token: parser.Token{
						Type:     parser.TkLiteral,
						Lexeme:   "0",
						Literal:  0,
						DataType: parser.DTNumber,
					},
				}
			case parser.DTString:
				stmt.Value = &parser.ExprLiteral{
					Token: parser.Token{
						Type:     parser.TkLiteral,
						Lexeme:   "",
						Literal:  "",
						DataType: parser.DTString,
					},
				}
			default:
				return a.newErrorTk("Unknown type.", stmt.Name)
			}
		}

		a.variables[stmt.Name.Lexeme] = variable
		if stmt.Value != nil {
			assign := &parser.StmtAssignment{
				Variable: stmt.Name,
				Operator: stmt.AssignToken,
				Value:    stmt.Value,
			}
			err := assign.Accept(a)
			if err != nil {
				delete(a.variables, stmt.Name.Lexeme)
				return err
			}
			variable.DataType = assign.Value.Type()
			a.variableInitializers = append(a.variableInitializers, assign)
		}

		if variable.DataType == "" {
			delete(a.variables, stmt.Name.Lexeme)
			return a.newErrorTk("Cannot infer the data type of the variable. Please explicitly provide type information.", stmt.Name)
		}

		variable.declared = true
	}
	return nil
}

func (a *analyzer) VisitConstDecl(stmt *parser.StmtConstDecl) error {
	if err := a.assertNotDeclared(stmt.Name); err != nil {
		return err
	}
	a.constants[stmt.Name.Lexeme] = &Constant{
		Name:  stmt.Name,
		Value: stmt.Value,
		Type:  stmt.Value.DataType,
	}
	return nil
}

func (a *analyzer) VisitFuncDecl(stmt *parser.StmtFuncDecl) error {
	if err := a.assertNotDeclared(stmt.Name); err != nil {
		return err
	}
	procCode := stmt.Name.Lexeme
	argumentIDs := make([]string, 0, len(stmt.Params))
	argumentNames := make([]string, 0, len(stmt.Params))
	for _, p := range stmt.Params {
		if slices.Contains(argumentNames, p.Name.Lexeme) {
			return a.newErrorTk("Duplicate parameter name.", p.Name)
		}

		id := uuid.NewString()
		argumentIDs = append(argumentIDs, id)
		argumentNames = append(argumentNames, p.Name.Lexeme)

		if p.Type.DataType == parser.DTBool {
			procCode += " %b"
		} else {
			if p.Type.DataType == parser.DTNumber {
				procCode += " %n"
			} else {
				procCode += " %s"
			}
		}
	}

	a.functions[stmt.Name.Lexeme] = &Function{
		Name:        stmt.Name,
		Params:      stmt.Params,
		ProcCode:    procCode,
		ArgumentIDs: argumentIDs,
		StartLine:   stmt.StartLine,
		EndLine:     stmt.EndLine,
	}

	a.currentFunction = a.functions[stmt.Name.Lexeme]
	for _, s := range stmt.Body {
		err := s.Accept(a)
		if err != nil {
			return err
		}
	}
	a.currentFunction = nil
	return nil
}

func (a *analyzer) assertNotDeclared(name parser.Token) error {
	if v, ok := a.variables[name.Lexeme]; ok {
		return a.newErrorTk(fmt.Sprintf("'%s' is already declared in line %d.", name.Lexeme, v.Name.Pos.Line+1), name)
	}
	if l, ok := a.lists[name.Lexeme]; ok {
		return a.newErrorTk(fmt.Sprintf("'%s' is already declared in line %d.", name.Lexeme, l.Name.Pos.Line+1), name)
	}
	if c, ok := a.constants[name.Lexeme]; ok {
		return a.newErrorTk(fmt.Sprintf("'%s' is already declared in line %d.", name.Lexeme, c.Name.Pos.Line+1), name)
	}
	if f, ok := a.functions[name.Lexeme]; ok {
		return a.newErrorTk(fmt.Sprintf("'%s' is already declared in line %d.", name.Lexeme, f.Name.Pos.Line+1), name)
	}
	return nil
}

func (a *analyzer) VisitEvent(stmt *parser.StmtEvent) error {
	ev, ok := Events[stmt.Name.Lexeme]
	if !ok {
		return a.newErrorStmt("Unknown event.", stmt)
	}
	if ev.Param == nil && stmt.Parameter != (parser.Token{}) {
		return a.newErrorTk("This event does not take a parameter.", stmt.Parameter)
	}
	if ev.Param != nil {
		if stmt.Parameter == (parser.Token{}) {
			return a.newErrorStmt(fmt.Sprintf("Please provide the %s parameter of type %s.", ev.Param.Name, ev.Param.Type), stmt)
		}
		var value any
		if stmt.Parameter.Type == parser.TkIdentifier {
			if constant, ok := a.constants[stmt.Parameter.Lexeme]; ok {
				if constant.Type != ev.Param.Type {
					return a.newErrorTk(fmt.Sprintf("Wrong data type. Expected '%s'.", ev.Param.Type), stmt.Parameter)
				}
				value = constant.Value.Literal
			} else {
				return a.newErrorTk("Unknown constant.", stmt.Parameter)
			}
		} else {
			if stmt.Parameter.DataType != ev.Param.Type {
				return a.newErrorTk(fmt.Sprintf("Wrong data type. Expected '%s'.", ev.Param.Type), stmt.Parameter)
			}
		}

		if ev.ParamOptions != nil {
			valid := false
			for _, o := range ev.ParamOptions {
				if value == o {
					valid = true
					break
				}
			}
			if !valid {
				strOptions := make([]string, len(ev.ParamOptions))
				for i, o := range ev.ParamOptions {
					strOptions[i] = fmt.Sprintf("%v", o)
				}
				return a.newErrorTk(fmt.Sprintf("Invalid value. Available options: %s", strings.Join(strOptions, ", ")), stmt.Parameter)
			}
		}
	}
	for _, s := range stmt.Body {
		err := s.Accept(a)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *analyzer) VisitFuncCall(stmt *parser.StmtFuncCall) error {
	if a.unreachable {
		a.newWarningStmt("Unreachable code.", stmt)
		return nil
	}

	if f, ok := a.functions[stmt.Name.Lexeme]; ok {
		f.used = true

		if len(stmt.Parameters) != len(f.Params) {
			return a.newErrorStmt("Wrong argument count.", stmt)
		}

		var err error
		for i, p := range stmt.Parameters {
			err = p.Accept(a)
			if err != nil {
				return err
			}
			if p.Type() != f.Params[i].Type.DataType {
				return a.newErrorExpr(fmt.Sprintf("Expected %s parameter '%s'.", f.Params[i].Type.DataType, f.Params[i].Name.Lexeme), p)
			}
		}
	} else {
		fn, ok := FuncCalls[stmt.Name.Lexeme]
		if !ok {
			if _, ok := ExprFuncCalls[stmt.Name.Lexeme]; ok {
				return a.newErrorStmt("Only functions which don't return a value are allowed in this context.", stmt)
			}
			return a.newErrorTk("Unknown function.", stmt.Name)
		}
		validSignature := false

		types := make([]string, len(stmt.Parameters))
		for i, p := range stmt.Parameters {
			err := p.Accept(a)
			if err != nil {
				return err
			}
			types[i] = string(p.Type())
		}
	signatures:
		for _, s := range fn.Signatures {
			if len(s.Params) != len(stmt.Parameters) {
				continue
			}
			for i, t := range types {
				if parser.DataType(t) != s.Params[i].Type {
					continue signatures
				}
			}
			validSignature = true
			break
		}
		if !validSignature {
			signatures := make([]string, len(fn.Signatures))
			for i, s := range fn.Signatures {
				sig := strings.Builder{}
				for j, p := range s.Params {
					sig.WriteString(string(p.Type))
					if j < len(s.Params)-1 {
						sig.WriteString(", ")
					}
				}
				signatures[i] = "(" + sig.String() + ")"
			}
			return a.newErrorStmt(fmt.Sprintf("Invalid arguments:\n  have: (%s)\n  want: %s", strings.Join(types, ", "), strings.Join(signatures, " or ")), stmt)
		}
	}
	return nil
}

func (a *analyzer) VisitAssignment(stmt *parser.StmtAssignment) error {
	if a.unreachable {
		a.newWarningStmt("Unreachable code.", stmt)
		return nil
	}

	if assignment, ok := Assignments[stmt.Variable.Lexeme]; ok {
		err := stmt.Value.Accept(a)
		if err != nil {
			return err
		}
		if stmt.Value.Type() != assignment.DataType {
			return a.newErrorExpr(fmt.Sprintf("Cannot assign %s value to %s variable.", stmt.Value.Type(), assignment.DataType), stmt.Value)
		}
	} else {
		v, ok := a.variables[stmt.Variable.Lexeme]
		if !ok {
			return a.newErrorTk("Unknown variable.", stmt.Variable)
		}
		err := stmt.Value.Accept(a)
		if err != nil {
			return err
		}
		if stmt.Value.Type() != v.DataType {
			return a.newErrorExpr(fmt.Sprintf("Cannot assign %s value to %s variable.", stmt.Value.Type(), v.DataType), stmt.Value)
		}
	}
	return nil
}

func (a *analyzer) VisitIf(stmt *parser.StmtIf) error {
	if a.unreachable {
		a.newWarningStmt("Unreachable code.", stmt)
		return nil
	}

	err := stmt.Condition.Accept(a)
	if err != nil {
		return err
	}
	if stmt.Condition.Type() != parser.DTBool {
		return a.newErrorExpr("Expected boolean condition.", stmt.Condition)
	}

	for _, s := range stmt.Body {
		err = s.Accept(a)
		if err != nil {
			return err
		}
	}

	for _, s := range stmt.ElseBody {
		err = s.Accept(a)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *analyzer) VisitLoop(stmt *parser.StmtLoop) error {
	if a.unreachable {
		a.newWarningStmt("Unreachable code.", stmt)
		return nil
	}
	switch stmt.Keyword.Type {
	case parser.TkWhile:
		err := stmt.Condition.Accept(a)
		if err != nil {
			return err
		}
		if stmt.Condition.Type() != parser.DTBool {
			return a.newErrorExpr("Expected boolean condition.", stmt.Condition)
		}
	case parser.TkFor:
		err := stmt.Condition.Accept(a)
		if err != nil {
			return err
		}
		if stmt.Condition.Type() != parser.DTNumber {
			return a.newErrorExpr("Expected number.", stmt.Condition)
		}
	default:
		return a.newErrorTk("Unknown loop type.", stmt.Keyword)
	}
	for _, s := range stmt.Body {
		err := s.Accept(a)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *analyzer) VisitIdentifier(expr *parser.ExprIdentifier) error {
	if a.currentFunction != nil {
		for _, p := range a.currentFunction.Params {
			if p.Name.Lexeme == expr.Name.Lexeme {
				expr.ReturnType = p.Type.DataType
				return nil
			}
		}
	}

	if v, ok := Variables[expr.Name.Lexeme]; ok {
		expr.ReturnType = v.DataType
		return nil
	}

	if variable, ok := a.variables[expr.Name.Lexeme]; ok {
		if !variable.declared {
			return a.newErrorTk("Cannot use variable in its own initializer.", expr.Name)
		}
		variable.used = true
		expr.ReturnType = variable.DataType
		return nil
	}

	if l, ok := a.lists[expr.Name.Lexeme]; ok {
		l.used = true
		a.variableIsList = true
		expr.ReturnType = l.DataType
		return nil
	}

	if c, ok := a.constants[expr.Name.Lexeme]; ok {
		c.used = true
		expr.ReturnType = c.Type
		return nil
	}

	return a.newErrorTk("Unknown identifier.", expr.Name)
}

func (a *analyzer) VisitExprFuncCall(expr *parser.ExprFuncCall) error {
	fn, ok := ExprFuncCalls[expr.Name.Lexeme]
	if !ok {
		if _, ok := FuncCalls[expr.Name.Lexeme]; ok {
			return a.newErrorExpr("Only functions which return a value are allowed in this context.", expr)
		}
		return a.newErrorTk("Unknown function.", expr.Name)
	}
	validSignature := false

	types := make([]string, len(expr.Parameters))
	for i, p := range expr.Parameters {
		err := p.Accept(a)
		if err != nil {
			return err
		}
		types[i] = string(p.Type())
	}
signatures:
	for _, s := range fn.Signatures {
		if len(s.Params) != len(expr.Parameters) {
			continue
		}
		for i, t := range types {
			if parser.DataType(t) != s.Params[i].Type {
				continue signatures
			}
		}
		validSignature = true
		expr.ReturnType = s.ReturnType
		break
	}
	if !validSignature {
		signatures := make([]string, len(fn.Signatures))
		for i, s := range fn.Signatures {
			sig := strings.Builder{}
			for j, p := range s.Params {
				sig.WriteString(string(p.Type))
				if j < len(s.Params)-1 {
					sig.WriteString(", ")
				}
			}
			signatures[i] = "(" + sig.String() + ")"
		}
		return a.newErrorExpr(fmt.Sprintf("Invalid arguments:\n  have: (%s)\n  want: %s", strings.Join(types, ", "), strings.Join(signatures, " or ")), expr)
	}
	return nil
}

func (a *analyzer) VisitTypeCast(expr *parser.ExprTypeCast) error {
	dataType := expr.Target.DataType
	err := expr.Value.Accept(a)
	if err != nil {
		return err
	}
	if expr.Target.DataType == parser.DTBool {
		return a.newErrorTk("Cannot cast to a boolean.", expr.Target)
	}
	if expr.Value.Type() == parser.DTBool {
		return a.newErrorExpr("Cannot cast from a boolean.", expr.Value)
	}

	if strings.HasSuffix(string(expr.Value.Type()), "[]") && expr.Target.DataType != parser.DTString {
		return a.newErrorExpr(fmt.Sprintf("Cannot cast list to %s.", expr.Target.DataType), expr.Value)
	}
	expr.ReturnType = dataType
	return nil
}

func (a *analyzer) VisitLiteral(expr *parser.ExprLiteral) error {
	expr.ReturnType = expr.Token.DataType
	return nil
}

func (a *analyzer) VisitListInitializer(expr *parser.ExprListInitializer) error {
	panic("Should never be called.")
}

func (a *analyzer) VisitUnary(expr *parser.ExprUnary) error {
	var dataType parser.DataType
	switch expr.Operator.Type {
	case parser.TkBang:
		dataType = parser.DTBool
	}
	err := expr.Right.Accept(a)
	if err != nil {
		return err
	}
	if expr.Right.Type() != dataType {
		return a.newErrorExpr(fmt.Sprintf("Expected operand of type %s.", dataType), expr.Right)
	}
	expr.ReturnType = dataType
	return nil
}

func (a *analyzer) VisitBinary(expr *parser.ExprBinary) error {
	retDataType := parser.DTBool
	if expr.Operator.Type == parser.TkPlus || expr.Operator.Type == parser.TkEqual {
		err := expr.Left.Accept(a)
		if err != nil {
			return err
		}
		leftType := expr.Left.Type()

		err = expr.Right.Accept(a)
		if err != nil {
			return err
		}
		rightType := expr.Right.Type()

		if leftType == parser.DTBool {
			return a.newErrorExpr("Expected number or string operand.", expr.Left)
		}
		if rightType == parser.DTBool {
			return a.newErrorExpr("Expected number or string operand.", expr.Right)
		}

		if expr.Operator.Type == parser.TkEqual {
			retDataType = parser.DTBool
		} else {
			if leftType == parser.DTString || rightType == parser.DTString {
				retDataType = parser.DTString
			} else {
				retDataType = parser.DTNumber
			}
		}
	} else {
		var operandDataType parser.DataType
		switch expr.Operator.Type {
		case parser.TkLess:
			operandDataType = parser.DTNumber
		case parser.TkGreater:
			operandDataType = parser.DTNumber
		case parser.TkAnd:
			operandDataType = parser.DTBool
		case parser.TkOr:
			operandDataType = parser.DTBool
		default:
			retDataType = parser.DTNumber
			operandDataType = parser.DTNumber
		}

		err := expr.Left.Accept(a)
		if err != nil {
			return err
		}
		if expr.Left.Type() != operandDataType {
			return a.newErrorExpr(fmt.Sprintf("Expected operand of type %s.", operandDataType), expr.Left)
		}

		err = expr.Right.Accept(a)
		if err != nil {
			return err
		}
		if expr.Right.Type() != operandDataType {
			return a.newErrorExpr(fmt.Sprintf("Expected operand of type %s.", operandDataType), expr.Right)
		}
	}

	expr.ReturnType = retDataType
	return nil
}

func (a *analyzer) VisitGrouping(expr *parser.ExprGrouping) error {
	return expr.Expr.Accept(a)
}

type AnalyzerError struct {
	Start   parser.Position
	End     parser.Position
	Message string
	Warning bool
}

func (e AnalyzerError) Error() string {
	if e.Warning {
		return "WARNING: " + e.Message
	} else {
		return "ERROR: " + e.Message
	}
}

func (a *analyzer) newErrorTk(message string, token parser.Token) error {
	end := token.Pos
	end.Column += len(token.Lexeme)
	if token.Type == parser.TkNewLine {
		end.Column += 1
	}
	return AnalyzerError{
		Start:   token.Pos,
		End:     end,
		Message: message,
	}
}

func (a *analyzer) newErrorExpr(message string, expr parser.Expr) error {
	start, end := expr.Position()
	return AnalyzerError{
		Start:   start,
		End:     end,
		Message: message,
	}
}

func (a *analyzer) newErrorStmt(message string, stmt parser.Stmt) error {
	start, end := stmt.Position()
	return AnalyzerError{
		Start:   start,
		End:     end,
		Message: message,
	}
}

func (a *analyzer) newWarningTk(message string, token parser.Token) {
	end := token.Pos
	end.Column += len(token.Lexeme)
	if token.Type == parser.TkNewLine {
		end.Column += 1
	}
	a.warnings = append(a.warnings, AnalyzerError{
		Start:   token.Pos,
		End:     end,
		Message: message,
		Warning: true,
	})
}

func (a *analyzer) newWarningExpr(message string, expr parser.Expr) {
	start, end := expr.Position()
	a.warnings = append(a.warnings, AnalyzerError{
		Start:   start,
		End:     end,
		Message: message,
		Warning: true,
	})
}

func (a *analyzer) newWarningStmt(message string, stmt parser.Stmt) {
	start, end := stmt.Position()
	a.warnings = append(a.warnings, AnalyzerError{
		Start:   start,
		End:     end,
		Message: message,
		Warning: true,
	})
}
