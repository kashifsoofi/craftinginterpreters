using cslox.Parser;
using cslox.Scanning;

namespace cslox.Interpreter;

class Void
{
    private Void() { }
}

class Interpreter : IExprVisitor<object?>, IStmtVisitor<Void?>
{
    private Environment environment = new Environment();

    public void Interpret(List<Stmt> statements)
    {
        try
        {
            foreach (var statement in statements)
            {
                Execute(statement);
            }
        }
        catch (RuntimeError error)
        {
            Lox.RuntimeError(error);
        }
    }

    public object? VisitBinaryExpr(Binary expr)
    {
        var left = Evaluate(expr.Left);
        var right = Evaluate(expr.Right);

        switch (expr.Operator.Type)
        {
            case TokenType.GREATER:
                CheckNumberOperands(expr.Operator, left, right);
                return (double)left! > (double)right!;
            case TokenType.GREATER_EQUAL:
                CheckNumberOperands(expr.Operator, left, right);
                return (double)left! >= (double)right!;
            case TokenType.LESS:
                CheckNumberOperands(expr.Operator, left, right);
                return (double)left! < (double)right!;
            case TokenType.LESS_EQUAL:
                CheckNumberOperands(expr.Operator, left, right);
                return (double)left! <= (double)right!;
            case TokenType.BANG_EQUAL:
                return !IsEqual(left, right);
            case TokenType.EQUAL_EQUAL:
                return IsEqual(left, right);
            case TokenType.MINUS:
                CheckNumberOperands(expr.Operator, left, right);
                return (double) left! - (double)right!;
            case TokenType.PLUS:
                if (left is double && right is double)
                {
                    return (double)left! + (double)right!;
                }

                if (left is string && right is string)
                {
                    return (string)left + (string)right;
                }
                throw new RuntimeError(expr.Operator, "Operands must be two numbers or two strings");
            case TokenType.SLASH:
                CheckNumberOperands(expr.Operator, left, right);
                if ((double)right! == 0)
                {
                    throw new RuntimeError(expr.Operator, "Division by zero");
                }
                return (double)left! / (double)right!;
            case TokenType.STAR:
                CheckNumberOperands(expr.Operator, left, right);
                return (double)left! * (double)right!;

        }

        return null;
    }

    public object? VisitGroupingExpr(Grouping expr)
    {
        return Evaluate(expr.Expression);
    }

    public object? VisitLiteralExpr(Literal expr)
    {
        return expr.Value;
    }

    public object? VisitUnaryExpr(Unary expr)
    {
        var right = Evaluate(expr.Right);
        switch (expr.Operator.Type)
        {
            case TokenType.BANG:
                return !IsTruthy(right);
            case TokenType.MINUS:
                CheckNumberOperand(expr.Operator, right);
                return -(double)right!;
        }

        return null;
    }

    public object? VisitVariableExpr(Variable expr)
    {
        return environment.Get(expr.Name);
    }

    public Void? VisitExpressionStmt(ExpressionStmt stmt)
    {
        Evaluate(stmt.Expression);
        return null;
    }

    public Void? VisitPrintStmt(Print stmt)
    {
        var value = Evaluate(stmt.Expression);
        Console.WriteLine(Stringify(value));
        return null;
    }

    public Void? VisitVarStmt(Var stmt)
    {
        object? value = null;
        if (stmt.Initializer != null)
        {
            value = Evaluate(stmt.Initializer);
        }

        environment.Define(stmt.Name.Lexeme, value);
        return null;
    }

    private object? Evaluate(Expr expr)
    {
        return expr.Accept(this);
    }

    private void Execute(Stmt stmt)
    {
        stmt.Accept(this);
    }

    private bool IsTruthy(object? obj)
    {
        if (obj == null)
        {
            return false;
        }

        if (obj is bool)
        {
            return (bool)obj;
        }

        return true;
    }

    private bool IsEqual(object? a, object? b)
    {
        if (a == null && b == null)
        {
            return true;
        }
        if (a == null)
        {
            return false;
        }

        return a.Equals(b);
    }

    private void CheckNumberOperand(Token @operator, object? operand)
    {
        if (operand is double)
        {
            return;
        }

        throw new RuntimeError(@operator, "Operand must be a number.");
    }

    private void CheckNumberOperands(Token @operator, object? left, object? right)
    {
        if (left is double && right is double)
        {
            return;
        }

        throw new RuntimeError(@operator, "Operands must be numbers.");
    }

    private string Stringify(object? value)
    {
        if (value == null)
        {
            return "nil";
        }

        if (value is double)
        {
            var text = ((double)value).ToString();
            if (text.EndsWith(".0"))
            {
                text = text.Substring(0, text.Length - 2);
            }
            return text;
        }

        return value.ToString()!;
    }
}

class RuntimeError : Exception
{
    public Token Token { get; }

    public RuntimeError(Token token, string message)
        : base(message)
    {
        Token = token;
    }
}
