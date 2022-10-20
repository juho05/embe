package parser

type ExprVisitor interface {
	VisitIdentifier(expr *ExprIdentifier) error
	VisitExprFuncCall(expr *ExprFuncCall) error
	VisitTypeCast(expr *ExprTypeCast) error
	VisitLiteral(expr *ExprLiteral) error
	VisitListInitializer(expr *ExprListInitializer) error
	VisitUnary(expr *ExprUnary) error
	VisitBinary(expr *ExprBinary) error
	VisitGrouping(expr *ExprGrouping) error
}

type Expr interface {
	Accept(visitor ExprVisitor) error
	Type() DataType
	Position() (start, end Position)
}

type ExprFuncCall struct {
	Name       Token
	Parameters []Expr
	ReturnType DataType
	CloseParen Token
}

func (e *ExprFuncCall) Accept(visitor ExprVisitor) error {
	return visitor.VisitExprFuncCall(e)
}

func (e *ExprFuncCall) Type() DataType {
	return e.ReturnType
}

func (e *ExprFuncCall) Position() (start, end Position) {
	return e.Name.Pos, e.CloseParen.Pos
}

type ExprTypeCast struct {
	Target     Token
	Value      Expr
	ReturnType DataType
	CloseParen Token
}

func (e *ExprTypeCast) Accept(visitor ExprVisitor) error {
	return visitor.VisitTypeCast(e)
}

func (e *ExprTypeCast) Type() DataType {
	return e.ReturnType
}

func (e *ExprTypeCast) Position() (start, end Position) {
	return e.Target.Pos, e.CloseParen.Pos
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

func (e *ExprIdentifier) Position() (start, end Position) {
	return e.Name.Pos, e.Name.EndPos
}

type ExprLiteral struct {
	Token      Token
	End        Position
	ReturnType DataType
}

func (e *ExprLiteral) Accept(visitor ExprVisitor) error {
	return visitor.VisitLiteral(e)
}

func (e *ExprLiteral) Type() DataType {
	return e.ReturnType
}

func (e *ExprLiteral) Position() (start, end Position) {
	if e.End.Line == 0 && e.End.Column == 0 {
		end = e.Token.EndPos
	} else {
		end = e.End
	}
	return e.Token.Pos, end
}

type ExprListInitializer struct {
	OpenBracket  Token
	CloseBracket Token
	Values       []Expr
	ReturnType   DataType
}

func (e *ExprListInitializer) Accept(visitor ExprVisitor) error {
	return visitor.VisitListInitializer(e)
}

func (e *ExprListInitializer) Type() DataType {
	return e.ReturnType
}

func (e *ExprListInitializer) Position() (start, end Position) {
	return e.OpenBracket.Pos, e.CloseBracket.Pos
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

func (e *ExprUnary) Position() (start, end Position) {
	_, end = e.Right.Position()
	return e.Operator.Pos, end
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

func (e *ExprBinary) Position() (start, end Position) {
	start, _ = e.Left.Position()
	_, end = e.Right.Position()
	return start, end
}

type ExprGrouping struct {
	OpenParen  Token
	CloseParen Token
	Expr       Expr
}

func (e *ExprGrouping) Accept(visitor ExprVisitor) error {
	return visitor.VisitGrouping(e)
}

func (e *ExprGrouping) Type() DataType {
	return e.Expr.Type()
}

func (e *ExprGrouping) Position() (start, end Position) {
	return e.OpenParen.Pos, e.CloseParen.Pos
}
