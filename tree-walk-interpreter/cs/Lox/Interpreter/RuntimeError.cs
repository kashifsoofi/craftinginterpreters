using System;
using Lox.Scanner;

namespace Lox.Interpreter;

class RuntimeError : Exception
{
    public Token Token { get; }

    public RuntimeError(Token token, string message)
        : base(message)
    {
        Token = token;
    }
}

