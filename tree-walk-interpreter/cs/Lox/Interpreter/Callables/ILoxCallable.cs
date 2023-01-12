using System;
namespace Lox.Interpreter.Callables;

internal interface ILoxCallable
{
    object? Call(Interpreter interpreter, List<object> arguments);
    int Arity();
}

