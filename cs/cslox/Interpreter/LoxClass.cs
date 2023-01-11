using System;

namespace cslox.Interpreter;

class LoxClass : ILoxCallable
{
	public string Name { get; }
    private Dictionary<string, LoxFunction> methods;

	public LoxClass(string name, Dictionary<string, LoxFunction> methods)
	{
		Name = name;
        this.methods = methods;
	}

    public int Arity()
    {
        var initializer = FindMethod("init");
        if (initializer == null)
        {
            return 0;
        }
        return initializer.Arity();
    }

    public object? Call(Interpreter interpreter, List<object> arguments)
    {
        var instance = new LoxInstance(this);
        var initializer = FindMethod("init");
        if (initializer != null)
        {
            initializer.Bind(instance).Call(interpreter, arguments);
        }

        return instance;
    }

    public LoxFunction? FindMethod(string name)
    {
        if (methods.ContainsKey(name))
        {
            return methods[name];
        }

        return null;
    }

    public override string ToString()
    {
        return Name;
    }
}

