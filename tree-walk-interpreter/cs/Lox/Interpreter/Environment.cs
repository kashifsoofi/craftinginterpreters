using System;
using Lox.Scanner;

namespace Lox.Interpreter;

public class Environment
{
	public Environment? Enclosing { get; }
	private Dictionary<string, object?> values = new Dictionary<string, object?>();

    public Environment() :
		this(null)
	{ }

    public Environment(Environment? enclosing)
	{
		this.Enclosing = enclosing;
	}

	public void Define(string name, object? value)
	{
		values[name] = value;
	}

	public object? Get(Token name)
	{
		if (values.ContainsKey(name.Lexeme))
		{
			return values[name.Lexeme];
		}

		if (Enclosing != null)
		{
			return Enclosing.Get(name);
		}

		throw new RuntimeError(name, $"Undefined variable '{name.Lexeme}'.");
	}

	public object? GetAt(int distance, string name)
	{
		return Ancestor(distance).values[name];
	}

	public void Assign(Token name, object? value)
	{
		if (values.ContainsKey(name.Lexeme))
		{
			values[name.Lexeme] = value;
			return;
		}

		if (Enclosing != null)
		{
			Enclosing.Assign(name, value);
			return;
		}

        throw new RuntimeError(name, $"Undefined variable '{name.Lexeme}'.");
    }

	public void AssignAt(int distance, Token name, object? value)
	{
		Ancestor(distance).values[name.Lexeme] = value;
	}

	private Environment Ancestor(int distance)
	{
		var environment = this;
		for (var i = 0; i < distance; i++)
		{
			environment = environment!.Enclosing;
		}

		return environment!;
	}
}

