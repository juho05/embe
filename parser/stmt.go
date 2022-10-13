package parser

type StmtVisitor interface {
	VisitVarDecl(stmt *StmtVarDecl) error
	VisitConstDecl(stmt *StmtConstDecl) error
	VisitFuncDecl(stmt *StmtFuncDecl) error
	VisitEvent(stmt *StmtEvent) error
	VisitFuncCall(stmt *StmtFuncCall) error
	VisitAssignment(stmt *StmtAssignment) error
	VisitIf(stmt *StmtIf) error
	VisitLoop(stmt *StmtLoop) error
}

type Stmt interface {
	Accept(visitor StmtVisitor) error
	Position() (start, end Position)
}

type StmtVarDecl struct {
	Name        Token
	DataType    DataType
	AssignToken Token
	Value       Expr
}

func (s *StmtVarDecl) Accept(visitor StmtVisitor) error {
	return visitor.VisitVarDecl(s)
}

func (s *StmtVarDecl) Position() (start, end Position) {
	start = s.Name.Pos
	_, end = s.Value.Position()
	return start, end
}

type StmtConstDecl struct {
	Name        Token
	AssignToken Token
	Value       Token
}

func (s *StmtConstDecl) Accept(visitor StmtVisitor) error {
	return visitor.VisitConstDecl(s)
}

func (s *StmtConstDecl) Position() (start, end Position) {
	start = s.Name.Pos
	end = s.Value.Pos
	end.Column += len(s.Value.Lexeme)
	return start, end
}

type FuncParam struct {
	Name Token
	Type Token
}

type StmtFuncDecl struct {
	Name       Token
	CloseParen Token
	Params     []FuncParam
	Body       []Stmt
	StartLine  int
	EndLine    int
}

func (s *StmtFuncDecl) Accept(visitor StmtVisitor) error {
	return visitor.VisitFuncDecl(s)
}

func (s *StmtFuncDecl) Position() (start, end Position) {
	return s.Name.Pos, s.CloseParen.Pos
}

type StmtEvent struct {
	At        Token
	Name      Token
	Parameter Token
	Body      []Stmt
}

func (s *StmtEvent) Accept(visitor StmtVisitor) error {
	return visitor.VisitEvent(s)
}

func (s *StmtEvent) Position() (start, end Position) {
	end = s.Name.Pos
	end.Column += len(s.Name.Lexeme)
	return s.At.Pos, end
}

type StmtFuncCall struct {
	Name       Token
	CloseParen Token
	Parameters []Expr
}

func (s *StmtFuncCall) Accept(visitor StmtVisitor) error {
	return visitor.VisitFuncCall(s)
}

func (s *StmtFuncCall) Position() (start, end Position) {
	return s.Name.Pos, end
}

type StmtAssignment struct {
	Variable Token
	Operator Token
	Value    Expr
}

func (s *StmtAssignment) Accept(visitor StmtVisitor) error {
	return visitor.VisitAssignment(s)
}

func (s *StmtAssignment) Position() (start, end Position) {
	_, end = s.Value.Position()
	return s.Variable.Pos, end
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

func (s *StmtIf) Position() (start, end Position) {
	_, end = s.Condition.Position()
	return s.Keyword.Pos, end
}

type StmtLoop struct {
	Keyword   Token
	Condition Expr
	Body      []Stmt
}

func (s *StmtLoop) Accept(visitor StmtVisitor) error {
	return visitor.VisitLoop(s)
}

func (s *StmtLoop) Position() (start, end Position) {
	_, end = s.Condition.Position()
	return s.Keyword.Pos, end
}
