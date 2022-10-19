package parser

import (
	"fmt"
	"strings"
)

type Constant struct {
	Token Token
}

type parser struct {
	tokens  []Token
	current int
	errors  []error
}

func Parse(tokens []Token) ([]Stmt, []error) {
	parser := &parser{
		tokens: tokens,
		errors: make([]error, 0),
	}
	return parser.parse()
}

func (p *parser) parse() ([]Stmt, []error) {
	statements := make([]Stmt, 0)
	for p.peek().Type != TkEOF {
		statements = append(statements, p.topLevel())
	}
	return statements, p.errors
}

func (p *parser) topLevel() Stmt {
	var err error
	var stmt Stmt
	switch p.peek().Type {
	case TkVar:
		stmt, err = p.varDecl()
	case TkConst:
		stmt, err = p.constDecl()
	case TkFunc:
		stmt, err = p.funcDecl()
	case TkAt:
		stmt, err = p.event()
	default:
		err = p.newError("Expected event or declaration.")
	}
	if err != nil {
		p.errors = append(p.errors, err)
		p.synchronize()
	}
	return stmt
}

func (p *parser) varDecl() (Stmt, error) {
	if !p.match(TkVar) {
		return nil, p.newError("Expected 'var' keyword.")
	}

	if !p.match(TkIdentifier) {
		return nil, p.newError("Expected variable name.")
	}
	name := p.previous()
	if strings.Contains(name.Lexeme, ".") {
		return nil, p.newErrorAt("Variable names cannot contain a dot.", name)
	}

	var dataType DataType
	if p.match(TkColon) {
		if !p.match(TkType) {
			return nil, p.newError("Expected type after ':'.")
		}
		var ok bool
		dataType, ok = types[p.previous().Lexeme]
		if !ok {
			if dataType, ok = types[strings.TrimSuffix(p.previous().Lexeme, "[]")]; !ok {
				return nil, p.newError("Unknown data type.")
			}
			dataType += "[]"
		}
		if dataType == DTBool {
			return nil, p.newErrorAt("Boolean variables are not supported.", p.previous())
		}
	}

	var value Expr
	var err error
	var assignToken Token
	if p.match(TkAssign) {
		assignToken = p.previous()
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	if !p.match(TkNewLine) {
		return nil, p.newError("Expected '\n' after variable declaration.")
	}

	return &StmtVarDecl{
		Name:        name,
		DataType:    dataType,
		AssignToken: assignToken,
		Value:       value,
	}, nil
}

func (p *parser) constDecl() (Stmt, error) {
	if !p.match(TkConst) {
		return nil, p.newError("Expected 'const' keyword.")
	}

	if !p.match(TkIdentifier) {
		return nil, p.newError("Expected constant name.")
	}
	name := p.previous()
	if strings.Contains(name.Lexeme, ".") {
		return nil, p.newErrorAt("Constant names cannot contain a dot.", name)
	}

	if p.match(TkColon) {
		if !p.match(TkType) {
			return nil, p.newError("Expected type after ':'.")
		}

		if _, ok := types[p.previous().Lexeme]; !ok {
			return nil, p.newError("Unknown data type.")
		}
	}

	if !p.match(TkAssign) {
		return nil, p.newError("Expected '=' after constant name.")
	}
	assignToken := p.previous()

	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(TkNewLine) {
		return nil, p.newError("Expected '\n' after constant declaration.")
	}

	return &StmtConstDecl{
		Name:        name,
		AssignToken: assignToken,
		Value:       value,
	}, nil
}

func (p *parser) funcDecl() (Stmt, error) {
	if !p.match(TkFunc) {
		return nil, p.newError("Expected 'func' keyword.")
	}

	if !p.match(TkIdentifier) {
		return nil, p.newError("Expected function name.")
	}
	name := p.previous()
	if strings.Contains(name.Lexeme, ".") {
		return nil, p.newErrorAt("Function names cannot contain a dot.", name)
	}

	if !p.match(TkOpenParen) {
		return nil, p.newError("Expected '(' after function name.")
	}

	parameters := make([]FuncParam, 0)
	for p.peek().Type != TkCloseParen && p.peek().Type != TkEOF {
		if !p.match(TkIdentifier) {
			return nil, p.newError("Expected parameter name.")
		}
		pName := p.previous()
		if !p.match(TkColon) {
			return nil, p.newError("Expected ':' after parameter name.")
		}
		if !p.match(TkType) {
			return nil, p.newError("Expected type after ':'.")
		}
		pType := p.previous()
		parameters = append(parameters, FuncParam{
			Name: pName,
			Type: pType,
		})
		if !p.match(TkComma) {
			break
		}
	}

	if !p.match(TkCloseParen) {
		return nil, p.newError("Expected ')' after parameter list.")
	}
	closeParen := p.previous()

	if !p.match(TkColon) {
		return nil, p.newError("Expected ':' after function declaration.")
	}

	if !p.match(TkNewLine) {
		return nil, p.newError("Expected '\n' after ':'.")
	}

	start := name.Pos.Line
	body := p.statements(name.Indent + 1)
	end := p.previous().Pos.Line + 1

	return &StmtFuncDecl{
		Name:       name,
		CloseParen: closeParen,
		Body:       body,
		Params:     parameters,
		StartLine:  start,
		EndLine:    end,
	}, nil
}

func (p *parser) event() (Stmt, error) {
	if !p.match(TkAt) {
		return nil, p.newError("Expected event.")
	}
	at := p.previous()

	if !p.match(TkIdentifier) {
		return nil, p.newError("Expected event name after '@'.")
	}
	name := p.previous()

	var parameter Expr
	var err error
	if p.peek().Type != TkColon && p.peek().Type != TkNewLine {
		parameter, err = p.expression()
		if err != nil {
			p.errors = append(p.errors, err)
			p.synchronize()
		}
	}

	if err == nil {
		if !p.match(TkColon) {
			p.errors = append(p.errors, p.newError("Expected ':' after parameter."))
			p.synchronize()
		} else if !p.match(TkNewLine) {
			p.errors = append(p.errors, p.newError("Expected '\n' after ':'."))
			p.synchronize()
		}
	}

	body := p.statements(name.Indent + 1)

	return &StmtEvent{
		At:        at,
		Name:      name,
		Parameter: parameter,
		Body:      body,
	}, nil
}

func (p *parser) statements(indent int) []Stmt {
	statements := make([]Stmt, 0, 10)
	for p.peek().Indent >= indent {
		stmt, err := p.statement()
		if err == nil {
			statements = append(statements, stmt)
		} else {
			p.errors = append(p.errors, err)
			p.synchronize()
		}
	}
	return statements
}

func (p *parser) statement() (Stmt, error) {
	switch p.peek().Type {
	case TkIf:
		return p.ifStmt()
	case TkWhile:
		return p.whileLoop()
	case TkFor:
		return p.forLoop()
	}

	if p.peekNext().Type == TkOpenParen {
		return p.funcCall()
	} else if p.peekNext().Type == TkAssign || p.peekNext().Type == TkPlusAssign || p.peekNext().Type == TkMinusAssign || p.peekNext().Type == TkMultiplyAssign || p.peekNext().Type == TkDivideAssign || p.peekNext().Type == TkModulusAssign {
		return p.assignment()
	}

	return nil, p.newError("Expected statement.")
}

func (p *parser) funcCall() (Stmt, error) {
	if !p.match(TkIdentifier) {
		return nil, p.newError("Expected identifier.")
	}
	name := p.previous()

	if !p.match(TkOpenParen) {
		return nil, p.newError("Expected '(' after identifier.")
	}

	parameters := make([]Expr, 0, 1)
	for p.peek().Type != TkCloseParen && p.peek().Type != TkEOF {
		param, err := p.expression()
		if err != nil {
			return nil, err
		}
		parameters = append(parameters, param)
		if !p.match(TkComma) {
			break
		}
	}

	if !p.match(TkCloseParen) {
		return nil, p.newError("Expected ')' after parameter list.")
	}
	closeParen := p.previous()

	if !p.match(TkNewLine) {
		return nil, p.newError("Expected '\n' after statement.")
	}

	return &StmtFuncCall{
		Name:       name,
		CloseParen: closeParen,
		Parameters: parameters,
	}, nil
}

func (p *parser) assignment() (Stmt, error) {
	if !p.match(TkIdentifier) {
		return nil, p.newError("Expected identifier.")
	}
	variable := p.previous()

	if !p.match(TkAssign, TkPlusAssign, TkMinusAssign, TkMultiplyAssign, TkDivideAssign, TkModulusAssign) {
		return nil, p.newError("Expected assignment operator after identifier.")
	}
	operator := p.previous()

	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(TkNewLine) {
		return nil, p.newError("Expected '\n' after statement.")
	}

	if operator.Type == TkMultiplyAssign || operator.Type == TkDivideAssign || operator.Type == TkModulusAssign {
		assignOp := operator
		assignOp.Type = TkAssign

		binOp := operator
		binOp.Type--

		return &StmtAssignment{
			Variable: variable,
			Operator: assignOp,
			Value: &ExprBinary{
				Left: &ExprIdentifier{
					Name: variable,
				},
				Operator: binOp,
				Right:    value,
			},
		}, nil
	}

	if operator.Type == TkMinusAssign {
		operator.Type = TkPlusAssign
		if v, ok := value.(*ExprLiteral); ok && v.Token.DataType == DTNumber {
			v.Token.Literal = -v.Token.Literal.(float64)
			return &StmtAssignment{
				Variable: variable,
				Operator: operator,
				Value:    v,
			}, nil
		} else {
			return &StmtAssignment{
				Variable: variable,
				Operator: operator,
				Value: &ExprBinary{
					Operator: Token{
						Type:   TkMultiply,
						Lexeme: operator.Lexeme,
						Pos:    operator.Pos,
						Indent: operator.Indent,
					},
					Left: &ExprLiteral{
						Token: Token{
							Type:     TkLiteral,
							Lexeme:   operator.Lexeme,
							Pos:      operator.Pos,
							Indent:   operator.Indent,
							DataType: DTNumber,
							Literal:  -1,
						},
					},
					Right: value,
				},
			}, nil
		}
	}

	return &StmtAssignment{
		Variable: variable,
		Operator: operator,
		Value:    value,
	}, nil
}

func (p *parser) ifStmt() (Stmt, error) {
	if !p.match(TkIf) {
		return nil, p.newError("Expected 'if' keyword.")
	}
	keyword := p.previous()

	condition, err := p.expression()
	if err != nil {
		p.errors = append(p.errors, err)
		p.synchronize()
	}

	if err == nil {
		if !p.match(TkColon) {
			p.errors = append(p.errors, p.newError("Expected ':' after if condition."))
			p.synchronize()
		} else if !p.match(TkNewLine) {
			p.errors = append(p.errors, p.newError("Expected '\\n' after ':'."))
			p.synchronize()
		}
	}

	body := p.statements(keyword.Indent + 1)

	stmt := &StmtIf{
		Keyword:   keyword,
		Condition: condition,
		Body:      body,
	}

	elifStmt := stmt
	for p.match(TkElif) {
		elifKeyword := p.previous()
		elifCondition, err := p.expression()
		if err != nil {
			p.errors = append(p.errors, err)
			p.synchronize()
		}

		if err == nil {
			if !p.match(TkColon) {
				p.errors = append(p.errors, p.newError("Expected ':' after if condition."))
				p.synchronize()
			} else if !p.match(TkNewLine) {
				p.errors = append(p.errors, p.newError("Expected '\\n' after ':'."))
				p.synchronize()
			}
		}

		elifBody := p.statements(elifKeyword.Indent + 1)
		s := &StmtIf{
			Keyword:   elifKeyword,
			Condition: elifCondition,
			Body:      elifBody,
		}
		elifStmt.ElseBody = []Stmt{s}
		elifStmt = s
	}

	if p.match(TkElse) {
		elseKeyword := p.previous()
		if !p.match(TkColon) {
			p.errors = append(p.errors, p.newError("Expected ':' after 'else'."))
			p.synchronize()
		} else if !p.match(TkNewLine) {
			p.errors = append(p.errors, p.newError("Expected '\\n' after ':'."))
			p.synchronize()
		}
		elifStmt.ElseBody = p.statements(elseKeyword.Indent + 1)
	}

	return stmt, nil
}

func (p *parser) whileLoop() (Stmt, error) {
	if !p.match(TkWhile) {
		return nil, p.newError("Expected 'while' keyword.")
	}
	keyword := p.previous()

	var condition Expr
	var err error
	if p.peek().Type != TkColon {
		condition, err = p.expression()
		if err != nil {
			p.errors = append(p.errors, err)
			p.synchronize()
		} else {
			// invert 'while' to 'until'
			if c, ok := condition.(*ExprUnary); ok && c.Operator.Type == TkBang {
				condition = c.Right
			} else {
				condition = &ExprUnary{
					Operator: Token{
						Type:   TkBang,
						Lexeme: keyword.Lexeme,
						Pos:    keyword.Pos,
						Indent: keyword.Indent,
					},
					Right: condition,
				}
			}
		}
	}

	if err == nil {
		if !p.match(TkColon) {
			p.errors = append(p.errors, p.newError("Expected ':' at the end of the while statement."))
			p.synchronize()
		} else if !p.match(TkNewLine) {
			p.errors = append(p.errors, p.newError("Expected '\\n' after ':'."))
			p.synchronize()
		}
	}

	body := p.statements(keyword.Indent + 1)

	return &StmtLoop{
		Keyword:   keyword,
		Condition: condition,
		Body:      body,
	}, nil
}

func (p *parser) forLoop() (Stmt, error) {
	if !p.match(TkFor) {
		return nil, p.newError("Expected 'for' keyword.")
	}
	keyword := p.previous()

	var condition Expr
	var err error
	if p.peek().Type != TkColon {
		condition, err = p.expression()
		if err != nil {
			p.errors = append(p.errors, err)
			p.synchronize()
		}
	}

	if err == nil {
		if !p.match(TkColon) {
			p.errors = append(p.errors, p.newError("Expected ':' at the end of the for statement."))
			p.synchronize()
		} else if !p.match(TkNewLine) {
			p.errors = append(p.errors, p.newError("Expected '\\n' after ':'."))
			p.synchronize()
		}
	}

	body := p.statements(keyword.Indent + 1)

	return &StmtLoop{
		Keyword:   keyword,
		Condition: condition,
		Body:      body,
	}, nil
}

func (p *parser) expression() (Expr, error) {
	return p.or()
}

func (p *parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		if !p.synchronize(TkOr) {
			p.current--
			return nil, err
		}
		p.errors = append(p.errors, err)
	}

	for p.match(TkOr) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			if !p.synchronize(TkOr) {
				p.current--
				return nil, err
			}
			p.errors = append(p.errors, err)
		}
		expr = &ExprBinary{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		if !p.synchronize(TkAnd) {
			p.current--
			return nil, err
		}
		p.errors = append(p.errors, err)
	}

	for p.match(TkAnd) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			if !p.synchronize(TkAnd) {
				p.current--
				return nil, err
			}
			p.errors = append(p.errors, err)
		}
		expr = &ExprBinary{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		if !p.synchronize(TkEqual, TkNotEqual) {
			p.current--
			return nil, err
		}
		p.errors = append(p.errors, err)
	}

	for p.match(TkEqual, TkNotEqual) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			if !p.synchronize(TkEqual, TkNotEqual) {
				p.current--
				return nil, err
			}
			p.errors = append(p.errors, err)
		}
		if operator.Type == TkNotEqual {
			operator.Type = TkEqual
			expr = &ExprUnary{
				Operator: Token{
					Type: TkBang,
				},
				Right: &ExprBinary{
					Operator: operator,
					Left:     expr,
					Right:    right,
				},
			}
		} else {
			expr = &ExprBinary{
				Operator: operator,
				Left:     expr,
				Right:    right,
			}
		}
	}

	return expr, nil
}

func (p *parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		if !p.synchronize(TkLess, TkLessEqual, TkGreater, TkGreaterEqual) {
			p.current--
			return nil, err
		}
		p.errors = append(p.errors, err)
	}

	for p.match(TkLess, TkLessEqual, TkGreater, TkGreaterEqual) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			if !p.synchronize(TkLess, TkLessEqual, TkGreater, TkGreaterEqual) {
				p.current--
				return nil, err
			}
			p.errors = append(p.errors, err)
		}
		if operator.Type == TkLessEqual || operator.Type == TkGreaterEqual {
			withoutEqual := TkLess
			if operator.Type == TkGreaterEqual {
				withoutEqual = TkGreater
			}
			expr = &ExprBinary{
				Operator: Token{
					Type:   TkOr,
					Lexeme: operator.Lexeme,
					Indent: operator.Indent,
					Pos:    operator.Pos,
				},
				Left: &ExprBinary{
					Operator: Token{
						Type:   withoutEqual,
						Lexeme: operator.Lexeme,
						Indent: operator.Indent,
						Pos:    operator.Pos,
					},
					Left:  expr,
					Right: right,
				},
				Right: &ExprBinary{
					Operator: Token{
						Type:   TkEqual,
						Lexeme: operator.Lexeme,
						Indent: operator.Indent,
						Pos:    operator.Pos,
					},
					Left:  expr,
					Right: right,
				},
			}
		} else {
			expr = &ExprBinary{
				Operator: operator,
				Left:     expr,
				Right:    right,
			}
		}
	}

	return expr, nil
}

func (p *parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		if !p.synchronize(TkPlus, TkMinus) {
			p.current--
			return nil, err
		}
		p.errors = append(p.errors, err)
	}

	for p.match(TkPlus, TkMinus) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			if !p.synchronize(TkPlus, TkMinus) {
				p.current--
				return nil, err
			}
			p.errors = append(p.errors, err)
		}
		expr = &ExprBinary{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		if !p.synchronize(TkMultiply, TkDivide, TkModulus) {
			p.current--
			return nil, err
		}
		p.errors = append(p.errors, err)
	}

	for p.match(TkMultiply, TkDivide, TkModulus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			if !p.synchronize(TkMultiply, TkDivide, TkModulus) {
				p.current--
				return nil, err
			}
			p.errors = append(p.errors, err)
		}
		expr = &ExprBinary{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) unary() (Expr, error) {
	if p.match(TkBang, TkMinus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		if operator.Type == TkMinus {
			if l, ok := right.(*ExprLiteral); ok && l.Token.DataType == DTNumber {
				l.Token.Literal = -l.Token.Literal.(float64)
				return l, nil
			}
			return &ExprBinary{
				Operator: Token{
					Type:   TkMultiply,
					Lexeme: operator.Lexeme,
					Pos:    operator.Pos,
					Indent: operator.Indent,
				},
				Left: &ExprLiteral{
					Token: Token{
						Type:     TkLiteral,
						Lexeme:   operator.Lexeme,
						Pos:      operator.Pos,
						Indent:   operator.Indent,
						DataType: DTNumber,
						Literal:  -1,
					},
				},
				Right: right,
			}, nil
		} else {
			return &ExprUnary{
				Operator: operator,
				Right:    right,
			}, nil
		}
	}

	return p.primary()
}

func (p *parser) primary() (Expr, error) {
	if p.match(TkIdentifier) {
		name := p.previous()
		if p.match(TkOpenParen) {
			parameters := make([]Expr, 0, 1)
			for p.peek().Type != TkCloseParen && p.peek().Type != TkEOF {
				param, err := p.expression()
				if err != nil {
					return nil, err
				}
				parameters = append(parameters, param)
				if !p.match(TkComma) {
					break
				}
			}

			if !p.match(TkCloseParen) {
				return nil, p.newError("Expected ')' after parameter list.")
			}

			return &ExprFuncCall{
				Name:       name,
				Parameters: parameters,
				CloseParen: p.previous(),
			}, nil
		}
		return &ExprIdentifier{
			Name: name,
		}, nil
	}

	if p.match(TkType) {
		token := p.previous()
		if !p.match(TkOpenParen) {
			return nil, p.newError("Expected '(' after type name for type cast.")
		}

		value, err := p.expression()
		if err != nil {
			return nil, err
		}

		if !p.match(TkCloseParen) {
			return nil, p.newError("Expected ')' after value for type cast.")
		}
		return &ExprTypeCast{
			Target:     token,
			Value:      value,
			CloseParen: p.previous(),
		}, nil
	}

	if p.match(TkOpenBracket) {
		openBracket := p.previous()
		values := make([]Expr, 0)
		for p.peek().Type != TkCloseBracket && p.peek().Type != TkEOF {
			expr, err := p.expression()
			if err != nil {
				if !p.synchronize(TkComma, TkCloseBracket) {
					p.current--
					return nil, err
				}
				p.errors = append(p.errors, err)
			}
			values = append(values, expr)
			if !p.match(TkComma) {
				break
			}
		}
		if !p.match(TkCloseBracket) {
			return nil, p.newError("Expected ']' after value list.")
		}
		return &ExprListInitializer{
			OpenBracket:  openBracket,
			CloseBracket: p.previous(),
			Values:       values,
		}, nil
	}

	if p.match(TkLiteral) {
		return &ExprLiteral{
			Token: p.previous(),
		}, nil
	}

	if p.match(TkOpenParen) {
		openParen := p.previous()
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if !p.match(TkCloseParen) {
			return nil, p.newError("Expected ')' after expression.")
		}
		closeParen := p.previous()
		return &ExprGrouping{
			OpenParen:  openParen,
			CloseParen: closeParen,
			Expr:       expr,
		}, nil
	}

	return nil, p.newError(fmt.Sprintf("Unexpected token '%s'", p.peek().Lexeme))
}

func (p *parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.peek().Type == t {
			p.current++
			return true
		}
	}
	return false
}

func (p *parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *parser) peek() Token {
	return p.tokens[p.current]
}

func (p *parser) peekNext() Token {
	return p.tokens[p.current+1]
}

func (p *parser) synchronize(tokens ...TokenType) bool {
	if p.peek().Type == TkEOF {
		return false
	}
	p.current++
	if p.previous().Type == TkNewLine {
		return false
	}
	for p.peek().Type != TkEOF {
		switch p.peek().Type {
		case TkNewLine:
			p.current++
			return false
		}
		for _, t := range tokens {
			if p.peek().Type == t {
				return true
			}
		}
		p.current++
	}
	return false
}

type ParseError struct {
	Token   Token
	Message string
}

func (p ParseError) Error() string {
	return "ERROR: " + p.Message
}

func (p *parser) newError(message string) error {
	return p.newErrorAt(message, p.peek())
}

func (p *parser) newErrorAt(message string, token Token) error {
	return ParseError{
		Token:   token,
		Message: message,
	}
}
