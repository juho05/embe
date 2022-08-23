package debug

import (
	"fmt"
	"strings"

	"github.com/Bananenpro/embe/parser"
)

func PrintAST(statements []parser.Stmt) {
	printer := &astPrinter{}
	for _, stmt := range statements {
		stmt.Accept(printer)
		fmt.Println()
	}
}

type astPrinter struct {
	indent int
}

func (p *astPrinter) VisitEvent(stmt *parser.StmtEvent) error {
	p.print("[ev] %s(%v):\n", stmt.Name.Lexeme, stmt.Parameter.Lexeme)
	p.indent++
	for _, s := range stmt.Body {
		s.Accept(p)
	}
	p.indent--
	return nil
}

func (p *astPrinter) VisitFuncCall(stmt *parser.StmtFuncCall) error {
	p.print("[fn] %s(", stmt.Name.Lexeme)
	for i, param := range stmt.Parameters {
		if i > 0 {
			fmt.Print(", ")
		}
		param.Accept(p)
	}
	fmt.Print(")\n")
	return nil
}

func (p *astPrinter) VisitAssignment(stmt *parser.StmtAssignment) error {
	p.print("[as] %s %s ", stmt.Variable.Lexeme, stmt.Operator.Lexeme)
	stmt.Value.Accept(p)
	fmt.Print("\n")
	return nil
}

func (p *astPrinter) VisitIf(stmt *parser.StmtIf) error {
	p.print("[if] ")
	stmt.Condition.Accept(p)
	fmt.Print(":\n")
	p.indent++
	for _, s := range stmt.Body {
		s.Accept(p)
	}
	p.indent--
	if stmt.ElseBody != nil {
		p.indent++
		fmt.Print("[el]:\n")
		for _, s := range stmt.ElseBody {
			s.Accept(p)
		}
		p.indent--
	}
	return nil
}

func (p *astPrinter) VisitLoop(stmt *parser.StmtLoop) error {
	if stmt.Keyword.Type == parser.TkWhile {
		p.print("[wh] ")
	} else {
		p.print("[fo] ")
	}
	stmt.Condition.Accept(p)
	fmt.Print(":\n")
	p.indent++
	for _, s := range stmt.Body {
		s.Accept(p)
	}
	p.indent--
	return nil
}

func (p *astPrinter) VisitIdentifier(expr *parser.ExprIdentifier) error {
	fmt.Print(expr.Name.Lexeme)
	return nil
}

func (p *astPrinter) VisitLiteral(expr *parser.ExprLiteral) error {
	fmt.Print(expr.Token.Literal)
	return nil
}

func (p *astPrinter) VisitUnary(expr *parser.ExprUnary) error {
	fmt.Print("(-")
	expr.Right.Accept(p)
	fmt.Print(")")
	return nil
}

func (p *astPrinter) VisitBinary(expr *parser.ExprBinary) error {
	fmt.Print("(")
	expr.Left.Accept(p)
	fmt.Print(" " + expr.Operator.Lexeme + " ")
	expr.Right.Accept(p)
	fmt.Print(")")
	return nil
}

func (p *astPrinter) print(format string, a ...any) {
	fmt.Printf("%s%s", strings.Repeat(" ", p.indent), fmt.Sprintf(format, a...))
}
