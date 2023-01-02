using cslox.Scanning;

namespace cslox.Parser;

interface IExprVisitor<T>
{
	T VisitBinaryExpr(Binary expr);
	T VisitGroupingExpr(Grouping expr);
	T VisitLiteralExpr(Literal expr);
	T VisitUnaryExpr(Unary expr);
	T VisitVariableExpr(Variable expr);
}

abstract class Expr
{
	public abstract T Accept<T>(IExprVisitor<T> visitor);
}

class Binary : Expr
{
	public Binary(Expr left, Token @operator, Expr right)
	{
		Left = left;
		Operator = @operator;
		Right = right;
	}

	public Expr Left { get; }
	public Token Operator { get; }
	public Expr Right { get; }

	public override T Accept<T>(IExprVisitor<T> visitor)
	{
		return visitor.VisitBinaryExpr(this);
	}
}

class Grouping : Expr
{
	public Grouping(Expr expression)
	{
		Expression = expression;
	}

	public Expr Expression { get; }

	public override T Accept<T>(IExprVisitor<T> visitor)
	{
		return visitor.VisitGroupingExpr(this);
	}
}

class Literal : Expr
{
	public Literal(object? value)
	{
		Value = value;
	}

	public object? Value { get; }

	public override T Accept<T>(IExprVisitor<T> visitor)
	{
		return visitor.VisitLiteralExpr(this);
	}
}

class Unary : Expr
{
	public Unary(Token @operator, Expr right)
	{
		Operator = @operator;
		Right = right;
	}

	public Token Operator { get; }
	public Expr Right { get; }

	public override T Accept<T>(IExprVisitor<T> visitor)
	{
		return visitor.VisitUnaryExpr(this);
	}
}

class Variable : Expr
{
	public Variable(Token name)
	{
		Name = name;
	}

	public Token Name { get; }

	public override T Accept<T>(IExprVisitor<T> visitor)
	{
		return visitor.VisitVariableExpr(this);
	}
}
