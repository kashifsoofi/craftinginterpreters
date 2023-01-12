using System;

namespace Lox.Interpreter;

class LoxClass : ILoxCallable
{
	public string Name { get; }
    private LoxClass? superclass;
    private Dictionary<string, LoxFunction> methods;

	public LoxClass(string name, LoxClass? superclass, Dictionary<string, LoxFunction> methods)
	{
		Name = name;
        this.superclass = superclass;
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

        if (superclass != null)
        {
            return superclass.FindMethod(name);
        }

        return null;
    }

    public override string ToString()
    {
        return Name;
    }
}

