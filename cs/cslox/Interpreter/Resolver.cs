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
        Initializer,
        Method,
    }

    private enum ClassType
    {
        None,
        Class,
        Subclass,
    }

	private readonly Interpreter interpreter;
    private readonly Stack<Dictionary<string, bool>> scopes = new Stack<Dictionary<string, bool>>();
    private FunctionType currentFunctionType = FunctionType.None;
    private ClassType currentClassType = ClassType.None;

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

    public Void? VisitGetExpr(Get expr)
    {
        Resolve(expr.Object);
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

    public Void? VisitSetExpr(Set expr)
    {
        Resolve(expr.Value);
        Resolve(expr.Object);
        return null;
    }

    public Void? VisitSuperExpr(Super expr)
    {
        if (currentClassType == ClassType.None)
        {
            Lox.Error(expr.Keyword, "Can't use 'super' outside of a class.");
        }
        else if (currentClassType != ClassType.Subclass)
        {
            Lox.Error(expr.Keyword, "Can't use 'super' in a class with no superclass.");
        }

        ResolveLocal(expr, expr.Keyword);
        return null;
    }

    public Void? VisitThisExpr(This expr)
    {
        if (currentClassType == ClassType.None)
        {
            Lox.Error(expr.Keyword, "Can't use 'this' outside of a class.");
            return null;
        }

        ResolveLocal(expr, expr.Keyword);
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
        var enclosingClassType = currentClassType;
        currentClassType = ClassType.Class;

        Declare(stmt.Name);
        Define(stmt.Name);

        if (stmt.Superclass != null &&
            stmt.Name.Lexeme == stmt.Superclass.Name.Lexeme)
        {
            Lox.Error(stmt.Superclass.Name, "A class can't inherit from itself.");
        }

        if (stmt.Superclass != null)
        {
            currentClassType = ClassType.Subclass;
            Resolve(stmt.Superclass);
        }

        if (stmt.Superclass != null)
        {
            BeginScope();
            scopes.Peek()["super"] = true;
        }

        BeginScope();
        scopes.Peek()["this"] = true;

        foreach (var method in stmt.Methods)
        {
            var declaration = FunctionType.Method;
            if (method.Name.Lexeme == "this")
            {
                declaration = FunctionType.Initializer;
            }
            ResolveFunction(method, declaration);
        }

        EndScope();

        if (stmt.Superclass != null)
        {
            EndScope();
        }

        currentClassType = enclosingClassType;
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
            if (currentFunctionType == FunctionType.Initializer)
            {
                Lox.Error(stmt.Keyword, "Can't return a value from an initializer.");
            }

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

