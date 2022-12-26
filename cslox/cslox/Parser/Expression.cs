using cslox.Scanning;

namespace cslox.Parser;

abstract class Expression
{
	public abstract T Accept<T>(IVisitor<T> visitor);
}

class Binary : Expression
{
	public Binary(Expression left, Token @operator, Expression right)
	{
		Left = left;
		Operator = @operator;
		Right = right;
	}

	public Expression Left { get; }
	public Token Operator { get; }
	public Expression Right { get; }

	public override T Accept<T>(IVisitor<T> visitor)
	{
		return visitor.VisitBinaryExpression(this);
	}
}

class Grouping : Expression
{
	public Grouping(Expression expression)
	{
		Expression = expression;
	}

	public Expression Expression { get; }

	public override T Accept<T>(IVisitor<T> visitor)
	{
		return visitor.VisitGroupingExpression(this);
	}
}

class Literal : Expression
{
	public Literal(object value)
	{
		Value = value;
	}

	public object Value { get; }

	public override T Accept<T>(IVisitor<T> visitor)
	{
		return visitor.VisitLiteralExpression(this);
	}
}

class Unary : Expression
{
	public Unary(Token @operator, Expression right)
	{
		Operator = @operator;
		Right = right;
	}

	public Token Operator { get; }
	public Expression Right { get; }

	public override T Accept<T>(IVisitor<T> visitor)
	{
		return visitor.VisitUnaryExpression(this);
	}
}

interface IVisitor<T>
{
	T VisitBinaryExpression(Binary expression);
	T VisitGroupingExpression(Grouping expression);
	T VisitLiteralExpression(Literal expression);
	T VisitUnaryExpression(Unary expression);
}
