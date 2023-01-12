using System;
namespace Lox.Interpreter;

class ReturnValue : Exception
{
    public object? Value { get; }

    public ReturnValue(object? value)
    {
        Value = value;
    }
}
