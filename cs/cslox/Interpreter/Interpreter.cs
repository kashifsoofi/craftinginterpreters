using cslox.Parser;
using cslox.Scanning;

namespace cslox.Interpreter;

class Void
{
    private Void() { }
}

class Interpreter : IExprVisitor<object?>, IStmtVisitor<Void?>
{
    public Environment Globals { get; } = new Environment();
    private Environment environment;
    private readonly Dictionary<Expr, int> locals = new Dictionary<Expr, int>();

    public Interpreter()
    {
        Globals.Define("clock", new Clock());
        environment = Globals;
    }

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

    public void Resolve(Expr expr, int depth)
    {
        locals[expr] = depth;
    }

    public object? VisitAssignExpr(Assign expr)
    {
        var value = Evaluate(expr.Value);
        if (locals.TryGetValue(expr, out var distance))
        {
            environment.AssignAt(distance, expr.Name, value);
        }
        else
        {
            Globals.Assign(expr.Name, value);
        }
        return value;
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

    public object? VisitCallExpr(Call expr)
    {
        var callee = Evaluate(expr.Callee);

        var arguments = new List<object>();
        foreach (var argument in expr.Arguments)
        {
            var argumentValue = Evaluate(argument);
            if (argumentValue != null)
            {
                arguments.Add(argumentValue);
            }
        }

        if (callee is not ILoxCallable)
        {
            throw new RuntimeError(expr.Paren, "Can only call functions and classes.");
        }
        var function = (ILoxCallable)callee!;
        if (arguments.Count != function.Arity())
        {
            throw new RuntimeError(expr.Paren, $"Expected {function.Arity()} arguments but got {arguments.Count}.");
        }

        return function.Call(this, arguments);
    }

    public object? VisitGroupingExpr(Grouping expr)
    {
        return Evaluate(expr.Expression);
    }

    public object? VisitLogicalExpr(Logical expr)
    {
        var left = Evaluate(expr.Left);

        if (expr.Operator.Type == TokenType.OR)
        {
            if (IsTruthy(left))
            {
                return left;
            }
        }
        else
        {
            if (!IsTruthy(left))
            {
                return left;
            }
        }

        return Evaluate(expr.Right);
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
        return LookupVariable(expr.Name, expr);
    }

    public Void? VisitBlockStmt(Block stmt)
    {
        ExecuteBlock(stmt.Statements, new Environment(environment));
        return null;
    }

    public Void? VisitExpressionStmt(ExpressionStmt stmt)
    {
        Evaluate(stmt.Expression);
        return null;
    }

    public Void? VisitFunctionStmt(Function stmt)
    {
        var function = new LoxFunction(stmt, environment);
        environment.Define(stmt.Name.Lexeme, function);
        return null;
    }

    public Void? VisitIfStmt(If stmt)
    {
        if (IsTruthy(Evaluate(stmt.Condition)))
        {
            Execute(stmt.ThenBranch);
        }
        else if (stmt.ElseBranch != null)
        {
            Execute(stmt.ElseBranch);
        }
        return null;
    }

    public Void? VisitPrintStmt(Print stmt)
    {
        var value = Evaluate(stmt.Expression);
        Console.WriteLine(Stringify(value));
        return null;
    }

    public Void? VisitReturnStmt(Return stmt)
    {
        object? value = null;
        if (stmt.Value != null)
        {
            value = Evaluate(stmt.Value);
        }

        throw new ReturnValue(value);
    }

    public Void? VisitWhileStmt(While stmt)
    {
        while (IsTruthy(stmt.Condition))
        {
            Execute(stmt.Body);
        }
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

    public void ExecuteBlock(List<Stmt> statements, Environment environment)
    {
        var previous = this.environment;
        try
        {
            this.environment = environment;

            foreach (var statement in statements)
            {
                Execute(statement);
            }
        }
        finally
        {
            this.environment = previous;
        }
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

    private object? LookupVariable(Token name, Expr expr)
    {
        if (locals.TryGetValue(expr, out var distance))
        {
            return environment.GetAt(distance, name.Lexeme);
        }

        return Globals.Get(name);
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

class ReturnValue : Exception
{
    public object? Value { get; }

    public ReturnValue(object? value)
    {
        Value = value;
    }
}