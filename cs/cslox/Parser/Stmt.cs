using cslox.Scanning;

namespace cslox.Parser;

interface IStmtVisitor<T>
{
	T VisitExpressionStmt(ExpressionStmt stmt);
	T VisitPrintStmt(Print stmt);
	T VisitVarStmt(Var stmt);
}

abstract class Stmt
{
	public abstract T Accept<T>(IStmtVisitor<T> visitor);
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
