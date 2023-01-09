using System;

namespace cslox.Interpreter;

class LoxClass : ILoxCallable
{
	public string Name { get; }

	public LoxClass(string name)
	{
		Name = name;
	}

    public int Arity()
    {
        return 0;
    }

    public object? Call(Interpreter interpreter, List<object> arguments)
    {
        var instance = new LoxInstance(this);
        return instance;
    }

    public override string ToString()
    {
        return Name;
    }
}

