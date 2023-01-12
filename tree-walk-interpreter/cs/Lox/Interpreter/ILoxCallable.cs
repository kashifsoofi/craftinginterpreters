using System;
namespace Lox.Interpreter;

interface ILoxCallable
{
    object? Call(Interpreter interpreter, List<object> arguments);
    int Arity();
}

