using System;
using System.Text;
using cslox.Scanning;

namespace cslox.Parser;

class AstPrinter : IExprVisitor<string>
{
    public string Print(Expr expr)
    {
        return expr.Accept(this);
    }

    public string VisitAssignExpr(Assign expr)
    {
        return Parenthesize2("=", expr.Name.Lexeme, expr.Value);
    }

    public string VisitBinaryExpr(Binary expr)
    {
        return Parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right);
    }

    public string VisitCallExpr(Call expr)
    {
        return Parenthesize2("call", expr.Callee, expr.Arguments);
    }

    public string VisitGroupingExpr(Grouping expr)
    {
        return Parenthesize("group", expr.Expression);
    }

    public string VisitLogicalExpr(Logical expr)
    {
        return Parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right);
    }

    public string VisitLiteralExpr(Literal expr)
    {
        if (expr.Value == null)
        {
            return "nil";
        }
        return expr.Value.ToString()!;
    }

    public string VisitUnaryExpr(Unary expr)
    {
        return Parenthesize(expr.Operator.Lexeme, expr.Right);
    }

    public string VisitVariableExpr(Variable expr)
    {
        return expr.Name.Lexeme;
    }

    private string Parenthesize(string name, params Expr[] exprs)
    {
        var builder = new StringBuilder();
        builder.Append("(").Append(name);
        foreach (var expr in exprs)
        {
            builder.Append(" ");
            builder.Append(expr.Accept(this));
        }
        builder.Append(")");

        return builder.ToString();
    }

    private string Parenthesize2(string name, params object[] parts)
    {
        var builder = new StringBuilder();
        builder.Append("(").Append(name);
        Transform(builder, parts);
        builder.Append(")");

        return builder.ToString();
    }

    private void Transform(StringBuilder builder, params object[] parts)
    {
        foreach (var part in parts)
        {
            builder.Append(" ");
            if (part is Expr)
            {
                builder.Append(((Expr)part).Accept(this));
            }
            else if (part is Stmt)
            {
                // builder.Append(((Stmt)part).Accept<string>(this));
            }
            else if (part is Token)
            {
                builder.Append(((Token)part).Lexeme);
            }
            //else if (part is List)
            //{
            //    builder.Append((List))part).ToArray());
            //}
            else
            {
                builder.Append(part);
            }
        }
    }
}

