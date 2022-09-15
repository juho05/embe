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
}

type ExprFuncCall struct {
	Name       Token
	Parameters []Expr
}

func (e *ExprFuncCall) Accept(visitor ExprVisitor) error {
	return visitor.VisitExprFuncCall(e)
}

type ExprTypeCast struct {
	Type  Token
	Value Expr
}

func (e *ExprTypeCast) Accept(visitor ExprVisitor) error {
	return visitor.VisitTypeCast(e)
}

type ExprIdentifier struct {
	Name Token
}

func (e *ExprIdentifier) Accept(visitor ExprVisitor) error {
	return visitor.VisitIdentifier(e)
}

type ExprLiteral struct {
	Token Token
}

func (e *ExprLiteral) Accept(visitor ExprVisitor) error {
	return visitor.VisitLiteral(e)
}

type ExprListInitializer struct {
	OpenBracket Token
	Values      []Token
}

func (e *ExprListInitializer) Accept(visitor ExprVisitor) error {
	return visitor.VisitListInitializer(e)
}

type ExprUnary struct {
	Operator Token
	Right    Expr
}

func (e *ExprUnary) Accept(visitor ExprVisitor) error {
	return visitor.VisitUnary(e)
}

type ExprBinary struct {
	Operator Token
	Left     Expr
	Right    Expr
}

func (e *ExprBinary) Accept(visitor ExprVisitor) error {
	return visitor.VisitBinary(e)
}
