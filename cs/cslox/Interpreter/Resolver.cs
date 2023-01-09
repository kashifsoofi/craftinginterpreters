using System;
using cslox.Parser;
using cslox.Scanning;

namespace cslox.Interpreter;

class Resolver : IExprVisitor<Void?>, IStmtVisitor<Void?>
{
    private enum FunctionType
    {
        None,
        Function,
    }

	private readonly Interpreter interpreter;
    private readonly Stack<Dictionary<string, bool>> scopes = new Stack<Dictionary<string, bool>>();
    private FunctionType currentFunctionType = FunctionType.None;

	public Resolver(Interpreter interpreter)
	{
		this.interpreter = interpreter;
	}

    public void Resolve(List<Stmt> statements)
    {
        foreach (var statement in statements)
        {
            Resolve(statement);
        }
    }

    public Void? VisitAssignExpr(Assign expr)
    {
        Resolve(expr.Value);
        ResolveLocal(expr, expr.Name);
        return null;
    }

    public Void? VisitBinaryExpr(Binary expr)
    {
        Resolve(expr.Left);
        Resolve(expr.Right);
        return null;
    }

    public Void? VisitCallExpr(Call expr)
    {
        Resolve(expr.Callee);

        foreach (var argument in expr.Arguments)
        {
            Resolve(argument);
        }

        return null;
    }

    public Void? VisitGroupingExpr(Grouping expr)
    {
        Resolve(expr.Expression);
        return null;
    }

    public Void? VisitLiteralExpr(Literal expr)
    {
        return null;
    }

    public Void? VisitLogicalExpr(Logical expr)
    {
        Resolve(expr.Left);
        Resolve(expr.Right);
        return null;
    }

    public Void? VisitUnaryExpr(Unary expr)
    {
        Resolve(expr.Right);
        return null;
    }

    public Void? VisitVariableExpr(Variable expr)
    {
        if (scopes.Count > 0 && scopes.Peek()[expr.Name.Lexeme] == false)
        {
            Lox.Error(expr.Name, "Can't read local variable in its own initializer.");
        }

        ResolveLocal(expr, expr.Name);
        return null;
    }

    public Void? VisitBlockStmt(Block stmt)
    {
        BeginScope();
        Resolve(stmt.Statements);
        EndScope();
        return null;
    }

    public Void? VisitClassStmt(Class stmt)
    {
        Declare(stmt.Name);
        Define(stmt.Name);
        return null;
    }

    public Void? VisitExpressionStmt(ExpressionStmt stmt)
    {
        Resolve(stmt.Expression);
        return null;
    }

    public Void? VisitFunctionStmt(Function stmt)
    {
        Declare(stmt.Name);
        Define(stmt.Name);

        ResolveFunction(stmt, FunctionType.Function);
        return null;
    }

    public Void? VisitIfStmt(If stmt)
    {
        Resolve(stmt.Condition);
        Resolve(stmt.ThenBranch);
        if (stmt.ElseBranch != null)
        {
            Resolve(stmt.ElseBranch);
        }
        return null;
    }

    public Void? VisitPrintStmt(Print stmt)
    {
        Resolve(stmt.Expression);
        return null;
    }

    public Void? VisitReturnStmt(Return stmt)
    {
        if (currentFunctionType == FunctionType.None)
        {
            Lox.Error(stmt.Keyword, "Can't return from top-level code.");
        }

        if (stmt.Value != null)
        {
            Resolve(stmt.Value);
        }

        return null;
    }

    public Void? VisitVarStmt(Var stmt)
    {
        Declare(stmt.Name);
        if (stmt.Initializer != null)
        {
            Resolve(stmt.Initializer);
        }
        Define(stmt.Name);
        return null;
    }

    public Void? VisitWhileStmt(While stmt)
    {
        Resolve(stmt.Condition);
        Resolve(stmt.Body);
        return null;
    }

    private void Resolve(Stmt stmt)
    {
        stmt.Accept(this);
    }

    private void Resolve(Expr expr)
    {
        expr.Accept(this);
    }

    private void BeginScope()
    {
        scopes.Push(new Dictionary<string, bool>());
    }

    private void EndScope()
    {
        scopes.Pop();
    }

    private void Declare(Token name)
    {
        if (scopes.Count == 0)
        {
            return;
        }

        var scope = scopes.Peek();
        if (scope.ContainsKey(name.Lexeme))
        {
            Lox.Error(name, "Already a vairable with this name in this scope.");
        }

        scope[name.Lexeme] = false;
    }

    private void Define(Token name)
    {
        if (scopes.Count == 0)
        {
            return;
        }
        scopes.Peek()[name.Lexeme] = true;
    }

    private void ResolveLocal(Expr expr, Token name)
    {
        for (var i = scopes.Count - 1; i >= 0; i--)
        {
            if (scopes.ElementAt(i).ContainsKey(name.Lexeme))
            {
                interpreter.Resolve(expr, scopes.Count - 1 - i);
                return;
            }
        }
    }

    private void ResolveFunction(Function function, FunctionType functionType)
    {
        var enclosingFunctionType = currentFunctionType;
        currentFunctionType = functionType;

        BeginScope();
        foreach (var param in function.Parameters)
        {
            Declare(param);
            Define(param);
        }
        Resolve(function.Body);
        EndScope();

        currentFunctionType = enclosingFunctionType;
    }
}

