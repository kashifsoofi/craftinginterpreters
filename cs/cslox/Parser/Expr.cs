using cslox.Scanning;

namespace cslox.Parser;

interface IExprVisitor<T>
{
	T VisitAssignExpr(Assign expr);
	T VisitBinaryExpr(Binary expr);
	T VisitCallExpr(Call expr);
	T VisitGetExpr(Get expr);
	T VisitGroupingExpr(Grouping expr);
	T VisitLiteralExpr(Literal expr);
	T VisitLogicalExpr(Logical expr);
	T VisitSetExpr(Set expr);
	T VisitUnaryExpr(Unary expr);
	T VisitVariableExpr(Variable expr);
}

abstract class Expr
{
	public abstract T Accept<T>(IExprVisitor<T> visitor);
}

class Assign : Expr
{
	public Assign(Token name, Expr value)
	{
		Name = name;
		Value = value;
	}

	public Token Name { get; }
	public Expr Value { get; }

	public override T Accept<T>(IExprVisitor<T> visitor)
	{
		return visitor.VisitAssignExpr(this);
	}
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

class Call : Expr
{
	public Call(Expr callee, Token paren, List<Expr> arguments)
	{
		Callee = callee;
		Paren = paren;
		Arguments = arguments;
	}

	public Expr Callee { get; }
	public Token Paren { get; }
	public List<Expr> Arguments { get; }

	public override T Accept<T>(IExprVisitor<T> visitor)
	{
		return visitor.VisitCallExpr(this);
	}
}

class Get : Expr
{
	public Get(Expr @object, Token name)
	{
		Object = @object;
		Name = name;
	}

	public Expr Object { get; }
	public Token Name { get; }

	public override T Accept<T>(IExprVisitor<T> visitor)
	{
		return visitor.VisitGetExpr(this);
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

class Logical : Expr
{
	public Logical(Expr left, Token @operator, Expr right)
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
		return visitor.VisitLogicalExpr(this);
	}
}

class Set : Expr
{
	public Set(Expr @object, Token name, Expr value)
	{
		Object = @object;
		Name = name;
		Value = value;
	}

	public Expr Object { get; }
	public Token Name { get; }
	public Expr Value { get; }

	public override T Accept<T>(IExprVisitor<T> visitor)
	{
		return visitor.VisitSetExpr(this);
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
