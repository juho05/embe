package parser

type ExprVisitor interface {
	VisitIdentifier(expr *ExprIdentifier) error
	VisitExprFuncCall(expr *ExprFuncCall) error
	VisitTypeCast(expr *ExprTypeCast) error
	VisitLiteral(expr *ExprLiteral) error
	VisitListInitializer(expr *ExprListInitializer) error
	VisitUnary(expr *ExprUnary) error
	VisitBinary(expr *ExprBinary) error
}

type Expr interface {
	Accept(visitor ExprVisitor) error
	Type() DataType
}

type ExprFuncCall struct {
	Name       Token
	Parameters []Expr
	ReturnType DataType
}

func (e *ExprFuncCall) Accept(visitor ExprVisitor) error {
	return visitor.VisitExprFuncCall(e)
}

func (e *ExprFuncCall) Type() DataType {
	return e.ReturnType
}

type ExprTypeCast struct {
	Target     Token
	Value      Expr
	ReturnType DataType
}

func (e *ExprTypeCast) Accept(visitor ExprVisitor) error {
	return visitor.VisitTypeCast(e)
}

func (e *ExprTypeCast) Type() DataType {
	return e.ReturnType
}

type ExprIdentifier struct {
	Name       Token
	ReturnType DataType
}

func (e *ExprIdentifier) Accept(visitor ExprVisitor) error {
	return visitor.VisitIdentifier(e)
}

func (e *ExprIdentifier) Type() DataType {
	return e.ReturnType
}

type ExprLiteral struct {
	Token      Token
	ReturnType DataType
}

func (e *ExprLiteral) Accept(visitor ExprVisitor) error {
	return visitor.VisitLiteral(e)
}

func (e *ExprLiteral) Type() DataType {
	return e.ReturnType
}

type ExprListInitializer struct {
	OpenBracket Token
	Values      []Token
	ReturnType  DataType
}

func (e *ExprListInitializer) Accept(visitor ExprVisitor) error {
	return visitor.VisitListInitializer(e)
}

func (e *ExprListInitializer) Type() DataType {
	return e.ReturnType
}

type ExprUnary struct {
	Operator   Token
	Right      Expr
	ReturnType DataType
}

func (e *ExprUnary) Accept(visitor ExprVisitor) error {
	return visitor.VisitUnary(e)
}

func (e *ExprUnary) Type() DataType {
	return e.ReturnType
}

type ExprBinary struct {
	Operator   Token
	Left       Expr
	Right      Expr
	ReturnType DataType
}

func (e *ExprBinary) Accept(visitor ExprVisitor) error {
	return visitor.VisitBinary(e)
}

func (e *ExprBinary) Type() DataType {
	return e.ReturnType
}
