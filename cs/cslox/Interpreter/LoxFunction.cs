using System;
using cslox.Parser;

namespace cslox.Interpreter;

class LoxFunction : ILoxCallable
{
    private readonly Function declaration;

    public LoxFunction(Function declaration)
    {
        this.declaration = declaration;
    }

    public int Arity()
    {
        return declaration.Parameters.Count;
    }

    public object? Call(Interpreter interpreter, List<object> arguments)
    {
        var environment = new Environment(interpreter.Globals);
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
            return returnValue.Value;
        }

        return null;
    }

    public override string ToString()
    {
        return $"<fn {declaration.Name.Lexeme}>";
    }
}

