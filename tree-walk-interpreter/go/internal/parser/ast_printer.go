package parser

import (
	"fmt"
	"strings"

	"github.com/kashifsoofi/go-lox/internal/scanner"
)

type AstPrinter struct{}

func (p *AstPrinter) VisitAssignExpr(expr *Assign) interface{} {
	return p.parenthesize2("=", expr.Name.Lexeme, expr.Value)
}

func (p *AstPrinter) VisitBinaryExpr(expr *Binary) interface{} {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *AstPrinter) VisitCallExpr(expr *Call) interface{} {
	return p.parenthesize2("call", expr.Callee, expr.Arguments)
}

func (p *AstPrinter) VisitGetExpr(expr *Get) interface{} {
	return p.parenthesize2(".", expr.Object, expr.Name.Lexeme)
}

func (p *AstPrinter) VisitGroupingExpr(expr *Grouping) interface{} {
	return p.parenthesize("group", expr.Expression)
}

func (p *AstPrinter) VisitLiteralExpr(expr *Literal) interface{} {
	if expr.Value == nil {
		return "nil"
	}

	return fmt.Sprintf("%v", expr.Value)
}

func (p *AstPrinter) VisitLogicalExpr(expr *Logical) interface{} {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *AstPrinter) VisitSetExpr(expr *Set) interface{} {
	return p.parenthesize2("=", expr.Object, expr.Name.Lexeme, expr.Value)
}

func (p *AstPrinter) VisitSuperExpr(expr *Super) interface{} {
	return p.parenthesize2("super", expr.Method)
}

func (p *AstPrinter) VisitThisExpr(expr *This) interface{} {
	return "this"
}

func (p *AstPrinter) VisitUnaryExpr(expr *Unary) interface{} {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p *AstPrinter) VisitVariableExpr(expr *Variable) interface{} {
	return expr.Name.Lexeme
}

func (p *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder
	builder.WriteString("(")
	for _, expr := range exprs {
		builder.WriteString(" ")
		v, _ := expr.Accept(p).(string)
		builder.WriteString(v)
	}
	builder.WriteString(")")

	return builder.String()
}

func (p *AstPrinter) parenthesize2(name string, parts ...interface{}) string {
	var builder strings.Builder
	builder.WriteString("(")
	builder.WriteString(name)
	p.transform(builder, parts)
	builder.WriteString(")")

	return builder.String()
}

func (p *AstPrinter) transform(builder strings.Builder, parts ...interface{}) {
	for _, part := range parts {
		builder.WriteString(" ")
		if expr, ok := part.(Expr); ok {
			v, _ := expr.Accept(p).(string)
			builder.WriteString(v)
		} else if _, ok := part.(Stmt); ok {
			// TODO
		} else if token, ok := part.(*scanner.Token); ok {
			builder.WriteString(token.Lexeme)
		} else {
			v, _ := part.(string)
			builder.WriteString(v)
		}
	}
}
