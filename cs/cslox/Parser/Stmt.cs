using cslox.Scanning;

namespace cslox.Parser;

interface IStmtVisitor<T>
{
	T VisitBlockStmt(Block stmt);
	T VisitExpressionStmt(ExpressionStmt stmt);
	T VisitFunctionStmt(Function stmt);
	T VisitIfStmt(If stmt);
	T VisitPrintStmt(Print stmt);
	T VisitReturnStmt(Return stmt);
	T VisitVarStmt(Var stmt);
	T VisitWhileStmt(While stmt);
}

abstract class Stmt
{
	public abstract T Accept<T>(IStmtVisitor<T> visitor);
}

class Block : Stmt
{
	public Block(List<Stmt> statements)
	{
		Statements = statements;
	}

	public List<Stmt> Statements { get; }

	public override T Accept<T>(IStmtVisitor<T> visitor)
	{
		return visitor.VisitBlockStmt(this);
	}
}

class ExpressionStmt : Stmt
{
	public ExpressionStmt(Expr expression)
	{
		Expression = expression;
	}

	public Expr Expression { get; }

	public override T Accept<T>(IStmtVisitor<T> visitor)
	{
		return visitor.VisitExpressionStmt(this);
	}
}

class Function : Stmt
{
	public Function(Token name, List<Token> parameters, List<Stmt> body)
	{
		Name = name;
		Parameters = parameters;
		Body = body;
	}

	public Token Name { get; }
	public List<Token> Parameters { get; }
	public List<Stmt> Body { get; }

	public override T Accept<T>(IStmtVisitor<T> visitor)
	{
		return visitor.VisitFunctionStmt(this);
	}
}

class If : Stmt
{
	public If(Expr condition, Stmt thenbranch, Stmt? elsebranch)
	{
		Condition = condition;
		ThenBranch = thenbranch;
		ElseBranch = elsebranch;
	}

	public Expr Condition { get; }
	public Stmt ThenBranch { get; }
	public Stmt? ElseBranch { get; }

	public override T Accept<T>(IStmtVisitor<T> visitor)
	{
		return visitor.VisitIfStmt(this);
	}
}

class Print : Stmt
{
	public Print(Expr expression)
	{
		Expression = expression;
	}

	public Expr Expression { get; }

	public override T Accept<T>(IStmtVisitor<T> visitor)
	{
		return visitor.VisitPrintStmt(this);
	}
}

class Return : Stmt
{
	public Return(Token keyword, Expr? value)
	{
		Keyword = keyword;
		Value = value;
	}

	public Token Keyword { get; }
	public Expr? Value { get; }

	public override T Accept<T>(IStmtVisitor<T> visitor)
	{
		return visitor.VisitReturnStmt(this);
	}
}

class Var : Stmt
{
	public Var(Token name, Expr? initializer)
	{
		Name = name;
		Initializer = initializer;
	}

	public Token Name { get; }
	public Expr? Initializer { get; }

	public override T Accept<T>(IStmtVisitor<T> visitor)
	{
		return visitor.VisitVarStmt(this);
	}
}

class While : Stmt
{
	public While(Expr condition, Stmt body)
	{
		Condition = condition;
		Body = body;
	}

	public Expr Condition { get; }
	public Stmt Body { get; }

	public override T Accept<T>(IStmtVisitor<T> visitor)
	{
		return visitor.VisitWhileStmt(this);
	}
}
