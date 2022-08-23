package parser

type ExprVisitor interface {
	VisitIdentifier(expr *ExprIdentifier) error
	VisitLiteral(expr *ExprLiteral) error
	VisitUnary(expr *ExprUnary) error
	VisitBinary(expr *ExprBinary) error
}

type Expr interface {
	Accept(visitor ExprVisitor) error
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
