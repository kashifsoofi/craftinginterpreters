using System;
using Lox.Parser;

namespace Lox.Interpreter.Callables;

class LoxFunction : ILoxCallable
{
    private readonly Function declaration;
    private readonly Environment closure;
    private readonly bool isInitializer;

    public LoxFunction(Function declaration, Environment closure, bool isInitializer)
    {
        this.declaration = declaration;
        this.closure = closure;
        this.isInitializer = isInitializer;
    }

    public int Arity()
    {
        return declaration.Parameters.Count;
    }

    public object? Call(Interpreter interpreter, List<object> arguments)
    {
        var environment = new Environment(closure);
        for (var i = 0; i < declaration.Parameters.Count; i++)
        {
            environment.Define(declaration.Parameters[i].Lexeme, arguments[i]);
        }

        try
        {
            interpreter.ExecuteBlock(declaration.Body, environment);
        }
        catch (ReturnValue returnValue)
        {
            if (isInitializer)
            {
                return closure.GetAt(0, "this");
            }

            return returnValue.Value;
        }

        if (isInitializer)
        {
            return closure.GetAt(0, "this");
        }

        return null;
    }

    public LoxFunction Bind(LoxInstance instance)
    {
        var environment = new Environment(closure);
        environment.Define("this", instance);
        return new LoxFunction(declaration, environment, isInitializer);
    }

    public override string ToString()
    {
        return $"<fn {declaration.Name.Lexeme}>";
    }
}

