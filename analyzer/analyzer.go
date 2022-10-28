package analyzer

import (
	"fmt"
	"path/filepath"
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
	changed  bool
}

type List struct {
	ID            string
	Name          parser.Token
	DataType      parser.DataType
	InitialValues []string
	used          bool
}

type Constant struct {
	Name      parser.Token
	ValueExpr parser.Expr
	Value     any
	Type      parser.DataType
	used      bool
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

type CustomEvent struct {
	ID        string
	Name      parser.Token
	triggered bool
	consumed  bool
}

type analyzer struct {
	variables map[string]*Variable
	lists     map[string]*List
	constants map[string]*Constant
	functions map[string]*Function
	events    map[string]*CustomEvent

	variableIsList bool

	currentFunction *Function

	errors   []error
	warnings []error

	unreachable bool

	launchEventCount     int
	variableInitializers []parser.Stmt
}

type Definitions struct {
	Variables map[string]*Variable
	Lists     map[string]*List
	Constants map[string]*Constant
	Functions map[string]*Function
	Events    map[string]*CustomEvent
}

type AnalyzerResult struct {
	Definitions Definitions
	Warnings    []error
	Errors      []error
}

func Analyze(statements []parser.Stmt) ([]parser.Stmt, AnalyzerResult) {
	a := &analyzer{
		variables:            make(map[string]*Variable),
		lists:                make(map[string]*List),
		constants:            make(map[string]*Constant),
		functions:            make(map[string]*Function),
		events:               make(map[string]*CustomEvent),
		errors:               make([]error, 0),
		warnings:             make([]error, 0),
		variableInitializers: make([]parser.Stmt, 0),
	}
	for _, stmt := range statements {
		err := stmt.Accept(a)
		if err != nil {
			a.errors = append(a.errors, err)
		}
	}

	if len(a.errors) == 0 {
		for _, v := range a.variables {
			if !v.used {
				a.newWarningTk("This variable is never used.", v.Name)
			} else if !v.changed {
				a.newWarningTk("The value of this variable is never used. Consider using 'const' instead.", v.Name)
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

		for _, e := range a.events {
			if !e.triggered {
				a.newWarningTk("This event is never triggered.", e.Name)
			} else if !e.consumed {
				a.newWarningTk("This event is never consumed.", e.Name)
			}
		}
	}

	if len(a.variableInitializers) > 0 {
		if a.launchEventCount > 1 {
			startID := uuid.NewString()
			a.variableInitializers = append(a.variableInitializers, &parser.StmtCall{
				Name: parser.Token{
					Lexeme: "$start",
				},
			})
			a.events["$start"] = &CustomEvent{
				ID: startID,
				Name: parser.Token{
					Lexeme: "$start",
				},
				triggered: true,
				consumed:  true,
			}
			for _, s := range statements {
				if e, ok := s.(*parser.StmtEvent); ok {
					if e.Name.Lexeme == "launch" {
						e.Name.Lexeme = "$start"
					}
				}
			}
		} else if a.launchEventCount == 1 {
			for i, s := range statements {
				if e, ok := s.(*parser.StmtEvent); ok {
					if e.Name.Lexeme == "launch" {
						a.variableInitializers = append(a.variableInitializers, e.Body...)
						statements[i] = statements[len(statements)-1]
						statements = statements[:len(statements)-1]
						break
					}
				}
			}
		}

		statements = append(statements, &parser.StmtEvent{
			Name: parser.Token{
				Lexeme: "launch",
			},
			Body: a.variableInitializers,
		})
	}

	definitions := Definitions{
		Variables: a.variables,
		Lists:     a.lists,
		Constants: a.constants,
		Functions: a.functions,
		Events:    a.events,
	}

	if len(a.errors) == 0 {
		cErrs, cWarns := CalculateConstants(statements, definitions)
		a.errors = append(a.errors, cErrs...)
		a.warnings = append(a.warnings, cWarns...)
	}

	return statements, AnalyzerResult{
		Definitions: definitions,
		Warnings:    a.warnings,
		Errors:      a.errors,
	}
}

func (a *analyzer) VisitVarDecl(stmt *parser.StmtVarDecl) error {
	if err := a.assertNotDeclared(stmt.Name); err != nil {
		return err
	}

	if _, ok := stmt.Value.(*parser.ExprListInitializer); ok || strings.HasSuffix(string(stmt.DataType), "[]") {
		list := &List{
			ID:       uuid.NewString(),
			Name:     stmt.Name,
			DataType: stmt.DataType,
		}

		if stmt.Value == nil {
			closeBracket := stmt.Name
			closeBracket.Pos.Column += len(stmt.Name.Lexeme) - 1
			stmt.Value = &parser.ExprListInitializer{
				OpenBracket:  stmt.Name,
				CloseBracket: closeBracket,
				Values:       make([]parser.Expr, 0),
			}
		}

		var init *parser.ExprListInitializer
		if init, ok = stmt.Value.(*parser.ExprListInitializer); !ok {
			return a.newErrorExpr("Expected a list initializer.", stmt.Value)
		}

		valueType := parser.DataType(strings.TrimSuffix(string(stmt.DataType), "[]"))
		for _, v := range init.Values {
			err := v.Accept(a)
			if err != nil {
				a.errors = append(a.errors, err)
				continue
			}
			if valueType == "" {
				valueType = v.Type()
			}
			if v.Type() != valueType {
				a.errors = append(a.errors, a.newErrorExpr(fmt.Sprintf("Wrong data type. Expected %s.", valueType), v))
				continue
			}
		}
		if valueType != "" {
			list.DataType = valueType + "[]"
		}

		if list.DataType == "" {
			return a.newErrorTk("Cannot infer the data type of the variable. Please explicitly provide type information.", stmt.Name)
		}

		if list.DataType == "boolean[]" {
			if start, _ := stmt.Value.Position(); start == stmt.Name.Pos {
				stmt.Value = nil
				return a.newErrorStmt("Boolean lists are not supported", stmt)
			}
			return a.newErrorExpr("Boolean lists are not supported.", stmt.Value)
		}
		if list.DataType == "image[]" {
			if start, _ := stmt.Value.Position(); start == stmt.Name.Pos {
				stmt.Value = nil
				return a.newErrorStmt("Image lists are not supported", stmt)
			}
			return a.newErrorExpr("Image lists are not supported.", stmt.Value)
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
					ReturnType: parser.DTNumber,
				}
			case parser.DTString:
				stmt.Value = &parser.ExprLiteral{
					Token: parser.Token{
						Type:     parser.TkLiteral,
						Lexeme:   "",
						Literal:  "",
						DataType: parser.DTString,
					},
					ReturnType: parser.DTString,
				}
			case parser.DTImage:
				stmt.Value = &parser.ExprTypeCast{
					Target: parser.Token{
						Type:     parser.TkType,
						Lexeme:   "image",
						DataType: parser.DTImage,
					},
					Value: &parser.ExprLiteral{
						ReturnType: parser.DTString,
						Token: parser.Token{
							Type:     parser.TkLiteral,
							Lexeme:   "",
							Literal:  "",
							DataType: parser.DTString,
						},
					},
					ReturnType: parser.DTImage,
				}
			default:
				return a.newErrorStmt(fmt.Sprintf("%s variables are not supported.", strings.ToTitle(string(variable.DataType))), stmt)
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

		if variable.DataType == parser.DTBool {
			return a.newErrorStmt("Boolean variables are not supported.", stmt)
		}

		variable.declared = true
	}
	return nil
}

func (a *analyzer) VisitConstDecl(stmt *parser.StmtConstDecl) error {
	if err := a.assertNotDeclared(stmt.Name); err != nil {
		return err
	}
	err := stmt.Value.Accept(a)
	if err != nil {
		return err
	}
	if stmt.Value.Type() == parser.DTImage {
		return a.newErrorStmt("Image constants are not suuported.", stmt)
	}
	if stmt.Value.Type() == parser.DTBool {
		return a.newErrorStmt("Boolean constants are not suuported.", stmt)
	}
	a.constants[stmt.Name.Lexeme] = &Constant{
		Name:  stmt.Name,
		Value: stmt.Value,
		Type:  stmt.Value.Type(),
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
			a.errors = append(a.errors, a.newErrorTk("Duplicate parameter name.", p.Name))
			continue
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
	a.visitBody(stmt.Body)
	a.currentFunction = nil
	return nil
}

func (a *analyzer) VisitEvent(stmt *parser.StmtEvent) error {
	if e, ok := a.events[stmt.Name.Lexeme]; ok {
		if stmt.Parameter != nil {
			a.errors = append(a.errors, a.newErrorExpr("This event does not take a parameter.", stmt.Parameter))
		}
		e.consumed = true
	} else if ev, ok := Events[stmt.Name.Lexeme]; ok {
		if ev.Param == nil && stmt.Parameter != nil {
			a.errors = append(a.errors, a.newErrorExpr("This event does not take a parameter.", stmt.Parameter))
		} else if ev.Param != nil {
			if stmt.Parameter == nil {
				a.errors = append(a.errors, a.newErrorStmt(fmt.Sprintf("Please provide the %s parameter of type %s.", ev.Param.Name, ev.Param.Type), stmt))
			} else {
				err := stmt.Parameter.Accept(a)
				if err != nil {
					a.errors = append(a.errors, err)
				} else if stmt.Parameter.Type() != ev.Param.Type {
					a.errors = append(a.errors, a.newErrorExpr(fmt.Sprintf("Wrong data type. Expected '%s'.", ev.Param.Type), stmt.Parameter))
				}
			}
		}
	} else {
		a.errors = append(a.errors, a.newErrorStmt("Unknown event.", stmt))
	}

	a.visitBody(stmt.Body)

	if stmt.Name.Lexeme == "launch" {
		a.launchEventCount++
	}
	return nil
}

func (a *analyzer) VisitEventDecl(stmt *parser.StmtEventDecl) error {
	if err := a.assertNotDeclared(stmt.Name); err != nil {
		return err
	}
	if _, ok := Events[stmt.Name.Lexeme]; ok {
		return a.newErrorTk("An event with this name already exists.", stmt.Name)
	}

	a.events[stmt.Name.Lexeme] = &CustomEvent{
		ID:   uuid.NewString(),
		Name: stmt.Name,
	}
	return nil
}

func (a *analyzer) assertNotDeclared(name parser.Token) error {
	if v, ok := a.variables[name.Lexeme]; ok {
		return a.newErrorTk(fmt.Sprintf("'%s' is already declared in `%s` line %d.", name.Lexeme, filepath.Base(v.Name.Pos.Path), v.Name.Pos.Line+1), name)
	}
	if l, ok := a.lists[name.Lexeme]; ok {
		return a.newErrorTk(fmt.Sprintf("'%s' is already declared in `%s`, line %d.", name.Lexeme, filepath.Base(l.Name.Pos.Path), l.Name.Pos.Line+1), name)
	}
	if c, ok := a.constants[name.Lexeme]; ok {
		return a.newErrorTk(fmt.Sprintf("'%s' is already declared in `%s` line %d.", name.Lexeme, filepath.Base(c.Name.Pos.Path), c.Name.Pos.Line+1), name)
	}
	if f, ok := a.functions[name.Lexeme]; ok {
		return a.newErrorTk(fmt.Sprintf("'%s' is already declared in `%s` line %d.", name.Lexeme, filepath.Base(f.Name.Pos.Path), f.Name.Pos.Line+1), name)
	}
	if e, ok := a.events[name.Lexeme]; ok {
		return a.newErrorTk(fmt.Sprintf("'%s' is already declared in `%s` line %d.", name.Lexeme, filepath.Base(e.Name.Pos.Path), e.Name.Pos.Line+1), name)
	}
	return nil
}

func (a *analyzer) VisitCall(stmt *parser.StmtCall) error {
	if a.unreachable {
		a.newWarningStmt("Unreachable code.", stmt)
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
				a.errors = append(a.errors, err)
				continue
			}
			if p.Type() != f.Params[i].Type.DataType {
				a.errors = append(a.errors, a.newErrorExpr(fmt.Sprintf("Expected %s parameter '%s'.", f.Params[i].Type.DataType, f.Params[i].Name.Lexeme), p))
			}
		}
	} else if fn, ok := FuncCalls[stmt.Name.Lexeme]; ok {
		types := make([]string, len(stmt.Parameters))
		var hadError bool
		for i, p := range stmt.Parameters {
			err := p.Accept(a)
			if err != nil {
				a.errors = append(a.errors, err)
				hadError = true
				continue
			}
			types[i] = string(p.Type())
		}
		if !hadError {
			validSignature := false
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
	} else if ev, ok := a.events[stmt.Name.Lexeme]; ok {
		ev.triggered = true
		if len(stmt.Parameters) > 0 {
			return a.newErrorStmt("Events don't take any arguments.", stmt)
		}
	} else {
		if _, ok := ExprFuncCalls[stmt.Name.Lexeme]; ok {
			return a.newErrorStmt("Only functions which don't return a value are allowed in this context.", stmt)
		}
		return a.newErrorTk("Unknown function.", stmt.Name)
	}

	endFuncs := []string{"script.stop", "script.stopAll"}
	if slices.Contains(endFuncs, stmt.Name.Lexeme) {
		a.unreachable = true
	}
	return nil
}

func (a *analyzer) VisitAssignment(stmt *parser.StmtAssignment) error {
	if a.unreachable {
		a.newWarningStmt("Unreachable code.", stmt)
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
			if _, ok := a.constants[stmt.Variable.Lexeme]; ok {
				return a.newErrorStmt("Cannot change the value of a constant. Consider using 'var' instead.", stmt)
			}
			return a.newErrorTk("Unknown variable.", stmt.Variable)
		}
		if v.declared {
			v.changed = true
		}
		err := stmt.Value.Accept(a)
		if err != nil {
			return err
		}
		if v.DataType != "" && stmt.Value.Type() != v.DataType {
			return a.newErrorExpr(fmt.Sprintf("Cannot assign %s value to %s variable.", stmt.Value.Type(), v.DataType), stmt.Value)
		}
	}
	return nil
}

func (a *analyzer) VisitIf(stmt *parser.StmtIf) error {
	if a.unreachable {
		a.newWarningStmt("Unreachable code.", stmt)
	}

	err := stmt.Condition.Accept(a)
	if err != nil {
		return err
	}
	if stmt.Condition.Type() != parser.DTBool {
		return a.newErrorExpr("Expected boolean condition.", stmt.Condition)
	}

	a.visitBody(stmt.Body)
	a.visitBody(stmt.ElseBody)
	return nil
}

func (a *analyzer) VisitLoop(stmt *parser.StmtLoop) error {
	if a.unreachable {
		a.newWarningStmt("Unreachable code.", stmt)
	}
	forever := stmt.Condition == nil
	if !forever {
		switch stmt.Keyword.Type {
		case parser.TkWhile:
			err := stmt.Condition.Accept(a)
			if err != nil {
				a.errors = append(a.errors, err)
			} else if stmt.Condition.Type() != parser.DTBool {
				a.errors = append(a.errors, a.newErrorExpr("Expected boolean condition.", stmt.Condition))
			}
		case parser.TkFor:
			err := stmt.Condition.Accept(a)
			if err != nil {
				a.errors = append(a.errors, err)
			} else if stmt.Condition.Type() != parser.DTNumber {
				return a.newErrorExpr("Expected number.", stmt.Condition)
			}
		default:
			a.errors = append(a.errors, a.newErrorTk("Unknown loop type.", stmt.Keyword))
		}
	}
	a.visitBody(stmt.Body)
	if forever {
		a.unreachable = true
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
		if variable.DataType == parser.DTImage {
			variable.changed = true
		}
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
	types := make([]string, len(expr.Parameters))
	var hadError bool
	for i, p := range expr.Parameters {
		err := p.Accept(a)
		if err != nil {
			a.errors = append(a.errors, err)
			hadError = true
			continue
		}
		types[i] = string(p.Type())
	}
	if !hadError {
		validSignature := false
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
	} else {
		expr.ReturnType = fn.Signatures[0].ReturnType
	}
	return nil
}

func (a *analyzer) VisitTypeCast(expr *parser.ExprTypeCast) error {
	err := expr.Value.Accept(a)
	if err != nil {
		return err
	}
	if expr.Target.DataType == parser.DTBool {
		return a.newErrorTk("Cannot cast to a boolean.", expr.Target)
	}
	if expr.Target.DataType == parser.DTImage {
		if expr.Value.Type() != parser.DTString {
			return a.newErrorExpr("Expected file path.", expr.Value)
		}
	}
	if expr.Value.Type() == parser.DTBool {
		return a.newErrorExpr("Cannot cast a boolean to another type.", expr.Value)
	}
	if expr.Value.Type() == parser.DTImage {
		return a.newErrorExpr("Cannot cast an image to another type.", expr.Value)
	}

	if strings.HasSuffix(string(expr.Value.Type()), "[]") && expr.Target.DataType != parser.DTString {
		return a.newErrorExpr(fmt.Sprintf("Cannot cast list to %s.", expr.Target.DataType), expr.Value)
	}
	expr.ReturnType = expr.Target.DataType
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
			a.errors = append(a.errors, err)
		}
		leftType := expr.Left.Type()

		err = expr.Right.Accept(a)
		if err != nil {
			a.errors = append(a.errors, err)
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
			a.errors = append(a.errors, err)
		} else if expr.Left.Type() != operandDataType {
			return a.newErrorExpr(fmt.Sprintf("Expected operand of type %s.", operandDataType), expr.Left)
		}

		err = expr.Right.Accept(a)
		if err != nil {
			a.errors = append(a.errors, err)
		} else if expr.Right.Type() != operandDataType {
			return a.newErrorExpr(fmt.Sprintf("Expected operand of type %s.", operandDataType), expr.Right)
		}
	}

	expr.ReturnType = retDataType
	return nil
}

func (a *analyzer) VisitGrouping(expr *parser.ExprGrouping) error {
	return expr.Expr.Accept(a)
}

func (a *analyzer) visitBody(body []parser.Stmt) {
	unreachable := a.unreachable
	for _, s := range body {
		err := s.Accept(a)
		if err != nil {
			a.errors = append(a.errors, err)
		}
	}
	a.unreachable = unreachable
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
	end.Column += len(token.Lexeme) - 1
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
	end.Column += len(token.Lexeme) - 1
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
