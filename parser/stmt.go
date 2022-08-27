package parser

type StmtVisitor interface {
	VisitVarDecl(stmt *StmtVarDecl) error
	VisitEvent(stmt *StmtEvent) error
	VisitFuncCall(stmt *StmtFuncCall) error
	VisitAssignment(stmt *StmtAssignment) error
	VisitIf(stmt *StmtIf) error
	VisitLoop(stmt *StmtLoop) error
}

type Stmt interface {
	Accept(visitor StmtVisitor) error
}

type StmtVarDecl struct {
	Name        Token
	DataType    DataType
	Initializer *StmtAssignment
}

func (s *StmtVarDecl) Accept(visitor StmtVisitor) error {
	return visitor.VisitVarDecl(s)
}

type StmtEvent struct {
	Name      Token
	Parameter Token
	Body      []Stmt
}

func (s *StmtEvent) Accept(visitor StmtVisitor) error {
	return visitor.VisitEvent(s)
}

type StmtFuncCall struct {
	Name       Token
	Parameters []Expr
}

func (s *StmtFuncCall) Accept(visitor StmtVisitor) error {
	return visitor.VisitFuncCall(s)
}

type StmtAssignment struct {
	Variable Token
	Operator Token
	Value    Expr
}

func (s *StmtAssignment) Accept(visitor StmtVisitor) error {
	return visitor.VisitAssignment(s)
}

type StmtIf struct {
	Keyword   Token
	Condition Expr
	Body      []Stmt
	ElseBody  []Stmt
}

func (s *StmtIf) Accept(visitor StmtVisitor) error {
	return visitor.VisitIf(s)
}

type StmtLoop struct {
	Keyword   Token
	Condition Expr
	Body      []Stmt
}

func (s *StmtLoop) Accept(visitor StmtVisitor) error {
	return visitor.VisitLoop(s)
}
