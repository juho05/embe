package analyzer

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Bananenpro/embe/parser"
)

type constCalculator struct {
	definitions Definitions
	errors      []error
	warnings    []error

	newExpr parser.Expr
}

func CalculateConstants(statements []parser.Stmt, definitions Definitions) (errors []error, warnings []error) {
	calc := &constCalculator{
		definitions: definitions,
		warnings:    make([]error, 0),
	}

	for _, s := range statements {
		err := s.Accept(calc)
		if err != nil {
			break
		}
	}

	return calc.errors, calc.warnings
}

func (c *constCalculator) newLiteral(value any, expr parser.Expr) parser.Expr {
	start, end := expr.Position()
	return &parser.ExprLiteral{
		Token: parser.Token{
			Type:     parser.TkLiteral,
			DataType: expr.Type(),
			Pos:      start,
			Literal:  value,
		},
		End:        end,
		ReturnType: expr.Type(),
	}
}

func (c *constCalculator) VisitIdentifier(expr *parser.ExprIdentifier) error {
	if cn, ok := c.definitions.Constants[expr.Name.Lexeme]; ok {
		c.newExpr = c.newLiteral(cn.Value, expr)
	} else {
		c.newExpr = expr
	}
	return nil
}

func (c *constCalculator) VisitExprFuncCall(expr *parser.ExprFuncCall) error {
	for i, p := range expr.Parameters {
		err := p.Accept(c)
		if err != nil {
			return err
		}
		expr.Parameters[i] = c.newExpr
	}
	c.newExpr = expr
	return nil
}

func (c *constCalculator) VisitTypeCast(expr *parser.ExprTypeCast) error {
	err := expr.Value.Accept(c)
	if err != nil {
		return err
	}
	expr.Value = c.newExpr

	if expr.Target.DataType == parser.DTImage {
		var path string
		loadEmpty := false
		if literal, ok := expr.Value.(*parser.ExprLiteral); ok {
			path = literal.Token.Literal.(string)
			if literal.Token.Lexeme == "" {
				loadEmpty = true
			}
		} else {
			return c.newErrorExpr("Expected a constant file path.", expr.Value)
		}
		var img string
		if loadEmpty {
			img = strings.TrimSuffix(strings.Repeat("#000,", 16*16), ",")
		} else {
			img, err = loadImage(path)
			if err != nil {
				return c.newErrorExpr("Couldn't load image. Please provide a valid path to a PNG, JPEG or GIF file.", expr.Value)
			}
		}
		token := expr.Target
		token.Type = parser.TkLiteral
		token.Literal = img
		token.Lexeme = "\"" + img + "\""
		expr.Value = &parser.ExprLiteral{
			Token:      token,
			ReturnType: parser.DTString,
		}
	}

	if l, ok := expr.Value.(*parser.ExprLiteral); ok {
		var newValue any
		switch expr.Type() {
		case parser.DTString:
			newValue = fmt.Sprintf("%v", l.Token.Literal)
		case parser.DTNumber:
			newValue, err = strconv.ParseFloat(fmt.Sprintf("%v", l.Token.Literal), 64)
			if err != nil {
				c.newErrorExpr(fmt.Sprintf("Cannot convert %v to a number.", l.Token.Literal), expr)
			}
		default:
			c.newExpr = expr
			return nil
		}
		c.newExpr = c.newLiteral(newValue, expr)
	} else {
		c.newExpr = expr
	}
	return nil
}

func (c *constCalculator) VisitLiteral(expr *parser.ExprLiteral) error {
	c.newExpr = expr
	return nil
}

func (c *constCalculator) VisitListInitializer(expr *parser.ExprListInitializer) error {
	for i, v := range expr.Values {
		err := v.Accept(c)
		if err != nil {
			return err
		}
		expr.Values[i] = c.newExpr
		if _, ok := expr.Values[i].(*parser.ExprLiteral); !ok {
			c.newErrorExpr("Values in a list initializer must be constant.", expr.Values[i])
		}
	}
	c.newExpr = expr
	return nil
}

func (c *constCalculator) VisitUnary(expr *parser.ExprUnary) error {
	err := expr.Right.Accept(c)
	if err != nil {
		return err
	}
	expr.Right = c.newExpr

	if l, ok := expr.Right.(*parser.ExprLiteral); ok {
		if expr.Operator.Type == parser.TkMinus {
			c.newExpr = c.newLiteral(-l.Token.Literal.(float64), expr)
			return nil
		}
	}

	c.newExpr = expr
	return nil
}

func (c *constCalculator) VisitBinary(expr *parser.ExprBinary) error {
	err := expr.Left.Accept(c)
	if err != nil {
		return err
	}
	expr.Left = c.newExpr
	err = expr.Right.Accept(c)
	if err != nil {
		return err
	}
	expr.Right = c.newExpr
	if ll, ok := expr.Left.(*parser.ExprLiteral); ok {
		if lr, ok := expr.Right.(*parser.ExprLiteral); ok {
			var value any
			switch expr.Operator.Type {
			case parser.TkPlus:
				if ll.Token.DataType == parser.DTString || lr.Token.DataType == parser.DTString {
					value = fmt.Sprintf("%v%v", ll.Token.Literal, lr.Token.Literal)
				} else {
					value = ll.Token.Literal.(float64) + lr.Token.Literal.(float64)
				}
			case parser.TkMinus:
				value = ll.Token.Literal.(float64) - lr.Token.Literal.(float64)
			case parser.TkMultiply:
				value = ll.Token.Literal.(float64) * lr.Token.Literal.(float64)
			case parser.TkDivide:
				if lr.Token.Literal.(float64) == 0 {
					return c.newErrorExpr("Cannot divide by zero.", expr.Right)
				}
				value = ll.Token.Literal.(float64) / lr.Token.Literal.(float64)
			case parser.TkModulus:
				if lr.Token.Literal.(float64) == 0 {
					return c.newErrorExpr("Cannot divide by zero.", expr.Right)
				}
				value = math.Mod(ll.Token.Literal.(float64), lr.Token.Literal.(float64))
			default:
				c.newExpr = expr
				return nil
			}
			c.newExpr = c.newLiteral(value, expr)
			return nil
		}
	}
	c.newExpr = expr
	return nil
}

func (c *constCalculator) VisitGrouping(expr *parser.ExprGrouping) error {
	err := expr.Expr.Accept(c)
	if err != nil {
		return err
	}
	expr.Expr = c.newExpr
	if l, ok := expr.Expr.(*parser.ExprLiteral); ok {
		c.newExpr = c.newLiteral(l.Token.Literal, expr)
	} else {
		c.newExpr = expr
	}
	return nil
}

func (c *constCalculator) VisitVarDecl(stmt *parser.StmtVarDecl) error {
	err := stmt.Value.Accept(c)
	if err != nil {
		return err
	}
	stmt.Value = c.newExpr

	if init, ok := stmt.Value.(*parser.ExprListInitializer); ok {
		if l, ok := c.definitions.Lists[stmt.Name.Lexeme]; ok {
			l.InitialValues = make([]string, len(init.Values))
			for i, v := range init.Values {
				if lit, ok := v.(*parser.ExprLiteral); ok {
					l.InitialValues[i] = fmt.Sprintf("%v", lit.Token.Literal)
				}
			}
		}
	}
	return nil
}

func (c *constCalculator) VisitConstDecl(stmt *parser.StmtConstDecl) error {
	err := stmt.Value.Accept(c)
	if err != nil {
		return err
	}
	stmt.Value = c.newExpr
	if l, ok := stmt.Value.(*parser.ExprLiteral); !ok {
		return c.newErrorExpr("Cannot assign a non-constant value to a constant.", stmt.Value)
	} else {
		c.definitions.Constants[stmt.Name.Lexeme].Value = l.Token.Literal
	}
	return nil
}

func (c *constCalculator) VisitFuncDecl(stmt *parser.StmtFuncDecl) error {
	for _, s := range stmt.Body {
		err := s.Accept(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *constCalculator) VisitEvent(stmt *parser.StmtEvent) error {
	if stmt.Parameter != nil {
		err := stmt.Parameter.Accept(c)
		if err != nil {
			return err
		}
		stmt.Parameter = c.newExpr
		if l, ok := stmt.Parameter.(*parser.ExprLiteral); !ok {
			return c.newErrorExpr("Event parameters must be constant.", stmt.Parameter)
		} else {
			ev := Events[stmt.Name.Lexeme]
			if ev.ParamOptions != nil {
				valid := false
				for _, o := range ev.ParamOptions {
					if l.Token.Literal == o {
						valid = true
						break
					}
				}
				if !valid {
					strOptions := make([]string, len(ev.ParamOptions))
					for i, o := range ev.ParamOptions {
						strOptions[i] = fmt.Sprintf("%v", o)
					}
					return c.newErrorExpr(fmt.Sprintf("Invalid value. Available options: %s", strings.Join(strOptions, ", ")), stmt.Parameter)
				}
			}
		}
	}

	for _, s := range stmt.Body {
		err := s.Accept(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *constCalculator) VisitFuncCall(stmt *parser.StmtFuncCall) error {
	for i, p := range stmt.Parameters {
		err := p.Accept(c)
		if err != nil {
			return err
		}
		stmt.Parameters[i] = c.newExpr
	}
	return nil
}

func (c *constCalculator) VisitAssignment(stmt *parser.StmtAssignment) error {
	err := stmt.Value.Accept(c)
	if err != nil {
		return err
	}
	stmt.Value = c.newExpr
	return nil
}

func (c *constCalculator) VisitIf(stmt *parser.StmtIf) error {
	err := stmt.Condition.Accept(c)
	if err != nil {
		return err
	}
	stmt.Condition = c.newExpr
	for _, s := range stmt.Body {
		err = s.Accept(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *constCalculator) VisitLoop(stmt *parser.StmtLoop) error {
	err := stmt.Condition.Accept(c)
	if err != nil {
		return err
	}
	stmt.Condition = c.newExpr
	for _, s := range stmt.Body {
		err = s.Accept(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *constCalculator) newErrorTk(message string, token parser.Token) error {
	end := token.Pos
	end.Column += len(token.Lexeme) - 1
	if token.Type == parser.TkNewLine {
		end.Column += 1
	}
	err := AnalyzerError{
		Start:   token.Pos,
		End:     end,
		Message: message,
	}
	c.errors = append(c.errors, err)
	return err
}

func (c *constCalculator) newErrorExpr(message string, expr parser.Expr) error {
	start, end := expr.Position()
	err := AnalyzerError{
		Start:   start,
		End:     end,
		Message: message,
	}
	c.errors = append(c.errors, err)
	return err
}

func (c *constCalculator) newErrorStmt(message string, stmt parser.Stmt) error {
	start, end := stmt.Position()
	err := AnalyzerError{
		Start:   start,
		End:     end,
		Message: message,
	}
	c.errors = append(c.errors, err)
	return err
}

func (c *constCalculator) newWarningTk(message string, token parser.Token) {
	end := token.Pos
	end.Column += len(token.Lexeme) - 1
	if token.Type == parser.TkNewLine {
		end.Column += 1
	}
	c.warnings = append(c.warnings, AnalyzerError{
		Start:   token.Pos,
		End:     end,
		Message: message,
		Warning: true,
	})
}

func (c *constCalculator) newWarningExpr(message string, expr parser.Expr) {
	start, end := expr.Position()
	c.warnings = append(c.warnings, AnalyzerError{
		Start:   start,
		End:     end,
		Message: message,
		Warning: true,
	})
}

func (c *constCalculator) newWarningStmt(message string, stmt parser.Stmt) {
	start, end := stmt.Position()
	c.warnings = append(c.warnings, AnalyzerError{
		Start:   start,
		End:     end,
		Message: message,
		Warning: true,
	})
}
