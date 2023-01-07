using System;
namespace cslox.Interpreter;

interface ILoxCallable
{
    object? Call(Interpreter interpreter, List<object> arguments);
    int Arity();
}

