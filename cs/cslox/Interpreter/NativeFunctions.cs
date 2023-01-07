using System;

namespace cslox.Interpreter;

class Clock : ILoxCallable
{
    public int Arity()
    {
        return 0;
    }

    public object? Call(Interpreter interpreter, List<object> arguments)
    {
        return (double)DateTimeOffset.Now.ToUnixTimeMilliseconds() / 1000.0;
    }

    public override string ToString()
    {
        return "<native fn>";
    }
}

