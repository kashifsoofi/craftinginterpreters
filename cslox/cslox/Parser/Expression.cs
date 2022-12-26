using cslox.Scanning;

namespace cslox.Parser;

abstract class Expression
{
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
}

class Grouping : Expression
{
	public Grouping(Expression expression)
	{
		Expression = expression;
	}

	public Expression Expression { get; }
}

class Literal : Expression
{
	public Literal(object value)
	{
		Value = value;
	}

	public object Value { get; }
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
}
